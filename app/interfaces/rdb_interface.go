package interfaces

import (
	"fiber-template/app/models"
)

type RdbQuery interface {
	Update(*models.Example) error
	SelectAll(*models.Example) ([]*models.ExampleTable, error)
	SSO(*models.Organization, *models.Account) error
	GetAccountInfo(*models.Account) error
	GetOrgList(*[]models.Organization) error
	UpdateAccount(*models.RegisterQuery) error
	Login(*models.Account) error
	ModifyAccount(*models.Modify, string) error
	WithdrawalAccount(*models.Withdrawal, *models.Account) error
}

func WithdrawalAccount(r RdbQuery, n *models.Withdrawal, m *models.Account) error {
	return r.WithdrawalAccount(n, m)
}
func ModifyAccount(r RdbQuery, m *models.Modify, userUuid string) error {
	return r.ModifyAccount(m, userUuid)
}
func Login(r RdbQuery, m *models.Account) error {
	return r.Login(m)
}
func UpdateAccount(r RdbQuery, m *models.RegisterQuery) error {
	return r.UpdateAccount(m)
}
func Update(r RdbQuery, e *models.Example) error {
	return r.Update(e)
}
func SelectAll(r RdbQuery, e *models.Example) ([]*models.ExampleTable, error) {
	return r.SelectAll(e)
}
func SSO(r RdbQuery, s *models.Organization, n *models.Account) error {
	return r.SSO(s, n)
}
func GetAccountInfo(r RdbQuery, s *models.Account) error {
	return r.GetAccountInfo(s)
}
func GetOrgList(r RdbQuery, s *[]models.Organization) error {
	return r.GetOrgList(s)
}
