package api

import (
	"database/sql"
	// "errors"
	// "fmt"
	// "log"
	"net/http"
	"strings"
	"strconv"
	// "math/rand"
	// "time"
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
  
					
// func generateID() string {
// 	rand.Seed(time.Now().UnixNano())
// 	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
// 	b := make([]byte, 6)
// 	for i := range b {
// 			b[i] = chars[rand.Intn(len(chars))]
// 	}
// 	return string(b)
// }

	
func PackageCreate(c *gin.Context) {
	var pkg models.Package
	if err := c.ShouldBindJSON(&pkg); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
	}


	db_i, ok := c.Get("db")
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.Error{Code: 500, Message: "Server error"})
		return
	}
	db, ok := db_i.(*sql.DB)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.Error{Code: 500, Message: "Server error"})
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

	if paramID == "" {
		paramID = "1"
	}
	if count > 0 {
		// Generate new ID if package ID already exists or if the id is not specified
		tempID, err := strconv.Atoi(paramID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, models.Error{Code: 500, Message: "Failed to check if package ID exists"})
			return
		}
		newID := strconv.Itoa(tempID)
		
		for {
			tempID, err = strconv.Atoi(newID)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, models.Error{Code: 500, Message: "Failed to check if package ID exists"})
				return
			}
			newID = strconv.Itoa(tempID + 1)
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
	// packageID, err := result.LastInsertId()
	// if err != nil {
	// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 		return
	// }

	c.JSON(http.StatusOK, models.PackageMetadata{Name: metadata.Name, Version: metadata.Version, ID: paramID})
}





// // CreateAuthToken -
// func CreateAuthToken(c *gin.Context) {
// 	username := c.PostForm("username")
// 	password := c.PostForm("password")

// 	// Check if the username and password are valid
// 	if !isValidUser(username, password) {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{Message: "Invalid username or password"})
// 			return
// 	}

// 	// Generate a new authentication token
// 	token, err := generateToken(username)
// 	if err != nil {
// 			c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: "Server error"})
// 			return
// 	}

// 	// Return the authentication token to the client
// 	c.JSON(http.StatusOK, gin.H{"token": token})
// }


// // PackagesList - Get the packages from the registry.
// func PackagesList(c *gin.Context) {
// 	db := c.MustGet("db").(*sql.DB)

// 	rows, err := db.Query("SELECT id, name, version, description FROM packages")
// 	if err != nil {
// 		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: "Server error"})
// 		return
// 	}
// 	defer rows.Close()

// 	packages := []Package{}
// 	for rows.Next() {
// 		var package_temp Package
// 		err = rows.Scan(&package_temp.id, &package_temp.metadata_id, &package_temp.data_id)
// 		if err != nil {
// 			c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: "Server error"})
// 			return
// 		}
// 		packages = append(packages, package_temp)
// 	}

// 	c.JSON(http.StatusOK, packages)
// }






// // PackageByNameDelete - Delete all versions of this package.

// func PackageByNameDelete(c *gin.Context) {
// 	packageName := c.Param("name")

// 	// Get database connection from context
// 	db, ok := c.Get("db")
// 	if !ok {
// 			c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: "Server error"})
// 			return
// 	}

// 	// Delete package from database
// 	stmt, err := db.(*sql.DB).Prepare("DELETE FROM packages WHERE name = ?")
// 	if err != nil {
// 			c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: "Server error"})
// 			return
// 	}
// 	defer stmt.Close()

// 	res, err := stmt.Exec(packageName)
// 	if err != nil {
// 			c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: "Server error"})
// 			return
// 	}

// 	// Check if package was deleted
// 	count, err := res.RowsAffected()
// 	if err != nil {
// 			c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: "Server error"})
// 			return
// 	}
// 	if count == 0 {
// 			c.AbortWithStatusJSON(http.StatusNotFound, ErrorResponse{Message: "Package not found"})
// 			return
// 	}

// 	// Return success response
// 	c.JSON(http.StatusOK, gin.H{"message": "Package deleted"})
// }

// // PacakgeByNameGet - 
// func PackageByNameGet(c *gin.Context) {
// 	// Return the history of this package (all versions).
// 	packageName := c.Param("name")

// 	db, ok := c.Get("db")
// 	if !ok {
// 		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: "Server error"})
// 		return
// 	}

// 	var package_temp Package
// 	// change here to verify with schema
// 	err := db.(*sql.DB).QueryRow("SELECT id, name, version, description FROM packages WHERE name = ?", packageName).Scan(&package_temp.id, &package_temp.metadata_id, &package_temp.data_id)
// 	if err != nil {
// 		if errors.Is(err, sql.ErrNoRows) {
// 			c.AbortWithStatusJSON(http.StatusNotFound, ErrorResponse{Message:"Package not found"})
// 		} else {
// 			c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: "Server error"})
// 		}
// 		return
// 	}
// 	c.JSON(http.StatusOK, package_temp)
// }

// // PackageByRegExGet - Get any packages fitting the regular expression.

