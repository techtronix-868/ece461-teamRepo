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

// func TestPackageCreate(t *testing.T) {
// 	// Set up a mock database connection
// 	db, mock, err := sqlmock.New()
// 	if err != nil {
// 			t.Fatalf("failed to set up mock database: %s", err)
// 	}
// 	defer db.Close()

// 	// Set up the expected query and result
// 	query := "INSERT INTO packages\\(name, version, content, url, js_program\\) VALUES\\(\\?, \\?, \\?, \\?, \\?\\)"
// 	result := sqlmock.NewResult(1, 1)

// 	// Expect the query to be executed with the given arguments
// 	mock.ExpectExec(query).
// 		WithArgs("package_name", "1.0.0", "package_content", "http://example.com", "js_program").
// 		WillReturnResult(result)

// 	// Call the function
// 	err = PackageCreate("package_name", "1.0.0", "package_content", "http://example.com", "js_program")
// 	if err != nil {
// 		t.Errorf("PackageCreate returned an error: %s", err)
// 	}

// 	// Verify that the query was executed as expected
// 	if err := mock.ExpectationsWereMet(); err != nil {
// 		t.Errorf("failed to meet expectations: %s", err)
// 	}
// }

