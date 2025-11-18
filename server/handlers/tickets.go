package handlers

import (
	"context"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

type Tickets struct {
	TicketIDs []int `json:"ticketIds"`
}

type TicketCartType struct {
		ID    int     `json:"id"`
		Name  string  `json:"name"`
		Price float64 `json:"price"`
}

func (h *Handler) GetTicketTypes(c *gin.Context) {
	var ticketIds Tickets;
	
	if err := c.BindJSON(&ticketIds); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return;
	}

	placeholders := make([]string, len(ticketIds.TicketIDs))
	args := make([]interface{}, len(ticketIds.TicketIDs))

	for i, id := range ticketIds.TicketIDs {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	query := fmt.Sprintf("SELECT id, name, price FROM ticket_types WHERE id IN (%s)", strings.Join(placeholders, ","))
	rows, err := h.DB.Query(context.Background(), query, args...)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	defer rows.Close()
	
	result := make(map[int]TicketCartType)

	for rows.Next() {
		var ticket TicketCartType
		err := rows.Scan(&ticket.ID, &ticket.Name, &ticket.Price)
		if (err != nil) {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		result[ticket.ID] = ticket
	}

	if err = rows.Err(); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	fmt.Printf("%+v",result);
	c.JSON(200, result)
}