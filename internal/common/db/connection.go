package db

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func GetConnection(c *gin.Context) *sqlx.DB {
	dbValue, exists := c.Get("db")
	if !exists {
		return nil
	}
	db, ok := dbValue.(*sqlx.DB)
	if !ok {
		return nil
	}
	return db
}
