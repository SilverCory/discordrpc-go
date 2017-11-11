package golang_discord_rpc

import (
	"github.com/natefinch/npipe"
	"strconv"
	"fmt"
)

type ConnectionWindows struct {
	ConnectionBase
	Conn *npipe.PipeConn
	Connected bool
}

func (c *ConnectionWindows) Open() bool {

	for i := 0; i < 10; i++ {
		con, err := npipe.Dial("\\\\.\\pipe\\discord-ipc-" + strconv.Itoa(i))
		if err == nil {
			c.Conn = con
			c.Connected = true
			return true
		}
	}

	return false
}