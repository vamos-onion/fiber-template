package models

import (
	"github.com/google/uuid"
)

type Login struct {
	SSO string `json:"sso"`
}

type Organization struct {
	Seq          uint16 `json:"seq"`
	Organization string `json:"organization"`
	Status       bool   `json:"status"`
}

type SsoUser struct {
	Organization string    `json:"organization"`
	Username     string    `json:"username"`
	UserUuid     uuid.UUID `json:"user_uuid"`
	UserIndex    uint64    `json:"user_index"`
	Status       bool      `json:"status"`
}
