package api

import (
	"database/sql"
	//"encoding/json"
	// "errors"
	"log"
	"net/http"
	"fmt"
	"strings"
	"strconv"
	"math/rand"
	"time"
	"os"
	"github.com/gin-gonic/gin"
	"github.com/mabaums/ece461-web/backend/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
	"github.com/blang/semver"
	_ "github.com/go-sql-driver/mysql"

)

type AuthenticationToken string

type UserCredentials struct {
	Name string `json:"username"`
	isAdmin bool `json:"isAdmin"`
	Password string `json:"password"`
}



/*  HELPER FUNCTIONS */
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

func ExtractUserInfoFromToken(tokenString string) (string, string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verify that the signing method is HMAC and the secret matches
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}

		secret_key := os.Getenv("SECRET_KEY")   
		return []byte(secret_key), nil
	})
	if err != nil {
		return "",  "", err
	}

	// Extract the "user_name" claim from the token
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		username, ok1 := claims["username"].(string)
		password, ok2 := claims["password"].(string)
		if !ok1 || !ok2 {
				return "", "", fmt.Errorf("invalid claims")
		}
		return username, password, nil
	}
	return "", "", fmt.Errorf("invalid token")
}

/*  API Endpoints */


func CreateUser(c *gin.Context) {
	db, ok := getDB(c)
	if !ok {
		return
	}

	var info UserCredentials
	if err := c.ShouldBindJSON(&info); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := db.Exec("INSERT INTO User (name, isAdmin) VALUES (?, ?)", info.Name, info.isAdmin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	userID, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	result, err = db.Exec("INSERT INTO UserAuthenticationInfo (user_id, password) VALUES (?, ?)", userID, info.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"Message": "created user"})
}


func CreateAuthToken(c *gin.Context) {
	db, ok := getDB(c)
	if !ok {
		return
	}

	//	Get authentication request from request body
	var authReq models.AuthenticationRequest
	if err := c.ShouldBindJSON(&authReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user authentication info is correct
	// For example, verify user's password against stored hash
	var dbAuthInfo models.UserAuthenticationInfo
	err := db.QueryRow("SELECT password FROM UserAuthenticationInfo INNER JOIN User ON User.id = UserAuthenticationInfo.user_id WHERE user.name = ?", authReq.User.Name).Scan(&dbAuthInfo.Password)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}


	// Verify password
	if (dbAuthInfo.Password != authReq.Secret.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"Message": "unauth"})
	}
	// Verify password
	// if !verifyPassword(authReq.Secret.Password, dbAuthInfo.Password) {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
	// 	return
	// }

	// Create the claims for the JWT token
	claims := jwt.MapClaims{}
	claims["username"] = authReq.User.Name
	claims["password"] = authReq.Secret.Password
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix() // token expires in 1 hour

	// Create the JWT token with HMAC SHA-256 signing
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with a secret key
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	secret_key := os.Getenv("SECRET_KEY")  
	secret := []byte(secret_key)
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return
	}
	c.String(http.StatusOK, "Bearer "+tokenString)
}


