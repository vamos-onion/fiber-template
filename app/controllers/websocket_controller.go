package controllers

import (
	"context"
	"encoding/json"
	"fiber-template/app/interfaces"
	"fiber-template/app/models"
	"fiber-template/app/queries"
	"fiber-template/pkg/middleware"
	log "fiber-template/pkg/utils/logger"
	"fiber-template/platform/database"
	"runtime/debug"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"gorm.io/gorm"
)

// Websocket APIs List up
//   - func (w *wsController) example() {...}
//   - func (w *wsController) keepAlive() {...}
type wsController struct {
	conn       *websocket.Conn
	apiRequest *models.API
	log        *log.User
	rdbQuery   interfaces.RdbQuery
	redisQuery interfaces.RedisQuery
}

func newWsClient(c *websocket.Conn) *wsController {
	c.UnderlyingConn().SetDeadline(time.Now().Add(wsTimeout * time.Second))
	l := log.NewUser()
	return &wsController{
		conn:       c,
		apiRequest: &models.API{},
		log:        l,
		rdbQuery: &queries.SqlQuery{
			Log: l,
		},
		redisQuery: &queries.RedisQuery{Log: l},
	}
}

func (w *wsController) txSet() *context.Context {
	db := database.DB.MariaDB.Begin()
	w.rdbQuery = &queries.SqlQuery{
		Log: w.log,
		DB:  db,
	}
	newCtx := context.WithValue(context.Background(), middleware.Tx("Tx"), db)
	return &newCtx
}

func (w *wsController) closeConn() {
	w.conn.Close()
	debug.FreeOSMemory()
}

func WsConn(c *websocket.Conn) {
	w := newWsClient(c)
	defer w.closeConn()

	/***
	* @ websocket.conn.RemoteAddr().String() returns IP_ADDRESS:PORT
	*	- need PORT to specify the user
	*	- to reduce memory usage, convert the value to type uint16
	*	- saving as client's port and reuse this value when you need
	**/
	w.log.Infoln("Websocket connected!!")
	var msg []byte
	var err error
	for {
		_, msg, err = c.ReadMessage()
		if err != nil && !websocket.IsUnexpectedCloseError(err, websocket.CloseAbnormalClosure) {
			w.log.Infoln("websocket timeout")
			break
		} else if err != nil && !websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure) {
			w.log.Infoln("client closed websocket")
			break
		} else if err != nil && !websocket.IsUnexpectedCloseError(err, websocket.CloseNoStatusReceived) {
			w.log.Errorln("no status received", err)
			break
		} else if err != nil {
			w.log.Errorln("unexpected close error:", err)
			break
		}
		json.Unmarshal(msg, &w.apiRequest)
		if len(w.apiRequest.Transaction) == 0 {
			w.conn.WriteJSON(&models.R{
				Status:      fiber.StatusBadRequest,
				Transaction: w.apiRequest.Transaction,
				Response:    "transaction mandatory",
			})
			continue
		}
		switch w.apiRequest.Request {
		case "example":
			go w.example(c)
		case "keep-alive":
			go w.keepAlive()
		default:
			w.conn.WriteJSON(&models.R{
				Status:      fiber.StatusNotFound,
				Transaction: w.apiRequest.Transaction,
				Response:    "requested function not found",
			})
		}
	}
}

func (w *wsController) example(c *websocket.Conn) {
	ctx := w.txSet()
	db := (*ctx).Value(middleware.Tx("Tx")).(*gorm.DB)
	defer db.Rollback()
	t := &models.Example{
		Payload: "example",
	}
	err := w.rdbQuery.Update(t)
	if err != nil {
		c.WriteJSON(&models.R{
			Status:      fiber.StatusInternalServerError,
			Transaction: w.apiRequest.Transaction,
			Response:    err.Error(),
		})
		return
	}
	rv, err := w.rdbQuery.SelectAll(t)
	if err != nil {
		c.WriteJSON(&models.R{
			Status:      fiber.StatusInternalServerError,
			Transaction: w.apiRequest.Transaction,
			Response:    err.Error(),
		})
		return
	}

	// API StatusOk
	db.Commit()
	c.WriteJSON(&models.R{
		Status:      fiber.StatusOK,
		Transaction: w.apiRequest.Transaction,
		Response:    rv,
	})
}

func (w *wsController) keepAlive() {
	w.conn.UnderlyingConn().SetDeadline(time.Now().Add(wsTimeout * time.Second))
	w.apiRequest.Body = nil
	w.conn.WriteJSON(&models.R{
		Status:      fiber.StatusOK,
		Transaction: w.apiRequest.Transaction,
		Response:    "keep-alive",
	})
}
