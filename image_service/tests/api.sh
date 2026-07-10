#!/usr/bin/env bash
set -euo pipefail

BASE_URL="${BASE_URL:-http://localhost:8080}"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

COOKIES_FILE="cookies.txt"
TEST_LOGIN="${TEST_LOGIN:-testuser}"
TEST_PASS="${TEST_PASS:-123456}"
IMAGE_ID=""

rm -f "$COOKIES_FILE"

echo "Base URL: $BASE_URL"
echo "--------------------------------"
echo ""

test_0_health() {
  echo "Тест №0 : проверка health"
  curl -X GET "$BASE_URL/api/health" \
    -w "\nстатус - %{http_code}\n" \
    -s
  echo ""
}

test_1_register() {
  echo "Тест №1 : регистрация пользователя"
  curl -X POST "$BASE_URL/api/auth/register" \
    -H "Content-Type: application/json" \
    -d "{\"login\":\"$TEST_LOGIN\",\"pass\":\"$TEST_PASS\"}" \
    -w "\nстатус - %{http_code}\nbody - " \
    -s
  echo ""
}

test_2_login() {
  echo "Тест №2 : вход пользователя"
  curl -X POST "$BASE_URL/api/auth/login" \
    -H "Content-Type: application/json" \
    -d "{\"login\":\"$TEST_LOGIN\",\"pass\":\"$TEST_PASS\"}" \
    -c "$COOKIES_FILE" \
    -w "\nстатус - %{http_code}\nbody - " \
    -s
  echo ""
}

test_3_refresh() {
  echo "Тест №3 : обновление JWT-токена"
  curl -X POST "$BASE_URL/api/auth/refresh" \
    -b "$COOKIES_FILE" \
    -c "$COOKIES_FILE" \
    -w "\nстатус - %{http_code}\nbody - " \
    -s
  echo ""
}

test_4_upload_image() {
  echo "Тест №4 : загрузка изображения"

  local image_base64
  image_base64=$(base64 -w 0 ./test.png)

  local response http_body http_code
  response=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/api/upload" \
    -H "Content-Type: application/json" \
    -d "{
      \"name\": \"test_image\",
      \"media_type\": \"jpeg\",
      \"data\": \"$image_base64\"
    }")

  http_body=$(echo "$response" | sed '$d')
  http_code=$(echo "$response" | tail -n 1)

  echo "$http_body"
  echo "статус - $http_code"

  IMAGE_ID=$(echo "$http_body" | jq -r '.data.id')

  if [ -z "$IMAGE_ID" ] || [ "$IMAGE_ID" = "null" ]; then
    echo "Ошибка: не удалось получить id изображения"
    exit 1
  fi

  echo "IMAGE_ID=$IMAGE_ID"
  echo ""
}

test_5_get_image_guest() {
  echo "Тест №5 : получение изображения (гостевой доступ)"
  curl -X GET "$BASE_URL/api/guest/image/$IMAGE_ID" \
    -w "\nстатус - %{http_code}\n" \
    -s -o /dev/null
  echo ""
}

test_6_get_image_admin() {
  echo "Тест №6 : получение изображения (администратор)"
  curl -X GET "$BASE_URL/api/admin/image/$IMAGE_ID" \
    -b "$COOKIES_FILE" \
    -w "\nстатус - %{http_code}\n" \
    -s -o /dev/null
  echo ""
}

test_7_get_meta() {
  echo "Тест №7 : метаданные изображения"
  curl -X GET "$BASE_URL/api/image/meta/$IMAGE_ID" \
    -w "\nстатус - %{http_code}\nbody - " \
    -s
  echo ""
}

test_8_unmoderated_list() {
  echo "Тест №8 : список немодерированных"
  curl -X GET "$BASE_URL/api/images/unmoderated" \
    -b "$COOKIES_FILE" \
    -w "\nстатус - %{http_code}\nbody - " \
    -s
  echo ""
}

test_9_unmoderated_pagination() {
  echo "Тест №9 : список немодерированных с пагинацией (page=0, page_size=2)"
  curl -X GET "$BASE_URL/api/images/unmoderated?page=0&page_size=2" \
    -b "$COOKIES_FILE" \
    -w "\nстатус - %{http_code}\nbody - " \
    -s
  echo ""
}

test_10_block_image() {
  echo "Тест №10 : блокировка изображения"
  curl -X PUT "$BASE_URL/api/images/$IMAGE_ID/block" \
    -b "$COOKIES_FILE" \
    -w "\nстатус - %{http_code}\nbody - " \
    -s
  echo ""
}

test_11_approve_image() {
  echo "Тест №11 : одобрение изображения"
  curl -X PUT "$BASE_URL/api/images/$IMAGE_ID/approve" \
    -b "$COOKIES_FILE" \
    -w "\nстатус - %{http_code}\nbody - " \
    -s
  echo ""
}

test_12_logout() {
  echo "Тест №12 : выход пользователя"
  curl -X POST "$BASE_URL/api/auth/logout" \
    -b "$COOKIES_FILE" \
    -c "$COOKIES_FILE" \
    -w "\nстатус - %{http_code}\nbody - " \
    -s
  echo ""
}

run_upload_flow() {
  test_4_upload_image
}

run_all_tests() {
  test_0_health
  test_1_register
  test_2_login
  test_3_refresh
  test_4_upload_image
  test_5_get_image_guest
  test_6_get_image_admin
  test_7_get_meta
  test_8_unmoderated_list
  test_9_unmoderated_pagination
  test_10_block_image
  test_11_approve_image
  test_12_logout
}

case "${1:-all}" in
  upload)
    run_upload_flow
    ;;
  all)
    run_all_tests
    ;;
  *)
    echo "Использование: $0 {upload|all}"
    exit 1
    ;;
esac
