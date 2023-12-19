package main

import (
	"education-website/api"
	authService2 "education-website/service/authService"
	"education-website/service/courseClass"
	"education-website/service/user"
	"education-website/store"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
)

func main() {
	// Set the log format to plain text
	f, err := os.OpenFile("code/batman.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	// Create a MultiWriter to write logs to both os.Stdout and the file
	wrt := io.MultiWriter(os.Stdout, f)
	// Set the log output to the MultiWriter
	log.SetOutput(wrt)
	//log.Println("Orders API Called")

	log.Info("Hello, service Batman is running")

	// Generate our config based on the config supplied
	// by the user in the flags
	cfgPath, err := api.ParseFlags()
	if err != nil {
		log.WithError(err).Errorf("Error setting path to config file: %s", err)
		log.Fatal(err)
	}

	// Load configuration from config.yml file
	cfg, err := api.NewConfig(cfgPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize the database connection
	db, err := InitDatabase(*cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	classService := courseClass.NewClassService(courseClass.ClassServiceCfg{
		ClassStore: store.NewClassManagementStore(store.ClassManagementStoreCfg{
			Db: db,
		}),
	})
	cfg.ClassService = classService

	userService := user.NewUserService(user.UserServiceCfg{
		UserStore: store.NewUserManagementStore(store.UserManagementStoreCfg{
			Db: db,
		}),
	})
	cfg.UserService = userService

	jwtService := authService2.NewJwtService(authService2.JwtServiceCfg{
		SecretKey: cfg.XApiKey,
	})
	cfg.JwtService = jwtService

	authService := authService2.NewAuthService(authService2.AuthServiceCfg{
		JwtService: jwtService,
	})
	cfg.AuthService = authService

	log.Printf("Successful connect to database")
	defer db.Close() // Close the database connection when finished

	apiCfg := api.Config{
		Server:       cfg.Server,
		Database:     cfg.Database,
		XApiKey:      cfg.XApiKey,
		UserService:  cfg.UserService,
		JwtService:   cfg.JwtService,
		AuthService:  cfg.AuthService,
		ClassService: cfg.ClassService,
	}
	api.Init(apiCfg)
	// Run the server
	apiCfg.Run()
}

func InitDatabase(config api.Config) (*sqlx.DB, error) {
	// Create a MySQL data source name (DSN) using the configuration
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s",
		config.Database.User,
		config.Database.Password,
		config.Database.Host,
		config.Database.Port,
		config.Database.DbName,
	)

	// Open a database connection
	log.Info("Open a database connection")
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
