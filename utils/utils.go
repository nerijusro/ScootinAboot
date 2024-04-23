package utils

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetClientIdFromRequest(c *gin.Context) (uuid.UUID, error) {
	clientId := c.GetHeader("Client-Id")
	if clientId == "" {
		return uuid.Nil, errors.New("client id is missing")
	}

	clientIdUUID, err := uuid.Parse(clientId)
	if err != nil {
		return uuid.Nil, errors.New("invalid client id")
	}

	return clientIdUUID, nil
}
