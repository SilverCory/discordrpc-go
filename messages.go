package golang_discord_rpc

type HandshakeMessage struct {
	Version int `json:"v"`
	ApplicationID string `json:"client_id"`
}

type CommandEventMessage struct {
	Command string `json:"cmd"`
	Event string `json:"evt"`
}