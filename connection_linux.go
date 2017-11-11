package golang_discord_rpc

import (
	"net"
	"strconv"
	"os"
)

type ConnectionLinux struct {
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

func (c *ConnectionLinux) Open() bool {
	path := GetTempPath()
	for i := 0; i < 10; i++ {
		con, err := net.Dial("unix", path + "/discord-ipc-" + strconv.Itoa(i))
		if err == nil {
			c.Conn = con
			c.Connected = true
			return true
		}
	}

	return true

}

func (c *ConnectionLinux) Write(data []byte, length uint) bool {
	_, err := c.Conn.Write(data)
	return err == nil
}

func (c *ConnectionLinux) Read(data []byte, length uint) bool {
	// TODO
	return false
}
