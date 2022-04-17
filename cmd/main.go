package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/cagrikilicoglu/shopping-basket/internal/models"
	"github.com/cagrikilicoglu/shopping-basket/internal/models/cart"
	"github.com/cagrikilicoglu/shopping-basket/internal/models/category"
	"github.com/cagrikilicoglu/shopping-basket/internal/models/item"
	"github.com/cagrikilicoglu/shopping-basket/internal/models/order"
	"github.com/cagrikilicoglu/shopping-basket/internal/models/product"
	"github.com/cagrikilicoglu/shopping-basket/internal/models/response"
	"github.com/cagrikilicoglu/shopping-basket/internal/models/user"
	"github.com/cagrikilicoglu/shopping-basket/pkg/auth"
	"github.com/cagrikilicoglu/shopping-basket/pkg/config"
	"github.com/cagrikilicoglu/shopping-basket/pkg/database"
	"github.com/cagrikilicoglu/shopping-basket/pkg/graceful"
	"github.com/cagrikilicoglu/shopping-basket/pkg/logging"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func main() {
	log.Println("Shopping basket service started")

	// Load environment to detect current app environemnt
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file.")
	}

	// Load configuration depending on app environment
	configFile := fmt.Sprintf("./pkg/config/%s", os.Getenv("APP_ENV"))
	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		log.Fatalf("Loadconfig failed, %v", err)
	}

	// Set a globalLogger
	logging.NewLogger(cfg)
	defer logging.Close()

	// Creating db
	db := database.Connect(cfg)

	// Close db connection
	sqlDb, err := db.DB()
	if err != nil {
		log.Fatal("Database connection cannot be closed.")
	}
	defer sqlDb.Close()

	log.Println("Postgress connected")

	router := gin.Default()
	// TODO custom format ekle -- middleware klasörüne al
	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// your custom format
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.ServerConfig.Port),
		ReadTimeout:  time.Duration(int64(cfg.ServerConfig.ReadTimeoutSecs) * int64(time.Second)),
		WriteTimeout: time.Duration(int64(cfg.ServerConfig.WriteTimeoutSecs) * int64(time.Second)),
		Handler:      router,
	}

	// TODO rooterları düzelt
	baseRouter := router.Group(cfg.ServerConfig.RoutePrefix)
	productRooter := baseRouter.Group("/products")
	categoryRouter := baseRouter.Group("/categories")
	cartRouter := baseRouter.Group("/cart")

	productRepo := product.NewProductRepository(db)
	productRepo.Migration()
	product.NewProductHandler(productRooter, productRepo, cfg)

	categoryRepo := category.NewCategoryRepository(db)
	categoryRepo.Migration()
	category.NewCategoryHandler(categoryRouter, categoryRepo, cfg)

	auth := auth.NewAuthenticator(cfg)
	fmt.Printf("auth: %s", cfg.JWTConfig.SecretKey)

	userRepo := user.NewUserRepository(db)
	userRepo.Migration()
	user.NewUserHandler(baseRouter, userRepo, auth) // TODO base routter değiştiilebilir

	// SampleQueries(*productRepo)

	cartRepo := cart.NewCartRepository(db)
	orderRepo := order.NewOrderRepository(db)
	itemRepo := item.NewItemRepository(db)
	cartRepo.Migration()
	orderRepo.Migration()

	itemRepo.Migration()
	itemService := item.NewItemService(itemRepo, *productRepo)

	cart.NewCartHandler(cartRouter, cartRepo, itemService, cfg)

	order.NewOrderHandler(baseRouter, orderRepo, cartRepo, itemService, cfg)

	// CreateAdmin(userRepo)
	// TODO aşağıdaki fonksiyonu kontrol et
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	baseRouter.GET("/health", checkHealth)
	// TODO aşağıyı anonymous func gibi handle etmeli?
	// baseRouter.GET("/ready", checkReady())
	graceful.Shutdown(srv, time.Duration(int64(cfg.ServerConfig.ShutdownTimeoutSecs)*int64(time.Second)))

}

func checkHealth(c *gin.Context) {
	response.RespondWithJson(c, http.StatusOK, nil)
}

func checkReady(c *gin.Context, db *gorm.DB) {
	DB, err := db.DB()
	if err != nil {
		zap.L().Fatal("cannot get sql database instance", zap.Error(err))
		response.RespondWithError(c, err)
		return
	}
	if err := DB.Ping(); err != nil {
		zap.L().Fatal("cannot ping database", zap.Error(err))
		response.RespondWithError(c, err)
		return
	}
	response.RespondWithJson(c, http.StatusOK, nil)
}

//////////////////////////------------------------------------//////////////////////////
// createAdmin creates an admin to manipulate the database
// the function is here for test purposes and should only be run in first usage
func CreateAdmin(userRepo *user.UserRepository) {
	admin := "admin1234.com"
	adminPass := getHash([]byte("admin1234"))
	u := models.User{
		Email:     &admin,
		Password:  &adminPass,
		FirstName: admin,
		LastName:  admin,
		Role:      "admin",
		ZipCode:   admin,
	}
	userRepo.Create(&u)
}

func getHash(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err) // TODO başka şekilde handle et
	}
	return string(hash)
}
