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

func main() {
	ctx := context.Background()
	a := app.NewApp()
	a.ConnectNats()
	a.ConnectFirebaseMessaging(ctx)

	a.GetMessage()
}