func PackageCreate(c *gin.Context) {
	// Get Database
	db, ok := getDB(c)
	if !ok {
		return
	}
	// Authentication
	authTokenHeader := c.Request.Header.Get("X-Authorization")
	if authTokenHeader == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Authentication token not found in request header"})
		return
	}
	username, password, err := ExtractUserInfoFromToken(authTokenHeader)
	var pass string 
	err = db.QueryRow("SELECT password FROM UserAuthenticationInfo WHERE user_id = (SELECT id FROM User WHERE name = ?)", username).Scan(&pass)
	if err != nil || pass != password {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"description": "There is missing field(s) in the PackageData/AuthenticationToken" + 
		"or it is formed improperly (e.g. Content and URL are both set), or the AuthenticationToken is invalid." })
		return
	}

	// Process Request
	var pkg models.Package
	if err := c.ShouldBindJSON(&pkg); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"description": "There is missing field(s) in the PackageData/AuthenticationToken" + 
		"or it is formed improperly (e.g. Content and URL are both set), or the AuthenticationToken is invalid." })
		return
	}
	data := pkg.Data
	metadata := pkg.Metadata

  // Verify that only one of Content and URl are set
	dataURLEmpty := len(data.URL) == 0
	dataContentEmpty := len(data.Content) == 0
	if (!((dataURLEmpty && !dataContentEmpty) || (!dataURLEmpty && dataContentEmpty))) { // XOR
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"description": "There is missing field(s) in the PackageData/AuthenticationToken" + 
		"or it is formed improperly (e.g. Content and URL are both set), or the AuthenticationToken is invalid." })
		return
	}

	// Check Rating
	
	
	// PackageMetadata
	paramID := strings.TrimLeft(c.Param("id"), "/")

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM PackageMetadata WHERE PackageID = ?", paramID).Scan(&count)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error5"})
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
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error6"})
				return
			}
			if count == 0 {
				paramID = newID
				break
			}
		}
	}
	metadata.ID = paramID
	// Insert PackageMetadata
	result, err := db.Exec("INSERT INTO PackageMetadata (Name, Version, PackageID) VALUES (?, ?, ?)", metadata.Name, metadata.Version, paramID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusConflict, gin.H{"description": "Package exists already." })
		return
	}

	metadataID, err := result.LastInsertId()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error7"})
		return
	}

	// Insert PackageData
	var dataID int64
	if (dataURLEmpty) {
		result, err := db.Exec("INSERT INTO PackageData (Content, JSProgram) VALUES (?, ?)", data.Content, data.JSProgram)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error1"})
			return
		}
		dataID, err = result.LastInsertId()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error2"})
			return
		}
	} else if (dataContentEmpty) {
		result, err := db.Exec("INSERT INTO PackageData (URL, JSProgram) VALUES (?, ?)", data.URL, data.JSProgram)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error3"})
			return
		}
		dataID, err = result.LastInsertId()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error4"})
			return
		}
	}
		
	// Insert Package
	result, err = db.Exec("INSERT INTO Package (metadata_id, data_id) VALUES (?, ?)", metadataID, dataID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error8"})
		return
	}

	// Insert PackageHistoryEntry
	var User_temp models.User
	User_temp.Name = username
	User_temp.IsAdmin = false
	packageHistoryEntry := models.PackageHistoryEntry{
		User: User_temp,
		Date: time.Now(),
		PackageMetadata: metadata,
		Action: "Create",
	}
	var user_table_id int
	err = db.QueryRow("SELECT user.id FROM User WHERE user.name= ?", User_temp.Name).Scan(&user_table_id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error9"})
		return
	}
	result, err = db.Exec("INSERT INTO PackageHistoryEntry (user_id, date, package_metadata_id, action) VALUES (?, ?, ?, ?)", user_table_id, packageHistoryEntry.Date, metadataID, packageHistoryEntry.Action)
	if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
	}

	// Successful Response
	c.JSON(http.StatusCreated, gin.H{
		"metadata": metadata,
		"data": data,
	})
}

// PackageUpdate - Update this content of the package.
// historyentry
func PackageUpdate(c *gin.Context) {
	db, ok := getDB(c)
	if !ok {
		return
	}

	// Authentication
	authTokenHeader := c.Request.Header.Get("X-Authorization")
	if authTokenHeader == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Authentication token not found in request header"})
		return
	}
	username, password, err := ExtractUserInfoFromToken(authTokenHeader)
	var pass string 
	err = db.QueryRow("SELECT password FROM UserAuthenticationInfo WHERE user_id = (SELECT id FROM User WHERE name = ?)", username).Scan(&pass)
	if err != nil || pass != password {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"description": "There is missing field(s) in the PackageData/AuthenticationToken" + 
		"or it is formed improperly (e.g. Content and URL are both set), or the AuthenticationToken is invalid." })
		return
	}

	var pkg models.Package
	if err := c.ShouldBindJSON(&pkg); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}


	metadata := pkg.Metadata
	var existingPackage models.Package
	var package_data_id int;
	var package_metadata_id int;

	err = db.QueryRow("SELECT p.data_id, pmd.id, pmd.Name, pmd.Version, pmd.PackageID "+
			"FROM Package p "+
			"JOIN PackageMetadata pmd ON p.metadata_id = pmd.id "+
			"WHERE pmd.Name = ? AND pmd.Version = ? AND pmd.PackageID = ?",
			metadata.Name, metadata.Version, metadata.ID).Scan(
			&package_data_id, &package_metadata_id, &existingPackage.Metadata.Name, &existingPackage.Metadata.Version, &existingPackage.Metadata.ID,
	)
	if err == sql.ErrNoRows {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"description": "Package does not exist" })
		return
	} else if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	// Update package data
	packageData := pkg.Data
	_, err = db.Exec("UPDATE PackageData pd SET Content = ?, URL = ?, JSProgram = ? WHERE pd.id = ? ",
		packageData.Content, packageData.URL, packageData.JSProgram, package_data_id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	// Insert PackageHistoryEntry
	var User_temp models.User
	User_temp.Name = username
	User_temp.IsAdmin = false
	packageHistoryEntry := models.PackageHistoryEntry{
		User: User_temp,
		Date: time.Now(),
		PackageMetadata: metadata,
		Action: "Update",
	}
	var user_table_id int
	err = db.QueryRow("SELECT user.id FROM User WHERE user.name= ?", User_temp.Name).Scan(&user_table_id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error9"})
		return
	}
	_, err = db.Exec("INSERT INTO PackageHistoryEntry (user_id, date, package_metadata_id, action) VALUES (?, ?, ?, ?)", user_table_id, packageHistoryEntry.Date, package_metadata_id, packageHistoryEntry.Action)
	if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
	}

	c.JSON(http.StatusOK, gin.H{"description": "Version is updated"})
}	

