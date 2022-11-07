package controllers

import (
	"fiber-template/app/interfaces"
	"fiber-template/app/models"
	"fiber-template/app/queries"
	"fiber-template/pkg/middleware"
	"fiber-template/pkg/utils"
	log "fiber-template/pkg/utils/logger"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type accountController struct {
	ctx        *fiber.Ctx
	log        *log.User
	account    *models.Account
	rdbQuery   interfaces.RdbQuery
	redisQuery interfaces.RedisQuery
}

func newAccountClient(c *fiber.Ctx) *accountController {
	l := log.NewUser()
	id := c.Locals("id")
	account := &models.Account{}
	s := &queries.SqlQuery{
		Log: l,
		DB:  c.UserContext().Value(middleware.Tx("RdbConnection")).(*gorm.DB),
	}
	if id == nil {
		l.Information = "no user info"
	} else {
		account.AccountUuid = id.(uuid.UUID)
		s.GetAccountInfo(account)
		l.Information = account.AccountName
	}
	return &accountController{
		ctx:        c,
		log:        l,
		account:    account,
		rdbQuery:   s,
		redisQuery: &queries.RedisQuery{Log: l},
	}
}

// POST /account/register
//   - params: account_id, account_pwd, account_name, account_email
//     Register process
//     1. Check account id and password
//     2. If account id and password are valid, register account
//     3. If account id and password are not valid, return error message
func Register(c *fiber.Ctx) error {
	n := newAccountClient(c)
	param := &models.Register{
		// AccountId    string    `json:"account_id"`
		// AccountPwd   string    `json:"account_pwd"`
		// AccountName  string    `json:"account_name"`
		// AccountEmail string    `json:"account_email"`
	}
	if err := n.ctx.BodyParser(param); err != nil {
		return n.ctx.Status(fiber.StatusBadRequest).JSON(&models.R{
			Status:   fiber.StatusBadRequest,
			Response: "failed; invalid body",
		})
	}
	if param.AccountId == "" || param.AccountPwd == "" {
		return n.ctx.Status(fiber.StatusBadRequest).JSON(&models.R{
			Status:   fiber.StatusBadRequest,
			Response: "failed; invalid body",
		})
	}
	s := strings.ContainsAny(param.AccountId, " !@#$%^&*()+{}|:\"<>?`~/,';[]\\=")
	if s {
		return n.ctx.Status(fiber.StatusBadRequest).JSON(&models.R{
			Status:   fiber.StatusBadRequest,
			Response: "failed; invalid account id",
		})
	}
	utils.HashPassword(&param.AccountPwd)
	m := &models.RegisterQuery{
		Account: models.Account{
			AccountId:    param.AccountId,
			AccountPwd:   param.AccountPwd,
			AccountName:  param.AccountName,
			AccountEmail: param.AccountEmail,
			AccountUuid:  utils.MakeUID(),
			Status:       true,
		},
		Organization: "next2us",
	}
	err := n.rdbQuery.UpdateAccount(m)
	if err != nil {
		return n.ctx.Status(fiber.StatusInternalServerError).JSON(&models.R{
			Status:   fiber.StatusInternalServerError,
			Response: "failed; register account " + err.Error(),
		})
	}
	return n.ctx.Status(fiber.StatusOK).JSON(&models.R{
		Status:   fiber.StatusOK,
		Response: "success",
	})
}

// POST /account/login
//   - params: account_id, account_pwd
//     Login process
//     1. Check account id and password
//     2. If account id and password are correct, return JWT token
//     3. If account id and password are incorrect, return error message
//     4. If account id and password are not exist, return error message
//     5. If account id and password are not valid, return error message
func Login(c *fiber.Ctx) error {
	n := newAccountClient(c)
	param := &models.Login{
		// 	AccountId  string `json:"account_id"`
		// 	AccountPwd string `json:"account_pwd"`
	}
	if err := n.ctx.BodyParser(param); err != nil {
		return n.ctx.Status(fiber.StatusBadRequest).JSON(&models.R{
			Status:   fiber.StatusBadRequest,
			Response: "failed; invalid body",
		})
	}
	if param.AccountId == "" || param.AccountPwd == "" {
		return n.ctx.Status(fiber.StatusBadRequest).JSON(&models.R{
			Status:   fiber.StatusBadRequest,
			Response: "failed; invalid body",
		})
	}
	s := strings.ContainsAny(param.AccountId, " !@#$%^&*()+{}|:\"<>?`~/,';[]\\=")
	if s {
		return n.ctx.Status(fiber.StatusBadRequest).JSON(&models.R{
			Status:   fiber.StatusBadRequest,
			Response: "failed; invalid account id",
		})
	}
	utils.HashPassword(&param.AccountPwd)
	m := &models.Account{
		AccountId:  param.AccountId,
		AccountPwd: param.AccountPwd,
	}
	err := n.rdbQuery.Login(m)
	if err != nil {
		return n.ctx.Status(fiber.StatusInternalServerError).JSON(&models.R{
			Status:   fiber.StatusInternalServerError,
			Response: "failed; login account " + err.Error(),
		})
	}

	// TODO: JWT token
	return n.ctx.Status(fiber.StatusOK).JSON(&models.R{
		Status:   fiber.StatusOK,
		Response: "success",
	})
}

// POST /account/modify
//   - params: account_id, account_pwd, account_name, account_email
//     Modify process
//     1. Check account validation by JWT
//     2. Check the param is valid
//     3. Modify account information
func Modify(c *fiber.Ctx) error {
	c.Locals("id", uuid.MustParse("57cd4e03-e40c-4af8-bb46-3d14e49c9313"))
	n := newAccountClient(c)
	param := &models.Modify{
		// AccountPwd   string `json:"account_pwd,omitempty"`
		// AccountName  string `json:"account_name,omitempty"`
		// AccountEmail string `json:"account_email,omitempty"`
	}

	if err := n.ctx.BodyParser(param); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&models.R{
			Status:   fiber.StatusBadRequest,
			Response: "failed; invalid body",
		})
	}

	if (*param == models.Modify{}) {
		return n.ctx.Status(fiber.StatusBadRequest).JSON(&models.R{
			Status:   fiber.StatusBadRequest,
			Response: "failed; invalid body",
		})
	}

	if err := n.rdbQuery.ModifyAccount(param, string(n.account.AccountUuid.String())); err != nil {
		return n.ctx.Status(fiber.StatusInternalServerError).JSON(&models.R{
			Status:   fiber.StatusInternalServerError,
			Response: "failed; modify account " + err.Error(),
		})
	}

	// TODO
	// Delete JWT token

	return n.ctx.Status(fiber.StatusOK).JSON(&models.R{
		Status:   fiber.StatusOK,
		Response: "success",
	})
}

