package app

import (
	restapi "api/src/api/rest"
	proccesses "api/src/processes"
	"api/src/shared/db"
	"api/src/shared/retry"
	"context"
	"os"
	"time"
)

func Start() {
	retry.Exec(context.TODO(), func() error {
		return db.Connect()
	}, retry.Options{Retries: -1, Delay: time.Second * 2})

	go proccesses.Start()

	restapi.Engine.Run(":" + os.Getenv("PORT"))
}
