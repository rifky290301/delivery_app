package controllers

import (
	"database/sql"
	"delivery_app/config"
	"delivery_app/middleware"
	"delivery_app/models"
	"delivery_app/utils"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

func GetUsers(w http.ResponseWriter, r *http.Request) {
	query := "SELECT * FROM Users"
	rows, err := config.DB.Query(query)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error querying database")
		return
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.Id,
			&user.UserName,
			&user.Email,
			&user.PasswordHash,
			&user.Address,
			&user.PhoneNumber,
			&user.Role,
			&user.Description,
			&user.ProfilePicture,
			&user.InstagramLink,
			&user.CreatedAt,
			&user.UpdatedAt,
		)

		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error scanning row"+err.Error())
			return
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error iterating rows")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, users)
}

func Register(w http.ResponseWriter, r *http.Request) {
	var user models.User

	// Decode the JSON body to user struct
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// Validate mandatory fields
	if user.UserName == "" || user.Email == "" || user.PasswordHash == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Please provide username: "+user.UserName+" email: "+user.Email+" password :"+user.PasswordHash)
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error hashing password")
		return
	}
	user.PasswordHash = string(hashedPassword)

	// Insert user into the database
	query := "INSERT INTO Users (user_name, email, password_hash) VALUES (?, ?, ?)"
	stmt, err := config.DB.Prepare(query)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error preparing query")
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.UserName, user.Email, user.PasswordHash)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error inserting user into database")
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, map[string]string{"message": "User registered successfully"})
}

func Login(w http.ResponseWriter, r *http.Request) {
	var loginRequest models.LoginRequest

	// Decode the JSON body to loginRequest struct
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&loginRequest); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// Validate mandatory fields
	if loginRequest.Email == "" || loginRequest.Password == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Please provide email and password")
		return
	}

	// Fetch user from the database
	var user models.User
	query := "SELECT id, user_name, password_hash FROM Users WHERE email = ?"
	err := config.DB.QueryRow(query, loginRequest.Email).Scan(&user.Id, &user.UserName, &user.PasswordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.RespondWithError(w, http.StatusUnauthorized, "Invalid email or password")
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error querying the database")
		}
		return
	}

	// Compare the hashed password with the password provided
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(loginRequest.Password))
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	// Create the JWT token
	claims := models.Claims{
		Id: user.Id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(), // Token expires in 24 hours
			Issuer:    "delivery_app",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(config.JWTSecret)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error generating token")
		return
	}

	// Send the token in the response
	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"token": tokenString})
}

func Logout(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		utils.RespondWithError(w, http.StatusUnauthorized, "Missing authorization header")
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	// Menambahkan token ke blacklist
	middleware.BlacklistedTokens[tokenString] = true

	// Menanggapi bahwa logout berhasil
	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Successfully logged out"})
}

func CompleteProfile(w http.ResponseWriter, r *http.Request) {
	id, ok := r.Context().Value(middleware.UserIDKey).(int)

	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized access")
		return
	}

	var user models.User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// Build the update query dynamically based on non-empty fields
	query := "UPDATE Users SET "
	args := []interface{}{}
	if user.Address.Valid {
		query += "address = ?, "
		args = append(args, user.Address.String)
	}
	if user.PhoneNumber.Valid {
		query += "phone_number = ?, "
		args = append(args, user.PhoneNumber.String)
	}
	if user.ProfilePicture.Valid {
		query += "profile_picture = ?, "
		args = append(args, user.ProfilePicture.String)
	}
	if user.InstagramLink.Valid {
		query += "instagram_link = ?, "
		args = append(args, user.InstagramLink.String)
	}
	if user.Description.Valid {
		query += "description = ?, "
		args = append(args, user.Description.String)
	}
	if user.UpdatedAt.Valid {
		query += "updated_at = ?, "
		args = append(args, time.Now())
	}

	// Remove trailing comma and space, and add WHERE clause
	query = query[:len(query)-2] + " WHERE id = ?"
	args = append(args, id)

	// Execute the update query
	if _, err := config.DB.Exec(query, args...); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error updating user profile")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Profile updated successfully"})
}
