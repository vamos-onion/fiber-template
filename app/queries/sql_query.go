package queries

import (
	"fiber-template/app/models"
	"fiber-template/pkg/configs"
	log "fiber-template/pkg/utils/logger"
	"strings"

	"gorm.io/gorm"
)

type SqlQuery struct {
	Log *log.User
	DB  *gorm.DB
}

func (s *SqlQuery) ModifyAccount(m *models.Modify, userUuid string) error {
	var pwd, name, email string
	if m.AccountPwd != "" {
		pwd = "account_pwd = '" + m.AccountPwd + "',"
	}
	if m.AccountName != "" {
		name = "account_name = '" + m.AccountName + "',"
	}
	if m.AccountEmail != "" {
		email = "account_email = '" + m.AccountEmail + "',"
	}
	return s.DB.Exec(`
		UPDATE account
		SET `+pwd+name+email+`
			updated_at = NOW()
		WHERE account_uuid = ?`,
		userUuid).Error
}

func (s *SqlQuery) UpdateAccount(m *models.RegisterQuery) error {
	err := s.DB.Exec(`
		insert into account # (auth_seq, account_id, account_pwd, account_name, account_email, account_uuid, created_at, updated_at)
		set
			auth_seq = (select seq
			from organization
			where organization = ?),
			account_id = ?,
			account_pwd = ?,
			account_name = ?,
			account_email = ?,
			account_uuid = ?,
			created_at = now(),
			updated_at = now();
	`, m.Organization, m.Account.AccountId, m.Account.AccountPwd,
		m.Account.AccountName, m.Account.AccountEmail, m.Account.AccountUuid).Error
	if err != nil && strings.Contains(strings.ToLower(err.Error()), "duplicate entry") {
		return configs.ErrDuplicatedExist
	} else if err != nil {
		s.Log.Errorln(err)
		return err
	}
	return err
}

func (s *SqlQuery) Login(m *models.Account) error {
	result := s.DB.Exec(`
		update account
		set connected_at =
			case when (select if (exists(
				select account_id
					from account
				where account_id = @ac), 1, 0)) = 1 then now() else null end
		where account_id = (@ac := ?) and account_pwd = ? and status = 1;
	`, m.AccountId, m.AccountPwd)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return configs.ErrRequestTooFast
	}
	return nil
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

func (s *SqlQuery) SSO(m *models.Organization, n *models.Account) error {
	return s.DB.Exec(`
	insert into account
	set auth_seq = (
		select seq
			from organization
		where organization = ?
		), account_id = ?, account_pwd = ?,
		account_name = ?, account_email = ?,
		account_uuid = ?, status = 1,
		created_at = now(), updated_at = now();
	`, m.Organization, n.AccountId, n.AccountPwd,
		n.AccountName, n.AccountEmail, n.AccountUuid).Error
}

func (s *SqlQuery) GetAccountInfo(m *models.Account) error {
	return s.DB.Raw(`
		select account_id, account_name, account_email
		from account
		where status = 1 and account_uuid =?
		limit 1;
	`, m.AccountUuid).Scan(&m).Error
}

func (s *SqlQuery) GetOrgList(m *[]models.Organization) error {
	return s.DB.Raw(`
		select *
		from organization
		where status = 1;
	`).Scan(m).Error
}

func (s *SqlQuery) WithdrawalAccount(n *models.Withdrawal, m *models.Account) error {
	return s.DB.Exec(`
		update account
		set status = 0
		where account_uuid = ? and account_pwd = ?;
	`, m.AccountUuid, n.AccountPwd).Error
}
