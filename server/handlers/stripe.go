package handlers

// func (h *Handler) CreateCheckoutSession(c *gin.Context) {
// 	var req struct {
//             EventID  int    `json:"event_id"`
//             Quantity int    `json:"quantity"`
//             Email    string `json:"email"`
//         }
//         if err := c.BindJSON(&req); err != nil {
//             c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//             return
//         }

//         s, err := session.New(&stripe.CheckoutSessionParams{
//             PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
//             LineItems: []*stripe.CheckoutSessionLineItemParams{
//                 {
//                     PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
//                         Currency: stripe.String("aud"),
//                         ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
//                             Name: stripe.String("Event Ticket"),
//                         },
//                         UnitAmount: stripe.Int64(2000), // $20.00
//                     },
//                     Quantity: stripe.Int64(int64(req.Quantity)),
//                 },
//             },
//             Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
//             SuccessURL: stripe.String("http://localhost:3000/success"),
//             CancelURL:  stripe.String("http://localhost:3000/cancel"),
//             CustomerEmail: stripe.String(req.Email),
//         })
//         if err != nil {
//             c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//             return
//         }

//         c.JSON(http.StatusOK, gin.H{"id": s.ID})
// }