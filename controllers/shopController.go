package controllers

import (
	"database/sql"
	"delivery_app/config"
	"delivery_app/middleware"
	"delivery_app/models"
	"delivery_app/utils"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func GetShop(w http.ResponseWriter, r *http.Request) {
	// Get the shop ID from the URL parameters
	vars := mux.Vars(r)
	shopID, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid shop ID")
		return
	}

	// Query the database for the shop
	var shop models.Shop
	query := "SELECT s.id, s.shop_name, s.shop_description, s.shop_address, u.id, u.user_name, u.email, u.address  FROM sellershops s JOIN users u ON s.user_id = u.id WHERE s.id = ?"
	// query := "SELECT * FROM sellershops s JOIN users u ON s.user_id = u.id WHERE s.id = ?"
	err = config.DB.QueryRow(query, shopID).Scan(
		&shop.Id,
		// &shop.SellerID,
		&shop.ShopName,
		&shop.ShopDescription,
		&shop.ShopAddress,
		// &shop.CreatedAt,
		&shop.Seller.Id,
		&shop.Seller.UserName,
		&shop.Seller.Email,
		&shop.Seller.Address,
		// &shop.Seller.PhoneNumber,
		// &shop.Seller.PasswordHash,
		// &shop.Seller.ProfilePicture,
		// &shop.Seller.InstagramLink,
		// &shop.Seller.Description,
		// &shop.Seller.CreatedAt,
		// &shop.Seller.Role,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.RespondWithError(w, http.StatusNotFound, "Shop not found")
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error retrieving shop")
		}
		return
	}

	// Respond with the shop data
	utils.RespondWithJSON(w, http.StatusOK, shop)
}

func GetAllShops(w http.ResponseWriter, r *http.Request) {
	// Query the database for all shops
	rows, err := config.DB.Query("SELECT * FROM sellershops")
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error retrieving shops")
		return
	}
	defer rows.Close()

	var shops []models.Shop

	// Iterate through the result set and populate the shops slice
	for rows.Next() {
		var shop models.Shop
		if err := rows.Scan(
			&shop.Id,
			&shop.SellerID,
			&shop.ShopName,
			&shop.ShopDescription,
			&shop.ShopAddress,
			&shop.CreatedAt,
		); err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error scanning shop")
			return
		}
		shops = append(shops, shop)
	}

	// Check for errors from iterating over rows
	if err = rows.Err(); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error iterating over shops")
		return
	}

	// Respond with the shops data
	utils.RespondWithJSON(w, http.StatusOK, shops)
}

func CreateShop(w http.ResponseWriter, r *http.Request) {
	var shop models.Shop

	// Decode request body into shop struct
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&shop); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// Get user ID from context
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	shop.SellerID = userID
	shop.CreatedAt = time.Now()

	// Insert shop into database
	query := `
        INSERT INTO sellershops (user_id, shop_name, shop_description, shop_address, created_at)
        VALUES (?, ?, ?, ?, ?)`
	result, err := config.DB.Exec(query, shop.SellerID, shop.ShopName, shop.ShopDescription, shop.ShopAddress, shop.CreatedAt)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error creating shop")
		return
	}

	shopID, err := result.LastInsertId()
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error retrieving last inserted ID")
		return
	}
	shop.Id = int(shopID)

	utils.RespondWithJSON(w, http.StatusOK, shop)
}

func UpdateShop(w http.ResponseWriter, r *http.Request) {
	// Get the shop ID from the URL parameters
	vars := mux.Vars(r)
	shopID, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid shop ID")
		return
	}

	var shop models.Shop

	// Decode request body into shop struct
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&shop); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// Update shop in the database
	query := `
			UPDATE sellershops SET shop_name = ?, shop_description = ?, shop_address = ?
			WHERE id = ?`
	_, err = config.DB.Exec(query, shop.ShopName, shop.ShopDescription, shop.ShopAddress, shopID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error updating shop")
		return
	}

	// Respond with the updated shop data
	shop.Id = shopID
	utils.RespondWithJSON(w, http.StatusOK, shop)
}

func DeleteShop(w http.ResponseWriter, r *http.Request) {
	// Get the shop ID from the URL parameters
	vars := mux.Vars(r)
	shopID, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid shop ID")
		return
	}

	// Delete shop from the database
	query := "DELETE FROM sellershops WHERE id = ?"
	_, err = config.DB.Exec(query, shopID)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.RespondWithError(w, http.StatusNotFound, "Shop not found")
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error deleting shop")
		}
		return
	}

	// Respond with success message
	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Shop deleted successfully"})
}
