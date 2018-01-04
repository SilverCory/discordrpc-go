package discordrpc

import (
	"errors"
	"io"
)

var (
	ErrorDiscordNotFound = errors.New("could not find discord")
	ErrorNoData          = errors.New("no data")
)

type ConnectionBase interface {
	io.ReadWriteCloser
	isOpen() bool
	Open() error
}

func NewConnection() ConnectionBase {
	return &Connection{}
}
