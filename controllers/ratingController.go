package controllers

import (
	"database/sql"
	"delivery_app/config"
	"delivery_app/models"
	"delivery_app/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// CreateRating handles creating a new rating
func CreateRating(w http.ResponseWriter, r *http.Request) {
	var rating models.Rating
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rating); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	rating.CreatedAt = time.Now()

	query := `
        INSERT INTO Ratings (order_id, buyer_id, shop_id, rating, feedback, created_at)
        VALUES (?, ?, ?, ?, ?, ?)`
	result, err := config.DB.Exec(query, rating.OrderID, rating.BuyerID, rating.ShopID, rating.Rating, rating.Feedback, rating.CreatedAt)
	if err != nil {
		fmt.Println(err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Error creating rating")
		return
	}

	ratingID, err := result.LastInsertId()
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error retrieving last insert ID")
		return
	}

	rating.Id = int(ratingID)
	utils.RespondWithJSON(w, http.StatusCreated, rating)
}

// GetRating handles retrieving a rating by its ID
func GetRating(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ratingID, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid rating ID")
		return
	}

	var rating models.Rating
	query := `
			SELECT * Ratings
			WHERE id = ?`
	err = config.DB.QueryRow(query, ratingID).Scan(
		&rating.Id,
		&rating.OrderID,
		&rating.BuyerID,
		&rating.ShopID,
		&rating.Rating,
		&rating.Feedback,
		&rating.CreatedAt,
		&rating.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.RespondWithError(w, http.StatusNotFound, "Rating not found")
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error retrieving rating")
		}
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, rating)
}

// GetAllRatings handles retrieving all ratings
func GetAllRatings(w http.ResponseWriter, r *http.Request) {
	var ratings []models.Rating
	query := `SELECT * FROM Ratings`
	rows, err := config.DB.Query(query)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error retrieving ratings")
		return
	}
	defer rows.Close()

	for rows.Next() {
		var rating models.Rating
		if err := rows.Scan(
			&rating.Id,
			&rating.OrderID,
			&rating.BuyerID,
			&rating.ShopID,
			&rating.Rating,
			&rating.Feedback,
			&rating.CreatedAt,
			&rating.UpdatedAt,
		); err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error scanning rating")
			return
		}
		ratings = append(ratings, rating)
	}

	utils.RespondWithJSON(w, http.StatusOK, ratings)
}

// UpdateRating handles updating an existing rating
func UpdateRating(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ratingID, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid rating ID")
		return
	}

	var rating models.Rating
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rating); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	rating.UpdatedAt = models.NullTime{NullTime: sql.NullTime{Time: time.Now(), Valid: true}}

	query := `
			UPDATE Ratings
			SET  order_id = ?,buyer_id = ?, shop_id = ?, rating = ?, feedback = ?, updated_at = ?
			WHERE id = ?`
	_, err = config.DB.Exec(query, rating.OrderID, rating.BuyerID, rating.ShopID, rating.Rating, rating.Feedback, rating.UpdatedAt, ratingID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error updating rating")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, rating)
}

// DeleteRating handles deleting a rating by its ID
func DeleteRating(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ratingID, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid rating ID")
		return
	}

	query := `
			DELETE FROM Ratings
			WHERE id = ?`
	_, err = config.DB.Exec(query, ratingID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error deleting rating")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}
