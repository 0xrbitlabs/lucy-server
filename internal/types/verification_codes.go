package types

import "time"

type VerificationCode struct {
	Id     string    `db:"id"`
	Code   string    `db:"code"`
	SentTo string    `db:"sent_to"`
	SentAt time.Time `db:"sent_at"`
	Used   bool      `db:"used"`
}
