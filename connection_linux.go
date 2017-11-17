package golang_discord_rpc

import (
	"net"
	"strconv"
	"os"
)

type Connection struct {
	ConnectionBase
	Conn net.Conn
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
		con, err := net.Dial("unix", path + "/discord-ipc-" + strconv.Itoa(i))
		if err == nil {
			c.Conn = con
			c.Connected = true
			return nil
		}
	}

	return ErrorDiscordNotFound

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
