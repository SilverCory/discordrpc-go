package discordrpc

import (
	"encoding/json"
	"os"
)

// TODO updateConnection
// TODO wait for io?
// TODO registering
// TODO reconnecting

type API struct {
	Connection *RPCConnection
}

func New(ApplicationID string) (*API, error) {
	api := &API{
		Connection: NewRPCConnection(ApplicationID),
	}

	return api, nil
}

func (a *API) Open() error {
	return a.Connection.Open()
}

func (a *API) IsOpen() bool {
	return a.Connection != nil && a.Connection.IsOpen()
}

func (a *API) GetLastErrorMessage() string {
	return a.Connection.lastErrorMessage
}

func (a *API) GetLastErrorCode() ErrorCode {
	return a.Connection.lastErrorCode
}

func (a *API) GetState() State {
	return a.Connection.State
}

func (a *API) SetRichPresence(activity *Activity) error {
	command := &CommandRichPresenceMessage{
		Args: &RichPresenceMessageArgs{
			Activity: activity,
			Pid:      os.Getpid(),
		},
		CommandMessage: CommandMessage{"SET_ACTIVITY"},
	}
	command.SetNonce()

	data, err := json.Marshal(command)
	if err != nil {
		return err
	}

	return a.Connection.Write(string(data))
}
