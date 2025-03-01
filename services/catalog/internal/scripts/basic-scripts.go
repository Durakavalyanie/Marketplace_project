package scripts

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/Dyrakavalyanie/Clothes_shop/services/catalog/internal/models"
)

func AddItem(ctx context.Context, poolConn *pgxpool.Pool, item *models.Item) (*models.ItemID, error) {
	addItemScript := `
		INSERT INTO Items(category, price, status, brand, color, size, sex, description, created_at, sold_at, seller_id, buyer_id)
		values($1, $2, $3, $4, $5, $6, $7, $8, COALESCE($9, NOW()), $10, $11, $12) RETURNING id;`

	var itemID models.ItemID
	err := poolConn.QueryRow(context.Background(), addItemScript, item.Category, item.Price,
		item.Status, item.Brand, item.Color, item.Size, item.Sex, item.Description, item.Created_at,
		item.Sold_at, item.Seller_id, item.Buyer_id).Scan(&itemID.UUID)

	if err != nil {
		return nil, fmt.Errorf("error adding item: %w", err)
	}

	return &itemID, nil
}

func UpdateItem(ctx context.Context, poolConn *pgxpool.Pool, item *models.ItemUpdate) error {
	addItemScript, args := generateUpdateQuery(item)

	_, err := poolConn.Exec(context.Background(), addItemScript, args...)

	if err != nil {
		return fmt.Errorf("error updating item: %w", err)
	}

	return nil
}

func DeleteItem(ctx context.Context, poolConn *pgxpool.Pool, item *models.ItemID) error {
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
	err := connPool.QueryRow(ctx, getItemScript, item.UUID).Scan(nil, &foundItem.Category, &foundItem.Price,
		&foundItem.Status, &foundItem.Brand, &foundItem.Color, &foundItem.Size, &foundItem.Sex,
		&foundItem.Description, &foundItem.Created_at, &foundItem.Sold_at, &foundItem.Seller_id, &foundItem.Buyer_id)
	if err != nil {
		return nil, fmt.Errorf("error in finding item: %w", err)
	}

	return &foundItem, nil
}

func AddPhoto(ctx context.Context, connPool *pgxpool.Pool, photo *models.Photo) (*models.PhotoID, error) {
	addPhotoScript := "INSERT INTO Items_photos(item_id, display_order) values($1, $2) RETURNING id"

	var photoID models.PhotoID
	err := connPool.QueryRow(ctx, addPhotoScript, photo.ItemUUID, photo.DisplayOrder).Scan(&photoID.UUID)
	if err != nil {
		return nil, fmt.Errorf("error while adding photo")
	}

	return &photoID, nil
}

func DeletePhoto(ctx context.Context, connPool *pgxpool.Pool, item *models.PhotoID) (error) {
	addPhotoScript := "DELETE FROM Items_photos WHERE id = $1"

	_, err := connPool.Exec(ctx, addPhotoScript, item.UUID)
	if err != nil {
		return fmt.Errorf("error while deleting photo")
	}

	return nil
}

func GetPhotoPath(ctx context.Context, connPool *pgxpool.Pool, item *models.PhotoID) (string, error) {
	addPhotoScript := "SELECT photo_path FROM Items_photos WHERE id = $1"

	var path string
	err := connPool.QueryRow(ctx, addPhotoScript, item.UUID).Scan(&path)
	if err != nil || path == "" {
		return "", fmt.Errorf("error while searching photo")
	}

	return path, nil
}


func generateUpdateQuery(item *models.ItemUpdate) (string, []interface{}) {
	setClauses := []string{}
	args := []interface{}{item.UUID}
	argCounter := 2

	if item.Category != nil {
		setClauses = append(setClauses, fmt.Sprintf("category = $%d", argCounter))
		args = append(args, *item.Category)
		argCounter++
	}

	if item.Price != nil {
		setClauses = append(setClauses, fmt.Sprintf("price = $%d", argCounter))
		args = append(args, *item.Price)
		argCounter++
	}

	if item.Status != nil {
		setClauses = append(setClauses, fmt.Sprintf("status = $%d", argCounter))
		args = append(args, *item.Status)
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

	if item.Sex != nil {
		setClauses = append(setClauses, fmt.Sprintf("sex = $%d", argCounter))
		args = append(args, *item.Sex)
		argCounter++
	}

	if item.Description != nil {
		setClauses = append(setClauses, fmt.Sprintf("description = $%d", argCounter))
		args = append(args, *item.Description)
		argCounter++
	}

	if item.Created_at != nil {
		setClauses = append(setClauses, fmt.Sprintf("created_at = $%d", argCounter))
		args = append(args, *item.Created_at)
		argCounter++
	}

	if item.Sold_at != nil {
		setClauses = append(setClauses, fmt.Sprintf("sold_at = $%d", argCounter))
		args = append(args, *item.Sold_at)
		argCounter++
	}

	if item.Seller_id != nil {
		setClauses = append(setClauses, fmt.Sprintf("seller_id = $%d", argCounter))
		args = append(args, *item.Seller_id)
		argCounter++
	}

	if item.Buyer_id != nil {
		setClauses = append(setClauses, fmt.Sprintf("buyer_id = $%d", argCounter))
		args = append(args, *item.Buyer_id)
		argCounter++
	}

	if len(setClauses) == 0 {
		return "", nil
	}

	query := "UPDATE Items SET " + strings.Join(setClauses, ", ") + " WHERE id = $1"
	return query, args
}
