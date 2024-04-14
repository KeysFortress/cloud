package middlewhere

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func Authorize() gin.HandlerFunc {
	return func(c *gin.Context) {

		fmt.Println("It works")
		c.Copy().Next()
	}
}

func Authenticate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		fmt.Println("Veryfing token")
		ctx.Next()
	}
}
