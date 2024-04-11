package app

import (
	"log"

	"gitlab.simbirsoft/verify/t.galeev/internal/consumer/config"
	"gitlab.simbirsoft/verify/t.galeev/internal/consumer/core"
)

func Run() {
	conf := config.NewConfig()
	app, err := core.NewConsumer(conf.Addresses)
	if err != nil {
		log.Fatal(err)
	}
	app.Run()
}
