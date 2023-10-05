package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.etcd.io/bbolt"
)

func handleCreateSecret(c *gin.Context, db *bbolt.DB) {
	var newSecret Secret
	if err := c.ShouldBindJSON(&newSecret); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	secretID := uuid.New().String()

	// Serialize Secret struct to JSON
	encodedSecret, err := json.Marshal(newSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encode secret"})
		return
	}

	// Store the secret in the database
	err = db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("secrets"))
		return b.Put([]byte(secretID), encodedSecret)
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store secret"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"secret_id": secretID})
}

func handleFetchAllURLs(c *gin.Context, db *bbolt.DB) {
	urlSet := make(map[string]bool)

	err := db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("secrets"))
		if b == nil {
			return fmt.Errorf("Bucket secrets not found!")
		}

		return b.ForEach(func(k, v []byte) error {
			var secret Secret
			if err := json.Unmarshal(v, &secret); err != nil {
				return err
			}

			// Add the URL to the set
			urlSet[secret.URL] = true
			return nil
		})
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch URLs"})
		return
	}

	// Convert the set to a list
	var urls []string
	for url := range urlSet {
		urls = append(urls, url)
	}

	c.JSON(http.StatusOK, gin.H{"urls": urls})
}

func handleFetchSecrets(c *gin.Context, db *bbolt.DB) {
	searchTerm := c.DefaultQuery("term", "")
	specificURL := c.DefaultQuery("url", "")
	var secrets []Secret

	err := db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("secrets"))
		if b == nil {
			return fmt.Errorf("Bucket secrets not found!")
		}

		return b.ForEach(func(k, v []byte) error {
			var secret Secret
			if err := json.Unmarshal(v, &secret); err != nil {
				return err
			}

			// Set the ID field with the UUID from the database
			secret.ID = string(k)

			// Perform text search on URL and username fields
			matchesTerm := searchTerm == "" ||
				strings.Contains(strings.ToLower(secret.URL), strings.ToLower(searchTerm)) ||
				strings.Contains(strings.ToLower(secret.Username), strings.ToLower(searchTerm))

			// Check for specific URL
			matchesURL := specificURL == "" || strings.ToLower(secret.URL) == strings.ToLower(specificURL)

			if matchesTerm && matchesURL {
				secrets = append(secrets, secret)
			}
			return nil
		})
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch secrets"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"secrets": secrets})
}
