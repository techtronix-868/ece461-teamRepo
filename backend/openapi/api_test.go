package openapi

import (
	"os"
	"testing"
	_ "github.com/go-sql-driver/mysql"
)

func TestConnectTCPSocket(t *testing.T) {
	// Set environment variables for testing
	os.Setenv("DB_USER", "test")
	os.Setenv("DB_PASS", "QmKw&vcFUZH2")
	os.Setenv("DB_NAME", "packagedir")
	os.Setenv("DB_PORT", "3306")
	os.Setenv("INSTANCE_HOST", "localhost")

	db, err := connectTCPSocket()
	if err != nil {
		t.Fatalf("Failed to connect to MySQL instance: %v", err)
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		t.Fatalf("Failed to ping MySQL instance: %v", err)
	}
}