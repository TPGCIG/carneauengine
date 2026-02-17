package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	redis "github.com/go-redis/redis/v8"
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
	cacheKey := "cache:events:summary"
	
	// 1. Try to get from Redis cache
	cachedEvents, err := h.Redis.Get(c.Request.Context(), cacheKey).Result()
	if err == nil {
		var events []SummaryEvent
		if err := json.Unmarshal([]byte(cachedEvents), &events); err == nil {
			log.Printf("Cache HIT for %s", cacheKey) // Log cache hit
			c.JSON(http.StatusOK, events)
			return
		}
		log.Printf("Error unmarshalling cached events: %v", err) // Log but continue to DB
	} else if err == redis.Nil {
		log.Printf("Cache MISS for %s", cacheKey) // Log cache miss
		// Proceed to DB query
	} else {
		log.Printf("Redis cache error for %s: %v", cacheKey, err) // Log other Redis errors
	}

	// 2. Get from DB
	rows, err := h.DB.Query(context.Background(), "SELECT e.id, e.title, o.name, e.description, i.url FROM events e JOIN organisations o ON e.organisation_id = o.id RIGHT JOIN event_images i ON e.id = i.id")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var events []SummaryEvent

	for rows.Next() {
		var e SummaryEvent
		err := rows.Scan(&e.ID, &e.Title, &e.OrganisationName, &e.Description, &e.ImageURL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		events = append(events, e)
	}

	if err = rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 3. Cache the result in Redis for 5 minutes
	eventsJSON, err := json.Marshal(events)
	if err != nil {
		log.Printf("Error marshalling events for cache: %v", err) // Log but don't fail request
	} else {
		if err := h.Redis.Set(c.Request.Context(), cacheKey, eventsJSON, 5*time.Minute).Err(); err != nil {
			log.Printf("Error setting events cache: %v", err) // Log but don't fail request
		}
	}

	c.JSON(http.StatusOK, events)
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
	cacheKey := fmt.Sprintf("cache:event:%s", id)

	// 1. Try to get from Redis cache
	cachedEvent, err := h.Redis.Get(c.Request.Context(), cacheKey).Result()
	if err == nil {
		var e Event
		if err := json.Unmarshal([]byte(cachedEvent), &e); err == nil {
			log.Printf("Cache HIT for %s", cacheKey) // Log cache hit
			c.JSON(http.StatusOK, e)
			return
		}
		log.Printf("Error unmarshalling cached event %s: %v", id, err) // Log but continue to DB
	} else if err == redis.Nil {
		log.Printf("Cache MISS for %s", cacheKey) // Log cache miss
		// Proceed to DB query
	} else {
		log.Printf("Redis cache error for %s: %v", cacheKey, err) // Log other Redis errors
	}

	// 2. Get from DB
	var e Event
	var ticketTypesJSON []byte

	err = h.DB.QueryRow(context.Background(), GetEventByID, id).Scan(
		&e.ID, &e.OrganisationName, &e.Title, &e.Description, &e.Location,
		&e.StartTime, &e.EndTime, &e.TotalCapacity, &e.ImageURLs, &ticketTypesJSON,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := json.Unmarshal(ticketTypesJSON, &e.TicketTypes); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 3. Cache the result in Redis for 5 minutes
	eventJSON, err := json.Marshal(e)
	if err != nil {
		log.Printf("Error marshalling event %s for cache: %v", id, err) // Log but don't fail request
	} else {
		if err := h.Redis.Set(c.Request.Context(), cacheKey, eventJSON, 5*time.Minute).Err(); err != nil {
			log.Printf("Error setting event %s cache: %v", id, err) // Log but don't fail request
		}
	}
	
	c.JSON(http.StatusOK, e)
}

func (h *Handler) UpdateEvent(c *gin.Context) {
	// Assuming event ID is passed in the URL, e.g., /api/events/:id
	eventID := c.Param("id")

	// Invalidate the cache for the specific event
	if err := h.Redis.Del(c.Request.Context(), fmt.Sprintf("cache:event:%s", eventID)).Err(); err != nil {
		log.Printf("Error invalidating cache for event %s: %v", eventID, err)
	} else {
		log.Printf("Invalidated cache for event %s", eventID)
	}

	// Invalidate the overall summarised events list cache
	if err := h.Redis.Del(c.Request.Context(), "cache:events:summary").Err(); err != nil {
		log.Printf("Error invalidating cache for summarised events: %v", err)
	} else {
		log.Println("Invalidated cache for summarised events")
	}

	// TODO: Add actual event update logic here

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Event %s updated (cache invalidated)", eventID),
	})
}
