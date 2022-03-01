package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/julis-bw-nw/julis-service-register-app/internal/app/register/regkey"
	"github.com/julis-bw-nw/julis-service-register-app/internal/app/register/user"
	"github.com/julis-bw-nw/julis-service-register-app/internal/pkg/data"
	"github.com/julis-bw-nw/julis-service-register-app/internal/pkg/ldap"
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

	if err := createConfigIfNotExist(configPath); err != nil && !errors.Is(err, os.ErrExist) {
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

	db := data.Service{
		Host:     cfg.Database.Host,
		Database: cfg.Database.Database,
		Username: cfg.Database.Username,
		Password: cfg.Database.Password,
		Timeout:  time.Second * 3,
	}

	if err := db.Open(); err != nil {
		log.Fatalf("Failed to connect to database; %s", err)
	}
	if err := db.Migrate(); err != nil {
		log.Fatalf("Failed to connect to database; %s", err)
	}

	regKeyService := regkey.Service{
		DataService: &db,
	}

	userService := user.Service{
		DataService: &db,
		LDAPService: ldap.Service{
			BaseURL: "http://lldap:17170",
			Client:  http.DefaultClient,
		},
	}

	fileServer(r, "/", "web/static")
	r.Route("/api", func(r chi.Router) {
		r.Mount("/users", userService)
		r.Mount("/regkeys", regKeyService)
	})

	srv := http.Server{
		Addr:    cfg.API.Bind,
		Handler: r,
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

func fileServer(r chi.Router, public string, static string) {
	if strings.ContainsAny(public, "{}*") {
		log.Println("FileServer does not permit any URL parameters.")
		return
	}

	root, _ := filepath.Abs(static)
	if _, err := os.Stat(root); os.IsNotExist(err) {
		log.Printf("Failed to find dir %q", root)
		return
	}

	fs := http.StripPrefix(public, http.FileServer(http.Dir(root)))

	if public != "/" && public[len(public)-1] != '/' {
		r.Get(public, http.RedirectHandler(public+"/", http.StatusMovedPermanently).ServeHTTP)
		public += "/"
	}

	r.Get(public+"*", func(w http.ResponseWriter, r *http.Request) {
		file := strings.Replace(r.RequestURI, public, "/", 1)

		if fileInfo, err := os.Stat(root + file); os.IsNotExist(err) || fileInfo.IsDir() {
			http.ServeFile(w, r, path.Join(root, "index.html"))
			return
		}
		fs.ServeHTTP(w, r)
	})
}
