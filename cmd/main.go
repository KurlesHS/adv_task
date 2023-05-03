package main

import (
	"context"
	"fmt"
	"kurles/adv_task/configs"
	"kurles/adv_task/pkg/model"
	"kurles/adv_task/pkg/repository/postgres"
	"os"
)

// https://pkg.go.dev/github.com/VolkovEgor/advertising-task#section-readme

// export DB_PASS=postgrespass && export DB_USER=postgres && export DB_NAME=test_db && export DB_HOST=localhost && export DB_PORT=5433 && make migrate_up
// set DB_PASS=postgrespass && set DB_USER=postgres && set DB_NAME=test_db && set DB_HOST=localhost && set DB_PORT=5433

func main() {
	conf, err := configs.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "read config error: %v\n", err)
		return
	}

	repo, err := postgres.New(conf.DBHost, conf.DBPort, conf.DBName, conf.DBUserName, conf.DBPassword)
	_ = repo
	fmt.Printf("%v\n", conf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error opening database connection: %v", err)
	}

	ctx := context.Background()
	res, err := repo.GetAdverts(ctx, 1, model.Date, true)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error selecting adverts: %v", err)
	}
	_ = res

	a, err := repo.GetAdvert(ctx, 1)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error selecting advert 1: %v", err)
	}
	_ = a

	a, err = repo.GetAdvert(ctx, 2)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error selecting advert 2: %v", err)
	}
	_ = a

	da := model.DetailedAdvert{
		Title:       "inserted adv",
		Description: "inserted adv descr",
		Price:       1000,
		Photos:      []string{"link1", "link2"},
	}
	advId, err := repo.InsertAdvert(ctx, da)

	if err != nil {
		fmt.Fprintf(os.Stderr, "error inserting advert: %v", err)
	}
	_ = advId
}
