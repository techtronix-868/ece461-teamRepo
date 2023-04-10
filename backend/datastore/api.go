package api

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"

)
 
type AuthenticationRequest struct {
	User User `json:"User"`

	Secret UserAuthenticationInfo `json:"Secret"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

// CreateAuthToken -
func CreateAuthToken(c *gin.Context) {
	// username := c.PostForm("username")
	// password := c.PostForm("password")

	// // Check if the username and password are valid
	// if !isValidUser(username, password) {
	// 		c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{Message: "Invalid username or password"})
	// 		return
	// }

	// // Generate a new authentication token
	// token, err := generateToken(username)
	// if err != nil {
	// 		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: "Server error"})
	// 		return
	// }

	// // Return the authentication token to the client
	// c.JSON(http.StatusOK, gin.H{"token": token})
}

// PackageByNameDelete - Delete all versions of this package.

func PackageByNameDelete(c *gin.Context) {
	packageName := c.Param("name")

	// Get database connection from context
	db, ok := c.Get("db")
	if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: "Server error"})
			return
	}

	// Delete package from database
	stmt, err := db.(*sql.DB).Prepare("DELETE FROM packages WHERE name = ?")
	if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: "Server error"})
			return
	}
	defer stmt.Close()

	res, err := stmt.Exec(packageName)
	if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: "Server error"})
			return
	}

	// Check if package was deleted
	count, err := res.RowsAffected()
	if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: "Server error"})
			return
	}
	if count == 0 {
			c.AbortWithStatusJSON(http.StatusNotFound, ErrorResponse{Message: "Package not found"})
			return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{"message": "Package deleted"})
}

// PacakgeByNameGet - 
func PackageByNameGet(c *gin.Context) {
	// Return the history of this package (all versions).
	packageName := c.Param("name")

	db, ok := c.Get("db")
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: "Server error"})
		return
	}

	var package models.Package
	// change here to verify with schema
	err := db.(*sql.DB).QueryRow("SELECT id, name, version, description FROM packages WHERE name = ?", packageName).Scan(&package.ID, &package.Name, &package.Version, &package.Description)
	if err != nil {
		if err.Is(err, sql.ErrNoRows) {
			c.AbortWithStatusJSON(http.StatusNotFound, ErrorResponse{Message:"Package not found"})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: "Server error"})
		}
		return
	}
	c.JSON(http.StatusOK, package)
}

// PackageByRegExGet - Get any packages fitting the regular expression.
func PackageByRegexGet(c *gin.Context) {
	packageNamePattern := c.Param("pattern")

	db, ok := c.Get("db")
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: "Server error"})
		return
	}

	rows, err := db.(*sql.DB).Query("SELECT id, name, version, description FROM packages WHERE name REGEXP ?", packageNamePattern)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: "Server error"})
		return
	}
	defer rows.Close()

	pkgs := []models.Package{}
	for rows.Next() {
		pkg := models.Package{}
		err := rows.Scan(&pkg.ID, &pkg.Name, &pkg.Version, &pkg.Description)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: "Server error"})
			return
		}
		pkgs = append(pkgs, pkg)
	}

	c.JSON(http.StatusOK, pkgs)
}

// PackageCreate -
func PackageCreate(c *gin.Context) {
	name := c.PostForm("name")
	version := c.PostForm("version")
	content := c.PostForm("content")
	url := c.PostForm("url")
	jsprogram := c.PostForm("jsprogram")

	db, ok := c.Get("db")
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: "Server error"})
		return
	}

	result, err := db.(*sql.DB).Exec("INSERT INTO packages (name, version, content, url, jsprogram) VALUES (?, ?, ?, ?, ?)", name, version, content, url, jsprogram)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: "Server error"})
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: "Server error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// PackageDelete - Delete this version of the package.
func PackageDelete(c *gin.Context) {
	name := c.Param("name")
	version := c.Param("version")

	db, ok := c.Get("db")
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: "Server error"})
		return
	}

	result, err := db.(*sql.DB).Exec("DELETE FROM packages WHERE name = ? AND version = ?", name, version)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: "Server error"})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: "Server error"})
		return
	}

	if rowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusNotFound, ErrorResponse{Message: "Package not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Package deleted"})
}

// PackageRate -
func PackageRate(c *gin.Context) {
	packageID := c.Param("id")

	db, ok := c.Get("db")
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: "Server error"})
		return
	}

	var ratings []models.Rating
	rows, err := db.(*sql.DB).Query("SELECT id, package_id, user_id, rating FROM ratings WHERE package_id = ?", packageID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: "Server error"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var rating models.Rating
		err = rows.Scan(&rating.ID, &rating.PackageID, &rating.UserID, &rating.Rating)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: "Server error"})
			return
		}
		ratings = append(ratings, rating)
	}

	if len(ratings) == 0 {
		c.AbortWithStatusJSON(http.StatusNotFound, ErrorResponse{Message: "Package not found"})
		return
	}

	c.JSON(http.StatusOK, ratings)
}

// PackageRetrieve - Interact with the package with this ID
func PackageRetrieve(c *gin.Context) {
	packageID := c.Param("id")

	db, ok := c.Get("db")
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: "Server error"})
		return
	}

	var package models.Package
	err := db.(*sql.DB).QueryRow("SELECT id, name, version, content, url, jsprogram FROM packages WHERE id = ?", packageID).Scan(&package.ID, &package.Name, &package.Version, &package.Content, &package.URL, &package.JSProgram)
	if err != nil {
		if err == sql.ErrNoRows {
			c.AbortWithStatusJSON(http.StatusNotFound, ErrorResponse{Message: "Package not found"})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: "Server error"})
		}
		return
	}

	c.JSON(http.StatusOK, package)
}

// PackageUpdate - Update this content of the package.
func PackageUpdate(c *gin.Context) {
	packageID := c.Param("id")

	name := c.PostForm("name")
	version := c.PostForm("version")
	content := c.PostForm("content")
	url := c.PostForm("url")
	jsprogram := c.PostForm("jsprogram")

	db, ok := c.Get("db")
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: "Server error"})
		return
	}

	result, err := db.(*sql.DB).Exec("UPDATE packages SET name = ?, version = ?, content = ?, url = ?, jsprogram = ? WHERE id = ?", name, version, content, url, jsprogram, packageID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: "Server error"})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: "Server error"})
		return
	}

	if rowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusNotFound, ErrorResponse{Message: "Package not found"})
		return
	}

	c.Status(http.StatusOK)
}

// PackagesList - Get the packages from the registry.
func PackagesList(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)

	rows, err := db.Query("SELECT id, name, version, description FROM packages")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: "Server error"})
		return
	}
	defer rows.Close()

	packages := []models.Package{}
	for rows.Next() {
		var pkg models.Package
		err = rows.Scan(&pkg.ID, &pkg.Name, &pkg.Version, &pkg.Description)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: "Server error"})
			return
		}
		packages = append(packages, pkg)
	}

	c.JSON(http.StatusOK, packages)
}

// RegistryReset - Reset the registry
func RegistryReset(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)

	_, err := db.Exec("DELETE FROM packages")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: "Server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Registry reset"})
}