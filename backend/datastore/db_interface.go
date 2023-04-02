package datastore

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"

)
 
var db *sql.DB

func initDB() error {
	var err error
	db, err = connectTCPSocket()
	if err != nil {
		return err
	}
	return nil
}

// connectTCPSocket initializes a TCP connection pool for a Cloud SQL
// instance of MySQL.
func connectTCPSocket() (*sql.DB, error) {
	mustGetenv := func(k string) string {
		v := os.Getenv(k)
		if v == "" {
			log.Fatalf("Fatal Error in connect_tcp.go: %s environment variable not set.", k)
		}
		return v
	}
	// Note: Saving credentials in environment variables is convenient, but not
	// secure - consider a more secure solution such as
	// Cloud Secret Manager (https://cloud.google.com/secret-manager) to help
	// keep secrets safe.
	var (
		dbUser    = mustGetenv("DB_USER")       // e.g. 'my-db-user'
		dbPwd     = mustGetenv("DB_PASS")       // e.g. 'my-db-password'
		dbName    = mustGetenv("DB_NAME")       // e.g. 'my-database'
		dbPort    = mustGetenv("DB_PORT")       // e.g. '3306'
		dbTCPHost = mustGetenv("INSTANCE_HOST") // e.g. '127.0.0.1' ('172.17.0.1' if deployed to GAE Flex)
	)

	dbURI := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		dbUser, dbPwd, dbTCPHost, dbPort, dbName)

	// dbPool is the pool of database connections.
	dbPool, err := sql.Open("mysql", dbURI)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %v", err)
	}

	return dbPool, nil
}

// CreateAuthToken -
func CreateAuthToken(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

// PackageByNameDelete - Delete all versions of this package.
func PackageByNameDelete(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

// PackageByNameGet -
func PackageByNameGet(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

// PackageByRegExGet - Get any packages fitting the regular expression.
func PackageByRegExGet(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

// PackageCreate -
func PackageCreate(name string, version string, content string, url string, jsProgram string) error {
	err := db.Ping()
	if err != nil {
		return fmt.Errorf("error verifying database connection: %w", err)
	}

	stmt, err := db.Prepare("INSERT INTO package(name, version, content, url, js_program) VALUES(?, ?, ?, ?, ?)")
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(name, version, content, url, jsProgram)
	if err != nil {
		return fmt.Errorf("failed to insert package: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get ID of newly created package: %w", err)
	}

	log.Printf("Inserted package with ID %d", id)

	return nil
}


// PackageDelete - Delete this version of the package.
func PackageDelete(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

// PackageRate -
func PackageRate(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

// PackageRetrieve - Interact with the package with this ID
func PackageRetrieve(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

// PackageUpdate - Update this content of the package.
func PackageUpdate(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

// PackagesList - Get the packages from the registry.
func PackagesList(c *gin.Context) {
}

// RegistryReset - Reset the registry
func RegistryReset(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}
