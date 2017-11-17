package golang_discord_rpc_test

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/SilverCory/golang-discord-rpc"
	"log"
	"net"
	"os"
	"testing"
)

func TestMeme(t *testing.T) {
	win := golang_discord_rpc.NewConnection()
	fmt.Println(win.Open())

	// TODO handshaking ect ect

}
