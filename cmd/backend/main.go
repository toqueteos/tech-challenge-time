package main

import (
	"os"
	"strings"
	"time"

	"pento-tech-challenge/backend/api"
	"pento-tech-challenge/backend/timer"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func main() {
	logger, _ := newLogger()
	defer logger.Sync()

	logger.Info(
		"BACKEND_START",
		zap.Any("when", time.Now()),
	)

	db := newDatabase(logger)
	runMigrations(logger, db)

	opt := api.Opts{
		Host:         "0.0.0.0",
		Port:         8080,
		Logger:       logger,
		DB:           db,
		TimerService: timer.NewService(timer.NewRepo(db)),
	}
	srv := api.NewServer(opt)

	logger.Info(
		"HTTP_START",
		zap.String("host", opt.Host),
		zap.Uint("port", opt.Port),
	)

	if err := srv.Run(); err != nil {
		logger.Fatal(
			"HTTP_STOPPED",
			zap.Error(err),
		)
	}
}

func newLogger() (*zap.Logger, error) {
	if env := strings.ToLower(os.Getenv("ENVIRONMENT")); env == "production" {
		return zap.NewProduction()
	}
	return zap.NewDevelopment()
}

func newDatabase(log *zap.Logger) *sqlx.DB {
	// DISCLAIMER: For simplicity I'm hardcoding the connection string here.
	addr := os.Getenv("DATABASE_ADDR")
	if len(addr) == 0 {
		addr = "postgres://postgres@localhost:5432/pento?sslmode=disable"
	}

	db, err := sqlx.Connect("postgres", addr)
	if err != nil {
		log.Fatal(
			"DATABASE_CONNECT_FAIL",
			zap.Error(err),
		)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(
			"DATABASE_PING_FAIL",
			zap.Error(err),
		)
	}

	// DISCLAIMER: These are magic numbers, which WOULD be configurable
	// in a real app.
	db.SetMaxIdleConns(16)
	db.SetMaxOpenConns(32)

	log.Info(
		"DATABASE_PING_OK",
		zap.String("addr", addr),
	)

	return db
}

func runMigrations(log *zap.Logger, db *sqlx.DB) {
	handleError := func(err error) {
		if err == migrate.ErrNoChange {
			return
		}

		if err != nil {
			log.Fatal(
				"DATABASE_MIGRATION_FAIL",
				zap.Error(err),
			)
		}
	}

	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	handleError(err)

	m, err := migrate.NewWithDatabaseInstance("file://backend/migrations", "postgres", driver)
	handleError(err)

	err = m.Up()
	handleError(err)

	log.Info("DATABASE_MIGRATIONS_OK")
}
