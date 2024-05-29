package common

import (
	"github.com/pusher/pusher-http-go/v5"
	"os"
)

var PusherClient pusher.Client

func NewPusherClient() {
	PusherClient = pusher.Client{
		AppID:   os.Getenv("PUSHER_APP_ID"),
		Key:     os.Getenv("PUSHER_APP_KEY"),
		Secret:  os.Getenv("PUSHER_APP_SECRET"),
		Cluster: os.Getenv("PUSHER_APP_CLUSTER"),
		Secure:  true,
	}

	//	data := map[string]string{"message": "hello world"}
	//
	//if err := common.PusherClient.Trigger("my-channel", "my-event", data); err != nil {
	//	fmt.Println(err.Error())
	//}
}
