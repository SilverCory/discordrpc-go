package golang_discord_rpc

import (
	"encoding/json"
	"fmt"
)

type RichPresenceMessageArgs struct {
	Pid      int       `json:"pid"`
	Activity *Activity `json:"activity"`
}

// A Game struct holds the name of the "playing .." game for a user
type Activity struct {
	Details    string     `json:"details,omitempty"`
	State      string     `json:"state,omitempty"`
	TimeStamps TimeStamps `json:"timestamps,omitempty"`
	Assets     Assets     `json:"assets,omitempty"`
	Instance   bool       `json:"instance,omitempty"`
}

// A TimeStamps struct contains start and end times used in the rich presence "playing .." Game
type TimeStamps struct {
	EndTimestamp   int64 `json:"end,omitempty"`
	StartTimestamp int64 `json:"start,omitempty"`
}

// UnmarshalJSON unmarshals JSON into TimeStamps struct
func (t *TimeStamps) UnmarshalJSON(b []byte) error {
	temp := struct {
		End   float64 `json:"end,omitempty"`
		Start float64 `json:"start,omitempty"`
	}{}
	err := json.Unmarshal(b, &temp)
	if err != nil {
		return err
	}
	t.EndTimestamp = int64(temp.End)
	t.StartTimestamp = int64(temp.Start)
	return nil
}

// An Assets struct contains assets and labels used in the rich presence "playing .." Game
type Assets struct {
	LargeImageID string `json:"large_image,omitempty"`
	SmallImageID string `json:"small_image,omitempty"`
	LargeText    string `json:"large_text,omitempty"`
	SmallText    string `json:"small_text,omitempty"`
}

// TODO implement party and secrets.

var nonceVal = 0

type Nonce struct {
	Nonce string `json:"nonce"`
}

func (n *Nonce) SetNonce() {
	nonceVal += 1
	n.Nonce = fmt.Sprintf("%032d", nonceVal)
}
