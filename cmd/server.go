package cmd

import (
	"evermos-hasby/ecommerce/internal/config"
	"evermos-hasby/ecommerce/internal/controllers"
	"evermos-hasby/ecommerce/internal/interfaces"
	"evermos-hasby/ecommerce/internal/repository"
	"evermos-hasby/ecommerce/internal/services"
	"evermos-hasby/ecommerce/utils"
	"log"
)

func ServeHTTP() {
	r := NewRouter()
	err := r.Listen(":" + utils.GetEnv("PORT", "8000"))
	if err != nil {
		log.Fatal(err)
	}
}

type Dependency struct {
	UserRepository interfaces.IUserRepository
	AuthAPI        *controllers.AuthController
	AccountAPI     *controllers.AccountController
	HealthCheckAPI *controllers.HealthCheckController
	StoreAPI       *controllers.StoreController
	AddressAPI     *controllers.AddressController
	CategoryAPI    *controllers.CategoryController
	ProductAPI     *controllers.ProductController
	TransactionAPI *controllers.TransactionController
	UploadAPI      *controllers.UploadController
	RegionAPI      *controllers.RegionController
}

func dependencyInject() Dependency {
	// 1. Repositories
	userRepo := &repository.UserRepository{DB: config.DB}
	storeRepo := &repository.StoreRepository{DB: config.DB}
	addressRepo := &repository.AddressRepository{DB: config.DB}
	categoryRepo := &repository.CategoryRepository{DB: config.DB}
	productRepo := &repository.ProductRepository{DB: config.DB}
	transactionRepo := &repository.TransactionRepository{DB: config.DB}

	// 2. Services
	healthCheckSvc := &services.HealthCheck{}
	authSvc := &services.AuthService{UserRepository: userRepo, DB: config.DB, Redis: config.Redis}
	accountSvc := &services.AccountService{Repository: userRepo, DB: config.DB}
	storeSvc := &services.StoreService{StoreRepository: storeRepo, DB: config.DB}
	addressSvc := &services.AddressService{AddressRepository: addressRepo}
	categorySvc := &services.CategoryService{CategoryRepository: categoryRepo, DB: config.DB}
	productSvc := &services.ProductService{Repository: productRepo, DB: config.DB, Redis: config.Redis}
	transactionSvc := &services.TransactionService{Repository: transactionRepo, DB: config.DB}
	uploadSvc := &services.UploadService{}
	regionSvc := &services.RegionService{}

	// 3. Controllers
	healthCheckAPI := &controllers.HealthCheckController{HealthCheckService: healthCheckSvc}
	authAPI := &controllers.AuthController{AuthService: authSvc}
	accountAPI := &controllers.AccountController{Service: accountSvc}
	storeAPI := &controllers.StoreController{StoreService: storeSvc}
	addressAPI := &controllers.AddressController{AddressService: addressSvc}
	categoryAPI := &controllers.CategoryController{CategoryService: categorySvc}
	productAPI := &controllers.ProductController{Service: productSvc}
	transactionAPI := &controllers.TransactionController{Service: transactionSvc}
	uploadAPI := &controllers.UploadController{Service: uploadSvc}
	regionAPI := &controllers.RegionController{RegionService: regionSvc}

	return Dependency{
		UserRepository: userRepo,
		AuthAPI:        authAPI,
		AccountAPI:     accountAPI,
		HealthCheckAPI: healthCheckAPI,
		StoreAPI:       storeAPI,
		AddressAPI:     addressAPI,
		CategoryAPI:    categoryAPI,
		ProductAPI:     productAPI,
		TransactionAPI: transactionAPI,
		UploadAPI:      uploadAPI,
		RegionAPI:      regionAPI,
	}
}
