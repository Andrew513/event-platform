package cryptoledger

type EventType string

const (
	Deposit EventType = "DEPOSIT"
	Withdraw EventType = "WITHDRAW"
)

type LedgerEvent struct {
	Account string `json:"account"`
	Amount float64 `json:"amount"`
}