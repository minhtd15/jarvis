package api

import (
	"context"
	"crypto/tls"
	batman "education-website"
	"education-website/client"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
	"go.elastic.co/apm/module/apmlogrus"
	"gopkg.in/yaml.v3"
	_ "gopkg.in/yaml.v3"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	XApiKeyHeader string = "X-Api-Key"
)

var (
	userService  batman.UserService
	authService  batman.AuthService
	jwtService   batman.JwtService
	classService batman.ClassService
	flashClient  client.FlashClient
	redisClient  client.RedisClient
)

func GetLoggerWithContext(ctx context.Context) *log.Entry {
	log := log.WithFields(apmlogrus.TraceContext(ctx))
	return log
}

// Config struct for webapp config
type Config struct {
	Server struct {
		// Host is the local machine IP Address to bind the HTTP Server to
		Host string `yaml:"host"`
		// Port is the local machine TCP Port to bind the HTTP Server to
		Port    string `yaml:"port"`
		Timeout struct {
			// Server is the general server timeout to use
			// for graceful shutdowns
			Server time.Duration `yaml:"server"`

			// Write is the amount of time to wait until an HTTP server
			// write opperation is cancelled
			Write time.Duration `yaml:"write"`

			// Read is the amount of time to wait until an HTTP server
			// read operation is cancelled
			Read time.Duration `yaml:"read"`

			// Read is the amount of time to wait
			// until an IDLE HTTP session is closed
			Idle time.Duration `yaml:"idle"`
		} `yaml:"timeout"`
	} `yaml:"server"`

	Database struct {
		Port     string `yaml:"port"`
		Host     string `yaml:"host"`
		User     string `yaml:"user"`
		Password string `yaml:"pass"`
		DbName   string `yaml:"dbName"`
	} `yaml:"database"`

	XApiKey        string `yaml:"XApiKey"`
	FlashClientTmp struct {
		Root                       string `yaml:"root"`
		GetCourseRevenueByCourseId string `yaml:"get-course-revenue-by-course-id"`
		GetYearlyRevenue           string `yaml:"get-yearly-revenue"`
	} `yaml:"flashClient"`

	RedisClientTmp struct {
		Addr     string `yaml:"addr"`
		Password string `yaml:"password"`
		Db       int    `yaml:"db"`
	} `yaml:"redisClient"`

	RedisClient  client.RedisClient
	FlashClient  client.FlashClient
	UserService  batman.UserService
	JwtService   batman.JwtService
	AuthService  batman.AuthService
	ClassService batman.ClassService
}

// NewConfig returns a new decoded Config struct
func NewConfig(configPath string) (*Config, error) {
	// Create config structure
	config := &Config{}

	// Open config file
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Init new YAML decode
	d := yaml.NewDecoder(file)

	// Start YAML decoding from file
	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}

// ValidateConfigPath just makes sure, that the path provided is a file,
// that can be read
func ValidateConfigPath(path string) error {
	s, err := os.Stat(path)
	if err != nil {
		return err
	}
	if s.IsDir() {
		return fmt.Errorf("'%s' is a directory, not a normal file", path)
	}
	return nil
}

// ParseFlags will create and parse the CLI flags
// and return the path to be used elsewhere
func ParseFlags() (string, error) {
	// String that contains the configured configuration path
	configFile := os.Getenv("CONFIG_FILE_PATH")
	// Use configFile in your application

	// Set up a CLI flag called "-config" to allow users
	// to supply the configuration file
	flag.StringVar(&configFile, "config", "conf/config.yml", "path to config file")

	// Actually parse the flags
	flag.Parse()

	// Validate the path first
	if err := ValidateConfigPath(configFile); err != nil {
		log.Infof("error validating config path: %v", err)
		return "", err
	}

	// Return the configuration path
	return configFile, nil
}

