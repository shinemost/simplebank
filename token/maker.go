package token

import "time"

type Maker interface {
	CreateToken(username string, duartion time.Duration) (string, *Payload, error)

	VerifyToken(token string) (*Payload, error)
}
