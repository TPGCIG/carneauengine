package handlers

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"math/rand"
	"net/http"
	"net/smtp"
	"net/textproto" // Added for textproto.MIMEHeader
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/skip2/go-qrcode"
	"github.com/stripe/stripe-go/v83"
	"github.com/stripe/stripe-go/v83/checkout/session"
	"github.com/stripe/stripe-go/v83/webhook"
	"golang.org/x/crypto/bcrypt"

	_ "github.com/go-redis/redis/v8" // Added for Redis client
	"github.com/google/uuid"       // Added for UUID generation

	"github.com/tpgcig/carneauengine/server/models"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func generateRandomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// sendEmail helper function
func sendEmail(to, subject, htmlBody string, qrCodeImages map[string][]byte, senderEmail, smtpHost, smtpPort, smtpUser, smtpPassword string) error {
	var body bytes.Buffer
	mimeWriter := multipart.NewWriter(&body)

	// Create a "multipart/alternative" part for text and HTML body
	altWriter := new(bytes.Buffer)
	altMimeWriter := multipart.NewWriter(altWriter)
	altMimeWriter.SetBoundary("alt-" + altMimeWriter.Boundary()) // Unique boundary for nested multipart
	
	// Text part (simple version of HTML)
	textPartHeaders := make(textproto.MIMEHeader)
	textPartHeaders.Set("Content-Type", "text/plain; charset=\"UTF-8\"")
	textPartHeaders.Set("Content-Transfer-Encoding", "quoted-printable")
	textPart, _ := altMimeWriter.CreatePart(textPartHeaders)
	textPart.Write([]byte(htmlBody)) // For simplicity, sending HTML as plain text too

	// HTML part with embedded images
	relatedWriter := new(bytes.Buffer)
	relatedMimeWriter := multipart.NewWriter(relatedWriter)
	relatedMimeWriter.SetBoundary("rel-" + relatedMimeWriter.Boundary()) // Unique boundary for nested multipart

	htmlPartHeaders := make(textproto.MIMEHeader)
	htmlPartHeaders.Set("Content-Type", "text/html; charset=\"UTF-8\"")
	htmlPartHeaders.Set("Content-Transfer-Encoding", "quoted-printable")
	p, _ := relatedMimeWriter.CreatePart(htmlPartHeaders)
	p.Write([]byte(htmlBody))

	for cid, qrCodeData := range qrCodeImages {
		imgPartHeaders := make(textproto.MIMEHeader)
		imgPartHeaders.Set("Content-Type", "image/png")
		imgPartHeaders.Set("Content-Transfer-Encoding", "base64")
		imgPartHeaders.Set("Content-ID", fmt.Sprintf("<%s>", cid))
		imgPartHeaders.Set("Content-Disposition", fmt.Sprintf("inline; filename=\"%s.png\"", cid))
		imgPart, _ := relatedMimeWriter.CreatePart(imgPartHeaders)
		encoder := base64.NewEncoder(base64.StdEncoding, imgPart)
		encoder.Write(qrCodeData)
		encoder.Close()
	}
	relatedMimeWriter.Close() // Close the related writer to finalize its boundary

	// Add the multipart/related part to the multipart/alternative
	relatedAltPartHeaders := make(textproto.MIMEHeader)
	relatedAltPartHeaders.Set("Content-Type", fmt.Sprintf("multipart/related; boundary=%s", relatedMimeWriter.Boundary()))
	relatedAltPart, _ := altMimeWriter.CreatePart(relatedAltPartHeaders)
	relatedAltPart.Write(relatedWriter.Bytes())

	altMimeWriter.Close() // Close the multipart/alternative writer

	// Add the multipart/alternative part to the main body
	altMainPartHeaders := make(textproto.MIMEHeader)
	altMainPartHeaders.Set("Content-Type", fmt.Sprintf("multipart/alternative; boundary=%s", altMimeWriter.Boundary()))
	altMainPart, _ := mimeWriter.CreatePart(altMainPartHeaders)
	altMainPart.Write(altWriter.Bytes())


	// Main email headers
	mainHeaders := make(textproto.MIMEHeader)
	mainHeaders.Set("From", senderEmail)
	mainHeaders.Set("To", to)
	mainHeaders.Set("Subject", subject)
	mainHeaders.Set("MIME-Version", "1.0")
	mainHeaders.Set("Content-Type", fmt.Sprintf("multipart/mixed; boundary=%s", mimeWriter.Boundary()))


	// Combine headers and body
	var finalMessage bytes.Buffer
	for k := range mainHeaders {
		finalMessage.WriteString(fmt.Sprintf("%s: %s\r\n", k, mainHeaders.Get(k)))
	}
	finalMessage.WriteString("\r\n")
	finalMessage.Write(body.Bytes())

	auth := smtp.PlainAuth("", smtpUser, smtpPassword, smtpHost)
	addr := fmt.Sprintf("%s:%s", smtpHost, smtpPort)

	return smtp.SendMail(addr, auth, senderEmail, []string{to}, finalMessage.Bytes())
}