// NewRouter generates the router used in the HTTP Server
func NewRouter(config Config) http.Handler {
	// Create router and define routes and return that router
	r := mux.NewRouter()
	//r.Use(AuthMiddleware())
	//r.Use(TenantMiddleware)

	// APIs that require token
	internalRouter := r.PathPrefix("/i/v1").Subrouter()
	internalRouter.HandleFunc("/user-verification", handlerUserAccount).Methods(http.MethodPost)
	internalRouter.HandleFunc("/change-password", handleChangePassword).Methods(http.MethodPut)

	// api that requires to interact with Flash service
	//backendRouter := r.PathPrefix("/b/v1").Subrouter()
	//backendRouter.HandleFunc("/callback", handleCallback).Methods(http.MethodGet)

	// APIs that does not require token
	externalRouter := r.PathPrefix("/e").Subrouter()

	externalRouter.HandleFunc("/login-user", handlerLogin).Methods(http.MethodPost)
	externalRouter.HandleFunc("/login-third-party", loginViaThirdParty).Methods(http.MethodGet)
	externalRouter.HandleFunc("/callback", handleCallback).Methods(http.MethodPost)
	externalRouter.HandleFunc("/login", handlerLoginUser).Methods(http.MethodPost)
	externalRouter.HandleFunc("/register", handlerRegisterUser).Methods(http.MethodPost)
	externalRouter.HandleFunc("/test", handleTest).Methods(http.MethodGet)

	externalRouter.HandleFunc("/sports", handleGetSports).Methods(http.MethodGet)
	externalRouter.HandleFunc("/upload-image", uploadImage).Methods(http.MethodPost)
	externalRouter.HandleFunc("/create-schema", createSchema).Methods(http.MethodPost)
	externalRouter.HandleFunc("/callback-oauth2", handleCallbackOAuth2).Methods(http.MethodGet)

	// API FOR BOOKING TICKET PROJECT
	externalRouter.HandleFunc("/queue", handleQueue).Methods(http.MethodGet)
	externalRouter.HandleFunc("/pdf", handlePushFileToQueue).Methods(http.MethodPost)
	externalRouter.HandleFunc("/gemini/{param}", handleGemini).Methods(http.MethodGet)
	externalRouter.HandleFunc("/gemini/image/{param}", handleGeminiImage).Methods(http.MethodGet)

	// Serving static files from the "./static" directory
	r.PathPrefix("/web/").Handler(http.StripPrefix("/web/", http.FileServer(http.Dir("dist"))))

	// Set up cron job to run sendDailyEmail at 7 AM daily
	//con := cron.New()
	//con.AddFunc("0 7 * * *", user.SendDailyEmail)
	//con.Start()

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{"*"},
	})

	handler := c.Handler(r)
	return handler
}

// Run will run the HTTP Server
func (config Config) Run() {
	// Set up a channel to listen to for interrupt signals
	var runChan = make(chan os.Signal, 1)

	// Set up a context to allow for graceful server shutdowns in the event
	// of an OS interrupt (defers the cancel just in case)
	ctx, cancel := context.WithTimeout(
		context.Background(),
		1200*time.Second,
	)
	defer cancel()

	// Define server options
	server := &http.Server{
		Addr:         config.Server.Host + ":" + config.Server.Port,
		Handler:      NewRouter(config),
		ReadTimeout:  config.Server.Timeout.Read * time.Second,
		WriteTimeout: config.Server.Timeout.Write * time.Second,
		IdleTimeout:  config.Server.Timeout.Idle * time.Second,
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS13,
		},
	}

	// Handle ctrl+c/ctrl+x interrupt
	signal.Notify(runChan, os.Interrupt, syscall.SIGTSTP)

	// Alert the user that the server is starting
	log.Printf("Server is starting on %s\n", server.Addr)
	//serverCrt := "/Users/minhtong/go-server.crt"
	//serverKey := "/Users/minhtong/go-server.key"

	// Run the server on a new goroutine

	go func() {
		log.Infof("Before Listen and Serve")
		if err := server.ListenAndServe(); err != nil {
			//if err := server.ListenAndServeTLS(serverCrt, serverKey); err != nil {
			if err == http.ErrServerClosed {
				// Normal interrupt operation, ignore
			} else {
				log.Fatalf("Server failed to start due to err: %v", err)
			}
		}
	}()

	log.Printf("After ListenAndServe")
	// Block on this channel listeninf for those previously defined syscalls assign
	// to variable so we can let the user know why the server is shutting down
	interrupt := <-runChan

	// If we get one of the pre-prescribed syscalls, gracefully terminate the server
	// while alerting the user
	log.Printf("Server is shutting down due to %+v\n", interrupt)
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server was unable to gracefully shutdown due to err: %+v", err)
	}
}

func respondWithJSON(w http.ResponseWriter, httpStatusCode int, data interface{}) {
	resp, err := json.Marshal(data)
	if err != nil {
		log.WithError(err).WithField("data", data).Error("failed to marshal data")
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)
	w.Write(resp)
	return
}

func Init(c Config) {
	userService = c.UserService
	authService = c.AuthService
	jwtService = c.JwtService
	classService = c.ClassService
	flashClient = c.FlashClient
	redisClient = c.RedisClient
}

//func HandleMessageFromQueue(message []byte) error {
//	// Xử lý message nhận được từ hàng đợi
//	log.Printf(" > Received message: %s\n", message)
//	var rq response.YearlyResponse
//	err := json.Unmarshal(message, &rq)
//	if err != nil {
//		log.WithError(err).Errorf("Error marshal request to fix course information")
//		return err
//	}
//	stringValue := fmt.Sprintf("%.6f", rq.TotalYearlyRevenue)
//	err = redisClient.Save(rq.Year, stringValue, context.Background())
//	tmp, err := redisClient.Get(rq.Year, context.Background())
//	log.Printf("Value from redis: %s", tmp)
//	if err != nil {
//		log.WithError(err).Errorf("Error save yearly revenue to redis")
//		log.WithError(err).Errorf("Error message: %s", message)
//		return err
//	}
//	return err
//}
