package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/julis-bw-nw/julis-service-register-app/backend/db"
	"github.com/julis-bw-nw/julis-service-register-app/backend/user"
)

const (
	prefixEnv     = "JULIS_REGISTER_APP_"
	configPathEnv = prefixEnv + "CONFIG_PATH"
)

var (
	configPath = "config.yml"
)

func getEnv(key string, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func init() {
	configPath = getEnv(configPathEnv, configPath)

	if err := createConfigIfNotExist(configPath); err != nil {
		log.Fatalf("Failed to create default config at %q; %s", configPath, err)
	}
}

func main() {
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	cfg, err := loadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to read config from %q, %s", configPath, err)
	}

	encrypter := encrypter{secret: cfg.Database.EncryptionSecret}
	db := db.DB{
		Host:     cfg.Database.Host,
		Database: cfg.Database.Database,
		Username: cfg.Database.Username,
		Password: cfg.Database.Password,
	}

	userHandler := user.Handler{
		EncryptionService: &encrypter,
		DataService:       &db,
	}

	r.Route("/api", func(r chi.Router) {
		r.Post("/register", userHandler.PostRegisterUserHandler())
	})

	srv := http.Server{
		Addr: cfg.API.Bind,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
