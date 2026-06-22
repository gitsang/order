package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"github.com/gitsang/order/internal/config"
	"github.com/gitsang/order/pkg/database"
	"github.com/gitsang/order/pkg/logger"
)

const passwordHashPlaceholder = "__ADMIN_PASSWORD_HASH__"

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	seedPath := flag.String("seed", "scripts/seed.sql", "path to seed SQL script")
	adminPassword := flag.String("admin-password", "hashed_password", "initial admin password")
	dbHost := flag.String("db-host", "", "database host override")
	dbPort := flag.Int("db-port", 0, "database port override")
	dbUser := flag.String("db-user", "", "database user override")
	dbPassword := flag.String("db-password", "", "database password override")
	dbName := flag.String("db-name", "", "database name override")
	dbSSLMode := flag.String("db-sslmode", "", "database sslmode override")
	flag.Parse()

	applyDatabaseFlags(&cfg.Database, databaseFlags{
		Host:     *dbHost,
		Port:     *dbPort,
		User:     *dbUser,
		Password: *dbPassword,
		DBName:   *dbName,
		SSLMode:  *dbSSLMode,
	})

	log, err := logger.New(cfg.LogLevel)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer log.Sync()

	db, err := database.NewPostgres(cfg.Database)
	if err != nil {
		log.Fatal("Failed to connect to database", zap.Error(err))
	}

	script, err := os.ReadFile(*seedPath)
	if err != nil {
		log.Fatal("Failed to read seed script", zap.String("path", *seedPath), zap.Error(err))
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(*adminPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("Failed to hash admin password", zap.Error(err))
	}

	sql := strings.ReplaceAll(string(script), passwordHashPlaceholder, string(hash))
	if strings.Contains(sql, passwordHashPlaceholder) {
		log.Fatal("Seed script still contains password hash placeholder")
	}

	result := db.Exec(sql)
	if result.Error != nil {
		log.Fatal("Failed to execute seed script", zap.Error(result.Error))
	}

	log.Info(
		"Seed data executed successfully",
		zap.String("script", *seedPath),
		zap.String("database", cfg.Database.DBName),
		zap.Int64("rows_affected", result.RowsAffected),
	)
}

type databaseFlags struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func applyDatabaseFlags(cfg *config.DatabaseConfig, flags databaseFlags) {
	if flags.Host != "" {
		cfg.Host = flags.Host
	}
	if flags.Port != 0 {
		cfg.Port = flags.Port
	}
	if flags.User != "" {
		cfg.User = flags.User
	}
	if flags.Password != "" {
		cfg.Password = flags.Password
	}
	if flags.DBName != "" {
		cfg.DBName = flags.DBName
	}
	if flags.SSLMode != "" {
		cfg.SSLMode = flags.SSLMode
	}
}
