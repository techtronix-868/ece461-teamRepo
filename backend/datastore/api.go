package api

import (
	"database/sql"
	// "errors"
	// "log"
	"net/http"
	"strings"
	// "strconv"
	"math/rand"
	"time"
	// "os"
	"github.com/gin-gonic/gin"
	"github.com/mabaums/ece461-web/backend/models"

	_ "github.com/go-sql-driver/mysql"

)

 
// PackageCreate -

// requestBody:
// 	content:
// 		application/json:
// 			schema:
// 				$ref: '#/components/schemas/Package'
// 	required: true
// responses:
// 	"201":
// 		content:
// 			application/json:
// 				schema:
// 					$ref: '#/components/schemas/PackageMetadata'
// 		description: Success. Check the ID in the returned metadata for the official
// 			ID.
// 	"403":
// 		description: Package exists already.
// 	"400":
// 		description: Malformed request.

// Package:
// 	required:
// 	- metadata
// 	- data
// 	type: object
// 	properties:
// 		metadata:
// 			$ref: '#/components/schemas/PackageMetadata'
// 			description: ""
// 		data:
// 			$ref: '#/components/schemas/PackageData'
// 			description: ""
// PackageMetadata:
// 	description: |-
// 		The "Name" and "Version" are used as a unique identifier pair when uploading a package.

// 		The "ID" is used as an internal identifier for interacting with existing packages.
// 	required:
// 	- Name
// 	- Version
// 	- ID
// 	type: object
// 	properties:
// 		Name:
// 			$ref: '#/components/schemas/PackageName'
// 			description: Package name
// 			example: my-package
// 		Version:
// 			description: Package version
// 			type: string
// 			example: 1.2.3
// 		ID:
// 			$ref: '#/components/schemas/PackageID'
// 			description: "Unique ID for use with the /package/{id} endpoint."
// 			example: "123567192081501"


// PackageData:
// 	description: |-
// 		This is a "union" type.
// 		- On package upload, either Content or URL should be set.
// 		- On package update, exactly one field should be set.
// 		- On download, the Content field should be set.
// 	type: object
// 	properties:
// 		Content:
// 			description: |-
// 				Package contents. This is the zip file uploaded by the user. (Encoded as text using a Base64 encoding).

// 				This will be a zipped version of an npm package's GitHub repository, minus the ".git/" directory." It will, for example, include the "package.json" file that can be used to retrieve the project homepage.

// 				See https://docs.npmjs.com/cli/v7/configuring-npm/package-json#homepage.
// 			type: string
// 		URL:
// 			description: Package URL (for use in public ingest).
// 			type: string
// 		JSProgram:
// 			description: A JavaScript program (for use with sensitive modules).
// 			type: string
  

// NOT AN API ENDPOINT
func getDB(c *gin.Context) (*sql.DB, bool) {
	db_i, ok := c.Get("db")
	if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, models.Error{Code: 500, Message: "Server error"})
			return nil, false
	}
	db, ok := db_i.(*sql.DB)
	if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, models.Error{Code: 500, Message: "Server error"})
			return nil, false
	}
	return db, true
}


func PackageCreate(c *gin.Context) {
	var pkg models.Package
	if err := c.ShouldBindJSON(&pkg); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
	}

	db, ok := getDB(c)
	if !ok {
		return
	}

	// Insert PackageMetadata
	metadata := pkg.Metadata
	paramID := strings.TrimLeft(c.Param("id"), "/")

	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM PackageMetadata WHERE PackageID = ?", paramID).Scan(&count)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to query database"})
		return
	}

	// Generate new ID if package ID already exists or if the id is not specified
	if count > 0 || paramID == "" {
		for {
			rand.Seed(time.Now().UnixNano())
			const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
			b := make([]byte, 6)
			for i := range b {
				b[i] = chars[rand.Intn(len(chars))]
			}
			newID := string(b)

			err = db.QueryRow("SELECT COUNT(*) FROM PackageMetadata WHERE PackageID = ?", newID).Scan(&count)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, models.Error{Code: 500, Message: "Failed to check if package ID exists"})
				return
			}
			if count == 0 {
				paramID = newID
				break
			}
		}
	}

	result, err := db.Exec("INSERT INTO PackageMetadata (Name, Version, PackageID) VALUES (?, ?, ?)", metadata.Name, metadata.Version, paramID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	metadataID, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Insert PackageData
	data := pkg.Data
	result, err = db.Exec("INSERT INTO PackageData (Content, URL, JSProgram) VALUES (?, ?, ?)", data.Content, data.URL, data.JSProgram)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	dataID, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Insert Package
	result, err = db.Exec("INSERT INTO Package (metadata_id, data_id) VALUES (?, ?)", metadataID, dataID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, models.PackageMetadata{Name: metadata.Name, Version: metadata.Version, ID: paramID})
}


