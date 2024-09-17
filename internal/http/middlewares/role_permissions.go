package middlewares

import (
	"fmt"
	"net/http"

	"github.com/TheAmirhosssein/room-reservation-api/internal/entity"
	"github.com/gin-gonic/gin"
)

func SupportOrAdminRoleMiddleware(context *gin.Context) {
	role := context.GetString("role")
	fmt.Println(role)
	if !(role == entity.SupportRole || role == entity.AdminRole) {
		context.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "forbidden"})
		return
	}
	context.Next()
}