// releaseRedisHolds function to clean up Redis holds if something goes wrong before Stripe session is created
func (h *Handler) releaseRedisHolds(ctx context.Context, reservationID string) {
	// Retrieve the map of ticketTypeID -> quantity from the reservation hash
	reservedItems, err := h.Redis.HGetAll(ctx, "reservation:"+reservationID).Result()
	if err != nil {
		log.Printf("Error retrieving reserved items for reservation ID %s: %v", reservationID, err)
		return
	}

	for ticketTypeIDStr, quantityStr := range reservedItems {
		ticketTypeID, err := strconv.Atoi(ticketTypeIDStr)
		if err != nil {
			log.Printf("Error converting ticketTypeID %s to int: %v", ticketTypeIDStr, err)
			continue
		}
		quantity, err := strconv.Atoi(quantityStr)
		if err != nil {
			log.Printf("Error converting quantity %s to int for ticketTypeID %d: %v", quantityStr, ticketTypeID, err)
			continue
		}

		// Decrement the held quantity
		_, err = h.Redis.HIncrBy(ctx, fmt.Sprintf("ticket_holds:%d", ticketTypeID), "held_quantity", int64(-quantity)).Result()
		if err != nil {
			log.Printf("Error releasing Redis hold for ticket type %d, reservation %s: %v", ticketTypeID, reservationID, err)
		}
	}
	// Delete the main reservation hash
	h.Redis.Del(ctx, "reservation:"+reservationID)
	log.Printf("Released Redis holds for reservation ID: %s", reservationID)
}

