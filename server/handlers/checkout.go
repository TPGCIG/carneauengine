package handlers

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

// func (h *Handler) CreateCheckoutSession(c *gin.Context) {
// 	var req CheckoutRequest;

// 	if err := c.BindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"});
// 		return
// 	}

// 	c.JSON(200, )

// }