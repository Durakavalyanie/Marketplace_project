package scripts

import (
	"context"
	"fmt"
	"strings"

	//"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/Dyrakavalyanie/Clothes_shop/services/catalog/internal/models"
)

func AddItem(ctx context.Context, poolConn *pgxpool.Pool, item *models.Item) (*models.ItemID, error) {
	addItemScript := `
		INSERT INTO Items(sex, category, brand, color, size, description, price)
		values($1, $2, $3, $4, $5, $6, $7) RETURNING id;`

	var itemID models.ItemID
	err := poolConn.QueryRow(context.Background(), addItemScript, item.Sex, item.Category,
		item.Brand, item.Color, item.Size, item.Description, item.Price).Scan(&itemID.UUID)

	if err != nil {
		return nil, fmt.Errorf("error adding item: %w", err)
	}

	return &itemID, nil
}

func UpdateItem(ctx context.Context, poolConn *pgxpool.Pool, item *models.ItemUpdate) (error) {
	addItemScript, args := generateUpdateQuery(item)

	_, err := poolConn.Exec(context.Background(), addItemScript, args...)

	if err != nil {
		return fmt.Errorf("error updating item: %w", err)
	}

	return nil
}

func DeleteItem(ctx context.Context, poolConn *pgxpool.Pool, item *models.ItemID) (error) {
	deleteItemScript := "DELETE FROM Items WHERE id = $1"

	_, err := poolConn.Exec(context.Background(), deleteItemScript, item.UUID)

	if err != nil {
		return fmt.Errorf("error in deleting item: %w", err)
	}

	return nil
}

func GetItem(ctx context.Context, connPool *pgxpool.Pool, item *models.ItemID) (*models.Item, error) {
	getItemScript := "SELECT * FROM Items WHERE id = $1"

	var foundItem models.Item
	err := connPool.QueryRow(ctx, getItemScript, item.UUID).Scan(nil, &foundItem.Sex, &foundItem.Category,
		&foundItem.Brand, &foundItem.Color, &foundItem.Size, &foundItem.Description, &foundItem.Price)
	if err != nil {
		return nil, fmt.Errorf("error in finding item: %w", err)
	}

	return &foundItem, nil
}


func generateUpdateQuery(item *models.ItemUpdate) (string, []interface{}) {
	setClauses := []string{}
	args := []interface{}{}
	argCounter := 1
   
	if item.Sex != nil {
	 	setClauses = append(setClauses, fmt.Sprintf("sex = $%d", argCounter))
		args = append(args, *item.Sex)
	 	argCounter++
	}
	if item.Category != nil {
	 	setClauses = append(setClauses, fmt.Sprintf("category = $%d", argCounter))
	 	args = append(args, *item.Category)
	 	argCounter++
	}
	if item.Brand != nil {
	 	setClauses = append(setClauses, fmt.Sprintf("brand = $%d", argCounter))
	 	args = append(args, *item.Brand)
	 	argCounter++
	}
	if item.Color != nil {
	 	setClauses = append(setClauses, fmt.Sprintf("color = $%d", argCounter))
	 	args = append(args, *item.Color)
	 	argCounter++
	}
	if item.Size != nil {
	 	setClauses = append(setClauses, fmt.Sprintf("size = $%d", argCounter))
	 	args = append(args, *item.Size)
	 	argCounter++
	}
	if item.Description != nil {
	 	setClauses = append(setClauses, fmt.Sprintf("description = $%d", argCounter))
	 	args = append(args, *item.Description)
	 	argCounter++
	}
	if item.Price != nil {
	 	setClauses = append(setClauses, fmt.Sprintf("price = $%d", argCounter))
	 	args = append(args, *item.Price)
	 	argCounter++
	}
   
	if len(setClauses) == 0 {
		return "", nil
	}

	setClauses = append(setClauses, fmt.Sprintf("uuid = $%d", argCounter))
	args = append(args, item.UUID)
   
	query := "UPDATE Items SET " + strings.Join(setClauses, ", ") + " WHERE uuid = $1"
	return query, args
}