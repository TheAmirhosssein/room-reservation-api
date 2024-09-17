package middlewares

import (
	"net/http"

	"github.com/TheAmirhosssein/room-reservation-api/internal/entity"
	"github.com/gin-gonic/gin"
)

func AdminRoleMiddleware(context *gin.Context) {
	role := context.GetString("role")
	if role != entity.AdminRole {
		context.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "forbidden"})
		return
	}
	context.Next()
}

func SupportRoleMiddleware(context *gin.Context) {
	role := context.GetString("role")
	if role != entity.SupportRole {
		context.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "forbidden"})
		return
	}
	context.Next()
}
