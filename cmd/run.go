package cmd

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"strings"
	"syscall"
	"time"

	_ "embed"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/julis-bw-nw/julis-service-register-app/internal/app/registerkey"
	"github.com/julis-bw-nw/julis-service-register-app/internal/app/user"
	"github.com/julis-bw-nw/julis-service-register-app/internal/pkg/data"
	"github.com/julis-bw-nw/julis-service-register-app/pkg/ldap/lldap"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//go:embed config.default.yml
var defaultConfig []byte

type Config struct {
	API      API
	Database Database
	LLDAP    LLDAP
}

type API struct {
	Bind string
}

type Database struct {
	Host     string
	Port     uint16
	Database string
	Username string
	Password string
}

func (cfg Database) DSN() string {
	return fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s",
		cfg.Host, cfg.Port, cfg.Database, cfg.Username, cfg.Password)
}

type LLDAP struct {
	Host     string `yaml:"host"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

var (
	configPath string

	runCmd = &cobra.Command{
		Use:   "run",
		Short: "Starts the API and file server",
		Run: func(_ *cobra.Command, _ []string) {
			log.Println("Starting Julis Service Register App")

			viper.SetEnvPrefix("JSR")
			viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
			viper.AutomaticEnv()
			configPath = viper.GetString("CONFIG_PATH")
			cfg, err := loadConfig()
			if err != nil {
				log.Fatalf("Failed to read config from %q, %s", configPath, err)
			}

			db, err := data.NewPostgres(cfg.Database.DSN())
			if err != nil {
				log.Fatalf("Failed connect to DB; %s", err)
			}

			regKeyService := registerkey.Service{
				DataService: db,
			}

			cli := &http.Client{
				Timeout: time.Second * 3,
			}

			ldapService := lldap.New(cli, cfg.LLDAP.Host, lldap.WithAuthenticatorTransport(cfg.LLDAP.Username, cfg.LLDAP.Password))

			userService := user.Service{
				DataService: db,
				LDAPService: ldapService,
			}

			r := chi.NewRouter()
			r.Use(middleware.RequestID)
			r.Use(middleware.RealIP)
			r.Use(middleware.Logger)
			r.Use(middleware.Recoverer)
			r.Use(middleware.Timeout(60 * time.Second))
			r.Mount("/", http.StripPrefix("/", http.FileServer(http.Dir("web/src"))))
			r.Route("/api", func(r chi.Router) {
				r.Mount("/users", userService.Handler())
				r.Mount("/register-keys", regKeyService.Handler())
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

			log.Println("Service is online")

			sc := make(chan os.Signal, 1)
			signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
			<-sc
		},
	}
)

func init() {
	runCmd.Flags().StringVarP(&configPath, "config", "c", "config.yml", "config path to look for configs")
	viper.BindPFlag("CONFIG_PATH", runCmd.Flags().Lookup("config"))
}

func loadConfig() (Config, error) {
	viper.SetConfigType("yml")
	if err := viper.ReadConfig(bytes.NewReader(defaultConfig)); err != nil {
		return Config{}, err
	}

	configPath = strings.TrimSpace(configPath)
	dir, file := path.Split(configPath)
	ext := path.Ext(file)
	fileName := strings.TrimSuffix(file, ext)
	if dir == "" {
		dir = "."
	}

	viper.SetConfigName(fileName)
	viper.AddConfigPath(dir)
	viper.SetConfigType(strings.TrimPrefix(ext, "."))

	_ = viper.SafeWriteConfigAs(configPath)

	if err := viper.ReadInConfig(); err != nil {
		return Config{}, err
	}

	var cfg Config
	return cfg, viper.Unmarshal(&cfg)
}
