package cryptoledger

type EventType string

const (
	Deposit EventType = "DEPOSIT"
	Withdrawal EventType = "WITHDRAWAL"
)

type LedgerEvent struct {
	Account string `json:"account"`
	Amount float64 `json:"amount"`
}