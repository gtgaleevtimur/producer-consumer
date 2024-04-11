package core

import "gitlab.simbirsoft/verify/t.galeev/pkg"

type Publisher interface {
	Run()
	Stop() error
	Send() error
	GenerateSync() pkg.Sync
}
