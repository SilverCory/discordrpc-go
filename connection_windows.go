package golang_discord_rpc

import (
	"errors"
	"fmt"
	"github.com/natefinch/npipe"
	"strconv"
)

type Connection struct {
	ConnectionBase
	Conn      *npipe.PipeConn
	Connected bool
}

func (c Connection) Open() error {

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

func (c Connection) Write(data []byte) (int, error) {
	tot, err := c.Conn.Write(data)
	if err != nil {
		return tot, err
	} else if tot <= 0 {
		// TODO c.Close()
		return tot, ErrorNoData
	}
	return tot, nil
}

func (c Connection) Read(data []byte) (int, error) {
	tot, err := c.Conn.Read(data)
	if err != nil {
		return tot, err
	} else if tot <= 0 {
		// TODO c.Close()
		return tot, ErrorNoData
	}
	return tot, nil
}

func (c *Connection) Close() error {
	err := c.Conn.Close()
	c.Conn = nil
	return err
}
