package database

import (
	"fiber-template/pkg/utils"
	log "fiber-template/pkg/utils/logger"
	"io"
	defualtLog "log"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *DBConn

type DBConn struct {
	MariaDB *gorm.DB
}

func InitDBConnection(a *fiber.App) {
	DB = &DBConn{}
	DB.MariaDB = rdbConnection()
}

func rdbConnection() *gorm.DB {
	dbUrl, err := utils.ConnectionURLBuilder("mariadb")
	if err != nil {
		log.Fatalf("ERROR MAKING RDB URL", err)
	}
	var writer io.Writer
	if os.Getenv("STAGE_STATUS") == "dev" {
		writer = io.MultiWriter(os.Stdout, log.QueryLogFile)
	} else {
		writer = log.QueryLogFile
	}
	l, err := strconv.Atoi(os.Getenv("DB_LOG_LEVEL"))
	if err != nil {
		log.Fatalf("ERROR RDB LOG LEVEL", err)
	}
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       dbUrl, // data source name
		DefaultStringSize:         256,   // default size for string fields
		DisableDatetimePrecision:  true,  // disable datetime precision, which not supported before MySQL 5.6
		DontSupportRenameIndex:    true,  // drop & create when rename index, rename index not supported before MySQL 5.7, MariaDB
		DontSupportRenameColumn:   true,  // `change` when rename column, rename column not supported before MySQL 8, MariaDB
		SkipInitializeWithVersion: false, // auto configure based on currently MySQL version
	}), &gorm.Config{
		Logger: logger.New(
			defualtLog.New(writer, "\r\n", defualtLog.LstdFlags), // io writer
			logger.Config{
				SlowThreshold:             time.Second,        // Slow SQL threshold
				LogLevel:                  logger.LogLevel(l), // Log level @ gorm.Config.Logger.LogLevel
				IgnoreRecordNotFoundError: true,               // Ignore ErrRecordNotFound error for logger
				Colorful:                  true,               // Disable color
			},
		),
	})
	if err != nil {
		log.Fatalf("ERROR RDB ", err)
	}
	conn, err := db.DB()
	if err != nil {
		log.Fatalf("ERROR RDB", err)
	}
	i, err := strconv.Atoi(os.Getenv("DB_MAX_IDLE_CONNECTIONS"))
	if err != nil {
		log.Fatalf("ERROR RDB CONFIG", err)
	}
	j, err := strconv.Atoi(os.Getenv("DB_MAX_OPEN_CONNECTIONS"))
	if err != nil {
		log.Fatalf("ERROR RDB CONFIG", err)
	}
	k, err := strconv.Atoi(os.Getenv("DB_CONNECTION_MAX_LIFETIME"))
	if err != nil {
		log.Fatalf("ERROR RDB CONFIG", err)
	}
	conn.SetMaxIdleConns(i)
	conn.SetMaxOpenConns(j)
	conn.SetConnMaxLifetime(time.Duration(k) * time.Minute)

	log.Println("DB RDB Connection OK.")
	return db
}

/***
* @ gorm.Config.Logger.LogLevel
**	const (
		// Silent silent log level
		Silent LogLevel = iota + 1
		// Error error log level
		Error
		// Warn warn log level
		Warn
		// Info info log level
		Info
	)
*/
