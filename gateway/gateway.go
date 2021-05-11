package gateway

type ClientRequest struct {
	AccountID   string      `json:"account_id"`
	TypeRequest TypeRequest `json:"type_request"`
}

type TypeRequest struct {
	Type   string `json:"type"`
	Amount uint64 `json:"amount"`
}

type Messenger interface {
	Publish(payload []byte) error
}
