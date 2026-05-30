# Autotest Result Storage (WIP)
Autotest result and artifacts storage system

## Current project tree
```
test_storage
├─ README.md
├─ cmd
│  └─ app
│     ├─ main.go
│     └─ main_test.go
├─ deployments
│  ├─ .dockerignore
│  ├─ Dockerfile
│  ├─ docker-compose.yml
│  └─ nginx.conf
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
│  │  │  ├─ analytics_test.go
│  │  │  ├─ artifacts.go
│  │  │  ├─ artifacts_test.go
│  │  │  ├─ download.go
│  │  │  ├─ download_test.go
│  │  │  ├─ health.go
│  │  │  ├─ health_test.go
│  │  │  ├─ index.go
│  │  │  ├─ index_test.go
│  │  │  ├─ login.go
│  │  │  ├─ login_test.go
│  │  │  ├─ logout.go
│  │  │  ├─ logout_test.go
│  │  │  ├─ mocks.go
│  │  │  ├─ preview.go
│  │  │  ├─ preview_test.go
│  │  │  ├─ upload.go
│  │  │  └─ upload_test.go
│  │  ├─ middleware
│  │  │  ├─ admin.go
│  │  │  ├─ admin_test.go
│  │  │  ├─ auth.go
│  │  │  ├─ auth_test.go
│  │  │  ├─ upload.go
│  │  │  └─ upload_test.go
│  │  ├─ router.go
│  │  └─ router_test.go
│  ├─ auth
│  │  ├─ auth_test.go
│  │  ├─ interface.go
│  │  ├─ repository
│  │  │  ├─ interface.go
│  │  │  ├─ repository.go
│  │  │  └─ repository_test.go
│  │  └─ service.go
│  ├─ config
│  │  ├─ config.go
│  │  └─ config_test.go
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
│  │  │  ├─ analytics_test.go
│  │  │  ├─ day_count.go
│  │  │  └─ status_count.go
│  │  ├─ artifact_info.go
│  │  ├─ build.go
│  │  ├─ component.go
│  │  ├─ file_type.go
│  │  ├─ models_test.go
│  │  ├─ result_status.go
│  │  ├─ run_status.go
│  │  ├─ test_artifact.go
│  │  ├─ test_run.go
│  │  ├─ test_suite.go
│  │  ├─ user.go
│  │  └─ user_group.go
│  ├─ preview
│  │  ├─ interface.go
│  │  ├─ preview_test.go
│  │  └─ service.go
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
│  └─ integration.go
└─ web
   ├─ static
   │  ├─ analytics.css
   │  ├─ analytics.js
   │  ├─ index.css
   │  ├─ login.css
   │  └─ pagination.js
   └─ templates
      ├─ analytics.html
      ├─ index.html
      └─ login.html
```
