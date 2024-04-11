package app

import (
	"log"

	"gitlab.simbirsoft/verify/t.galeev/internal/producer/config"
	"gitlab.simbirsoft/verify/t.galeev/internal/producer/core"
)

func Run() {
	conf := config.NewConfig()
	app, err := core.NewProducer(conf.Addresses)
	if err != nil {
		log.Fatal(err)
	}
	app.Run()
}
