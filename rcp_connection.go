package golang_discord_rpc

import (
	"errors"
	"encoding/binary"
	"bytes"
)

const MaxRpcFrameSize = 64 * 1024

var (
	ErrorInvalidState = errors.New("invalid state on read/write")
	ErrorReadCorrupt = errors.New("read corrupted")
	ErrorWriteEncode = errors.New("write encode error")
)

type ErrorCode int8

const (
	CodeSuccess ErrorCode = iota
	CodePipeClosed
	CodeReadCorrupt
)

type OpCode uint32
const (
	OpCodeHandshake OpCode = iota
	OpCodeFrame
	OpCodeClose
	OpCodePing
	OpCodePong
)

type MessageFrameHeader struct {
	OpCode OpCode
	Length uint32
}

type MessageFrame struct {
	MessageFrameHeader
	Message string
}

type State uint32
const (
	StateDisconnected State = iota
	StateSentHandshake
	StateAwaitingResponse
	StateConnected
)

type RcpConnection struct {
	Connection ConnectionBase
	State State
	ApplicationID string
	lastErrorCode ErrorCode
	lastErrorMessage string
}

func New(ApplicationID string) *RcpConnection {
	return &RcpConnection{
		Connection: NewConnection(),
		State: StateDisconnected,
		ApplicationID: ApplicationID,
		lastErrorCode: 0,
		lastErrorMessage: "",
	}
}

func (r *RcpConnection) Open() error {
	if r.State == StateConnected {
		return nil
	}

	if r.State == StateDisconnected {
		if err := r.Connection.Open(); err != nil {
			return err
		}
	}

	if r.State == StateSentHandshake {

	}

	return errors.New("unimplemented")

}

func (r *RcpConnection) Read() (string, error) {

	if r.State != StateConnected && r.State != StateSentHandshake {
		return "", ErrorInvalidState
	}

	var frame MessageFrame
	for {
		data, err := r.readData(8) // TODO sizeof MessageFrameHeader
		if err != nil {
			return "", err
		}

		if err := binary.Read(bytes.NewReader(data), binary.LittleEndian, &frame.MessageFrameHeader); err != nil {
			return "", err
		}

		if frame.Length > 0 {
			data, err := r.readData(frame.Length)
			if err != nil {
				r.lastErrorCode = CodeReadCorrupt
				r.lastErrorMessage = "Partial data in frame"
				return "", err
			}
			frame.Message = string(data)
		}

		switch frame.OpCode {
			case OpCodeClose:
				// TODO close message.
				// TODO close.
				return "", nil
			case OpCodeFrame:
				// TODO parse frame.
				return "", nil
			case OpCodePing:
				frame.OpCode = OpCodePong
				// TODO writeFrame
				break
			case OpCodePong:
				break
			case OpCodeHandshake:
			default:
				r.lastErrorMessage = "Bad ipc frame"
				r.lastErrorCode = CodeReadCorrupt
				// TODO r.Close()
				return "", ErrorReadCorrupt

		}

	}

	return frame.Message, nil

}

func (r *RcpConnection) readData(length uint32) (data []byte, err error) {
	data = make([]byte, length)
	err = r.Connection.Read(data)
	if err != nil {
		if !r.Connection.isOpen() {
			// TODO r.Close()
			r.lastErrorCode = CodePipeClosed
			r.lastErrorMessage = "Pipe closed"
		}
		return
	}
	return
}

func (r *RcpConnection) Write(data string) error {

	header := MessageFrameHeader{
		OpCode:OpCodeFrame,
		Length: uint32(len(data) + 8),
	}

	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, header); err != nil {
		return err
	}

	if _, err := buf.WriteString(data); err != nil {
		return err
	}

	if err := r.Connection.Write(buf.Bytes()); err != nil {
		// TODO r.Close()
		return ErrorWriteEncode
	}

	return nil

}