package main

import (
	"github.com/gin-gonic/gin"
)

// CheckRegisteredUser function checks if the cookie token is valid
func CheckRegisteredUser(c *gin.Context) {
	_, err := c.Request.Cookie("authorisation_token")
	if err != nil {
		println("Token has expired.")
		c.Abort()
		ShowLoggingForm(c)
		return
	}
}
