package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type CheckoutRequest struct {
	EventID int `json:"event_id"`
	Tickets []struct {
		TicketTypeID int `json:"ticket_type_id"`
		Quantity     int `json:"quantity"`
	} `json:"tickets"`
}

type CheckoutResponse struct {
	URL string `json:"url"`
}

// Response models for GetPurchaseByStripeSessionID
type PurchasedTicketDetail struct {
	ID        int    `json:"id"`
	QRCode    string `json:"qr_code"`
	Status    string `json:"status"`
	TicketTypeName string `json:"ticket_type_name"`
	TicketTypePrice float64 `json:"ticket_type_price"`
}

type PurchaseEventDetails struct {
	ID            int       `json:"id"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	Location      string    `json:"location"`
	StartTime     time.Time `json:"start_time"`
	EndTime       time.Time `json:"end_time"`
	ImageURLs     []string  `json:"image_urls"`
}

type FullPurchaseDetails struct {
	PurchaseID     int                     `json:"purchase_id"`
	TotalAmount    float64                 `json:"total_amount"`
	PaymentStatus  string                  `json:"payment_status"`
	PurchaserEmail string                  `json:"purchaser_email"`
	Event          PurchaseEventDetails    `json:"event"`
	Tickets        []PurchasedTicketDetail `json:"tickets"`
	CreatedAt      time.Time               `json:"created_at"`
}

// GetPurchaseByStripeSessionID retrieves full details of a purchase by its Stripe Session ID.
func (h *Handler) GetPurchaseByStripeSessionID(c *gin.Context) {
	stripeSessionID := c.Param("stripeSessionId")

	var purchase FullPurchaseDetails
	var eventJSON []byte
	var ticketsJSON []byte

	// Query to get purchase, event, and associated tickets
	// This query is complex as it gathers details from multiple tables and aggregates tickets
	query := `
		SELECT
			p.id, p.total_amount, p.payment_status, p.created_at,
			u.email AS purchaser_email,
			JSON_BUILD_OBJECT(
				'id', e.id,
				'title', e.title,
				'description', e.description,
				'location', e.location,
				'start_time', e.start_time,
				'end_time', e.end_time,
				'image_urls', COALESCE(ARRAY_AGG(ei.url) FILTER (WHERE ei.url IS NOT NULL), '{}')
			) AS event_details,
			COALESCE(JSON_AGG(
				JSON_BUILD_OBJECT(
					'id', t.id,
					'qr_code', t.qr_code,
					'status', t.status,
					'ticket_type_name', tt.name,
					'ticket_type_price', tt.price
				)
			) FILTER (WHERE t.id IS NOT NULL), '[]') AS tickets_details
		FROM purchases p
		JOIN users u ON p.user_id = u.id
		JOIN events e ON p.event_id = e.id
		LEFT JOIN event_images ei ON e.id = ei.event_id
		LEFT JOIN tickets t ON p.id = t.purchase_id
		LEFT JOIN ticket_types tt ON t.ticket_type_id = tt.id
		WHERE p.stripe_payment_id = $1
		GROUP BY p.id, u.email, e.id
	`

	err := h.DB.QueryRow(c.Request.Context(), query, stripeSessionID).Scan(
		&purchase.PurchaseID, &purchase.TotalAmount, &purchase.PaymentStatus, &purchase.CreatedAt,
		&purchase.PurchaserEmail,
		&eventJSON,
		&ticketsJSON,
	)

	if err == pgx.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Purchase not found"})
		return
	}
	if err != nil {
		log.Printf("Error querying purchase details: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve purchase details"})
		return
	}

	if err := json.Unmarshal(eventJSON, &purchase.Event); err != nil {
		log.Printf("Error unmarshalling event details: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process event details"})
		return
	}
	if err := json.Unmarshal(ticketsJSON, &purchase.Tickets); err != nil {
		log.Printf("Error unmarshalling tickets details: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process ticket details"})
		return
	}

	c.JSON(http.StatusOK, purchase)
}
