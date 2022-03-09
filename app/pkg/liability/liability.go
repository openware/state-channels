package liability

import (
	"app/models/broker"
	"bytes"
	"encoding/gob"
	"time"

	"github.com/shopspring/decimal"
)

type Type string

const (
	REQ     Type = "request"
	ACK     Type = "acknowledged"
	ERR     Type = "error"
	Unknown Type = "?"
)

type Liability struct {
	From            *broker.Broker
	To              *broker.Broker
	Type            Type
	ToCurrencyArray []string
	ToAmountArray   []decimal.Decimal
	CreatedAt       time.Time
}

func (l *Liability) EncodeToBytes() ([]byte, error) {
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)

	err := enc.Encode(l)
	if err != nil {
		return []byte{}, err
	}

	return buf.Bytes(), nil
}
