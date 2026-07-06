package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
	"traineesheep/notifyservice/email"
	"traineesheep/notifyservice/handlers"
	"traineesheep/notifyservice/tgbot"

	"github.com/go-telegram/bot"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

// Структура SmtpDTO нужна для передачи параметров почты, с которой будет
// вестись рассылка
type SmtpDTO struct {
	Email    string
	Password string
	Host     string
	Port     string
}

// Функция main это главная функция для работа программы
func main() {
	godotenv.Load("./config/data.env")

	var channel_size, _ = strconv.Atoi(os.Getenv("CHANNEL_SIZE"))

	pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_CONNECT"))
	if err != nil {
		panic(err)
	}
	defer pool.Close()

	opts := []bot.Option{
		bot.WithDefaultHandler(tgbot.Handler),
	}

	bot, err := bot.New(os.Getenv("BOT_TOKEN"), opts[0])
	if err != nil {
		panic(err)
	}
	smtpDto := email.NewSmtpDTO(os.Getenv("smtpEmail"), os.Getenv("smtpPassword"), os.Getenv("smtpHost"), os.Getenv("smtpPort"))

	grtChan := make(chan int, channel_size)
	wg := &sync.WaitGroup{}
	d := handlers.NewDTO(bot, pool, smtpDto, grtChan, wg, os.Getenv("JWT_SECRET"))

	r := mux.NewRouter()
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/notify", d.HandleNotify).Methods("OPTIONS", "POST")
	api.HandleFunc("/notify_types", d.HandleSaveSettingsCheckmarks).Methods("OPTIONS", "POST")
	api.HandleFunc("/get_notify_settings", d.HandleGetNotifySettings).Methods("OPTIONS", "GET")
	api.HandleFunc("/moderator_login", d.HandleModeratorLogin).Methods("OPTIONS", "POST")

	srv := &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: r,
	}

	go func() {
		log.Println("Telegram bot has been started")
	}()

	go func() {
		log.Println("Server has been started!")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT)
	<-stop

	log.Println("Shutting down server")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Println("HTTP server shutdown error:", err)
	}

	log.Println("Waiting for active requests to finish...")
	wg.Wait()

	if _, err := bot.Close(context.Background()); err != nil {
		log.Println("Tg bot close error:", err)
	}
}
