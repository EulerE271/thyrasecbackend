// In "helpers/DBContext.go"
package helpers

import (
	"thyra/internal/common/db" // Import the package with the correct name

	"github.com/gin-gonic/gin"
)

func DBContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the database connection
		database := db.GetDB()

		// Use the database connection as needed
		// For example, you can set it in the context
		c.Set("db", database)
		c.Next()
	}
}
