package dto

import (
	"order-service/constants"
	"time"

	"github.com/google/uuid"
)

type PaymentData struct {
	OrderID   uuid.UUID                     `json:"orderID"`
	PaymentID uuid.UUID                     `json:"paymentID"`
	Status    constants.PaymentStatusString `json:"status"`
	ExpiredAt *time.Time                    `json:"expiredAt"`
	PaidAt    *time.Time                    `json:"paidAt"`
}

type PaymentContent struct {
	Event    KafkaEvent             `json:"event"`
	Metadata KafkaMetaData          `json:"metadata"`
	Body     KafkaBody[PaymentData] `json:"body"`
}
