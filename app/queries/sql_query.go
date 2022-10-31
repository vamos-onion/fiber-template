package queries

import (
	"fiber-template/app/models"
	log "fiber-template/pkg/utils/logger"

	"gorm.io/gorm"
)

type SqlQuery struct {
	Log *log.User
	DB  *gorm.DB
}

func (s *SqlQuery) Update(t *models.Example) error {
	return s.DB.Table("example").Create(t).Error
}

func (s *SqlQuery) SelectAll(t *models.Example) (rv []*models.ExampleTable, err error) {
	err = s.DB.Raw(`
		select *
		from example
		where payload = ?;
	`, t.Payload).Scan(&rv).Error
	return
}

func (s *SqlQuery) SSO(m *models.SsoUser) error {
	var cnt int64
	_ = s.DB.Raw(`
		select count(sso_user.seq)
		from sso_user
		where username = ? and organization = ?;
	`, m.Username, m.Organization).Count(&cnt)
	if cnt > 0 {
		return gorm.ErrRegistered
	}
	return s.DB.Table("sso_user").Create(m).Error
}

func (s *SqlQuery) GetAccountInfo(m *models.SsoUser) error {
	return s.DB.Raw(`
		select organization, username, user_index
		from sso_user
		where status = 1 and user_uuid =?
		limit 1;`,
		m.UserUuid).Scan(&m).Error
}

func (s *SqlQuery) GetOrgList(m *[]models.Organization) error {
	return s.DB.Raw(`
		select *
		from organization
		where status = 1;
	`).Scan(m).Error
}
