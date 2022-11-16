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
    organization varchar(50)       not null,
    org_uuid     varchar(50)       null,
    status       tinyint default 0 not null,
    constraint organization_org_uuid_uindex
        unique (org_uuid),
    constraint organization_organization_uindex
        unique (organization)
);

create table account
(
    seq           int unsigned auto_increment
        primary key,
    auth_seq      int(50)           not null,
    account_id    varchar(20)       not null,
    account_pwd   varchar(100)      not null,
    account_name  varchar(20)       null,
    account_email varchar(50)       not null,
    account_uuid  varchar(40)       not null,
    status        tinyint default 1 not null,
    created_at    datetime          null,
    updated_at    datetime          null,
    connected_at  datetime          null,
    constraint users_account_email_uindex
        unique (account_email),
    constraint users_account_id_uindex
        unique (account_id),
    constraint users_account_uuid_uindex
        unique (account_uuid),
    constraint account_organization_seq_fk
        foreign key (auth_seq) references organization (seq)
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
