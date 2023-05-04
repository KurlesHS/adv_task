package main

import (
	"context"
	"fmt"
	"kurles/adv_task/configs"
	"kurles/adv_task/pkg/handler"
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

	ctx := context.Background()

	h := handler.New(conf.ServicePort)
	err = h.Start()
	_ = err
	_ = ctx

}
