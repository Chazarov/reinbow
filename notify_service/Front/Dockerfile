# ================== ЭТАП 1: СБОРКА ==================
FROM node:18-alpine AS builder

WORKDIR /app

# Копируем файлы зависимостей
COPY package*.json ./

# Устанавливаем зависимости
RUN npm install

# Копируем исходники (включая webpack.config.js, static, etc.)
COPY . .

# Собираем проект (создаёт папку dist)
RUN npm run build

# ================== ЭТАП 2: ЗАПУСК ЧЕРЕЗ NGINX ==================
FROM nginx:alpine

# Копируем собранную статику из этапа сборки
COPY --from=builder /app/dist /usr/share/nginx/html

# Копируем конфиг Nginx с прокси для API
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Открываем порт (внутри контейнера Nginx слушает 80)
EXPOSE 80

# Запускаем Nginx
CMD ["nginx", "-g", "daemon off;"]