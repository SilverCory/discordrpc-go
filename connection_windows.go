package golang_discord_rpc

import (
	"github.com/natefinch/npipe"
	"strconv"
	"fmt"
	"errors"
)

type Connection struct {
	ConnectionBase
	Conn *npipe.PipeConn
	Connected bool
}

func (c *Connection) Open() error {

	for i := 0; i < 10; i++ {
		con, err := npipe.Dial("\\\\.\\pipe\\discord-ipc-" + strconv.Itoa(i))
		if err == nil {
			c.Conn = con
			c.Connected = true
			return nil
		}
	}

	return
}

func (c *Connection) Write(data []byte) error {
	tot, err := c.Conn.Write(data)
	if err != nil {
		return err
	} else if tot <= 0 {
		// TODO c.Close()
		return ErrorNoData
	}
	return nil
}

func (c *Connection) Read(data []byte) error {
	tot, err := c.Conn.Read(data)
	if err != nil {
		return err
	} else if tot <= 0 {
		// TODO c.Close()
		return ErrorNoData
	}
	return nil
}