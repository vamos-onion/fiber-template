package models

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
	AuthSeq      int64     `json:"auth_seq" gorm:"column:auth_seq"`
	AccountId    string    `json:"account_id" gorm:"column:account_id"`
	AccountPwd   string    `json:"account_pwd" gorm:"column:account_pwd"`
	AccountName  string    `json:"account_name" gorm:"column:account_name"`
	AccountEmail string    `json:"account_email" gorm:"column:account_email"`
	AccountUuid  uuid.UUID `json:"account_uuid" gorm:"column:account_uuid"`
	CreatedAt    time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"column:updated_at"`
	ConnectedAt  time.Time `json:"connected_at" gorm:"column:connected_at"`
	Status       bool      `json:"status" gorm:"column:status"`
}

type Register struct {
	AccountId    string `json:"account_id"`
	AccountPwd   string `json:"account_pwd"`
	AccountName  string `json:"account_name"`
	AccountEmail string `json:"account_email"`
}

type RegisterQuery struct {
	Account             //Account //`json:"account"`
	Organization string `json:"organization"`
}

type Login struct {
	AccountId  string `json:"account_id"`
	AccountPwd string `json:"account_pwd"`
}

type SSO struct {
	SSO string `json:"sso"`
}

type Organization struct {
	Seq          int64  `json:"seq"`
	Organization string `json:"organization"`
	OrgUuid      string `json:"org_uuid"`
	Status       bool   `json:"status"`
}

type SsoQuery struct {
	Account      Account      `json:"account"`
	Organization Organization `json:"organization"`
}

type Modify struct {
	AccountPwd   string `json:"account_pwd,omitempty"`
	AccountName  string `json:"account_name,omitempty"`
	AccountEmail string `json:"account_email,omitempty"`
}

type Withdrawal struct {
	AccountPwd string `json:"account_pwd"`
}
