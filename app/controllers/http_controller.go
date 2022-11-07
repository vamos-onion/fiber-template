package controllers

import (
	"encoding/json"
	"fiber-template/app/interfaces"
	"fiber-template/app/models"
	"fiber-template/app/queries"
	"fiber-template/pkg/middleware"
	log "fiber-template/pkg/utils/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type restController struct {
	ctx        *fiber.Ctx
	apiRequest *models.API
	log        *log.User
	account    *models.Account
	rdbQuery   interfaces.RdbQuery
	redisQuery interfaces.RedisQuery
}

func newClient(c *fiber.Ctx) *restController {
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
	return &restController{
		ctx:        c,
		log:        l,
		account:    account,
		rdbQuery:   s,
		redisQuery: &queries.RedisQuery{Log: l},
	}
}

func Get(c *fiber.Ctx) error {
	/* JWT Token validation check */
	if status, err := tokenValidCheck(c); err != nil {
		return c.Status(status).JSON(&models.R{
			Status:   uint16(status),
			Response: err.Error(),
		})
	}
	j := newClient(c)
	switch j.ctx.Params("param") {
	case "org-list":
		return j.getOrgList()
	default:
		return c.Status(fiber.StatusInternalServerError).JSON(&models.R{
			Status:   fiber.StatusNotFound,
			Response: "failed; no such request",
		})
	}
}

func Post(c *fiber.Ctx) error {
	/* JWT Token validation check */
	if status, err := tokenValidCheck(c); err != nil {
		return c.Status(status).JSON(&models.R{
			Status:   uint16(status),
			Response: err.Error(),
		})
	}
	n := newClient(c)
	if err := c.BodyParser(n.apiRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&models.R{
			Status:   fiber.StatusBadRequest,
			Response: "failed; invalid body",
		})
	}
	switch n.apiRequest.Request {
	case "update":
		return n.update()
	default:
		return c.Status(fiber.StatusInternalServerError).JSON(&models.R{
			Status:   fiber.StatusNotFound,
			Response: "failed; no such request",
		})
	}
}

func InternalRequest(c *fiber.Ctx) error {
	n := newClient(c)
	/* Comment out while DEV testing */
	if xRealIp := c.Get("X-Real-Ip"); len(xRealIp) != 0 {
		n.log.Errorln(xRealIp)
		return c.Status(fiber.StatusUnauthorized).JSON(&models.R{
			Status:   fiber.StatusUnauthorized,
			Response: "failed; unauthorized request",
		})
	}
	if err := c.BodyParser(n.apiRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&models.R{
			Status:   fiber.StatusBadRequest,
			Response: "failed; invalid body",
		})
	}
	n.log.Information = "Request from internal server"
	return nil
}

func (r *restController) update() error {
	m := &models.Example{}
	itoj, err := json.Marshal(r.apiRequest.Body)
	if err != nil {
		return r.ctx.Status(fiber.StatusBadRequest).JSON(&models.R{
			Status:   fiber.StatusBadRequest,
			Response: "failed; " + err.Error(),
		})
	}
	json.Unmarshal(itoj, m)
	if err := r.rdbQuery.Update(m); err != nil {
		return r.ctx.Status(fiber.StatusBadRequest).JSON(&models.R{
			Status:   fiber.StatusBadRequest,
			Response: "failed; " + err.Error(),
		})
	}
	return r.ctx.Status(fiber.StatusOK).JSON(&models.R{
		Status:   fiber.StatusOK,
		Response: r.apiRequest,
	})
}

func (r *restController) getOrgList() error {
	var m []models.Organization
	if err := r.rdbQuery.GetOrgList(&m); err != nil {
		return r.ctx.Status(fiber.StatusBadRequest).JSON(&models.R{
			Status:   fiber.StatusBadRequest,
			Response: "failed; " + err.Error(),
		})
	}
	return r.ctx.Status(fiber.StatusOK).JSON(&models.R{
		Status:   fiber.StatusOK,
		Response: m,
	})
}
