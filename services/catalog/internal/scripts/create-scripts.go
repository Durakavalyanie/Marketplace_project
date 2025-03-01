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

	connPool, err := pgxpool.Connect(context.Background(), connStr)
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

			IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'status_var') THEN
        		CREATE TYPE status_var AS ENUM ('В наличии', 'Продано');
    		END IF;
		END $$;

		CREATE TABLE IF NOT EXISTS Items (
    		id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
			category category_var NOT NULL,
			price INTEGER NOT NULL,
			status status_var NOT NULL,
    		brand TEXT NOT NULL,
    		color TEXT NOT NULL,
    		size TEXT NOT NULL,
			sex sex_var NOT NULL,
    		description TEXT DEFAULT NULL,
			created_at TIMESTAMP DEFAULT NOW(),
			sold_at TIMESTAMP DEFAULT NULL,
			seller_id UUID NOT NULL,
			buyer_id UUID DEFAULT NULL
		);

		CREATE TABLE IF NOT EXISTS Items_photos(
    		id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    		item_id UUID REFERENCES Items(id) ON DELETE CASCADE,
    		photo_path TEXT NOT NULL,
			display_order INTEGER NOT NULL
		);
		
		CREATE OR REPLACE FUNCTION generate_photo_path() 
		RETURNS TRIGGER AS $$
		BEGIN
    		NEW.photo_path := CONCAT(NEW.item_id, '/', NEW.id, '.jpg');
			RETURN NEW;
		END;
		$$ LANGUAGE plpgsql;

		DROP TRIGGER IF EXISTS set_photo_path on Items_photos;
		
		CREATE TRIGGER set_photo_path
		BEFORE INSERT ON Items_photos
		FOR EACH ROW
		EXECUTE FUNCTION generate_photo_path();`
	_, err := pool.Exec(context.Background(), createCatalogTablesScript)
	if err != nil {
		return fmt.Errorf("error in creating catalog tables: %w", err)
	}
	return nil
}
