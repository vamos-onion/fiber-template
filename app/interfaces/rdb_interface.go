package interfaces

import (
	"fiber-template/app/models"
)

type RdbQuery interface {
	Update(interface{}, *models.Example) error
	SelectAll(interface{}, *models.Example) ([]*models.ExampleTable, error)
	SSO(interface{}, *models.SsoUser) error
	GetAccountInfo(interface{}, *models.SsoUser) error
	GetOrgList(interface{}, *[]models.Organization) error
}

func Update(r RdbQuery, db interface{}, e *models.Example) error {
	return r.Update(db, e)
}
func SelectAll(r RdbQuery, db interface{}, e *models.Example) ([]*models.ExampleTable, error) {
	return r.SelectAll(db, e)
}
func SSO(r RdbQuery, db interface{}, s *models.SsoUser) error {
	return r.SSO(db, s)
}
func GetAccountInfo(r RdbQuery, db interface{}, s *models.SsoUser) error {
	return r.GetAccountInfo(db, s)
}
func GetOrgList(r RdbQuery, db interface{}, s *[]models.Organization) error {
	return r.GetOrgList(db, s)
}
