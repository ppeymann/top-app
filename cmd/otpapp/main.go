package main

import (
	"fmt"
	"log"
	"os"
	"time"

	kitLog "github.com/go-kit/log"
	"github.com/ppeymann/top-app.git/config"
	"github.com/ppeymann/top-app.git/env"
	"github.com/ppeymann/top-app.git/server"
	pg "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	now := time.Now().UTC()

	base := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC).Unix()
	start := time.Date(now.Year(), now.Month(), now.Day(), 7, 35, 0, 0, time.UTC).Unix()
	end := time.Date(now.Year(), now.Month(), now.Day(), 23, 30, 0, 0, time.UTC).Unix()

	fmt.Println("date:", base, "start:", start, "end:", end)

	config, err := config.NewConfiguration()
	if err != nil {
		log.Fatal(err)

		return
	}

	_, err = gorm.Open(pg.Open(env.GetEnv("DSN", "")), &gorm.Config{SkipDefaultTransaction: true})
	if err != nil {
		log.Fatal(err)

		return
	}

	// configuration logger
	var logger kitLog.Logger
	logger = kitLog.NewJSONLogger(kitLog.NewSyncWriter(os.Stderr))
	logger = kitLog.With(logger, "ts", kitLog.DefaultTimestampUTC)

	// Service Logger
	sl := kitLog.With(logger, "component", "http")

	// Server instance
	svr := server.NewServer(sl, config)

	// listen and serve...
	svr.Listen()

}
