package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type SummaryEvent struct {
	ID 		int 	`json:"id"`
	Title 		string 	`json:"title"`
	OrganisationName  string 	`json:"organisation_name"`
	Description 	string 	`json:"description"`
	ImageURL 	string 	`json:"image_url"`
}	

const GetEventByID = `
SELECT 
    e.id,
    o.name,
    e.title,
    e.description,
    e.location,
    e.start_time,
    e.end_time,
    e.total_capacity,
    ARRAY_AGG(i.url) AS image_urls,
    COALESCE(
        JSON_AGG(
            JSON_BUILD_OBJECT(
                'id', t.id,
                'name', t.name,
                'price', t.price,
                'total_quantity', t.total_quantity
            )
        ) FILTER (WHERE t.id IS NOT NULL),
        '[]'
    ) AS ticket_types
FROM events e
JOIN organisations o ON e.organisation_id = o.id
LEFT JOIN event_images i ON e.id = i.event_id
LEFT JOIN ticket_types t ON e.id = t.event_id
WHERE e.id = $1
GROUP BY e.id, o.name, e.title, e.description, e.location, e.start_time, e.end_time, e.total_capacity
`

func (h *Handler) GetSummarisedEvents(c *gin.Context) {
	rows, err := h.DB.Query(context.Background(), "SELECT e.id, e.title, o.name, e.description, i.url FROM events e JOIN organisations o ON e.organisation_id = o.id RIGHT JOIN event_images i ON e.id = i.id")

	if (err != nil) {
		c.JSON(500, gin.H{"error": err.Error()})
		return;
	}
	defer rows.Close()

	var events []SummaryEvent;

	for rows.Next() {
		var e SummaryEvent
		err := rows.Scan(&e.ID, &e.Title, &e.OrganisationName, &e.Description, &e.ImageURL)
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

type TicketType struct {
	ID	int	`json:"id"`
	Name	string	`json:"name"`
	Price	float32	`json:"price"`
	TotalQuantity	int	`json:"total_quantity"`

}

type Event struct {
	ID 		int 	`json:"id"`
	OrganisationName  string 	`json:"organisation_name"`
	Title 		string 	`json:"title"`
	Description 	string 	`json:"description"`
	Location 	string 	`json:"location"`
	StartTime time.Time 	 `json:"start_time"`
	EndTime time.Time 		`json:"end_time"`
	TotalCapacity int 	`json:"total_capacity"`
	ImageURLs 	[]string 	`json:"image_urls"`
	TicketTypes []TicketType `json:"ticket_types"`

}

func (h *Handler) GetEvent(c *gin.Context) {
	id := c.Param("id")
	var e Event
	var ticketTypesJSON []byte

	err := h.DB.QueryRow(context.Background(), GetEventByID, id).Scan(
		&e.ID, &e.OrganisationName, &e.Title, &e.Description, &e.Location,
		&e.StartTime, &e.EndTime, &e.TotalCapacity, &e.ImageURLs, &ticketTypesJSON,
	)

	if (err != nil) {
		c.JSON(500, gin.H{"error": err.Error()})
		return;
	}

	if err := json.Unmarshal(ticketTypesJSON, &e.TicketTypes); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(200, e)
}

func (h *Handler) UpdateEvent(c *gin.Context) {
	// Return JSON response
	c.JSON(http.StatusOK, gin.H{
		"message": "AddEvent",
	})
}
