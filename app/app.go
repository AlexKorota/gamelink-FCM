package app

import (
	"FCMTestClient/config"
	"context"
	"firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"fmt"
	"github.com/nats-io/go-nats"
	"log"
	"google.golang.org/api/option"
	"sync"
	"github.com/gogo/protobuf/proto"
	push "gamelink-go/proto_nats_msg"
)

type App struct {
	nc *nats.Conn
	fbm *messaging.Client
	mchan chan push.PushMsgStruct
}

func NewApp() App {
	mchan := make(chan push.PushMsgStruct)
	return App{mchan:mchan}
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

func(a *App) GetMessage()  {
	wg := sync.WaitGroup{}
	wg.Add(1)
	var msgStruct push.PushMsgStruct
	// Subscribe to updates
	_, err := a.nc.Subscribe(config.NatsChan, func(m *nats.Msg) {
		err := proto.Unmarshal(m.Data, &msgStruct)
		if err != nil {
			return
		}
		fmt.Println("msgStruct", msgStruct)
		a.mchan <- msgStruct
	})
	if err != nil {
		log.Fatal(err)
	}
	for {
		m := <- a.mchan
		go a.prepareMsg(m)
	}
	wg.Wait()
}

func (a *App)prepareMsg(msg push.PushMsgStruct) {
	for _, v := range msg.UserInfo {
		message := &messaging.Message{
			Data: map[string]string{
				"message": msg.Message,
			},
			Token:v.DeviceID,
		}
		fmt.Println("prepared", message)
		//resp, err := a.SendMessage(ctx, message)
		//if err != nil {
		//	log.Warn(err)
		//}
		// Response is a message ID string.
		//fmt.Println("Successfully sent message:", resp)
	}
}

func (a *App) SendMessage(ctx context.Context, msg *messaging.Message) (string, error) {
	resp, err := a.fbm.Send(ctx, msg)
	if err != nil {
		return "", nil
	}
	return resp, nil
}
