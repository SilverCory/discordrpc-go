package go_discordrpc

import (
	"fmt"
	"net"
	"os"
	"strconv"
)

type Connection struct {
	ConnectionBase
	Conn      net.Conn
	Connected bool
}

func GetTempPath() string {
	temp := os.Getenv("XDG_RUNTIME_DIR")
	if temp != "" {
		return temp
	}

	temp = os.Getenv("TMPDIR")
	if temp != "" {
		return temp
	}

	temp = os.Getenv("TMP")
	if temp != "" {
		return temp
	}

	temp = os.Getenv("TEMP")
	if temp != "" {
		return temp
	}

	return "/tmp"
}

func (c *Connection) Open() error {
	path := GetTempPath()
	for i := 0; i < 10; i++ {
		con, err := net.Dial("unix", path+"/discord-ipc-"+strconv.Itoa(i))
		if err == nil {
			c.Conn = con
			c.Connected = true
			return nil
		}
	}

	return ErrorDiscordNotFound

}

func (c *Connection) Write(data []byte) (int, error) {
	fmt.Printf("%X\n\n", data)
	tot, err := c.Conn.Write(data)
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
	return err
}