// POST /account/withdrawal
//   - params: account_id, account_pwd
//     Withdrawal process
//     1. Check account validation by JWT
//     2. Check the param is valid
//     3. Withdrawal account information
func Withdrawal(c *fiber.Ctx) error {
	n := newAccountClient(c)
	param := &models.Withdrawal{
		// AccountPwd string `json:"account_pwd"`
	}

	if err := n.ctx.BodyParser(param); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&models.R{
			Status:   fiber.StatusBadRequest,
			Response: "failed; invalid body",
		})
	}

	if param.AccountPwd == "" {
		return n.ctx.Status(fiber.StatusBadRequest).JSON(&models.R{
			Status:   fiber.StatusBadRequest,
			Response: "failed; invalid body",
		})
	}

	utils.HashPassword(&param.AccountPwd)
	if err := n.rdbQuery.WithdrawalAccount(param, n.account); err != nil {
		return n.ctx.Status(fiber.StatusInternalServerError).JSON(&models.R{
			Status:   fiber.StatusInternalServerError,
			Response: "failed; delete account " + err.Error(),
		})
	}

	// TODO
	// Delete JWT token

	return n.ctx.Status(fiber.StatusOK).JSON(&models.R{
		Status:   fiber.StatusOK,
		Response: "success",
	})
}
