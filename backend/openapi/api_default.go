/*
 * ECE 461 - Spring 2023 - Project 2
 *
 * API for ECE 461/Spring 2023/Project 2: A Trustworthy Module Registry
 *
 * API version: 2.0.0
 * Contact: davisjam@purdue.edu
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

import (
	"database/sql"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"time"

	"github.com/mabaums/ece461-web/backend/datastore"
	"github.com/mabaums/ece461-web/backend/models"

	"github.com/gin-gonic/gin"
)

var db *sql.DB
var ds datastore.InMemoryDatstore

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
	pkg_history := []models.PackageHistoryEntry{
		{
			User: models.User{
				Name:    "James Davis",
				IsAdmin: true,
			},
			Date: time.Now(),
			PackageMetadata: models.PackageMetadata{
				Name:    "Underscore",
				Version: "1.0.0",
				ID:      "underscore",
			},
			Action: "DOWNLOAD",
		},
		{
			User: models.User{
				Name:    "James Davis",
				IsAdmin: true,
			},
			Date: time.Now(),
			PackageMetadata: models.PackageMetadata{
				Name:    "Underscore",
				Version: "1.0.0",
				ID:      "underscore",
			},
			Action: "DOWNLOAD",
		},
	}
	c.JSON(http.StatusOK, pkg_history)
}

// PackageByRegExGet - Get any packages fitting the regular expression.
func PackageByRegExGet(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

// PackageCreate -
func PackageCreate(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

// PackageDelete - Delete this version of the package.
func PackageDelete(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{})
}

// PackageRate -
func PackageRate(c *gin.Context) {
	pkg_rating := models.PackageRating{
		BusFactor:            1,
		RampUp:               1,
		Correctness:          1,
		ResponsiveMaintainer: 1,
		GoodPinningPractice:  1,
		NetScore:             1,
		PullRequest:          1,
		LicenseScore:         1,
	}
	c.JSON(http.StatusOK, pkg_rating)
}

// PackageRetrieve - Interact with the package with this ID
func PackageRetrieve(c *gin.Context) {
	id, _ := c.Params.Get("id")
	pkg, err := ds.GetPackage(id)
	if err != nil {
		c.JSON(int(err.Code), err)
		return
	}
	c.JSON(http.StatusOK, pkg)

}

// PackageUpdate - Update this content of the package.
func PackageUpdate(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

// PackagesList - Get the packages from the registry.
func PackagesList(c *gin.Context) {
	query := []models.PackageQuery{}
	//_, _ := c.Params.Get("offset")
	err := c.BindJSON(&query)
	if err != nil {
		log.Fatal("Query binding error")
	}
	c.JSON(http.StatusOK, ds.ListPackages(0, 0, query[0].Name, query[0].Version))
}

// RegistryReset - Reset the registry
func RegistryReset(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}