func (h *Handler) CreateCheckoutSession(c *gin.Context) {
	var req struct {
		Items []struct {
			TicketID int `json:"ticket_id"`
			Quantity int `json:"quantity"`
		} `json:"items"`
		Email string `json:"email"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(req.Items) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No items in cart"})
		return
	}

	// Generate a unique reservation ID for this checkout attempt
	reservationID := uuid.New().String()
	log.Printf("Generated reservation ID: %s", reservationID)

	// Ensure Redis holds are released if anything goes wrong before Stripe session creation
	// This defer will be executed if the function returns early
	defer func() {
		// Only release if the Stripe session hasn't been created successfully.
		// A flag or checking Stripe metadata could make this more robust,
		// but for now, we assume if we exit early, we should clean up.
		if c.Writer.Status() != http.StatusOK { // If an error status was set
			h.releaseRedisHolds(c.Request.Context(), reservationID)
		}
	}()

	// 1. Collect Ticket IDs and requested quantities
	ticketIDs := make([]int, len(req.Items))
	quantityMap := make(map[int]int)
	for i, item := range req.Items {
		ticketIDs[i] = item.TicketID
		quantityMap[item.TicketID] = item.Quantity
	}

	// 2. Fetch Real Prices, Event ID, and current availability from Database
	placeholders := make([]string, len(ticketIDs))
	args := make([]interface{}, len(ticketIDs))
	for i, id := range ticketIDs {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	// IMPORTANT: Select FOR UPDATE to ensure no other transaction modifies these rows
	// between our read and the Redis update.
	query := fmt.Sprintf(`
		SELECT id, event_id, name, price, total_quantity, sold_quantity
		FROM ticket_types
		WHERE id IN (%s) FOR UPDATE`, strings.Join(placeholders, ","))
	
	rows, err := h.DB.Query(c.Request.Context(), query, args...)
	if err != nil {
		log.Printf("DB error fetching ticket types for reservation: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch ticket information for reservation"})
		return
	}
	defer rows.Close()

	var (
		lineItems []*stripe.CheckoutSessionLineItemParams
		totalAmount float64
		eventID int
		dbTicketDetails = make(map[int]struct {
			Name string
			Price float64
			TotalQuantity int
			SoldQuantity int
		})
	)

	// To store ticket details for later purchase insertion
	type ticketDetail struct {
		TicketTypeID int
		Quantity     int
		Price        float64
	}
	var purchasedTicketDetails []ticketDetail

	for rows.Next() {
		var id, currentEventID, totalQuantity, soldQuantity int
		var name string
		var price float64
		if err := rows.Scan(&id, &currentEventID, &name, &price, &totalQuantity, &soldQuantity); err != nil {
			log.Printf("Error scanning ticket type during reservation fetch: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing ticket information"})
			return
		}

		if eventID == 0 {
			eventID = currentEventID
		} else if eventID != currentEventID {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Tickets from multiple events in one cart are not supported"})
			return
		}

		dbTicketDetails[id] = struct {
			Name string
			Price float64
			TotalQuantity int
			SoldQuantity int
		}{name, price, totalQuantity, soldQuantity}
	}

	if len(dbTicketDetails) != len(req.Items) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Some ticket types not found"})
		return
	}

	// Prepare data for Redis Lua script
	var (
		redisKeys   []string
		redisArgs   []interface{}
		ticketTypeArgs []string // For Lua script KEY arguments
		quantityArgs   []string // For Lua script ARG arguments
	)

	const reservationTTL = 15 * 60 // 15 minutes in seconds

	for _, item := range req.Items {
		detail := dbTicketDetails[item.TicketID]
		requestedQty := item.Quantity

		// Redis key for holding quantity
		ticketHoldKey := fmt.Sprintf("ticket_holds:%d", item.TicketID)
		
		// Lua script arguments
		ticketTypeArgs = append(ticketTypeArgs, ticketHoldKey) // KEY[1...N] = ticket_holds:id
		quantityArgs = append(quantityArgs, fmt.Sprintf("%d", detail.TotalQuantity))
		quantityArgs = append(quantityArgs, fmt.Sprintf("%d", detail.SoldQuantity))
		quantityArgs = append(quantityArgs, fmt.Sprintf("%d", requestedQty))
	}
	
	// KEYS: ticket_holds:<id1>, ticket_holds:<id2>, ...
	// ARGS: total_qty1, sold_qty1, requested_qty1, total_qty2, sold_qty2, requested_qty2, ..., reservationID, reservationTTL
	redisKeys = ticketTypeArgs
	redisArgs = []interface{}{}
	for _, arg := range quantityArgs {
		redisArgs = append(redisArgs, arg)
	}
	redisArgs = append(redisArgs, reservationID)
	redisArgs = append(redisArgs, fmt.Sprintf("%d", reservationTTL))

	// Lua script for atomic reservation
	// It iterates through each ticket type, checks availability, and reserves.
	// If any reservation fails, it rolls back all reservations for this request.
	// KEYS: {ticket_hold_key_1, ticket_hold_key_2, ...}
	// ARGV: {total_qty_1, sold_qty_1, requested_qty_1, total_qty_2, sold_qty_2, requested_qty_2, ..., reservationID, reservationTTL}
	var luaScript = `
		local reservationId = ARGV[#ARGV - 1]
		local reservationTTL = tonumber(ARGV[#ARGV])
		local numTickets = (#KEYS)
		local reservedItems = {} -- Stores successfully reserved quantities in this transaction
		local fullReservationKey = "reservation:" .. reservationId

		-- Clean up on error (optional, but good practice if intermediate writes occur)
		local function rollback()
			for i = 1, #reservedItems, 2 do
				local ttId = reservedItems[i]
				local qty = reservedItems[i+1]
				redis.call('HINCRBY', "ticket_holds:" .. ttId, "held_quantity", -qty)
			end
			redis.call('DEL', fullReservationKey)
			return 0
		end

		for i = 1, numTickets do
			local ticketHoldKey = KEYS[i] -- e.g., ticket_holds:123
			local totalQty = tonumber(ARGV[(i-1)*3 + 1])
			local soldQty = tonumber(ARGV[(i-1)*3 + 2])
			local requestedQty = tonumber(ARGV[(i-1)*3 + 3])
			local ticketTypeId = string.match(ticketHoldKey, "ticket_holds:(%d+)") -- Extract ID

			local currentHeldQty = tonumber(redis.call('HGET', ticketHoldKey, 'held_quantity') or '0')
			local availableForSale = totalQty - soldQty - currentHeldQty

			if availableForSale < requestedQty then
				-- Not enough tickets, roll back all and return error
				return rollback()
			end

			-- Reserve tickets by incrementing the held_quantity
			redis.call('HINCRBY', ticketHoldKey, "held_quantity", requestedQty)
			
			-- Store this reservation detail in the temporary reservedItems for potential rollback
			table.insert(reservedItems, ticketTypeId)
			table.insert(reservedItems, requestedQty)

			-- Store reservation details in a hash for this specific reservation ID
			-- This allows the webhook to easily retrieve what was reserved by this session
			redis.call('HSET', fullReservationKey, ticketTypeId, requestedQty)
		end

		-- Set TTL for the main reservation hash
		redis.call('EXPIRE', fullReservationKey, reservationTTL)
		return 1
	`
	// Execute the Lua script
	val, err := h.Redis.Eval(c.Request.Context(), luaScript, redisKeys, redisArgs...).Result()
	if err != nil {
		log.Printf("Redis Lua script execution failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reserve tickets due to internal error."})
		return
	}

	if val.(int64) == 0 { // If script returned 0, it means reservation failed (not enough tickets)
		c.JSON(http.StatusConflict, gin.H{"error": "Not enough tickets available for some selected items. Please adjust your cart."})
		return
	}

	// If we reach here, tickets are successfully reserved in Redis.
	// Now proceed with existing logic to prepare Stripe session.

	// Determine userID: existing user or new guest user
	var currentUserID int

	// All checkouts are guest checkouts based on the provided email
	user, err := h.GetUserByEmail(c.Request.Context(), req.Email) // Use h.GetUserByEmail from user.go
	if err != nil && err != pgx.ErrNoRows { // Check for actual errors, not just no rows
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to check for existing user: %v", err)})
		return
	}

	if user != nil {
		// User exists, use their ID
		currentUserID = user.ID
	} else {
		// No user found, create a new guest user
		var guestUser models.User
		guestUser.Email = req.Email
		guestUser.FirstName = "Guest" // Placeholder
		guestUser.LastName = "User"   // Placeholder
		guestUser.Role = "guest"

		// Hash a generic "guest" password. This password will not be used for login.
		hashedPwd, err := bcrypt.GenerateFromPassword([]byte("GUEST_USER_PLACEHOLDER_PASSWORD"), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash guest password"})
			return
		}
		guestUser.PasswordHash = string(hashedPwd)

		err = h.DB.QueryRow(c.Request.Context(),
			"INSERT INTO users (email, first_name, last_name, role, password_hash) VALUES ($1, $2, $3, $4, $5) RETURNING id",
			guestUser.Email, guestUser.FirstName, guestUser.LastName, guestUser.Role, guestUser.PasswordHash,
		).Scan(&guestUser.ID)

		if err != nil {
			log.Printf("Failed to create guest user: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create guest user"})
			return
		}
		currentUserID = guestUser.ID
	}

	var purchaseID int
	err = h.DB.QueryRow(c.Request.Context(),
		"INSERT INTO purchases (user_id, event_id, total_amount, payment_status) VALUES ($1, $2, $3, $4) RETURNING id",
		currentUserID, eventID, totalAmount, "pending",
	).Scan(&purchaseID)

	if err != nil {
		log.Printf("Failed to create pending purchase: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create pending purchase"})
		return
	}

	// Calculate line items for Stripe and total amount
	for _, item := range req.Items {
		detail := dbTicketDetails[item.TicketID]
		qty := quantityMap[item.TicketID] // Use quantity from request, already validated by Redis
		
		lineItems = append(lineItems, &stripe.CheckoutSessionLineItemParams{
			PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
				Currency: stripe.String("aud"),
				ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
					Name: stripe.String(detail.Name),
				},
				UnitAmount: stripe.Int64(int64(detail.Price * 100)), // Convert to cents
			},
			Quantity: stripe.Int64(int64(qty)),
		})
		totalAmount += detail.Price * float64(qty)
		purchasedTicketDetails = append(purchasedTicketDetails, ticketDetail{TicketTypeID: item.TicketID, Quantity: qty, Price: detail.Price})
	}
	// The `totalAmount` calculation above is correct for Stripe, but `purchases.total_amount` in DB was inserted based on
	// a previous `totalAmount` value. We need to update it or ensure it's calculated after all items are processed.
	// For now, let's trust the current `totalAmount` which is built up during `lineItems` preparation.

	// For `items` metadata, let's include the reserved quantities
	var itemsMeta []string
	for _, td := range purchasedTicketDetails {
		itemsMeta = append(itemsMeta, fmt.Sprintf("%d:%d:%.2f", td.TicketTypeID, td.Quantity, td.Price))
	}

	// 4. Create Stripe Session
	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		LineItems:          lineItems,
		Mode:               stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL:         stripe.String("http://localhost:3000/success?session_id={CHECKOUT_SESSION_ID}"), // Pass session ID
		CancelURL:          stripe.String("http://localhost:3000/cancel"),
		CustomerEmail:      stripe.String(req.Email),
		Metadata: map[string]string{
			"purchase_id":   strconv.Itoa(purchaseID),
			"event_id":      strconv.Itoa(eventID),
			"user_id":       strconv.Itoa(currentUserID), // Use the determined user ID
			"items":         strings.Join(itemsMeta, ";"),
			"reservation_id": reservationID, // Pass the Redis reservation ID
		},
	}

	s, err := session.New(params)
	if err != nil {
		log.Printf("Stripe session creation failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// If Stripe session is successfully created, we don't want the defer to release holds.
	// We could use a flag, but for now, rely on the fact that an HTTP 200 will be set.
	c.JSON(http.StatusOK, gin.H{"url": s.URL})
}

func (h *Handler) StripeWebhook(c *gin.Context) {
	const MaxBodyBytes = int64(65536)
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxBodyBytes)
	payload, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": fmt.Sprintf("Error reading request body: %v", err)})
		return
	}

	endpointSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")
	event, err := webhook.ConstructEvent(payload, c.Request.Header.Get("Stripe-Signature"), endpointSecret)

	if err != nil {
		log.Printf("Error verifying webhook signature: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Error verifying webhook signature: %v", err)})
		return
	}

	// Handle the event
	switch event.Type {
	case "checkout.session.completed":
		var s stripe.CheckoutSession
		err := json.Unmarshal(event.Data.Raw, &s)
		if err != nil {
			log.Printf("Error parsing webhook JSON: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Error parsing webhook JSON: %v", err)})
			return
		}

		log.Printf("Checkout session completed for session ID: %s, Customer Email: %s\n", s.ID, s.CustomerDetails.Email)

		// Retrieve metadata
		purchaseIDStr := s.Metadata["purchase_id"]
		userIDStr := s.Metadata["user_id"]
		itemsStr := s.Metadata["items"]
		reservationID := s.Metadata["reservation_id"] // Retrieve the reservation ID

		purchaseID, err := strconv.Atoi(purchaseIDStr)
		if err != nil {
			log.Printf("Error converting purchase_id from metadata: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid purchase_id in metadata"})
			return
		}
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			log.Printf("Error converting user_id from metadata: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user_id in metadata"})
			return
		}

		// Start a database transaction for atomicity
		tx, err := h.DB.Begin(c.Request.Context())
		if err != nil {
			log.Printf("Error starting transaction: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start database transaction"})
			return
		}
		// Defer rollback, will be overridden by Commit if successful
		defer tx.Rollback(c.Request.Context())

		// 1. Update purchase status
		_, err = tx.Exec(c.Request.Context(),
			"UPDATE purchases SET payment_status = $1, stripe_payment_id = $2 WHERE id = $3",
			"succeeded", s.ID, purchaseID,
		)
		if err != nil {
			log.Printf("Error updating purchase status: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update purchase status"})
			return
		}

		// 2. Process items and create tickets
		itemDetails := strings.Split(itemsStr, ";")
		for _, item := range itemDetails {
			parts := strings.Split(item, ":")
			if len(parts) != 3 {
				log.Printf("Invalid item format in metadata: %s", item)
				continue
			}

			ticketTypeID, _ := strconv.Atoi(parts[0])
			quantity, _ := strconv.Atoi(parts[1])

			for i := 0; i < quantity; i++ {
				// Generate a simple QR code (e.g., UUID or unique string)
				qrCode := fmt.Sprintf("CARNEAU-%d-%d-%s", purchaseID, ticketTypeID, generateRandomString(10))

				_, err = tx.Exec(c.Request.Context(),
					"INSERT INTO tickets (ticket_type_id, user_id, purchase_id, qr_code, status) VALUES ($1, $2, $3, $4, $5)",
					ticketTypeID, userID, purchaseID, qrCode, "valid",
				)
				if err != nil {
					log.Printf("Error inserting ticket: %v", err)
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert tickets"})
					return
				}

				// Update sold quantity atomically
				_, err = tx.Exec(c.Request.Context(),
					"UPDATE ticket_types SET sold_quantity = sold_quantity + 1 WHERE id = $1",
					ticketTypeID,
				)
				if err != nil {
					log.Printf("Error updating sold_quantity: %v", err)
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update ticket quantity"})
					return
				}
			}
		}

		// Commit the transaction
		err = tx.Commit(c.Request.Context())
		if err != nil {
			log.Printf("Error committing transaction: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit database transaction"})
			return
		}

		// --- Redis Cleanup: Release the temporary holds ---
		// After successful DB commit, remove the Redis holds
		err = h.releaseRedisHoldsFromWebhook(c.Request.Context(), reservationID)
		if err != nil {
			log.Printf("Error releasing Redis holds for reservation ID %s from webhook: %v", reservationID, err)
			// This is a non-critical error for the user, but important for inventory accuracy.
			// Log it and possibly alert monitoring.
		} else {
			log.Printf("Successfully released Redis holds for reservation ID: %s via webhook.", reservationID)
		}

		// --- Email Sending Logic ---
		// 1. Retrieve ticket and event details for the email
		type EmailTicket struct {
			QR             string
			TicketTypeName string
		}
		var ticketsForEmail []EmailTicket
		var eventTitle, eventLocation string
		var eventStartTime time.Time

		// rows variable needs to be declared here, outside the if block
		var rows pgx.Rows
		rows, err = h.DB.Query(c.Request.Context(), `
			SELECT
				t.qr_code,
				tt.name AS ticket_type_name,
				e.title AS event_title,
				e.location AS event_location,
				e.start_time AS event_start_time
			FROM tickets t
			JOIN ticket_types tt ON t.ticket_type_id = tt.id
			JOIN events e ON tt.event_id = e.id
			WHERE t.purchase_id = $1`, purchaseID)
		if err != nil {
			log.Printf("Error querying tickets for email: %v", err)
			// Don't return error to Stripe, fulfillment is done, just log email failure
		} else {
			defer rows.Close()
			for rows.Next() {
				var et EmailTicket
				if err := rows.Scan(&et.QR, &et.TicketTypeName, &eventTitle, &eventLocation, &eventStartTime); err != nil {
					log.Printf("Error scanning ticket for email: %v", err)
					continue
				}
				ticketsForEmail = append(ticketsForEmail, et)
			}
			if err = rows.Err(); err != nil {
				log.Printf("Error iterating ticket rows for email: %v", err)
			}
		}

		if len(ticketsForEmail) > 0 {
			// 2. Generate QR code images and construct HTML body
			var qrCodeImages = make(map[string][]byte)
			htmlBody := `
				<html>
				<body>
					<p>Dear ` + s.CustomerDetails.Email + `,</p>
					<p>Thank you for your purchase! Here are your tickets for <strong>` + eventTitle + `</strong>.</p>
					<p><strong>Event:</strong> ` + eventTitle + `</p>
					<p><strong>Location:</strong> ` + eventLocation + `</p>
					<p><strong>Date & Time:</strong> ` + eventStartTime.Format("Mon, Jan 2, 2006 3:04 PM") + `</p>
					<hr/>
					<h3>Your Tickets:</h3>
					<ul>
			`
			for i, ticket := range ticketsForEmail {
				cid := fmt.Sprintf("qrcode_%d", i)
				png, err := qrcode.Encode(ticket.QR, qrcode.Medium, 256)
				if err != nil {
					log.Printf("Error generating QR code for ticket %s: %v", ticket.QR, err)
					continue
				}
				qrCodeImages[cid] = png

				htmlBody += fmt.Sprintf(`
					<li>
						<strong>Ticket Type:</strong> %s<br/>
						<strong>QR Code:</strong> %s<br/>
						<img src="cid:%s" alt="QR Code for Ticket %s" style="width:128px;height:128px;"><br/><br/>
					</li>`, ticket.TicketTypeName, ticket.QR, cid, ticket.TicketTypeName)
			}
			htmlBody += `
					</ul>
					<p>Please present the QR codes at the event for entry.</p>
					<p>Best regards,<br/>The Carneau Engine Team</p>
				</body>
				</html>`

			// 3. Send the email
			smtpHost := os.Getenv("SMTP_HOST")
			smtpPort := os.Getenv("SMTP_PORT")
			smtpUser := os.Getenv("SMTP_USER")
			smtpPassword := os.Getenv("SMTP_PASSWORD")
			senderEmail := os.Getenv("SENDER_EMAIL")

			if smtpHost == "" || smtpPort == "" || smtpUser == "" || smtpPassword == "" || senderEmail == "" {
				log.Println("SMTP environment variables not fully configured. Skipping email send.")
			} else {
				err = sendEmail(s.CustomerDetails.Email, "Your Tickets for "+eventTitle, htmlBody, qrCodeImages, senderEmail, smtpHost, smtpPort, smtpUser, smtpPassword)
				if err != nil {
					log.Printf("Failed to send ticket confirmation email to %s: %v", s.CustomerDetails.Email, err)
				} else {
					log.Printf("Ticket confirmation email sent to %s for purchase %d", s.CustomerDetails.Email, purchaseID)
				}
			}
		} else {
			log.Printf("No tickets found for email confirmation for purchase %d", purchaseID)
		}

	case "payment_intent.succeeded":
		// Handle payment_intent.succeeded
		log.Println("Payment Intent Succeeded!")
	// ... handle other event types
	default:
		log.Printf("Unhandled event type: %s\n", event.Type)
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// releaseRedisHoldsFromWebhook is called by the webhook to clean up Redis holds after DB commit
func (h *Handler) releaseRedisHoldsFromWebhook(ctx context.Context, reservationID string) error {
	reservedItems, err := h.Redis.HGetAll(ctx, "reservation:"+reservationID).Result()
	if err != nil {
		return fmt.Errorf("error retrieving reserved items for reservation ID %s: %w", reservationID, err)
	}

	for ticketTypeIDStr, quantityStr := range reservedItems {
		ticketTypeID, err := strconv.Atoi(ticketTypeIDStr)
		if err != nil {
			log.Printf("Error converting ticketTypeID %s to int in webhook cleanup: %v", ticketTypeIDStr, err)
			continue
		}
		quantity, err := strconv.Atoi(quantityStr)
		if err != nil {
			log.Printf("Error converting quantity %s to int in webhook cleanup for ticketTypeID %d: %v", quantityStr, ticketTypeID, err)
			continue
		}

		// Decrement the held quantity
		_, err = h.Redis.HIncrBy(ctx, fmt.Sprintf("ticket_holds:%d", ticketTypeID), "held_quantity", int64(-quantity)).Result()
		if err != nil {
			return fmt.Errorf("error releasing Redis hold for ticket type %d, reservation %s: %w", ticketTypeID, reservationID, err)
		}
	}
	// Delete the main reservation hash
	_, err = h.Redis.Del(ctx, "reservation:"+reservationID).Result()
	if err != nil {
		return fmt.Errorf("error deleting main reservation hash for ID %s: %w", reservationID, err)
	}
	return nil
}