// func PackageByRegExGet(c *gin.Context) {
// 	packageNamePattern := c.Param("pattern")

// 	db, ok := c.Get("db")
// 	if !ok {
// 		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: "Server error"})
// 		return
// 	}

// 	rows, err := db.(*sql.DB).Query("SELECT id, metadata_id, data_id, rating_id FROM packages WHERE name REGEXP ?", packageNamePattern)
// 	if err != nil {
// 		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: "Server error"})
// 		return
// 	}
// 	defer rows.Close()

// 	pkgs := []Package{}
// 	for rows.Next() {
// 		package_temp := Package{}
// 		err := rows.Scan(&package_temp.id, &package_temp.metadata_id, &package_temp.data_id)
// 		if err != nil {
// 			c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: "Server error"})
// 			return
// 		}
// 		pkgs = append(pkgs, package_temp)
// 	}

// 	c.JSON(http.StatusOK, pkgs)
// }

// // PackageDelete - Delete this version of the package.
// func PackageDelete(c *gin.Context) {
// 	name := c.Param("name")
// 	version := c.Param("version")

// 	db, ok := c.Get("db")
// 	if !ok {
// 		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: "Server error"})
// 		return
// 	}

// 	result, err := db.(*sql.DB).Exec("DELETE FROM packages WHERE name = ? AND version = ?", name, version)
// 	if err != nil {
// 		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: "Server error"})
// 		return
// 	}

// 	rowsAffected, err := result.RowsAffected()
// 	if err != nil {
// 		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: "Server error"})
// 		return
// 	}

// 	if rowsAffected == 0 {
// 		c.AbortWithStatusJSON(http.StatusNotFound, ErrorResponse{Message: "Package not found"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Package deleted"})
// }

// // PackageRate -
// func PackageRate(c *gin.Context) {
// 	packageID := c.Param("id")

// 	db, ok := c.Get("db")
// 	if !ok {
// 		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: "Server error"})
// 		return
// 	}

// 	var ratings []Rating
// 	rows, err := db.(*sql.DB).Query("SELECT id, package_id, user_id, rating FROM ratings WHERE package_id = ?", packageID)
// 	if err != nil {
// 		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: "Server error"})
// 		return
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		var rating Rating

// 		err = rows.Scan(&rating.id, &rating.bus_factor, &rating.correctness, &rating.ramp_up, &rating.responsive_maintainer, &rating.license_score, &rating.good_pinning_practice, &rating.pull_request, &rating.net_score)
// 		if err != nil {
// 			c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: "Server error"})
// 			return
// 		}
// 		ratings = append(ratings, rating)
// 	}

// 	if len(ratings) == 0 {
// 		c.AbortWithStatusJSON(http.StatusNotFound, ErrorResponse{Message: "Package not found"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, ratings)
// }

// // PackageRetrieve - Interact with the package with this ID
// func PackageRetrieve(c *gin.Context) {
// 	packageID := c.Param("id")

// 	db, ok := c.Get("db")
// 	if !ok {
// 		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: "Server error"})
// 		return
// 	}

// 	var package_temp Package
// 	err := db.(*sql.DB).QueryRow("SELECT id, name, version, content, url, jsprogram FROM packages WHERE id = ?", packageID).Scan(&package_temp.id, &package_temp.metadata_id, &package_temp.data_id)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			c.AbortWithStatusJSON(http.StatusNotFound, ErrorResponse{Message: "Package not found"})
// 		} else {
// 			c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: "Server error"})
// 		}
// 		return
// 	}

// 	c.JSON(http.StatusOK, package_temp)
// }

// // PackageUpdate - Update this content of the package.
// func PackageUpdate(c *gin.Context) {
// 	packageID := c.Param("id")

// 	name := c.PostForm("name")
// 	version := c.PostForm("version")
// 	content := c.PostForm("content")
// 	url := c.PostForm("url")
// 	jsprogram := c.PostForm("jsprogram")

// 	db, ok := c.Get("db")
// 	if !ok {
// 		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: "Server error"})
// 		return
// 	}

// 	result, err := db.(*sql.DB).Exec("UPDATE packages SET name = ?, version = ?, content = ?, url = ?, jsprogram = ? WHERE id = ?", name, version, content, url, jsprogram, packageID)
// 	if err != nil {
// 		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: "Server error"})
// 		return
// 	}

// 	rowsAffected, err := result.RowsAffected()
// 	if err != nil {
// 		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: "Server error"})
// 		return
// 	}

// 	if rowsAffected == 0 {
// 		c.AbortWithStatusJSON(http.StatusNotFound, ErrorResponse{Message: "Package not found"})
// 		return
// 	}

// 	c.Status(http.StatusOK)
// }


// // RegistryReset - Reset the registry
// func RegistryReset(c *gin.Context) {
// 	db := c.MustGet("db").(*sql.DB)

// 	_, err := db.Exec("DELETE FROM packages")
// 	if err != nil {
// 		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: "Server error"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Registry reset"})
// }