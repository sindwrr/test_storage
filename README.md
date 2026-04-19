# Autotest Result Storage (WIP)
Autotest result and artifacts storage system

## Project tree
```
test_storage
├─ README.md
├─ cmd
│  └─ app
│     └─ main.go
├─ go.mod
├─ internal
│  ├─ api
│  │  ├─ handlers
│  │  │  ├─ hello.go
│  │  │  ├─ login.go
│  │  │  └─ logout.go
│  │  └─ router.go
│  └─ auth
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
