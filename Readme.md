# Telegram Mini App Form Backend

Бэкенд для мини-приложения Telegram и бота для обработки форм обратной связи. Приложение позволяет пользователям оставлять заявки через Telegram Mini App или напрямую через бота.

## Функциональность

- Прием и обработка форм обратной связи через Mini App
- Автоматическая отправка уведомлений о новых заявках в указанные Telegram чаты
- Сохранение заявок в PostgreSQL
- Валидация данных и защита от спама

## API Endpoints

### POST /api/v1/form
Создание новой заявки

**Headers:**
- `Authorization`: initData от Telegram Mini App (обязательный)

**Request Body:**
```json
{
  "name": "Имя пользователя",     // обязательное, 1-128 символов
  "feedback": "Способ связи",      // опциональное, 0-256 символов
  "comment": "Комментарий"        // опциональное, 0-512 символов
}
```

## Структура проекта

```
.
├── cmd/                    # Точки входа
│   ├── api/               # API сервер
│   └── migrate/           # Утилита для миграций
├── internal/              # Внутренняя логика
│   ├── api/              # API слой
│   ├── service/          # Бизнес-логика
│   ├── repository/       # Работа с БД
│   ├── model/           # Модели данных
│   └── tg/              # Интеграция с Telegram
└── pkg/                  # Переиспользуемые пакеты
    ├── config/          # Утилиты конфигурации
    └── tg/              # Общие утилиты для Telegram
```

## Требования

- Go 1.21+
- PostgreSQL 14+
- Telegram Bot Token

## Конфигурация

Создайте `.env` файл:
```env
# API
API_ADDR=localhost
API_PORT=3000
LIMITER_RATE=5
LIMITER_BURST=10

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASS=postgres
DB_NAME=tgform
DB_SSLMODE=disable

# Telegram
TG_TOKEN=your_bot_token
TG_MESSAGE_CHATS=123456789,-987654321  # ID чатов для уведомлений
TG_EXPIRATION_HOURS=24
TG_CLEANUP_INTERVAL_MINUTES=60
```

## Запуск

1. Миграции:
```bash
go run cmd/migrate/main.go
```

2. API сервер:
```bash
go run cmd/api/main.go
```

## В разработке

- [ ] Валидация и DTO для API
- [ ] CI/CD конфигурация
- [ ] Кэширование
- [ ] Метрики и мониторинг
