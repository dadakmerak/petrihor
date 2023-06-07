package main

import (
	"fmt"

	"github.com/dadakmerak/petrihor/pkg/config"
	"github.com/dadakmerak/petrihor/pkg/database"
	"github.com/dadakmerak/petrihor/pkg/querier"
	"github.com/dadakmerak/petrihor/pkg/sqlx"
	"github.com/dadakmerak/petrihor/routes"
	"github.com/dadakmerak/petrihor/services"
)

func main() {

	config, err := config.LoadConfig()
	if err != nil {
		fmt.Println("err:", err)
	}

	db := new(database.DB).InitDatabase(&config)
	query := querier.NewAPIGen()
	repos := sqlx.NewRepository(db, query)
	services := services.NewService(repos)
	handlers := routes.NewHandler(services)

	router := handlers.SetupRouter()

	fmt.Printf("test :%v\n", config.Port)
	router.Run()
}
