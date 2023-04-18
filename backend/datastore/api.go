package api

import (
	"database/sql"
	//"encoding/json"
	// "errors"
	"log"
	"net/http"
	"strings"
	// "strconv"
	"math/rand"
	"time"
	// "fmt"
	"os"
	"github.com/gin-gonic/gin"
	"github.com/mabaums/ece461-web/backend/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
	_ "github.com/go-sql-driver/mysql"

)

type AuthenticationToken string

type UserCredentials struct {
	Name string `json:"username"`
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

// func VerifyPassword(username string, password string) (*User, error) {
// 	db, ok := getDB(c)
// 	if !ok {
// 		return
// 	}
// 	// Retrieve the user from the database
// 	var user User
// 	err := db.QueryRow("SELECT id, name, isAdmin FROM User WHERE name = ?", username).
// 		Scan(&user.ID, &user.Name, &user.IsAdmin)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return nil, nil // user not found
// 		}
// 		return nil, err // other error
// 	}

// 	// Retrieve the user's authentication info from the database
// 	var authInfo UserAuthenticationInfo
// 	err = db.QueryRow("SELECT id, user_id, password FROM UserAuthenticationInfo WHERE user_id = ?", user.ID).
// 		Scan(&authInfo.ID, &authInfo.UserID, &authInfo.Password)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return nil, nil // authentication info not found
// 		}
// 		return nil, err // other error
// 	}

// 	// Verify the password against the stored hash
// 	err = bcrypt.CompareHashAndPassword([]byte(authInfo.Password), []byte(password))
// 	if err != nil {
// 		return nil, nil // incorrect password
// 	}

// 	// Password is correct, return the user
// 	return &user, nil
// }

// func ExtractUserIDFromToken(tokenString string) (string, error) {
// 	// Parse the JWT token string
// 	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 		// Validate the signing method
// 		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 			return nil, jwt.NewValidationError("Unexpected signing method", jwt.ValidationErrorSignatureInvalid)
// 		}
// 		// Get the secret key used to sign the token
// 		secretKey := []byte("your-secret-key")
// 		return secretKey, nil
// 	})
// 	if err != nil {
// 		return "", err
// 	}
// 	// Extract the user ID from the token claims
// 	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
// 		if userID, ok := claims["user_id"].(string); ok {
// 			return userID, nil
// 		}
// 	}
// 	return "", jwt.NewValidationError("User ID not found in token", jwt.ValidationErrorClaimsInvalid)
// }

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
	result, err := db.Exec("INSERT INTO User (name, isAdmin) VALUES (?, ?)", info.Name, false)
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

	// Create JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": authReq.User.Name,
		// Add any other relevant user info to the token
	})

	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	secret_key := os.Getenv("SECRET_KEY")      

	// Sign the token with a secret key
	secretKey := []byte(secret_key) // Replace with your own secret key
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//Return the token as a response
	c.JSON(http.StatusOK, AuthenticationToken(tokenString))
}



	// historyentry
