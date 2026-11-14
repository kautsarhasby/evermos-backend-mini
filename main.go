package main

import (
	"evermos-hasby/ecommerce/cmd"
	"evermos-hasby/ecommerce/internal/config"
	"evermos-hasby/ecommerce/utils"
)

func main() {
	utils.SetupLogger()
	utils.SetupEnv()
	config.SetupDatabase()
	config.SeedData()
	config.SetupRedis()

	cmd.ServeHTTP()
}
