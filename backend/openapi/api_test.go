package openapi

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mabaums/ece461-web/backend/models"
	"github.com/stretchr/testify/assert"
)

func TestConnectTCPSocket(t *testing.T) {
	// Load environment variables
	/*err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")
	instanceHost := os.Getenv("INSTANCE_HOST")

	os.Setenv("DB_USER", dbUser)
	os.Setenv("DB_PASS", dbPass)
	os.Setenv("DB_NAME", dbName)
	os.Setenv("DB_PORT", dbPort)
	os.Setenv("INSTANCE_HOST", instanceHost)

	db, err := connectTCPSocket()
	if err != nil {
		t.Fatalf("Failed to connect to MySQL instance: %v", err)
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		t.Fatalf("Failed to ping MySQL instance: %v", err)
	} */
}

func TestPackageList(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = &http.Request{}
	packageQuery := models.PackageQuery{
		Name:    "Name",
		Version: "Version",
	}
	queries := []models.PackageQuery{packageQuery}
	buf, err := json.Marshal(queries)
	if err != nil {
		log.Fatal("Error jsoning package query")
	}

	c.Request.Body = io.NopCloser(bytes.NewBuffer(buf))
	PackagesList(c)

	assert.Equal(t, w.Code, http.StatusOK, "Error")
}
