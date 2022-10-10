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
)

type restController struct {
	ctx        *fiber.Ctx
	apiRequest *models.API
	log        *log.User
	account    *models.SsoUser
	rdbQuery   interfaces.RdbQuery
	redisQuery interfaces.RedisQuery
}

func newClient(c *fiber.Ctx) *restController {
	l := log.NewUser()
	id := c.Locals("id")
	account := &models.SsoUser{}
	if id == nil {
		l.Information = "no user info"
	} else {
		account.UserUuid = id.(uuid.UUID)
		(&queries.SqlQuery{Log: l}).GetAccountInfo(c.UserContext().Value(middleware.Tx("RdbConnection")), account)
		l.Information = account.Username
	}
	return &restController{
		ctx:        c,
		apiRequest: &models.API{},
		log:        l,
		account:    account,
		rdbQuery:   &queries.SqlQuery{Log: l},
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
	db := r.ctx.UserContext().Value(middleware.Tx("RdbConnection"))
	m := &models.Example{}
	itoj, err := json.Marshal(r.apiRequest.Body)
	if err != nil {
		return r.ctx.Status(fiber.StatusBadRequest).JSON(&models.R{
			Status:   fiber.StatusBadRequest,
			Response: "failed; " + err.Error(),
		})
	}
	json.Unmarshal(itoj, m)
	if err := r.rdbQuery.Update(db, m); err != nil {
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
	db := r.ctx.UserContext().Value(middleware.Tx("RdbConnection"))
	var m []models.Organization
	if err := r.rdbQuery.GetOrgList(db, &m); err != nil {
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