// PackageDelete - Delete this version of the package. given packageid
func PackageDelete(c *gin.Context) {
	db, ok := getDB(c)
	if !ok {
		return
	}

	// Authentication
	authTokenHeader := c.Request.Header.Get("X-Authorization")
	if authTokenHeader == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Authentication token not found in request header"})
		return
	}
	username, password, err := ExtractUserInfoFromToken(authTokenHeader)
	var pass string 
	err = db.QueryRow("SELECT password FROM UserAuthenticationInfo WHERE user_id = (SELECT id FROM User WHERE name = ?)", username).Scan(&pass)
	if err != nil || pass != password {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"description": "There is missing field(s) in the PackageData/AuthenticationToken" + 
		"or it is formed improperly (e.g. Content and URL are both set), or the AuthenticationToken is invalid." })
		return
	}

	packageID := strings.TrimLeft(c.Param("id"), "/")
	var metadataID int
	err = db.QueryRow("SELECT id FROM PackageMetadata WHERE PackageID = ?", packageID).Scan(&metadataID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"description": "Package does not exist"})
		return
	}
	_, err = db.Exec("DELETE FROM PackageHistoryEntry WHERE package_metadata_id = ?", metadataID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	_, err = db.Exec("DELETE  pmd, pd, p FROM Package p "+
		"LEFT JOIN PackageMetaData pmd ON p.metadata_id = pmd.id "+
		"LEFT JOIN PackageData pd ON p.data_id = pd.id "+
		"WHERE p.metadata_id = ?", metadataID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"description": "Package is deleted"})
}

// PackageByNameDelete - Delete all versions of this package. return string of package name
func PackageByNameDelete(c *gin.Context) {
	db, ok := getDB(c)
	if !ok {
		return
	}

	// Authentication
	authTokenHeader := c.Request.Header.Get("X-Authorization")
	if authTokenHeader == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Authentication token not found in request header"})
		return
	}
	username, password, err := ExtractUserInfoFromToken(authTokenHeader)
	var pass string 
	err = db.QueryRow("SELECT password FROM UserAuthenticationInfo WHERE user_id = (SELECT id FROM User WHERE name = ?)", username).Scan(&pass)
	if err != nil || pass != password {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"description": "There is missing field(s) in the PackageData/AuthenticationToken" + 
		"or it is formed improperly (e.g. Content and URL are both set), or the AuthenticationToken is invalid." })
		return
	}


	packageName := strings.TrimLeft(c.Param("name"), "/")
	var metadataID int
	err = db.QueryRow("SELECT id FROM PackageMetadata WHERE Name = ?", packageName).Scan(&metadataID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"description": "Package does not exist"})
		return
	}

	_, err = db.Exec("DELETE FROM PackageHistoryEntry WHERE package_metadata_id = ?", metadataID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	_, err = db.Exec("DELETE p, pmd, pd, FROM Package p "+
		"LEFT JOIN PackageMetaData pmd ON p.metadata_id = pmd.id "+
		"LEFT JOIN PackageData pd ON p.data_id = pd.id "+
		"WHERE p.metadata_id = ?", metadataID)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"description": "Package  is deleted"})
}

// PackageRetrieve - Interact with the package with this ID
	// historyentry
