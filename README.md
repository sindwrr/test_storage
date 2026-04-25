# Autotest Result Storage (WIP)
Autotest result and artifacts storage system

## Current project tree
```
test_storage
├─ README.md
├─ cmd
│  └─ app
│     └─ main.go
├─ go.mod
├─ go.sum
├─ internal
│  ├─ api
│  │  ├─ handlers
│  │  │  ├─ health.go
│  │  │  ├─ hello.go
│  │  │  ├─ login.go
│  │  │  └─ logout.go
│  │  ├─ middleware
│  │  │  └─ auth.go
│  │  └─ router.go
│  ├─ auth
│  │  └─ service.go
│  ├─ config
│  │  └─ config.go
│  └─ health
│     ├─ interface.go
│     └─ service.go
├─ migrations
│  ├─ 001_init.down.sql
│  ├─ 001_init.up.sql
│  ├─ 002_indexes.down.sql
│  └─ 002_indexes.up.sql
└─ web
   └─ templates
      ├─ index.html
      └─ login.html
```
