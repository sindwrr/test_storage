# Autotest Result Storage (WIP)
Autotest result and artifacts storage system

## Current project tree
```
test_storage
├─ README.md
├─ cmd
│  └─ app
│     └─ main.go
├─ deployments
│  ├─ .dockerignore
│  ├─ Dockerfile
│  └─ docker-compose.yml
├─ docs
│  ├─ docs.go
│  ├─ swagger.json
│  └─ swagger.yaml
├─ go.mod
├─ go.sum
├─ internal
│  ├─ api
│  │  ├─ handlers
│  │  │  ├─ health.go
│  │  │  ├─ hello.go
│  │  │  ├─ login.go
│  │  │  ├─ logout.go
│  │  │  └─ upload.go
│  │  ├─ middleware
│  │  │  ├─ auth.go
│  │  │  └─ cors.go
│  │  └─ router.go
│  ├─ auth
│  │  └─ service.go
│  ├─ config
│  │  └─ config.go
│  ├─ health
│  │  ├─ interface.go
│  │  └─ service.go
│  ├─ metadata
│  │  ├─ interface.go
│  │  ├─ repository
│  │  │  ├─ interface.go
│  │  │  └─ repository.go
│  │  └─ service.go
│  ├─ models
│  │  ├─ build.go
│  │  ├─ component.go
│  │  ├─ file_type.go
│  │  ├─ result_status.go
│  │  ├─ run_status.go
│  │  ├─ test_artifact.go
│  │  ├─ test_run.go
│  │  ├─ test_suite.go
│  │  ├─ user.go
│  │  └─ user_group.go
│  └─ storage
│     ├─ interface.go
│     └─ service.go
├─ migrations
│  ├─ 001_init.down.sql
│  ├─ 001_init.up.sql
│  ├─ 002_indexes.down.sql
│  ├─ 002_indexes.up.sql
│  ├─ 003_seed.down.sql
│  └─ 003_seed.up.sql
└─ web
   ├─ static
   │  └─ login.css
   └─ templates
      ├─ index.html
      └─ login.html
```