func PackageRetrieve(c *gin.Context) {
	db, ok := getDB(c)
	if !ok {
		return
	}

	// Authentication
	authTokenHeader := c.Request.Header.Get("X-Authorization")
	if authTokenHeader == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Authentication token not found in request header"})
		return
	}
	username, password, err := ExtractUserInfoFromToken(authTokenHeader)
	var pass string 
	err = db.QueryRow("SELECT password FROM UserAuthenticationInfo WHERE user_id = (SELECT id FROM User WHERE name = ?)", username).Scan(&pass)
	if err != nil || pass != password {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"description": "There is missing field(s) in the PackageData/AuthenticationToken" + 
		"or it is formed improperly (e.g. Content and URL are both set), or the AuthenticationToken is invalid." })
		return
	}
		
	var packageID string
	packageID = strings.TrimLeft(c.Param("id"), "/")

	var packageName, packageVersion, packageContent, packageURL, packageJSProgram string
	fmt.Print(packageID)
	err = db.QueryRow("SELECT m.Name, m.Version, d.Content, d.URL, d.JSProgram " + 
	"FROM Package p " + 
	"INNER JOIN PackageMetadata m ON p.metadata_id = m.id " + 
	"INNER JOIN PackageData d ON p.data_id = d.id " + 
	"WHERE m.PackageID = ?;", packageID).Scan(&packageName, &packageVersion, &packageContent, &packageURL, &packageJSProgram)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err})
		// c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"description":"Package does not exist."})
		return
	}
	
	if err == sql.ErrNoRows {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"description":"Package does not exist."})
		return	
	}


	metadata := models.PackageMetadata {
		ID: packageID,
		Name: packageName,
		Version: packageVersion,
	}
	data := models.PackageData {
		Content: packageContent,
		URL: packageURL,
		JSProgram: packageJSProgram,
	}

	c.JSON(http.StatusOK, gin.H{
		"metadata": metadata,
		"data": data,
	})
}

// RegistryReset - Reset the registry
func RegistryReset(c *gin.Context) {
	db, ok := getDB(c)
	if !ok {
		return
	}

	// Authentication
	authTokenHeader := c.Request.Header.Get("X-Authorization")
	if authTokenHeader == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Authentication token not found in request header"})
		return
	}
	username, password, err := ExtractUserInfoFromToken(authTokenHeader)
	var pass string 
	err = db.QueryRow("SELECT password FROM UserAuthenticationInfo WHERE user_id = (SELECT id FROM User WHERE name = ?)", username).Scan(&pass)
	if err != nil || pass != password {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"description": "There is missing field(s) in the PackageData/AuthenticationToken" + 
		"or it is formed improperly (e.g. Content and URL are both set), or the AuthenticationToken is invalid." })
		return
	}

	// verify admin status	
	var isAdmin bool 
	err = db.QueryRow("SELECT isAdmin FROM User WHERE name ?", username).Scan(&isAdmin)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	if (!isAdmin) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"description": "You do not have permission to reset the registry"})
		return
	}

	// Delete all data from Package table
	_, err = db.Exec("DELETE FROM Package")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	// Delete all data from PackageData table
	_, err = db.Exec("DELETE FROM PackageData")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	// Delete all data from PackageHistoryEntry table
	_, err = db.Exec("DELETE FROM PackageHistoryEntry")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	// Delete all data from PackageMetadata table
	_, err = db.Exec("DELETE FROM PackageMetadata")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"description": "Registry is reset"})
}

// PacakgeByNameGet - 
func PackageByNameGet(c *gin.Context) {
	// Return the history of this package (all versions).
	db, ok := getDB(c)
	if !ok {
		return
	}

	// Authentication
	authTokenHeader := c.Request.Header.Get("X-Authorization")
	if authTokenHeader == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Authentication token not found in request header"})
		return
	}
	username, password, err := ExtractUserInfoFromToken(authTokenHeader)
	var pass string 
	err = db.QueryRow("SELECT password FROM UserAuthenticationInfo WHERE user_id = (SELECT id FROM User WHERE name = ?)", username).Scan(&pass)
	if err != nil || pass != password {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"description": "There is missing field(s) in the PackageData/AuthenticationToken" + 
		"or it is formed improperly (e.g. Content and URL are both set), or the AuthenticationToken is invalid." })
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

	if rows == nil || err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"description": "No such package."})
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
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
		packageMetadata.ID = "" // Set ID to "" since it's not needed
		packageHistoryEntry.PackageMetadata = packageMetadata
		packageHistoryEntries = append(packageHistoryEntries, packageHistoryEntry)
	}

	// Check for errors during iteration
	if err := rows.Err(); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, packageHistoryEntries)
}



