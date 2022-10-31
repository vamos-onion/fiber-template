package interfaces

import (
	"fiber-template/app/models"
)

type RdbQuery interface {
	Update(*models.Example) error
	SelectAll(*models.Example) ([]*models.ExampleTable, error)
	SSO(*models.SsoUser) error
	GetAccountInfo(*models.SsoUser) error
	GetOrgList(*[]models.Organization) error
}

func Update(r RdbQuery, e *models.Example) error {
	return r.Update(e)
}
func SelectAll(r RdbQuery, e *models.Example) ([]*models.ExampleTable, error) {
	return r.SelectAll(e)
}
func SSO(r RdbQuery, s *models.SsoUser) error {
	return r.SSO(s)
}
func GetAccountInfo(r RdbQuery, s *models.SsoUser) error {
	return r.GetAccountInfo(s)
}
func GetOrgList(r RdbQuery, s *[]models.Organization) error {
	return r.GetOrgList(s)
}
