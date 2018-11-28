package app

import (
	"context"
	"firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"fmt"
	"gamelink-fcm/config"
	push "gamelink-go/proto_nats_msg"
	"github.com/gogo/protobuf/proto"
	"github.com/nats-io/go-nats"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/option"
)

type App struct {
	nc    *nats.Conn
	fbm   *messaging.Client
	mchan chan push.PushMsgStruct
}

func NewApp() App {
	mchan := make(chan push.PushMsgStruct)
	return App{mchan: mchan}
}

func (a *App) ConnectNats() {
	nc, err := nats.Connect(config.NatsDialAddress)
	if err != nil {
		log.Fatal(err)
	}
	a.nc = nc
}

func (a *App) ConnectFirebaseMessaging(ctx context.Context) {
	opt := option.WithCredentialsFile(config.ServiceKeyPath)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}
	client, err := app.Messaging(ctx)
	a.fbm = client
}

func (a *App) GetMessage() {
	var msgStruct push.PushMsgStruct
	// Subscribe to updates
	_, err := a.nc.Subscribe(config.NatsFirebaseChan, func(m *nats.Msg) {
		err := proto.Unmarshal(m.Data, &msgStruct)
		if err != nil {
			log.Fatal(err)
			return
		}
		a.mchan <- msgStruct
	})
	if err != nil {
		log.Fatal(err)
	}
	for {
		m := <-a.mchan
		go a.prepareAndSend(m)
	}
}

func (a *App) prepareAndSend(msg push.PushMsgStruct) {
	ctx, _ := context.WithCancel(context.Background())
	message := &messaging.Message{
		Data: map[string]string{
			"message": msg.Message,
		},

		Token: msg.UserInfo.DeviceID,
	}
	fmt.Println("prepared", message)
	resp, err := a.fbm.Send(ctx, message)
	if err != nil {
		log.Warn(err)
	}
	//Response is a message ID string.
	fmt.Println("Successfully sent message:", resp)
}
