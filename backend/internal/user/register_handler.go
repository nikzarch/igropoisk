package user

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type RegisterHandler struct {
	service *Service
}

func NewRegisterHandler(service *Service) *RegisterHandler {
	return &RegisterHandler{service: service}
}

type request struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (r *RegisterHandler) HandleRegistration(c *gin.Context) {
	var req request
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if len(req.Username) < 3 || len(req.Username) > 32 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username must be between 3 and 32 characters long"})
		return
	}

	if ok, err := validatePassword(req.Password); !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := r.service.Register(req.Username, req.Password)
	if err != nil {
		log.Println("Unable to register: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to register"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func validatePassword(pass string) (bool, error) {
	if len(pass) < 8 || len(pass) > 32 {
		return false, errors.New("password must be between 8 and 32 characters long")
	}

	return true, nil
}
