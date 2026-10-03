// @title           Toko Emas API
// @version         1.0
// @description     API sistem manajemen toko emas dengan PostgreSQL backend.
// @host            localhost:8000
// @BasePath        /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

package main

import (
	"flag"

	"toko-emas/internal/app"
)

func main() {
	confDir := flag.String("conf", "configs", "config directory")
	flag.Parse()
	app.Run(*confDir)
}
