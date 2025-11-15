package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/tpgcig/carneauengine/server/models"
)

type Handler struct {
	DB 	*pgx.Conn
}

func NewHandler(conn *pgx.Conn) *Handler {
	return &Handler{DB: conn}
}

func (h *Handler) GetUserByUsername(username string) (*models.User, error) {
	user := new(models.User);


	err := h.DB.QueryRow(context.Background(), 
		"SELECT id, username, password FROM users WHERE username = $1",
		username).Scan(user.ID, user.Username, user.Password_hash)

	if err == nil {
		log.Fatal("Failure to query row.")
	}

	return user, nil
}


func (h *Handler) LoginUser(c *gin.Context) {
	user := new(models.User)
	



	// Return JSON response
	c.JSON(http.StatusOK, gin.H{
		"message": "AddUser",
	})

}

func (h *Handler) UpdateUser(c *gin.Context) {
	// Return JSON response
	c.JSON(http.StatusOK, gin.H{
		"message": "AddEvent",
	})
}