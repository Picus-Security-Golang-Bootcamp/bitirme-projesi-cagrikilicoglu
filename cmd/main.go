package main

import (
	"fmt"
	"log"
	"os"

	"github.com/cagrikilicoglu/shopping-basket/pkg/config"
	"github.com/cagrikilicoglu/shopping-basket/pkg/logging"
	"github.com/joho/godotenv"
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

	// DB := db.Connect(cfg)

}
