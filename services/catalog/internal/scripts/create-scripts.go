package scripts

import (
	"context"
	"fmt"
	//"log"

	//"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/Dyrakavalyanie/Clothes_shop/services/catalog/internal/config"
)

func CreateDBConnection(cfg config.Config) (*pgxpool.Pool, error) {
	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.SSLMode)

	connPool, err := pgxpool.Connect(context.Background(), connStr);
	if err != nil {
		return nil, fmt.Errorf("error in connecting to db: %w", err)
	}

	return connPool, nil
}

func CreateCatalogTables(pool *pgxpool.Pool) error {
	createCatalogTablesScript := `
		CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

		DO $$
		BEGIN
    		IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'sex_var') THEN
        		CREATE TYPE sex_var AS ENUM ('М', 'Ж', 'У');
    		END IF;
    
    		IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'category_var') THEN
        		CREATE TYPE category_var AS ENUM ('Головной убор', 'Верх', 'Низ', 'Обувь', 'Аксессуар', 'Прочее');
    		END IF;
		END $$;

		CREATE TABLE IF NOT EXISTS Items (
    		id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    		sex sex_var NOT NULL,
    		category category_var NOT NULL,
    		brand TEXT NOT NULL,
    		color TEXT NOT NULL,
    		size TEXT NOT NULL,
    		description TEXT,
    		price INTEGER NOT NULL
		);

		CREATE TABLE IF NOT EXISTS Items_photos(
    		id SERIAL PRIMARY KEY,
    		item_id UUID REFERENCES Items(id) ON DELETE CASCADE,
    		photo_path TEXT NOT NULL
		);`
	_, err := pool.Exec(context.Background(), createCatalogTablesScript)
	if err != nil {
		return fmt.Errorf("error in creating catalog tables: %w", err)
	}
	return nil
}