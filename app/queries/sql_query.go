package queries

import (
	"fiber-template/app/models"
	log "fiber-template/pkg/utils/logger"

	"gorm.io/gorm"
)

type SqlQuery struct {
	Log *log.User
}

func (s *SqlQuery) Update(db interface{}, t *models.Example) error {
	return db.(*gorm.DB).Table("example").Create(t).Error
}

func (s *SqlQuery) SelectAll(db interface{}, t *models.Example) (rv []*models.ExampleTable, err error) {
	err = db.(*gorm.DB).Raw(`
		select *
		from example
		where payload = ?;
	`, t.Payload).Scan(&rv).Error
	return
}

func (s *SqlQuery) SSO(db interface{}, m *models.SsoUser) error {
	var cnt int64
	_ = db.(*gorm.DB).Raw(`
		select count(sso_user.seq)
		from sso_user
		where username = ? and organization = ?;
	`, m.Username, m.Organization).Count(&cnt)
	if cnt > 0 {
		return gorm.ErrRegistered
	}
	return db.(*gorm.DB).Table("sso_user").Create(m).Error
}

func (s *SqlQuery) GetAccountInfo(db interface{}, m *models.SsoUser) error {
	return db.(*gorm.DB).Raw(`
		select organization, username, user_index
		from sso_user
		where status = 1 and user_uuid =?
		limit 1;`,
		m.UserUuid).Scan(&m).Error
}

func (s *SqlQuery) GetOrgList(db interface{}, m *[]models.Organization) error {
	return db.(*gorm.DB).Raw(`
		select *
		from organization
		where status = 1;
	`).Scan(m).Error
}
