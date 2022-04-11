package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/cagrikilicoglu/shopping-basket/internal/models"
	"github.com/cagrikilicoglu/shopping-basket/internal/models/cart"
	"github.com/cagrikilicoglu/shopping-basket/internal/models/category"
	"github.com/cagrikilicoglu/shopping-basket/internal/models/product"
	"github.com/cagrikilicoglu/shopping-basket/internal/models/response"
	"github.com/cagrikilicoglu/shopping-basket/internal/models/user"
	"github.com/cagrikilicoglu/shopping-basket/pkg/auth"
	"github.com/cagrikilicoglu/shopping-basket/pkg/config"
	"github.com/cagrikilicoglu/shopping-basket/pkg/database"
	"github.com/cagrikilicoglu/shopping-basket/pkg/logging"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func main() {
	log.Println("Shopping basket service started")

	// TODO aşağıdaki fonksiyonu hangi environment'ta olduğumu anlamak için kullanıyorum bunun için daha iyi bir yol olaiblir.
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file.")
	}
	// // //

	configFile := fmt.Sprintf("./pkg/config/%s", os.Getenv("APP_ENV"))
	cfg, err := config.LoadConfig(configFile)

	if err != nil {
		log.Fatalf("Loadconfig failed, %v", err)
	}

	log.Println(cfg)

	// Set globalLogger
	logging.NewLogger(cfg)
	defer logging.Close()

	// Creating db

	db := database.Connect(cfg)

	// TODO farklı şekilde handle edilebilir
	// Close db connection
	sqlDb, err := db.DB()
	if err != nil {
		log.Fatal("Database connection cannot be closed.")
	}
	defer sqlDb.Close()

	log.Println("Postgress connected")

	//////------//////////

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
	baseRooter := router.Group(cfg.ServerConfig.RoutePrefix)
	productRooter := baseRooter.Group("/products")
	categoryRouter := baseRooter.Group("/categories")
	// authRouter := baseRooter.Group("/user") // TODO başından user'ı sil

	productRepo := product.NewProductRepository(db)
	productRepo.Migration()
	product.NewProductHandler(productRooter, productRepo)

	categoryRepo := category.NewCategoryRepository(db)
	categoryRepo.Migration()
	category.NewCategoryHandler(categoryRouter, categoryRepo)

	// auth.NewAuthHandler(authRouter, cfg)

	auth := auth.NewAuthenticator(cfg)
	fmt.Printf("auth: %s", cfg.JWTConfig.SecretKey)

	userRepo := user.NewUserRepository(db)
	userRepo.Migration()
	user.NewUserHandler(baseRooter, userRepo, auth) // TODO base routter değiştiilebilir

	cartRepo := cart.NewCartRepository(db)
	cartRepo.Migration()

	// SampleQueries(*productRepo)

	// TODO aşağıdaki fonksiyonu kontrol et
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	baseRooter.GET("/health", checkHealth)
	// TODO aşağıyı anonymous func gibi handle etmeli?
	// baseRooter.GET("/ready", checkReady())
	GracefulShutdown(srv, 15*time.Second)
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

func SampleQueries(productRepo product.ProductRepository) {
	nameStr := "nikeAirForce1"
	catNameStr := "Shoes"
	nikeAirForce1 := models.Product{
		Name:  &nameStr,
		Price: 30,
		// Price: models.Price{Amount: 90, CurrencyCode: "USD"},
		Stock: models.Stock{
			SKU:    "1234F",
			Number: 15,
			Status: "decreasing",
		},
		CategoryName: &catNameStr,
	}
	nameStr2 := "logiMouse"
	catNameStr2 := "Technology"
	logitechMouse := models.Product{
		Name:  &nameStr2,
		Price: 12.4,
		// Price: models.Price{Amount: 12, CurrencyCode: "USD"},
		Stock: models.Stock{
			SKU:    "9874U",
			Number: 90,
			Status: "enough",
		},
		CategoryName: &catNameStr2,
	}
	productRepo.Create(&nikeAirForce1)
	productRepo.Create(&logitechMouse)

	// result1, _ := productRepo.GetAll()

	// for _, v := range result1 {
	// 	log.Println(*v)
	// }

	// log.Println(productRepo.GetByName("Mo"))
	// result1, _ := productRepo.GetBySKU("1234F")
	// log.Println(*result1)
	// log.Println(productRepo.GetBySKU("1234F"))
	// log.Println(productRepo.Delete("9874U"))

	// nikeAirForce1.Price.Subtract(money.New(20, "USD"))
	// log.Println(productRepo.Update(&nikeAirForce1))
	// log.Println(productRepo.GetAll)

}

// TODO başka bir yere taşı
func GracefulShutdown(srv *http.Server, timeout time.Duration) {
	c := make(chan os.Signal, 1)

	// when there is a interrupt signal, relay it to the channel
	signal.Notify(c, os.Interrupt)

	// block until any signal is received by the channel
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// wait until the timeout deadline and shutdown the server if there is no connections. if there is no connection shutdown immediately
	srv.Shutdown(ctx)

	log.Println("shutting down the server")
	os.Exit(0)
}
