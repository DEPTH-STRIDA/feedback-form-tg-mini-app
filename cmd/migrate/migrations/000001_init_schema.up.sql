-- Сначала создаем таблицу groups без owner_id
CREATE TABLE groups (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(64) NOT NULL UNIQUE,
    title VARCHAR(256) NOT NULL,
    participants BIGINT[] DEFAULT '{}',
    alternating_weeks BOOLEAN NOT NULL DEFAULT false,
    odd_monday DATE,
    odd_week JSONB NOT NULL DEFAULT '{}'::jsonb,
    even_week JSONB NOT NULL DEFAULT '{}'::jsonb,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Затем создаем таблицу users
CREATE TABLE users (
    id BIGINT PRIMARY KEY,
    first_name VARCHAR(64),
    last_name VARCHAR(64),
    username VARCHAR(32),
    member_of BIGINT REFERENCES groups(id) ON DELETE SET NULL,
    owned_groups BIGINT[] DEFAULT '{}',
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- И только потом добавляем owner_id в groups
ALTER TABLE groups 
ADD COLUMN owner_id BIGINT REFERENCES users(id) ON DELETE CASCADE; 

