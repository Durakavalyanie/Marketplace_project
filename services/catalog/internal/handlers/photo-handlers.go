package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"github.com/Dyrakavalyanie/Clothes_shop/services/catalog/internal/models"
	"github.com/Dyrakavalyanie/Clothes_shop/services/catalog/internal/scripts"
	"github.com/disintegration/imaging"
)


func (handler *Handler)AddPhotos(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	r.ParseMultipartForm(10 << 20)

	var itemID models.ItemID
	metaJSON := r.FormValue("json")

	err := json.Unmarshal([]byte(metaJSON), &itemID)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
  		return
	}

	saveDir := handler.Cfg.PhotosStoragePath + "/" + itemID.UUID.String()
	err = os.MkdirAll(saveDir, os.ModePerm)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	files := r.MultipartForm.File["files"]
	var filesCounter int32 = 1
	var photo models.Photo
	photo.ItemUUID = &itemID.UUID
	photo.DisplayOrder = &filesCounter

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
		 	http.Error(w, "Failed to open file", http.StatusInternalServerError)
		 	return
		}
		defer file.Close()

		photoID, err := scripts.AddPhoto(ctx, handler.ConnPool, &photo)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		*photo.DisplayOrder++

		outFile, err := os.Create(saveDir + "/" + photoID.UUID.String() + ".jpg")
		if err != nil {
		 	http.Error(w, "Failed to prepare file", http.StatusInternalServerError)
		 	return
		}
		defer outFile.Close()
	  
		err = SaveFileToJPG(outFile, &file)
		if err != nil {
			err = scripts.DeletePhoto(ctx, handler.ConnPool, photoID)
			if (err != nil) {
				http.Error(w, "Failed to save file to jpg and to clear database", http.StatusInternalServerError)
			}
			http.Error(w, "Failed to save file to jpg", http.StatusInternalServerError)
		 	return
		}
	}
}

func (handler *Handler)DeletePhoto(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "error", http.StatusBadRequest)
		return
	}

	var photoID models.PhotoID
	err = json.Unmarshal(body, &photoID)
	if err != nil {
		http.Error(w, "Json parsing error", http.StatusBadRequest)
		return
	}

	var photoPath string
	photoPath, err = scripts.GetPhotoPath(ctx, handler.ConnPool, &photoID)
	if err != nil {
		http.Error(w, "error while handling photo data", http.StatusInternalServerError)
		return
	}

	err = scripts.DeletePhoto(ctx, handler.ConnPool, &photoID)
	if err != nil {
		http.Error(w, "error while deleting photo data", http.StatusInternalServerError)
		return
	}

	err = os.Remove(handler.Cfg.PhotosStoragePath + "/" + photoPath)
	if err != nil {
		http.Error(w, "error while deleting photo file", http.StatusInternalServerError)
		return
	}

	err = os.Remove(handler.Cfg.PhotosStoragePath + "/" + strings.SplitN(photoPath, "/", 2)[0])
	if err != nil && !strings.Contains(err.Error(), "directory not empty"){
		http.Error(w, "error while deleting empty directory", http.StatusInternalServerError)
		return
	}
}


func SaveFileToJPG(output *os.File, input *multipart.File) (error) {
	img, _, err := image.Decode(*input)
	if err != nil {
		return fmt.Errorf("error in decoding photo: %w", err)
	}

	err = imaging.Encode(output, img, imaging.JPEG)
	if err != nil  {
		return fmt.Errorf("error in encoding photo: %w", err)
	}
	return nil
}