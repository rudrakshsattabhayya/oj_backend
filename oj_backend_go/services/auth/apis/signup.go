package apis

import "github.com/gin-gonic/gin"

type SignUpParams struct {
	Email string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	MobileNo string `json:"mobile_no" binding:"optional"`
}

func SignUp(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Request recieved!",
	})
}
