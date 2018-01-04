package discordrpc_test

import (
	"encoding/json"
	"fmt"
	"github.com/SilverCory/golang_discord_rpc"
	"os"
	"testing"
	"time"
)

func TestMeme(t *testing.T) {

	time.Sleep(time.Second * 3)

	//316245861074206730
	win := golang_discord_rpc.NewRPCConnection("368924690946850817")
	err := win.Open()
	if err != nil {
		fmt.Println(err)
	}

	str, err := win.Read()
	fmt.Println(err)
	fmt.Println(str)

	time.Sleep(time.Second * 3)

	for {
		fmt.Println(os.Getpid())
		presence := &golang_discord_rpc.CommandRichPresenceMessage{
			CommandMessage: golang_discord_rpc.CommandMessage{Command: "SET_ACTIVITY"},
			Args: &golang_discord_rpc.RichPresenceMessageArgs{
				Pid: os.Getpid(),
				Activity: &golang_discord_rpc.Activity{
					Details:  "Dean",
					State:    "Proud To Be A Developer ",
					Instance: false,
					Assets: &golang_discord_rpc.Assets{
						LargeText:    "Unknown Album",
						LargeImageID: "unknown",
						SmallText:    "Dank Memes",
						SmallImageID: "default",
					},
				},
			},
		}

		presence.SetNonce()
		data, err := json.Marshal(presence)

		if err != nil {
			fmt.Println(err)
			continue
		}

		err = win.Write(string(data))
		if err != nil {
			fmt.Println(err)
			continue
		}

		str, err := win.Read()
		fmt.Println(err)
		fmt.Println(str)

		fmt.Println("---\nDone?")
		time.Sleep(time.Second * 20)
	}

}