func PackageCreate(c *gin.Context) {
	// process request
	var pkg models.Package
	if err := c.ShouldBindJSON(&pkg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// connect to DB
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

	// authTokenHeader := c.Request.Header.Get("X-Authorization")
	// if authTokenHeader == "" {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Authentication token not found in request header"})
	// 	return
	// }
	// // Parse the authentication token
	// var authToken AuthenticationToken
	// if err := c.ShouldBindHeader(&authToken); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }
	// // Extract the user ID from the token
	// userID, err := ExtractUserIDFromToken(authToken.Token)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }

	// packageHistoryEntry := models.PackageHistoryEntry{
	// 	User: userID
	// 	Date: time.Now().Unix(),
	// 	PackageMetaData: paramID,
	// 	Action: "Create",
	// }

	// result, err = db.Exec("INSERT INTO PackageHistoryEntry (user_id, date, package_metadata_id, action) VALUES (?, ?, ?,", packageHistoryEntry.PackageID, packageHistoryEntry.Action, packageHistoryEntry.Time)
	// if err != nil {
	// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 		return
	// }

	c.JSON(http.StatusCreated, models.PackageMetadata{Name: metadata.Name, Version: metadata.Version, ID: paramID})
}

// PackageUpdate - Update this content of the package.
// historyentry
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
	// historyentry
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
	// // Parse the request body as an array of PackageQuery objects
	// var queries []models.PackageQuery
	// err := c.ShouldBindJSON(&queries)
	// if err != nil {
	// 		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
	// 		return
	// }

	// // Get the offset parameter from the query string
	// offsetStr := c.Query("offset")
	// var offset int
	// if offsetStr == "" {
	// 		offset = 0
	// } else {
	// 		var err error
	// 		offset, err = strconv.Atoi(offsetStr)
	// 		if err != nil {
	// 				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid offset parameter"})
	// 				return
	// 		}
	// }

	// db, ok := getDB(c)
	// if !ok {
	// 	return
	// }

	// // Build the SQL query based on the queries and offset
	// var query strings.Builder
	// query.WriteString("SELECT * FROM PackageMetadata")
	// if len(queries) > 0 {
	// 		query.WriteString(" WHERE ")
	// 		for i, q := range queries {
	// 				if i > 0 {
	// 						query.WriteString(" AND ")
	// 				}
	// 				query.WriteString(q.ToSQL())
	// 		}
	// }
	// query.WriteString(fmt.Sprintf(" LIMIT %d, 10", offset))

	// // Execute the SQL query and retrieve the package metadata
	// rows, err := db.Query(query.String())
	// if err != nil {
	// 		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to query packages from database"})
	// 		return
	// }
	// defer rows.Close()

	// packages := make([]models.PackageMetadata, 0)
	// for rows.Next() {
	// 		var packageMetadata models.PackageMetadata
	// 		err := rows.Scan(&packageMetadata.ID, &packageMetadata.Name, &packageMetadata.Version)
	// 		if err != nil {
	// 				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan package metadata"})
	// 				return
	// 		}
	// 		packages = append(packages, packageMetadata)
	// }

	// // Set the offset header in the response
	// nextOffset := offset + len(packages)
	// c.Header("offset", strconv.Itoa(nextOffset))

	// // Return the package metadata as a JSON array
	// c.JSON(http.StatusOK, packages)
}

// PackageByRegExGet - Get any packages fitting the regular expression.
func PackageByRegExGet(c *gin.Context) {
	// // Parse the request body as a PackageRegEx object
	// var query models.PackageRegEx
	// err := c.ShouldBindJSON(&query)
	// if err != nil {
	// 		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
	// 		return
	// }

	// // Search for packages that match the regular expression
	// db, ok := getDB(c)
	// if !ok {
	// 	return
	// }

	// var packages []models.PackageMetadata
	// rows, err := db.Query("SELECT version, name, id FROM packages WHERE name REGEXP ? OR readme REGEXP ?", query.RegEx, query.RegEx)
	// if err != nil {
	// 		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
	// 		return
	// }
	// defer rows.Close()
	// for rows.Next() {
	// 		var pkg models.PackageMetadata
	// 		err := rows.Scan(&pkg.Version, &pkg.Name, &pkg.ID)
	// 		if err != nil {
	// 				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
	// 				return
	// 		}
	// 		packages = append(packages, pkg)
	// }
	// if err := rows.Err(); err != nil {
	// 		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
	// 		return
	// }

	// // Return the packages as a JSON array
	// if len(packages) > 0 {
	// 		c.JSON(http.StatusOK, packages)
	// } else {
	// 		c.AbortWithStatus(http.StatusNotFound)
	// }
	// c.Json(http.StatusOK, gin.H{"message":"good"})
}

// PackageRate -
	// historyentry
func PackageRate(c *gin.Context) {
	// var rating Rating	
	// c.JSON(http.StatusOK, ratings)
}