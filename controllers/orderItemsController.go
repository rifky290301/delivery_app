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

// CreateOrderItem handles creating a new order item
func CreateOrderItem(w http.ResponseWriter, r *http.Request) {
	var orderItem models.OrderItem
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&orderItem); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	orderItem.CreatedAt = time.Now()

	query := `
        INSERT INTO OrderItems (order_id, product_id, quantity, price, created_at)
        VALUES (?, ?, ?, ?, ?)`
	result, err := config.DB.Exec(query, orderItem.OrderID, orderItem.ProductID, orderItem.Quantity, orderItem.Price, orderItem.CreatedAt)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error creating order item")
		return
	}

	orderItemID, err := result.LastInsertId()
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error retrieving last insert ID")
		return
	}

	orderItem.OrderItemID = int(orderItemID)
	utils.RespondWithJSON(w, http.StatusCreated, orderItem)
}

// GetOrderItem handles retrieving an order item by its ID, including order and product details
func GetOrderItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderItemID, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid order item ID")
		return
	}

	var orderItem models.OrderItem
	query := `
        SELECT 
            oi.order_item_id, oi.order_id, oi.product_id, oi.quantity, oi.price, oi.created_at, oi.updated_at,
            o.id, o.buyer_id, o.amount, o.status, o.created_at, o.updated_at,
            p.id, p.shop_id, p.product_name, p.description, p.price, p.stock, p.created_at, p.updated_at
        FROM 
            OrderItems oi
        INNER JOIN 
            Orders o ON oi.order_id = o.id
        INNER JOIN 
            Products p ON oi.product_id = p.id
        WHERE 
            oi.order_item_id = ?`
	err = config.DB.QueryRow(query, orderItemID).Scan(
		&orderItem.OrderItemID,
		&orderItem.OrderID,
		&orderItem.ProductID,
		&orderItem.Quantity,
		&orderItem.Price,
		&orderItem.CreatedAt,
		&orderItem.UpdatedAt,
		&orderItem.Order.Id,
		&orderItem.Order.BuyerId,
		&orderItem.Order.Amount,
		&orderItem.Order.Status,
		&orderItem.Order.CreatedAt,
		&orderItem.Order.UpdatedAt,
		&orderItem.Product.Id,
		&orderItem.Product.ShopID,
		&orderItem.Product.Name,
		&orderItem.Product.Description,
		&orderItem.Product.Price,
		&orderItem.Product.Stock,
		&orderItem.Product.CreatedAt,
		&orderItem.Product.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.RespondWithError(w, http.StatusNotFound, "Order item not found")
		} else {
			fmt.Println(err)
			utils.RespondWithError(w, http.StatusInternalServerError, "Error retrieving order item")
		}
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, orderItem)
}

// GetAllOrderItems handles retrieving all order items
func GetAllOrderItems(w http.ResponseWriter, r *http.Request) {
	var orderItems []models.OrderItem
	// SELECT OrderItemID, OrderID, ProductID, Quantity, Price, CreatedAt, UpdatedAt
	query := `
			SELECT *
			FROM OrderItems`
	rows, err := config.DB.Query(query)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error retrieving order items")
		return
	}
	defer rows.Close()

	for rows.Next() {
		var orderItem models.OrderItem
		if err := rows.Scan(
			&orderItem.OrderItemID,
			&orderItem.OrderID,
			&orderItem.ProductID,
			&orderItem.Quantity,
			&orderItem.Price,
			&orderItem.CreatedAt,
			&orderItem.UpdatedAt,
		); err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error scanning order item")
			return
		}
		orderItems = append(orderItems, orderItem)
	}

	utils.RespondWithJSON(w, http.StatusOK, orderItems)
}

// UpdateOrderItem handles updating an existing order item
func UpdateOrderItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderItemID, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid order item ID")
		return
	}

	var orderItem models.OrderItem
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&orderItem); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	orderItem.UpdatedAt = models.NullTime{NullTime: sql.NullTime{Time: time.Now(), Valid: true}}

	query := `
			UPDATE OrderItems
			SET order_id = ?, product_id = ?, quantity = ?, price = ?, UpdatedAt = ?
			WHERE order_item_id = ?`
	_, err = config.DB.Exec(query, orderItem.OrderID, orderItem.ProductID, orderItem.Quantity, orderItem.Price, orderItem.UpdatedAt, orderItemID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error updating order item")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, orderItem)
}

// DeleteOrderItem handles deleting an order item by its ID
func DeleteOrderItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderItemID, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid order item ID")
		return
	}

	query := `
			DELETE FROM OrderItems
			WHERE order_item_id = ?`
	_, err = config.DB.Exec(query, orderItemID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error deleting order item")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}
