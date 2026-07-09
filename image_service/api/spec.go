// Package docs содержит Swagger/OpenAPI документацию internal API сервиса изображений.
//
// @title           Image Service Internal API
// @version         1.0
// @description     Internal HTTP API сервиса изображений.
// @host            localhost:8080
// @BasePath        /api
//
// @securityDefinitions.apikey AccessToken
// @in              cookie
// @name            access_token
package api

//go:generate go run github.com/swaggo/swag/cmd/swag@v1.16.4 init -g spec.go -o . -d .,../infrastructure/ports/api/v1 --parseInternal --instanceName internal
