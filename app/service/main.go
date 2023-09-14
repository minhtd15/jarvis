package main

import (
	"education-website/api"
	"education-website/service"
	"education-website/store"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"log"
)

func main() {
	logrus.Info("Hello, service Batman is running")

	// Generate our config based on the config supplied
	// by the user in the flags
	cfgPath, err := api.ParseFlags()
	if err != nil {
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

	userService := service.NewUserService(service.UserServiceCfg{
		UserStore: store.NewUserManagementStore(store.UserManagementStoreCfg{
			Db: db,
		}),
	})
	cfg.UserService = userService

	log.Printf("Successful connect to database")
	defer db.Close() // Close the database connection when finished

	apiCfg := api.Config{
		Server:      cfg.Server,
		Database:    cfg.Database,
		UserService: cfg.UserService,
	}
	api.Init(apiCfg)
	// Run the server
	cfg.Run()
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
	logrus.Info("Open a database connection")
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
