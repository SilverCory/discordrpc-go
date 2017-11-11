package golang_discord_rpc

type ConnectionBase interface {
	isOpen() bool
	Open() bool
	Close() bool
	Write(data []byte, length uint) bool
	Read(data []byte, length uint) bool
}
