// routes.go
package routes

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/kuyjajan/kuyjajan-backend/controllers"
)

// SetupRoutes sets up all the routes for the application
func SetupRoutes() http.Handler {
	r := mux.NewRouter()

	// Authentication routes
	r.HandleFunc("/register", controllers.Register).Methods("POST")
	r.HandleFunc("/login", controllers.Login).Methods("POST")

	// Region routes
	r.HandleFunc("/regions", controllers.CreateRegion).Methods("POST")
	r.HandleFunc("/regions", controllers.GetRegions).Methods("GET")

	// Product routes
	r.HandleFunc("/products", controllers.CreateProduct).Methods("POST")
	r.HandleFunc("/products", controllers.GetProducts).Methods("GET")
	r.HandleFunc("/products/region", controllers.GetProductsByRegion).Methods("GET")

	// Cart routes
	r.HandleFunc("/cart", controllers.AddToCart).Methods("POST")
	r.HandleFunc("/cart", controllers.GetCart).Methods("GET")
	r.HandleFunc("/cart/item", controllers.UpdateCartItem).Methods("PUT")
	r.HandleFunc("/cart/item", controllers.RemoveCartItem).Methods("DELETE")

	// Configure CORS
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://127.0.0.1:5501"}), // Allow Live Server origin
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"}), // Allowed methods
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	return corsHandler(r)
}
