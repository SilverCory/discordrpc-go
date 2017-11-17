package golang_discord_rpc

import "errors"

var (
	ErrorDiscordNotFound = errors.New("could not find discord")
	ErrorNoData = errors.New("no data")
)

type ConnectionBase interface {
	isOpen() bool
	Open() error
	Close() error
	Write(data []byte) error
	Read(data []byte) error
}

func NewConnection() ConnectionBase {
	return Connection{}.ConnectionBase
}