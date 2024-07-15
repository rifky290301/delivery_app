package controllers

import (
	"database/sql"
	"delivery_app/config"
	"delivery_app/models"
	"delivery_app/utils"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// CreateProduct handles creating a new product
func CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product models.Product

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&product); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	query := `
        INSERT INTO Products (shop_id, product_name, description, price, stock, created_at)
        VALUES (?, ?, ?, ?, ?, ?)`
	result, err := config.DB.Exec(query, product.ShopID, product.Name, product.Description, product.Price, product.Stock, time.Now())
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error creating product")
		return
	}

	productID, err := result.LastInsertId()
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error retrieving product ID")
		return
	}

	product.Id = int(productID)
	product.CreatedAt = time.Now()

	utils.RespondWithJSON(w, http.StatusCreated, product)
}

// GetProduct handles retrieving a product by its ID
func GetProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	var product models.Product
	query := "SELECT * FROM products WHERE id = ?"
	err = config.DB.QueryRow(query, productID).Scan(
		&product.Id,
		&product.ShopID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.Stock,
		&product.CreatedAt,
		&product.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.RespondWithError(w, http.StatusNotFound, "Product not found")
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error retrieving product")
		}
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, product)
}

// GetAllProducts handles retrieving all products
func GetAllProducts(w http.ResponseWriter, r *http.Request) {
	rows, err := config.DB.Query("SELECT * FROM products")
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error retrieving products")
		return
	}
	defer rows.Close()

	var products []models.Product

	for rows.Next() {
		var product models.Product
		if err := rows.Scan(
			&product.Id,
			&product.ShopID,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.Stock,
			&product.CreatedAt,
			&product.UpdatedAt,
		); err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error scanning product")
			return
		}
		products = append(products, product)
	}

	if err = rows.Err(); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error iterating over products")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, products)
}

// UpdateProduct handles updating a product by its ID
func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	var product models.Product
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&product); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	var timeNow = time.Now()
	query := `
			UPDATE products SET shop_id = ?, product_name = ?, description = ?, price = ?, stock = ?, updated_at = ?
			WHERE id = ?`
	_, err = config.DB.Exec(query, product.ShopID, product.Name, product.Description, product.Price, product.Stock, timeNow, productID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error updating product")
		return
	}

	// product.UpdatedAt = timeNow;
	product.Id = productID
	utils.RespondWithJSON(w, http.StatusOK, product)
}

// DeleteProduct handles deleting a product by its ID
func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	query := "DELETE FROM products WHERE id = ?"
	_, err = config.DB.Exec(query, productID)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.RespondWithError(w, http.StatusNotFound, "Product not found")
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error deleting product")
		}
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Product deleted successfully"})
}