func PackagesList(c *gin.Context) {
	db, ok := getDB(c)
	if !ok {
		return
	}

	// Authentication
	authTokenHeader := c.Request.Header.Get("X-Authorization")
	if authTokenHeader == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Authentication token not found in request header"})
		return
	}
	username, password, err := ExtractUserInfoFromToken(authTokenHeader)
	var pass string 
	err = db.QueryRow("SELECT password FROM UserAuthenticationInfo WHERE user_id = (SELECT id FROM User WHERE name = ?)", username).Scan(&pass)
	if err != nil || pass != password {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"description": "There is missing field(s) in the PackageData/AuthenticationToken" + 
		"or it is formed improperly (e.g. Content and URL are both set), or the AuthenticationToken is invalid." })
		return
	}

	// Parse request body
	var packageQueries []models.PackageQuery
	err = c.BindJSON(&packageQueries)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var queryStrings []string
    for _, packageQuery := range packageQueries {
        if packageQuery.Name == "" {
            c.JSON(http.StatusBadRequest, gin.H{"error": "PackageQuery object must have non-empty 'Name' field"})
            return
        }
        queryString := packageQuery.Name
        if packageQuery.Version != "" {
					semverRange, err := semver.ParseRange(packageQuery.Version)
					if err != nil {
							c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
							return
					}
					queryString += fmt.Sprintf("@%v", semverRange)
			}
        queryStrings = append(queryStrings, queryString)
    }



		offset := c.Query("offset")
		if offset == "" {
				offset = "0"
		}
		offsetInt, err := strconv.Atoi(offset)
		if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'offset' parameter"})
				return
		}

		limit := 100
    packages, err := getPackages(c, queryStrings, offsetInt, limit+1)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    if len(packages) > limit {
        c.Header("offset", strconv.Itoa(offsetInt+limit))
        packages = packages[:limit]
    } else {
        c.Header("offset", "0")
    }

    c.JSON(http.StatusOK, packages)
}


func getPackages(c *gin.Context, queryStrings []string, offset, limit int) ([]models.PackageMetadata, error) {
	db, ok := getDB(c)
	if !ok {
		return nil, nil
	}


	var whereClause string
	if len(queryStrings) > 0 {
			whereClause = "WHERE " + strings.Join(queryStrings, " AND ")
	}

	rows, err := db.Query("SELECT Name, Version, PackageID FROM packagemetadata "+whereClause+" LIMIT ? OFFSET ?", limit+1, offset)
	if err != nil {
			return nil, err
	}
	defer rows.Close()

	packages := make([]models.PackageMetadata, 0, limit)
	for rows.Next() {
			var pkg models.PackageMetadata
			err := rows.Scan(&pkg.Name, &pkg.Version, &pkg.ID)
			if err != nil {
					return nil, err
			}
			packages = append(packages, pkg)
	}
	if err := rows.Err(); err != nil {
			return nil, err
	}

	return packages, nil
}



// PackageByRegExGet - Get any packages fitting the regular expression.
func PackageByRegExGet(c *gin.Context) {
	// Search for packages that match the regular expression
	db, ok := getDB(c)
	if !ok {
		return
	}

	// Authentication
	authTokenHeader := c.Request.Header.Get("X-Authorization")
	if authTokenHeader == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Authentication token not found in request header"})
		return
	}
	username, password, err := ExtractUserInfoFromToken(authTokenHeader)
	var pass string 
	err = db.QueryRow("SELECT password FROM UserAuthenticationInfo WHERE user_id = (SELECT id FROM User WHERE name = ?)", username).Scan(&pass)
	if err != nil || pass != password {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"description": "There is missing field(s) in the PackageData/AuthenticationToken" + 
		"or it is formed improperly (e.g. Content and URL are both set), or the AuthenticationToken is invalid." })
		return
	}

	// Parse the request body as a PackageRegEx object
	var query string
	err = c.ShouldBindJSON(&query)
	if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
	}	

	var packages []models.PackageMetadata
	rows, err := db.Query("SELECT version, name, id FROM packagesmetadata WHERE name REGEXP ?", query)
	if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
	}
	defer rows.Close()
	for rows.Next() {
			var pkg models.PackageMetadata
			err := rows.Scan(&pkg.Version, &pkg.Name, &pkg.ID)
			if err != nil {
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
					return
			}
			packages = append(packages, pkg)
	}
	if err := rows.Err(); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
	}

	// Return the packages as a JSON array
	if len(packages) > 0 {
			c.JSON(http.StatusOK, packages)
	} else {
			c.AbortWithStatus(http.StatusNotFound)
	}
	c.JSON(http.StatusOK, gin.H{"message":"good"})
}

// PackageRate -
	// historyentry
func PackageRate(c *gin.Context) {
	// var rating Rating	
	// c.JSON(http.StatusOK, ratings)
}