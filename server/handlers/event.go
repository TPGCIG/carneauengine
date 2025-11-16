package handlers

import (

	"context"
	"net/http"
	"github.com/gin-gonic/gin"
)

type SummaryEvent struct {
	ID 		int 	`json:"id"`
	Title 		string 	`json:"title"`
	OrganisationID  *int 	`json:"organisation_id"`
	Description 	string 	`json:"description"`
	ImageURL 	string 	`json:"image_url"`
}	

func (h *Handler) GetEvents(c *gin.Context) {
	rows, err := h.DB.Query(context.Background(), "SELECT id, organisation_id, title, description, image_url FROM events")

	if (err != nil) {
		c.JSON(500, gin.H{"error": err.Error()})
		return;
	}
	defer rows.Close()

	var events []SummaryEvent;

	for rows.Next() {
		var e SummaryEvent
		err := rows.Scan(&e.ID, &e.OrganisationID, &e.Title, &e.Description, &e.ImageURL)
		if (err != nil) {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		events = append(events, e)
	}

	if err = rows.Err(); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, events)
}

func (h *Handler) AddEvent(c *gin.Context) {
	// Return JSON response
	c.JSON(http.StatusOK, gin.H{
		"message": "AddEvent",
	})
}

func (h *Handler) UpdateEvent(c *gin.Context) {
	// Return JSON response
	c.JSON(http.StatusOK, gin.H{
		"message": "AddEvent",
	})
}
