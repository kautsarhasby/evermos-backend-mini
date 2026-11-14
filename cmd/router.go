package cmd

import (
	"evermos-hasby/ecommerce/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"time"
)

func NewRouter() *fiber.App {
	dependency := dependencyInject()

	r := fiber.New()
	r.Use(cors.New(cors.Config{
		AllowOrigins:  "*",
		AllowMethods:  "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders:  "Origin, Content-Type, Accept, Authorization",
		ExposeHeaders: "Content-Length",
	}))

	provcity := r.Group("/provcity")
	provcity.Get("/listprovincies", dependency.RegionAPI.GetListProvincies)
	provcity.Get("/listcities/:prov_id", dependency.RegionAPI.GetListCities)
	provcity.Get("/detailprovince/:prov_id", dependency.RegionAPI.GetDetailProvince)
	provcity.Get("/detailcity/:city_id", dependency.RegionAPI.GetDetailCity)

	auth := r.Group("/auth")
	auth.Post("/register", limiter.New(limiter.Config{Max: 10, Expiration: 1 * time.Minute}), dependency.AuthAPI.Register)
	auth.Post("/login", limiter.New(limiter.Config{Max: 10, Expiration: 1 * time.Minute}), dependency.AuthAPI.Login)
	auth.Post("/refresh", middleware.RequireRefreshAuth, dependency.AuthAPI.Refresh)
	auth.Post("/logout", middleware.RequireAuth, dependency.AuthAPI.Logout)

	// Profile
	auth.Get("/profile", middleware.RequireAuth, dependency.AuthAPI.GetProfile)
	auth.Put("/profile", middleware.RequireAuth, dependency.AuthAPI.UpdateProfile)

	api := r.Group("/api", middleware.RequireAuth)
	// Store
	api.Post("/store", dependency.StoreAPI.CreateStore)
	api.Get("/store", dependency.StoreAPI.GetStore)
	api.Put("/store", dependency.StoreAPI.UpdateStore)

	// Address
	api.Post("/address", dependency.AddressAPI.CreateAddress)
	api.Get("/addresses", dependency.AddressAPI.List)
	api.Put("/addresses/:id", dependency.AddressAPI.Update)
	api.Delete("/addresses/:id", dependency.AddressAPI.Delete)

	// Products
	api.Post("/products", dependency.ProductAPI.Create)
	api.Get("/products", dependency.ProductAPI.List)
	api.Get("/products/:id", dependency.ProductAPI.Get)
	api.Put("/products/:id", dependency.ProductAPI.Update)
	api.Delete("/products/:id", dependency.ProductAPI.Delete)

	// Transactions
	api.Post("/transactions", dependency.TransactionAPI.Create)
	api.Get("/transactions", dependency.TransactionAPI.List)

	// Upload
	api.Post("/upload", dependency.UploadAPI.Upload)

	// Categories
	admin := api.Group("", middleware.RequireRole("admin"))
	admin.Get("/accounts", dependency.AccountAPI.List)
	admin.Get("/accounts/:id", dependency.AccountAPI.Get)
	admin.Post("/accounts", dependency.AccountAPI.Create)
	admin.Put("/accounts/:id", dependency.AccountAPI.Update)
	admin.Delete("/accounts/:id", dependency.AccountAPI.Delete)

	admin.Post("/categories", dependency.CategoryAPI.CreateCategory)
	admin.Put("/categories/:id", dependency.CategoryAPI.UpdateCategory)
	admin.Delete("/categories/:id", dependency.CategoryAPI.DeleteCategory)

	return r
}
