# Autotest Result Storage
GitHub: https://github.com/sindwrr/test_storage

Система централизованного хранения результатов и артефактов автоматизированных тестов.  
Разработана как легковесное решение для команд тестирования, работающих в условиях ограниченных вычислительных ресурсов.

## Содержание

- [1. Ключевые возможности](#1-ключевые-возможности)
- [2. Стек технологий](#2-стек-технологий)
- [3. Быстрый старт](#3-быстрый-старт)
- [4. Архитектура](#4-архитектура)
- [5. Тестирование](#5-тестирование)
- [6. API-эндпоинты](#6-api-эндпоинты)
- [7. Скриншоты веб-интерфейса](#7-скриншоты-веб-интерфейса)
- [8. Дерево проекта](#8-дерево-проекта)
- [9. Лицензия](#9-лицензия)

## 1. Ключевые возможности

- Загрузка и скачивание артефактов любых форматов (логи, изображения, видео, PDF)
- Хранение метаданных в PostgreSQL с быстрой фильтрацией по компонентам, сборкам, наборам тестов и датам
- Авторизация через LDAP и локальных пользователей
- Прямая отдача файлов через Nginx (sendfile) без нагрузки на приложение
- Потребление RAM < 2 ГБ при пиковых нагрузках

## 2. Стек технологий

| Компонент        | Технология           |
|------------------|----------------------|
| Бэкенд           | Go 1.25              |
| База данных      | PostgreSQL 16        |
| Прокси-сервер    | Nginx                |
| Контейнеризация  | Docker Compose       |
| Документация API | Swagger (OpenAPI)    |

## 3. Быстрый старт

1. **Клонировать репозиторий**
   ```bash
   git clone https://github.com/sindwrr/test_storage.git
   cd test_storage
   ```

2. **Создать `.env` файл в корне проекта**

   Настройте значения переменных под свое окружение:
   ```env
   # ===== База данных =====
   DB_HOST=db                  # Имя сервиса в docker-compose
   DB_PORT=5432                # Порт PostgreSQL (внутри контейнера)
   DB_USER=postgres            # Пользователь БД
   DB_PASSWORD=postgres        # Пароль БД
   DB_NAME=test_storage        # Имя базы данных
   
   # ===== Хранилище артефактов =====
   ARTIFACT_VOLUME=/app/artifacts  # Путь внутри контейнера приложения
   
   # ===== LDAP-авторизация =====
   LDAP_ADDR=ldap.example.com:389         # Адрес LDAP-сервера
   LDAP_BASE_DN=dc=example,dc=com         # Базовый DN для поиска
   LDAP_USER=cn=admin,dc=example,dc=com   # Сервисная учётная запись для bind (опционально)
   LDAP_PASSWORD=admin                    # Пароль сервисной учётной записи (опционально)
   
   # ===== Ограничения =====
   MAX_FILE_BYTES=31457280    # Максимальный размер файла артефакта в байтах (30 МБ)
   ```

   *⚠️ Система протестирована на артефактах размером не более 30 МБ. Использование с файлами артефактов большего размера возможно, но не рекомендуется.*

3. **Поднять окружение**
   ```bash
   docker compose -f deployments/docker-compose.yml up -d
   ```
4. **Открыть интерфейс в браузере**

   - Веб-интерфейс: http://localhost:8080
   - Swagger UI: http://localhost:8080/docs/ (доступен админам)

5. **Авторизоваться через LDAP на странице логина**

   При необходимости можно добавить тестовых пользователей в миграции `003_seed.up.sql`.

## 4. Архитектура

Три Docker-контейнера:
- app – Go-приложение
- db – PostgreSQL
- nginx – обратный прокси и раздача статики

Файлы артефактов хранятся в именованном томе artifacts.
При скачивании Nginx отдаёт их напрямую через X-Accel-Redirect, минуя бизнес-логику.

Приложение разделено на независимые сервисы, каждый со своей зоной ответственности:

### Бизнес-сервисы
- **Metadata Service** - управление метаданными артефактов и тестовых прогонов. Сохраняет и извлекает записи о компонентах, сборках, наборах тестов и результатах.
- **Storage Service** - работа с файловой системой. Принимает поток загружаемого файла, генерирует уникальное имя и сохраняет артефакт в томе `artifacts`.
- **Auth Service** - аутентификация через LDAP. Выполняет bind к LDAP-серверу, поиск пользователя по uid и повторный bind с паролем.
- **Analytics Service** - сбор статистики: количество артефактов по дням и распределение результатов тестов.
- **Preview Service** - предпросмотр текстовых и графических файлов прямо в браузере без скачивания.

### Системные сервисы
- **Config** - загрузка конфигурации из переменных окружения (параметры БД, LDAP, лимиты).
- **Health Service** - проверка готовности приложения и подключения к базе данных (эндпоинты `/health/alive` и `/health/ready`).
- **Background Worker** - фоновый пул задач. Периодически проверяет целостность: удаляет из БД записи артефактов, файлы которых были удалены из хранилища.

## 5. Тестирование

- Юнит-тесты: `go test ./...` (покрытие ~82%)
- Интеграционный тест: `go run tests/integration.go`
- Нагрузочное тестирование: скрипты в `stress/`

## 6. API-эндпоинты

- User - авторизованные пользователи
- Admin - администраторы (доступ ко всем эндпоинтам)
- ATF - фреймворк автотестирования

| Метод | Эндпоинт                        | Доступ      | Описание                                                                 |
|-------|---------------------------------|-------------|--------------------------------------------------------------------------|
| GET   | `/`                             | User        | Главная страница с таблицей артефактов                                   |
| GET   | `/login`                        | Все         | Страница авторизации                                                     |
| POST  | `/login`                        | Все         | Аутентификация (принимает `username` и `password`)                       |
| POST  | `/logout`                       | Все         | Выход из системы (очистка cookie)                                        |
| GET   | `/artifacts`                    | User        | JSON-массив с информацией об артефактах (поддержка фильтрации)           |
| GET   | `/artifact/download/{id}`       | User        | Скачивание файла артефакта по ID                                         |
| GET   | `/preview?id={id}`              | User        | Предпросмотр содержимого файла в браузере                                |
| GET   | `/analytics`                    | User        | Страница с графической аналитикой                                        |
| GET   | `/analytics/artifacts-per-day`  | User        | JSON: количество загруженных артефактов по дням                          |
| GET   | `/analytics/status-distribution`| User        | JSON: распределение результатов тестов                                   |
| GET   | `/health/alive`                 | Все         | Проверка работоспособности приложения                                    |
| GET   | `/health/ready`                 | Все         | Проверка готовности приложения и БД                                      |
| GET   | `/docs`                         | Admin       | Документация Swagger UI                                                  |
| POST  | `/upload`                       | Admin, ATF  | Загрузка файла артефакта с метаданными (multipart/form-data)             |

**Особенности доступа:**
- Все эндпоинты, кроме `/login`, `/logout` и `/health/*`, требуют авторизации (cookie `session`).
- `/docs` доступен только администраторам.
- `/upload` доступен администраторам и фреймворку автотестирования (ATF).

## 7. Скриншоты веб-интерфейса

- Страница авторизации:
<img width="947" height="538" alt="image" src="https://github.com/user-attachments/assets/a739509a-5e18-47f5-a812-ff40e3652a62" />


- Страница артефактов:
<img width="1221" height="645" alt="image" src="https://github.com/user-attachments/assets/677a4aa3-6f0d-4eab-be33-e4e0dc338e32" />


- Страница аналитики:
<img width="928" height="682" alt="image" src="https://github.com/user-attachments/assets/1f9044d9-134e-4fdc-8cd4-ba8a32a0dba2" />

## 8. Дерево проекта

```
test_storage
├─ README.md
├─ cmd                          # точка входа
│  └─ app
│     ├─ main.go
│     └─ main_test.go
├─ deployments                  # развертывание
│  ├─ .dockerignore
│  ├─ Dockerfile
│  ├─ docker-compose.yml
│  └─ nginx.conf
├─ docs                         # документация API
│  ├─ docs.go
│  ├─ swagger.json
│  └─ swagger.yaml
├─ go.mod
├─ go.sum
├─ internal                    
│  ├─ analytics                 # сервис аналитики
│  │  ├─ analytics_test.go
│  │  ├─ interface.go
│  │  ├─ repository
│  │  │  ├─ interface.go
│  │  │  ├─ repository.go
│  │  │  └─ repository_test.go
│  │  └─ service.go
│  ├─ api
│  │  ├─ handlers               # API-обработчики
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
│  │  ├─ middleware             # middleware для авторизации
│  │  │  ├─ admin.go
│  │  │  ├─ admin_test.go
│  │  │  ├─ auth.go
│  │  │  ├─ auth_test.go
│  │  │  ├─ upload.go
│  │  │  └─ upload_test.go
│  │  ├─ router.go              # маршрутизатор API-эндпоинтов
│  │  └─ router_test.go
│  ├─ auth                      # сервис авторизации
│  │  ├─ auth_test.go
│  │  ├─ interface.go
│  │  ├─ repository
│  │  │  ├─ interface.go
│  │  │  ├─ repository.go
│  │  │  └─ repository_test.go
│  │  └─ service.go
│  ├─ config                    # конфигурация приложения
│  │  ├─ config.go
│  │  └─ config_test.go
│  ├─ health                    # сервис проверки работоспособности
│  │  ├─ health_test.go
│  │  ├─ interface.go
│  │  └─ service.go
│  ├─ metadata                  # сервис метаданных
│  │  ├─ interface.go
│  │  ├─ metadata_test.go
│  │  ├─ repository
│  │  │  ├─ interface.go
│  │  │  ├─ repository.go
│  │  │  └─ repository_test.go
│  │  └─ service.go
│  ├─ models                    # структуры моделей и сущностей
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
│  ├─ preview                   # сервис предпросмотра артефактов
│  │  ├─ interface.go
│  │  ├─ preview_test.go
│  │  └─ service.go
│  ├─ storage                   # сервис хранения артефактов
│  │  ├─ interface.go
│  │  ├─ service.go
│  │  └─ storage_test.go
│  └─ worker                    # пул фоновых воркеров (проверка целостности файлов)
│     ├─ pool.go
│     ├─ tasks.go
│     └─ worker_test.go
├─ migrations                   # SQL-миграции
│  ├─ 001_init.down.sql
│  ├─ 001_init.up.sql
│  ├─ 002_indexes.down.sql
│  ├─ 002_indexes.up.sql
│  ├─ 003_seed.down.sql
│  └─ 003_seed.up.sql
├─ stress                       # материалы нагрузочного тестирования
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
├─ tests                         # интеграционное тестирование
│  └─ integration.go
└─ web                           # веб-интерфейс (фронтенд)
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

## 9. Лицензия

MIT License. Подробнее в файле [LICENSE](LICENSE).
