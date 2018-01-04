package discordrpc

import (
	"gopkg.in/natefinch/npipe.v2"
	"strconv"
	"time"
)

type Connection struct {
	ConnectionBase
	Conn      *npipe.PipeConn
	Connected bool
}

func (c *Connection) Open() error {

	for i := 0; i < 10; i++ {
		con, err := npipe.DialTimeout(`\\.\pipe\discord-ipc-`+strconv.Itoa(i), 1*time.Second)
		if err == nil {
			c.Conn = con
			c.Connected = true
			return nil
		}
	}

	return ErrorDiscordNotFound
}

func (c *Connection) isOpen() bool {
	return c.Connected
}

func (c *Connection) Write(data []byte) (int, error) {
	conn := c.Conn
	tot, err := conn.Write(data)
	if err != nil {
		return tot, err
	} else if tot <= 0 {
		c.Close()
		return tot, ErrorNoData
	}
	return tot, nil
}

func (c *Connection) Read(data []byte) (int, error) {
	tot, err := c.Conn.Read(data)
	if err != nil {
		return tot, err
	} else if tot <= 0 {
		c.Close()
		return tot, ErrorNoData
	}
	return tot, nil
}

func (c *Connection) Close() error {
	err := c.Conn.Close()
	c.Conn = nil
	c.Connected = false
	return err
}
