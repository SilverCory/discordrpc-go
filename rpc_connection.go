package go_discordrpc

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
)

const RpcVersion = 1
const MaxRpcFrameSize = 64 * 1024

var (
	ErrorInvalidState = errors.New("invalid state on read/write")
	ErrorReadCorrupt  = errors.New("read corrupted")
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
	Message [MaxRpcFrameSize]byte
}

func (m *MessageFrame) GetMessage() string {
	return string(m.Message[:])
}

func (m *MessageFrame) SetMessage(str string) {
	copy(m.Message[:], str)
}

type State uint32

const (
	StateDisconnected State = iota
	StateSentHandshake
	StateAwaitingResponse
	StateConnected
)

type RPCConnection struct {
	io.Closer
	Connection       ConnectionBase
	State            State
	ApplicationID    string
	lastErrorCode    ErrorCode
	lastErrorMessage string
}

func NewRPCConnection(ApplicationID string) *RPCConnection {
	return &RPCConnection{
		Connection:       NewConnection(),
		State:            StateDisconnected,
		ApplicationID:    ApplicationID,
		lastErrorCode:    0,
		lastErrorMessage: "",
	}
}

func (r *RPCConnection) Open() error {
	if r.State == StateConnected {
		return nil
	}

	if r.State == StateDisconnected {
		if err := r.Connection.Open(); err != nil {
			return err
		}
	}

	if r.State == StateSentHandshake {
		str, err := r.Read()
		if err != nil {
			return err
		}

		msg := &CommandEventMessage{}
		if err := json.Unmarshal([]byte(str), msg); err != nil {
			return err
		}

		if strings.EqualFold(msg.Command, "DISPATCH") && strings.EqualFold(msg.Event, "READY") {
			r.State = StateConnected
			// TODO r.onConnect();
		}
	} else {
		data, err := json.Marshal(&HandshakeMessage{
			Version:       RpcVersion,
			ApplicationID: r.ApplicationID,
		})

		if err != nil {
			return err
		}

		if err := r.writeFrame(OpCodeHandshake, string(data)); err != nil {
			r.Close()
			return err
		} else {
			r.State = StateSentHandshake
			fmt.Println("Sent handshake.")
		}
	}

	return nil

}

func (r *RPCConnection) IsOpen() bool {
	return r.Connection != nil && r.State == StateConnected && r.Connection.isOpen()
}

func (r *RPCConnection) Read() (string, error) {

	if r.State != StateConnected && r.State != StateSentHandshake {
		return "", ErrorInvalidState
	}

	var frame MessageFrame
	for { // this is blocking?
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
			frame.SetMessage(string(data))
		}

		switch frame.OpCode {
		case OpCodeClose:
			return frame.GetMessage(), r.Close()
		case OpCodeFrame:
			// TODO parse frame.
			return frame.GetMessage(), nil
		case OpCodePing:
			r.writeFrame(OpCodePong, frame.GetMessage())
			break
		case OpCodePong:
			break
		case OpCodeHandshake:
		default:
			r.lastErrorMessage = "Bad ipc frame"
			r.lastErrorCode = CodeReadCorrupt
			r.Close()
			return frame.GetMessage(), ErrorReadCorrupt
		}

	}

	return frame.GetMessage(), nil

}

func (r *RPCConnection) readData(length uint32) (data []byte, err error) {
	data = make([]byte, length)
	_, err = r.Connection.Read(data)
	if err != nil {
		if !r.Connection.isOpen() {
			r.Close()
			r.lastErrorCode = CodePipeClosed
			r.lastErrorMessage = "Pipe closed"
		}
		return
	}
	return
}

func (r *RPCConnection) Write(data string) error {
	return r.writeFrame(OpCodeFrame, data)
}

func (r *RPCConnection) writeFrame(code OpCode, data string) error {
	header := MessageFrame{
		MessageFrameHeader: MessageFrameHeader{
			OpCode: code,
			Length: uint32(len(data)),
		},
	}

	header.SetMessage(data)

	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, header); err != nil {
		return err
	}

	if _, err := r.Connection.Write(buf.Bytes()[:header.Length+8]); err != nil {
		r.Close()
		return err
	}

	return nil
}

func (r *RPCConnection) Close() error {
	err := r.Connection.Close()
	r.State = StateDisconnected
	r.Connection = nil
	return err
}
