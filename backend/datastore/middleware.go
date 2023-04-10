package api

func VerifyConnection(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
			err := db.Ping()
			if err != nil {
					err = db.Close()
					if err != nil {
							c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to close database connection"})
							return
					}

					db, err = sql.Open("mysql", "<username>:<password>@tcp(<host>:<port>)/<dbname>")
					if err != nil {
							c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to reconnect to database"})
							return
					}

					// Try the ping again
					err = db.Ping()
					if err != nil {
							c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to ping database"})
							return
					}
			}

			// Store the database connection in the context
			c.Set("db", db)

			// Continue handling the request
			c.Next()
	}
}
