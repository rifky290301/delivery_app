package routes

import (
	"delivery_app/controllers"
	"delivery_app/middleware"

	"github.com/gorilla/mux"
)

func InitRoutes() *mux.Router {
	router := mux.NewRouter()

	// Public routes
	router.HandleFunc("/register", controllers.Register).Methods("POST")
	router.HandleFunc("/login", controllers.Login).Methods("POST")

	// Private routes with JWT Middleware
	api := router.PathPrefix("/api").Subrouter()
	// Apply JWT Middleware
	api.Use(middleware.JWTMiddleware)

	api.HandleFunc("/users", controllers.GetUsers).Methods("GET")
	api.HandleFunc("/logout", controllers.Logout).Methods("POST")
	api.HandleFunc("/complete-profile", controllers.CompleteProfile).Methods("PUT")

	api.HandleFunc("/shop", controllers.CreateShop).Methods("POST")
	api.HandleFunc("/shop/{id}", controllers.GetShop).Methods("GET")
	api.HandleFunc("/shop", controllers.GetAllShops).Methods("GET")
	api.HandleFunc("/shop/{id}", controllers.UpdateShop).Methods("PUT")
	api.HandleFunc("/shop/{id}", controllers.DeleteShop).Methods("DELETE")

	api.HandleFunc("/product", controllers.CreateProduct).Methods("POST")
	api.HandleFunc("/product/{id}", controllers.GetProduct).Methods("GET")
	api.HandleFunc("/product", controllers.GetAllProducts).Methods("GET")
	api.HandleFunc("/product/{id}", controllers.UpdateProduct).Methods("PUT")
	api.HandleFunc("/product/{id}", controllers.DeleteProduct).Methods("DELETE")

	api.HandleFunc("/order", controllers.CreateOrder).Methods("POST")
	api.HandleFunc("/order/{id}", controllers.GetOrder).Methods("GET")
	api.HandleFunc("/order", controllers.GetAllOrders).Methods("GET")
	api.HandleFunc("/order/{id}", controllers.UpdateOrder).Methods("PUT")
	api.HandleFunc("/order/{id}", controllers.DeleteOrder).Methods("DELETE")

	api.HandleFunc("/orderitems", controllers.CreateOrderItem).Methods("POST")
	api.HandleFunc("/orderitems/{id}", controllers.GetOrderItem).Methods("GET")
	api.HandleFunc("/orderitems", controllers.GetAllOrderItems).Methods("GET")
	api.HandleFunc("/orderitems/{id}", controllers.UpdateOrderItem).Methods("PUT")
	api.HandleFunc("/orderitems/{id}", controllers.DeleteOrderItem).Methods("DELETE")

	api.HandleFunc("/ratings", controllers.CreateRating).Methods("POST")
	api.HandleFunc("/ratings/{id}", controllers.GetRating).Methods("GET")
	api.HandleFunc("/ratings", controllers.GetAllRatings).Methods("GET")
	api.HandleFunc("/ratings/{id}", controllers.UpdateRating).Methods("PUT")
	api.HandleFunc("/ratings/{id}", controllers.DeleteRating).Methods("DELETE")

	return router
}
