package golang_discord_rpc_test

import (
	"fmt"
	"github.com/SilverCory/golang-discord-rpc"
	"testing"
)

func TestMeme(t *testing.T) {
	win := &golang_discord_rpc.ConnectionWindows{}
	fmt.Println(win.Open())

	// TODO handshaking ect ect

}