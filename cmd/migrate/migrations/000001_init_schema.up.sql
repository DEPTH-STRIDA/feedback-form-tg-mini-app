-- Создаем таблицу пользователей
CREATE TABLE users (
    id BIGINT PRIMARY KEY,                                    -- ID пользователя из Telegram
    first_name VARCHAR(64) NOT NULL,                         -- Имя из Telegram
    last_name VARCHAR(64),                                   -- Фамилия из Telegram
    username VARCHAR(32) UNIQUE,                             -- Username из Telegram
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Создаем таблицу форм обратной связи
CREATE TABLE forms (
    id BIGSERIAL PRIMARY KEY,                               -- Уникальный ID формы
    user_id BIGINT NOT NULL REFERENCES users(id),           -- ID пользователя, оставившего форму
    name VARCHAR(128) NOT NULL,                            -- Имя пользователя
    feedback VARCHAR(256),                                 -- Предпочтительный способ обратной связи
    comment VARCHAR(512),                                  -- Комментарий к заявке
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);