// PackageUpdate - Update this content of the package.
func PackageUpdate(c *gin.Context) {
	var pkg models.Package
	if err := c.ShouldBindJSON(&pkg); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
	}

	db, ok := getDB(c)
	if !ok {
		return
	}

	metadata := pkg.Metadata
	var existingPackage models.Package
	var package_id int;

	err := db.QueryRow("SELECT p.id, pmd.Name, pmd.Version, pmd.PackageID "+
			"FROM Package p "+
			"JOIN PackageMetadata pmd ON p.metadata_id = pmd.id "+
			"WHERE pmd.Name = ? AND pmd.Version = ? AND pmd.PackageID = ?",
			metadata.Name, metadata.Version, metadata.ID).Scan(
			&package_id, &existingPackage.Metadata.Name, &existingPackage.Metadata.Version, &existingPackage.Metadata.ID,
	)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Package not found1", "name": metadata.Name, "version":metadata.Version, "id":metadata.ID})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Update package data
	packageData := pkg.Data
	_, err = db.Exec("UPDATE PackageData pd SET Content = ?, URL = ?, JSProgram = ? WHERE pd.id = ? ",
		packageData.Content, packageData.URL, packageData.JSProgram, package_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Package updated"})
}	


// PackageDelete - Delete this version of the package. given packageid
func PackageDelete(c *gin.Context) {
	db, ok := getDB(c)
	if !ok {
		return
	}

	packageID := strings.TrimLeft(c.Param("id"), "/")
	var metadataID int
	err := db.QueryRow("SELECT id FROM PackageMetadata WHERE PackageID = ?", packageID).Scan(&metadataID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error1"})
		return
	}
	_, err = db.Exec("DELETE p, pmd, pd, ph FROM Package p "+
		"LEFT JOIN PackageMetaData pmd ON p.metadata_id = pmd.id "+
		"LEFT JOIN PackageData pd ON p.data_id = pd.id "+
		"LEFT JOIN PackageHistoryEntry ph ON p.metadata_id = ph.package_metadata_id "+
		"WHERE p.metadata_id = ?", metadataID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error2"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Package deleted"})
}


// PackageByNameDelete - Delete all versions of this package. return string of package name
func PackageByNameDelete(c *gin.Context) {
	db, ok := getDB(c)
	if !ok {
		return
	}

	packageName := strings.TrimLeft(c.Param("name"), "/")
	var metadataID int
	err := db.QueryRow("SELECT id FROM PackageMetadata WHERE Name = ?", packageName).Scan(&metadataID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error1"})
		return 
	}

	_, err = db.Exec("DELETE p, pmd, pd, ph FROM Package p "+
		"LEFT JOIN PackageMetaData pmd ON p.metadata_id = pmd.id "+
		"LEFT JOIN PackageData pd ON p.data_id = pd.id "+
		"LEFT JOIN PackageHistoryEntry ph ON p.metadata_id = ph.package_metadata_id "+
		"WHERE p.metadata_id = ?", metadataID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error1"})
		return 
	}

	c.JSON(http.StatusOK, gin.H{"message": "Package deleted"})
}


