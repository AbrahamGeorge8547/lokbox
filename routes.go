package main

import (
	"github.com/gin-gonic/gin"
	"go.etcd.io/bbolt"
)

func setupRouter(db *bbolt.DB) *gin.Engine {
	r := gin.Default()

	r.POST("/secret", func(c *gin.Context) {
		handleCreateSecret(c, db)
	})

	r.GET("/secrets", func(c *gin.Context) {
		handleFetchSecrets(c, db)
	})

	r.GET("/urls", func(c *gin.Context) {
		handleFetchAllURLs(c, db)
	})
	return r
}
