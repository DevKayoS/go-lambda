package database

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

func Init(ctx context.Context) error {
	database := os.Getenv("DATABASE_URL")

	if database == "" {
		return fmt.Errorf("DATABASE_URL is not configured")
	}

	config, err := pgxpool.ParseConfig(database)
	if err != nil {
		return err
	}

	config.MaxConns = 2                      // Lambda = poucos workers, não precisa de muito
	config.MinConns = 1                      // Mínimo 1 para manter conexão viva
	config.MaxConnLifetime = 5 * time.Minute // Recicla conexões antigas
	config.MaxConnIdleTime = 1 * time.Minute // Fecha conexões ociosas
	config.HealthCheckPeriod = 1 * time.Minute

	Pool, err = pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return fmt.Errorf("Erro ao realizar o pool: ", err.Error())
	}

	if err = Pool.Ping(ctx); err != nil {
		return fmt.Errorf("Erro ao tentar realizar o ping: ", err.Error())
	}

	return nil
}

func Close() {
	if Pool != nil {
		Pool.Close()
	}
}
