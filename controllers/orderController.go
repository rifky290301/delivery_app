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

// CreateOrder handles creating a new order
func CreateOrder(w http.ResponseWriter, r *http.Request) {
	var order models.Order

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&order); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	query := `
        INSERT INTO Orders (buyer_id, amount, status, created_at)
        VALUES (?, ?, ?, ?)`
	result, err := config.DB.Exec(query, order.BuyerId, order.Amount, "pending", time.Now())
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error creating order")
		return
	}

	orderID, err := result.LastInsertId()
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error retrieving order ID")
		return
	}

	order.Id = int(orderID)
	order.Status = "pending"
	order.CreatedAt = time.Now()

	utils.RespondWithJSON(w, http.StatusCreated, order)
}

// GetOrder handles retrieving an order by its ID
func GetOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid order ID")
		return
	}

	var order models.Order
	query := `SELECT 
	o.id, o.buyer_id, o.amount, o.status, o.created_at, o.updated_at, 
	u.id, u.user_name, u.email, u.address, u.phone_number, u.instagram_link, u.profile_picture, u.description, u.created_at, u.role
	FROM orders o 
	JOIN users u 
	ON o.buyer_id = u.id 
	WHERE o.id = ?`
	err = config.DB.QueryRow(query, orderID).Scan(
		&order.Id,
		&order.BuyerId,
		&order.Amount,
		&order.Status,
		&order.CreatedAt,
		&order.UpdatedAt,
		&order.Buyer.Id,
		&order.Buyer.UserName,
		&order.Buyer.Email,
		&order.Buyer.Address,
		&order.Buyer.PhoneNumber,
		&order.Buyer.InstagramLink,
		&order.Buyer.ProfilePicture,
		&order.Buyer.Description,
		&order.Buyer.CreatedAt,
		&order.Buyer.Role,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.RespondWithError(w, http.StatusNotFound, "Order not found")
		} else {
			fmt.Println(err)
			utils.RespondWithError(w, http.StatusInternalServerError, "Error retrieving order")
		}
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, order)
}

// GetAllOrders handles retrieving all orders
func GetAllOrders(w http.ResponseWriter, r *http.Request) {
	rows, err := config.DB.Query("SELECT * FROM orders")
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error retrieving orders")
		return
	}
	defer rows.Close()

	var orders []models.Order

	for rows.Next() {
		var order models.Order
		if err := rows.Scan(
			&order.Id,
			&order.BuyerId,
			&order.Amount,
			&order.Status,
			&order.CreatedAt,
			&order.UpdatedAt,
		); err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error scanning order")
			return
		}
		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error iterating over orders")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, orders)
}

// UpdateOrder handles updating an order by its ID
func UpdateOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid order ID")
		return
	}

	var order models.Order
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&order); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	query := `
			UPDATE orders SET buyer_id = ?, amount = ?, status = ?, updated_at = ?
			WHERE id = ?`
	_, err = config.DB.Exec(query, order.BuyerId, order.Amount, order.Status, time.Now(), orderID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error updating order")
		return
	}

	order.Id = orderID
	order.UpdatedAt = models.NullTime{NullTime: sql.NullTime{Time: time.Now(), Valid: true}}
	utils.RespondWithJSON(w, http.StatusOK, order)
}

// DeleteOrder handles deleting an order by its ID
func DeleteOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid order ID")
		return
	}

	query := "DELETE FROM orders WHERE id = ?"
	_, err = config.DB.Exec(query, orderID)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.RespondWithError(w, http.StatusNotFound, "Order not found")
		} else {
			fmt.Println(err)
			utils.RespondWithError(w, http.StatusInternalServerError, "Error deleting order")
		}
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Order deleted successfully"})
}
