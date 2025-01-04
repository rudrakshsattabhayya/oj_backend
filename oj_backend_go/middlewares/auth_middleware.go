package middlewares

import (
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rudrakshsattabhayya/oj_backend_go/helpers"
	models "github.com/rudrakshsattabhayya/oj_backend_go/models/session"
)

func AuthMiddleware(c *gin.Context) {
	var headers models.Headers

	err := c.ShouldBindHeader(&headers)
	if err != nil || headers.Authorization == "" {
		log.Print("Error", "Token is missing")
		c.JSON(401, gin.H{"error": "Token is missing"})
		c.Abort()
		return
	}

	splitArr := strings.Split(headers.Authorization, " ")
	if len(splitArr) != 2 {
		log.Print("Error", "Token is invalid")
		c.JSON(401, gin.H{"error": "Token is invalid"})
		c.Abort()
		return
	}

	JwtToken := splitArr[1]

	res, err := helpers.DecodeSessionToken(JwtToken)

	if err != nil {
		log.Print("Error in DecodeSessionToken:", err.Error())
		c.JSON(400, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	if err = helpers.AuthenticateWithToken(res["performed_by_id"].(string), res["token"].(string)); err != nil {
		log.Print("Error in AuthenticateWithToken:", err.Error())
		c.JSON(400, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	session := models.Session{
		Email:         res["email"].(string),
		PerformedByID: res["performed_by_id"].(string),
	}

	c.Set("session", session)
	log.Print("Success: Session set in context")
	c.Next()
}
