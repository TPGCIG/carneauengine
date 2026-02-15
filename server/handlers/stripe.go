package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v83"
	"github.com/stripe/stripe-go/v83/checkout/session"
)

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

	// 1. Collect Ticket IDs
	ticketIDs := make([]int, len(req.Items))
	quantityMap := make(map[int]int)
	for i, item := range req.Items {
		ticketIDs[i] = item.TicketID
		quantityMap[item.TicketID] = item.Quantity
	}

	// 2. Fetch Real Prices from Database
	placeholders := make([]string, len(ticketIDs))
	args := make([]interface{}, len(ticketIDs))
	for i, id := range ticketIDs {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	query := fmt.Sprintf("SELECT id, name, price FROM ticket_types WHERE id IN (%s)", strings.Join(placeholders, ","))
	rows, err := h.DB.Query(context.Background(), query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch ticket prices"})
		return
	}
	defer rows.Close()

	// 3. Build Stripe Line Items
	var lineItems []*stripe.CheckoutSessionLineItemParams

	for rows.Next() {
		var id int
		var name string
		var price float64
		if err := rows.Scan(&id, &name, &price); err != nil {
			continue
		}

		qty := quantityMap[id]
		if qty > 0 {
			lineItems = append(lineItems, &stripe.CheckoutSessionLineItemParams{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String("aud"),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String(name),
					},
					UnitAmount: stripe.Int64(int64(price * 100)), // Convert to cents
				},
				Quantity: stripe.Int64(int64(qty)),
			})
		}
	}

	if len(lineItems) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid items or prices"})
		return
	}

	// 4. Create Stripe Session
	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		LineItems:          lineItems,
		Mode:               stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL:         stripe.String("http://localhost:3000/success"), // Create this page
		CancelURL:          stripe.String("http://localhost:3000/cancel"),  // Create this page
		CustomerEmail:      stripe.String(req.Email),
	}

	s, err := session.New(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": s.URL})
}