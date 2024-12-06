package url

import (
	controllers "github.com/ChekoutGobiz/BackendChekout/controller"
	"github.com/ChekoutGobiz/BackendChekout/middleware"
	"github.com/gofiber/fiber/v2"
)

// SetupRoutes mendefinisikan semua rute aplikasi
func SetupRoutes(app *fiber.App) {
	// Grouping untuk lebih rapi
	api := app.Group("/api")

	// Authentication routes
	api.Post("/register", controllers.Register)
	api.Post("/login", controllers.Login)
	app.Post("/logout", controllers.Logout)  // Menambahkan route logout

	// Region routes
	api.Post("/regions", controllers.CreateRegion)
	api.Get("/regions", controllers.GetRegions)

	// Product routes - Protected by JWT middleware
	api.Use(middleware.VerifyJWT) // Hanya proteksi rute produk
	api.Post("/products", controllers.CreateProduct)
	api.Get("/products", controllers.GetProducts)

	// Cart routes - Protected by JWT middleware
	// Menambahkan rute untuk menambahkan item ke keranjang dengan verifikasi JWT
	api.Post("/cart", middleware.VerifyJWT, controllers.AddItemToCart)

	// Menambahkan rute untuk mengambil keranjang pengguna dengan verifikasi JWT
	api.Get("/cart", middleware.VerifyJWT, controllers.GetCart)

	// Menambahkan rute untuk memperbarui jumlah item dalam keranjang dengan verifikasi JWT
	api.Put("/cart", middleware.VerifyJWT, controllers.UpdateCartItem)

	// Menambahkan rute untuk menghapus item dari keranjang dengan verifikasi JWT
	api.Delete("/cart/:product_id", middleware.VerifyJWT, controllers.RemoveItemFromCart)
}
