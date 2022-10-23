[![License](https://img.shields.io/badge/license-Unlicense-blue.svg)](https://github.com/vamos-onion/fiber-template/blob/master/LICENSE)
# ⚡️Fiber-Template
Simple Golang Fiber framework template

<br>

## 🏗️Dependencies
- Go >= 1.17
- RDB
  - MariaDB or MySQL
- Redis
- If you don't wanna use query related things, comment out the lines in main.go that initializing DB connection

<br>

## 🗂️Used packages
- `go-fiber` : main framework
  - https://github.com/gofiber/fiber
  - https://gofiber.io/
- `go-fiber websocket` : websocket
  - https://github.com/gofiber/websocket
- `logrus` : logger pkg
  - https://github.com/sirupsen/logrus
- `gorm` : go ORM
  - https://github.com/go-gorm/gorm
- `go-redis` : go Redis
  - https://github.com/go-redis/redis
- `godotenv` : to read and set .env variables
  - https://github.com/joho/godotenv
- `google-uuid` : to generate uuid
  - https://github.com/google/uuid
- `golang-jwt` : JWT
  - https://github.com/golang-jwt/jwt
- `robfig-cron` : cron (not used yet)
  - https://github.com/robfig/cron
  
<br>

## 🚀How to run
- DB / Redis settings (If you wanna use)
- generate `.env` file using `env-sample`
- `git clone https://github.com/vamos-onion/fiber-template.git`
- `go mod tidy` & `go run main.go` (or you can build the project and then run the binary executable file)

<br>

## 🗄️DB settings
### Create MariaDB Table
```sql
create table example
(
    seq     int(50) auto_increment
        primary key,
    payload varchar(50) default '' not null
);

create table organization
(
    seq          int(50) auto_increment
        primary key,
    organization varchar(50) default '' not null,
    status       tinyint     default 0  not null,
    constraint organization_organization_uindex
        unique (organization)
);

create table sso_user
(
    seq          int(50) auto_increment
        primary key,
    organization varchar(50) default '' null,
    username     varchar(30) default '' not null,
    user_uuid    varchar(50)            null,
    user_index   int(50)                null,
    status       tinyint     default 0  null,
    constraint sso_user_user_index_uindex
        unique (user_index),
    constraint sso_user_user_uuid_uindex
        unique (user_uuid),
    constraint sso_user_organization_organization_fk
        foreign key (organization) references organization (organization)
);
```

<br>

## 🆘Known issues
- `logrus.Logger.SetReportCaller(true)`
  - only printing the point in `logger.go` file
  - if you want to print out original logrus report caller, should use logrus object directly

## ⚠️License
The packages that I used include MIT, BSD-2, and BSD-3 Licenses. \
Please be aware of use it.

#### 😊 Thx 👍
