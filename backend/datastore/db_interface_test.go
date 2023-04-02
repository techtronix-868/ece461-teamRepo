package datastore

import (
	"testing"
	"os"
	"log"
	// "github.com/DATA-DOG/go-sqlmock"
	// _ "github.com/go-sql-driver/mysql"
  "github.com/joho/godotenv"
	
)

func TestConnectTCPSocket(t *testing.T) {
	// Load environment variables
	err := godotenv.Load()

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
	} 
}

func TestPackageCreate(t *testing.T) {
	
	// Load environment variables
	err := godotenv.Load()

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

	initDB()

	// Call the function
	err = PackageCreate("package_name", "1.0.0", "package_content", "http://example.com", "js_program")
	if err != nil {
		t.Errorf("PackageCreate returned an error: %s", err)
	}
}

