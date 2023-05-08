package main

import (
	"fmt"
	"kurles/adv_task/configs"
	"kurles/adv_task/pkg/handler"
	"kurles/adv_task/pkg/repository/postgres"
	"log"
	"os"
)

// https://pkg.go.dev/github.com/VolkovEgor/advertising-task#section-readme

// export DB_PASS=postgrespass && export DB_USER=postgres && export DB_NAME=test_db && export DB_HOST=localhost && export DB_PORT=5433 && make migrate_up
// set DB_PASS=postgrespass && set DB_USER=postgres && set DB_NAME=test_db && set DB_HOST=localhost && set DB_PORT=5433

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Println("panic occurred:", err)
			os.Exit(1)
		}
	}()
	conf, err := configs.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "read config error: %v\n", err)
		return
	}

	repo, err := postgres.New(conf.DBHost, conf.DBPort, conf.DBName, conf.DBUserName, conf.DBPassword)
	if err != nil {
		fmt.Fprintf(os.Stderr, "create repository error: %v\n", err)
		return
	}
	h := handler.New(&repo, conf.ServicePort)
	log.Fatal(h.Start())
}
