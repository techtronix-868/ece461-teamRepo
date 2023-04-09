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
func PackageByNameDelete(name string) error {
	err := db.Ping()
	if err != nil {
		return fmt.Errorf("error verifying database connection: %w", err)
	}

	// Prepare the delete statement
	stmt, err := db.Prepare("DELETE FROM Packages WHERE Name = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the delete statement
	res, err := stmt.Exec(name)
	if err != nil {
		return err
	}

	// Check if any rows were affected by the delete statement
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		// If no rows were deleted, return a custom error
		return fmt.Errorf("no package with name %s found", name)
	}

	return nil
}


// PackageByNameGet -
func PackageByNameGet(name string) {
	// Prepare the select statement
	stmt, err := db.Prepare("SELECT * FROM Packages WHERE Name = ?")
	if err != nil {
			return Package{}, err
	}
	defer stmt.Close()

	// Execute the select statement and scan the result into a Package struct
	var packageInfo Package
	err = stmt.QueryRow(name).Scan(&packageInfo.ID, &packageInfo.Name, &packageInfo.Version, &packageInfo.Description)
	if err != nil {
			return Package{}, err
	}

	return packageInfo, nil
}

// PackageByRegExGet - Get any packages fitting the regular expression.
func PackageByRegExGet(regex string) {
	// Query for packages that match the regular expression
	rows, err := db.Query("SELECT Name FROM Packages WHERE Name REGEXP ?", regex)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Collect the results into a slice of strings
	var packageNames []string
	for rows.Next() {
		var packageName string
		if err := rows.Scan(&packageName); err != nil {
			return nil, err
		}
		packageNames = append(packageNames, packageName)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// If no packages were found, return an error
	if len(packageNames) == 0 {
		return nil, fmt.Errorf("no packages found matching the regular expression '%s'", regex)
	}

	return packageNames, nil
}

// PackageCreate -
func PackageCreate(name string, version string, content string, url string, jsProgram string) error {
	err := db.Ping()
	if err != nil {
		return fmt.Errorf("error verifying database connection: %w", err)
	}

	stmt, err := db.Prepare("INSERT INTO packages(name, version, content, url, jsprogram) VALUES(?, ?, ?, ?, ?)")
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
func PackageDelete(name string, version string) {
	// Prepare the delete statement
	stmt, err := db.Prepare("DELETE FROM Packages WHERE Name = ? AND Version = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the delete statement
	res, err := stmt.Exec(name, version)
	if err != nil {
		return err
	}

	// Check if any rows were affected by the delete statement
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no rows deleted")
	}

	return nil
}

// PackageRate -
func PackageRate(name string) {
	// Prepare the SELECT statement
	stmt, err := db.Prepare("SELECT AVG(Rating) FROM Ratings WHERE PackageName = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// Execute the SELECT statement
	row := stmt.QueryRow(name)

	// Scan the result into a PackageRating struct
	var rating PackageRating
	rating.Name = name
	err = row.Scan(&rating.Rating)
	if err != nil {
		return nil, err
	}

	return &rating, nil
}

// PackageRetrieve - Interact with the package with this ID
func PackageRetrieve(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

// PackageUpdate - Update this content of the package.
func PackageUpdate(name string, version string, content string, url string, jsProgram string) {
	// Prepare the update statement
	stmt, err := db.Prepare("UPDATE Packages SET Content = ?, URL = ?, JSProgram = ? WHERE Name = ? AND Version = ?")
	if err != nil {
			return err
	}
	defer stmt.Close()

	// Execute the update statement
	res, err := stmt.Exec(content, url, jsProgram, name, version)
	if err != nil {
			return err
	}

	// Check if any rows were affected by the update statement
	rowsAffected, err := res.RowsAffected()
	if err != nil {
			return err
	}

	if rowsAffected == 0 {
			// If no rows were updated, return an error
			return fmt.Errorf("no rows updated")
	}

	return nil
}


// PackagesList - Get the packages from the registry.
func PackagesList(c *gin.Context) {
}

// RegistryReset - Reset the registry
func RegistryReset(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}
