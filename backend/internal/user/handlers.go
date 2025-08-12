package user

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

type request struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (r *Handler) HandleRegistration(c *gin.Context) {
	var req request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := validateRequest(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := r.service.Register(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
func validateRequest(req request) error {
	if len(req.Username) < 3 || len(req.Username) > 32 {
		return errors.New("username must be between 3 and 32 characters long")
	}
	if len(req.Password) < 8 || len(req.Password) > 32 {
		return errors.New("password must be between 8 and 32 characters long")
	}
	return nil
}

func (r *Handler) HandleLogin(c *gin.Context) {
	var req request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := validateRequest(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := r.service.Login(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})

}
