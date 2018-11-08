package main

import (
	"FCMTestClient/app"
	"FCMTestClient/config"
	"context"
	log "github.com/sirupsen/logrus"
)

func init() {
	config.LoadEnvironment()
	if config.IsDevelopmentEnv() {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.WarnLevel)
	}
}

//const app = "test-431ca"

func main() {
	ctx := context.Background()
	a := app.NewApp()
	a.ConnectNats()						//Нужны ли горутины?
	a.ConnectFirebaseMessaging(ctx)		//Нужны ли горутины?

	a.GetMessage()
}
