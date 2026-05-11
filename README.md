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
│  ├─ analytics
│  │  ├─ analytics_test.go
│  │  ├─ interface.go
│  │  ├─ repository
│  │  │  ├─ interface.go
│  │  │  ├─ repository.go
│  │  │  └─ repository_test.go
│  │  └─ service.go
│  ├─ api
│  │  ├─ handlers
│  │  │  ├─ analytics.go
│  │  │  ├─ artifacts.go
│  │  │  ├─ download.go
│  │  │  ├─ health.go
│  │  │  ├─ index.go
│  │  │  ├─ login.go
│  │  │  ├─ logout.go
│  │  │  └─ upload.go
│  │  ├─ middleware
│  │  │  ├─ auth.go
│  │  │  └─ cors.go
│  │  └─ router.go
│  ├─ auth
│  │  ├─ auth_test.go
│  │  ├─ interface.go
│  │  └─ service.go
│  ├─ config
│  │  └─ config.go
│  ├─ health
│  │  ├─ health_test.go
│  │  ├─ interface.go
│  │  └─ service.go
│  ├─ metadata
│  │  ├─ interface.go
│  │  ├─ metadata_test.go
│  │  ├─ repository
│  │  │  ├─ interface.go
│  │  │  ├─ repository.go
│  │  │  └─ repository_test.go
│  │  └─ service.go
│  ├─ models
│  │  ├─ analytics
│  │  │  ├─ day_count.go
│  │  │  └─ status_count.go
│  │  ├─ artifact_info.go
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
│  ├─ storage
│  │  ├─ interface.go
│  │  ├─ service.go
│  │  └─ storage_test.go
│  └─ worker
│     ├─ pool.go
│     ├─ tasks.go
│     └─ worker_test.go
├─ migrations
│  ├─ 001_init.down.sql
│  ├─ 001_init.up.sql
│  ├─ 002_indexes.down.sql
│  ├─ 002_indexes.up.sql
│  ├─ 003_seed.down.sql
│  └─ 003_seed.up.sql
├─ stress
│  ├─ download
│  │  ├─ latency_graph.py
│  │  ├─ latency_graph_lin.png
│  │  ├─ latency_graph_log.png
│  │  ├─ loader.go
│  │  ├─ ram_graph.png
│  │  └─ ram_graph.py
│  ├─ filter
│  │  ├─ targets.txt
│  │  └─ vegeta_log.txt
│  └─ upload
│     ├─ latency_graph.py
│     ├─ latency_graph_lin.png
│     ├─ latency_graph_log.png
│     ├─ loader.go
│     ├─ ram_graph.png
│     └─ ram_graph.py
├─ tests
└─ web
   ├─ static
   │  ├─ analytics.css
   │  ├─ analytics.js
   │  ├─ index.css
   │  └─ login.css
   └─ templates
      ├─ analytics.html
      ├─ index.html
      └─ login.html
```
