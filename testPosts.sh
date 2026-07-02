#!/bin/bash

DB="postgres"
USER="postgres"
HOST="localhost"
PORT="5433"

echo "Очистка таблиц и создание 30 пользователей + 500 постов..."
echo "База: $DB, пользователь: $USER"

export PGPASSWORD="123"

psql -h "$HOST" -p "$PORT" -U "$USER" -d "$DB" <<'EOF'
-- Очищаем всё связанное каскадно
TRUNCATE post, image_post, refresh_tokens, users RESTART IDENTITY CASCADE;

DO $$
DECLARE
    i INTEGER;
    v_user_id INTEGER;
    v_post_id INTEGER;
    v_num_images INTEGER;
    j INTEGER;
BEGIN
    -- 1. Создаём 30 пользователей
    FOR i IN 1..30 LOOP
        INSERT INTO users (username, password, email, tg_chat_id)
        VALUES (
            'user' || i,
            'password',                     -- заглушка, пароль не хэширован
            'user' || i || '@example.com',
            NULL                             -- tg_chat_id не обязателен
        );
    END LOOP;

    -- 2. Создаём 500 постов, привязанных к случайным пользователям (id 1..30)
    FOR i IN 1..500 LOOP
        -- Случайный user_id от 1 до 30
        v_user_id := floor(random() * 30 + 1)::INT;

        INSERT INTO post (user_id, title, description)
        VALUES (v_user_id, 'Post ' || i, 'Description for post ' || i)
        RETURNING id INTO v_post_id;

        -- Случайное число изображений от 1 до 3
        v_num_images := floor(random() * 3 + 1)::INT;

        -- Вставка изображений со случайным image_id (1–1000)
        FOR j IN 1..v_num_images LOOP
            INSERT INTO image_post (post_id, image_id)
            VALUES (v_post_id, floor(random() * 1000 + 1)::INT);
        END LOOP;
    END LOOP;
END;
$$;
EOF

if [ $? -eq 0 ]; then
    echo "Готово! 30 пользователей и 500 постов с 1–3 изображениями вставлены."
else
    echo "Ошибка при выполнении запроса."
fi

read -p "Нажмите Enter, чтобы закрыть..."