// PackageRetrieve - Interact with the package with this ID
func PackageRetrieve(c *gin.Context) {
	db, ok := getDB(c)
	if !ok {
		return
	}

	packageID := strings.TrimLeft(c.Param("id"), "/")

	var packageName, packageVersion, packageContent, packageURL, packageJSProgram string

	err := db.QueryRow(`
			SELECT m.Name, m.Version, d.Content, d.URL, d.JSProgram
			FROM Package p
			INNER JOIN PackageMetadata m ON p.metadata_id = m.id
			INNER JOIN PackageData d ON p.data_id = d.id
			WHERE m.PackageID = ?;
	`, packageID).Scan(&packageName, &packageVersion, &packageContent, &packageURL, &packageJSProgram)

	if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
	}

	c.JSON(http.StatusOK, gin.H{
			"package_id": packageID,
			"name": packageName,
			"version": packageVersion,
			"content": packageContent,
			"url": packageURL,
			"js_program": packageJSProgram,
	})
}


// RegistryReset - Reset the registry
func RegistryReset(c *gin.Context) {
	db, ok := getDB(c)
	if !ok {
		return
	}

	// Delete all data from Package table
	_, err := db.Exec("DELETE FROM Package")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Delete all data from PackageData table
	_, err = db.Exec("DELETE FROM PackageData")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Delete all data from PackageHistoryEntry table
	_, err = db.Exec("DELETE FROM PackageHistoryEntry")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Delete all data from PackageMetadata table
	_, err = db.Exec("DELETE FROM PackageMetadata")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Registry reset complete"})
}


// PacakgeByNameGet - 
func PackageByNameGet(c *gin.Context) {
	// Return the history of this package (all versions).
	db, ok := getDB(c)
	if !ok {
		return
	}

	// Get name from query parameter
	name := strings.TrimLeft(c.Param("name"), "/")

	// Query for all packages with matching name
	rows, err := db.Query("SELECT pmd.Name, pmd.Version, pmd.PackageID, phe.user, phe.date, phe.action " +
		"FROM Package p " +
		"JOIN PackageMetadata pmd ON p.metadata_id = pmd.id " +
		"JOIN PackageHistoryEntry phe ON p.metadata_id = phe.package_metadata_id " +
		"WHERE pmd.Name = ?", name)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	// Store results in a slice of PackageHistoryEntry structs
	packageHistoryEntries := make([]models.PackageHistoryEntry, 0)
	for rows.Next() {
		var packageHistoryEntry models.PackageHistoryEntry
		var packageMetadata models.PackageMetadata

		err := rows.Scan(&packageMetadata.Name, &packageMetadata.Version, &packageMetadata.ID,
			&packageHistoryEntry.User, &packageHistoryEntry.Date, &packageHistoryEntry.Action)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		packageMetadata.ID = "" // Set ID to "" since it's not needed
		packageHistoryEntry.PackageMetadata = packageMetadata
		packageHistoryEntries = append(packageHistoryEntries, packageHistoryEntry)
	}

	// Check for errors during iteration
	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, packageHistoryEntries)
}




func PackagesList(c *gin.Context) {
	db, ok := getDB(c)
	if !ok {
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}


// CreateAuthToken -
func CreateAuthToken(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	// Check if the username and password are valid
	if !isValidUser(username, password) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{Message: "Invalid username or password"})
			return
	}

	// Generate a new authentication token
	token, err := generateToken(username)
	if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: "Server error"})
			return
	}

	// Return the authentication token to the client
	c.JSON(http.StatusOK, gin.H{"token": token})
}













// PackageByRegExGet - Get any packages fitting the regular expression.

func PackageByRegExGet(c *gin.Context) {
	packageNamePattern := c.Param("pattern")

	db, ok := c.Get("db")
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: "Server error"})
		return
	}

	rows, err := db.(*sql.DB).Query("SELECT id, metadata_id, data_id, rating_id FROM packages WHERE name REGEXP ?", packageNamePattern)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: "Server error"})
		return
	}
	defer rows.Close()

	pkgs := []Package{}
	for rows.Next() {
		package_temp := Package{}
		err := rows.Scan(&package_temp.id, &package_temp.metadata_id, &package_temp.data_id)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: "Server error"})
			return
		}
		pkgs = append(pkgs, package_temp)
	}

	c.JSON(http.StatusOK, pkgs)
}






// PackageRate -
func PackageRate(c *gin.Context) {
	var rating Rating

	c.JSON(http.StatusOK, ratings)
}



