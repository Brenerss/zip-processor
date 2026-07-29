package processor

import (
	"archive/zip"
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"path/filepath"
	"strings"
	"sync"

	"github.com/processorsystem/internal/storage"
)

func ProcessZip(filepath string, fileID int64, db *storage.DynamoStorage) error {
	archive, err := zip.OpenReader(filepath)
	if err != nil {
		return fmt.Errorf("failed to open zip: %v", err)
	}

	defer archive.Close()

	var wg sync.WaitGroup

	for _, f := range archive.File {
		if f.FileInfo().IsDir() {
			continue
		}

		wg.Add(1)

		go func(file *zip.File) {
			defer wg.Done()

			err := processSingleFile(file, fileID, db)
			if err != nil {
				log.Printf("error an image %s: %v", file.Name, err)
			}
		}(f)
	}

	wg.Wait()

	return nil
}

func processSingleFile(file *zip.File, fileID int64, db *storage.DynamoStorage) error {
	log.Printf("extracting and processing: %s", file.Name)

	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	ext := strings.ToLower(filepath.Ext(file.Name))

	metadata := storage.FileMetadata{
		FileID:    fileID,
		Filename:  file.Name,
		Extension: ext,
		Size:      file.FileInfo().Size(),
	}

	switch ext {
	case ".png":
		metadata.Type = "image"
	case ".csv", ".txt":
		metadata.Type = "text"
	case ".mp4":
		metadata.Type = "video"
	default:
		log.Printf("File ignored, format not supported: %s", file.Name)
		return nil
	}

	err = db.SaveMetadata(metadata)
	if err != nil {
		log.Printf("an error ocurred to save on DB: %v", err)
	}

	log.Printf("File %s sucessfully resized", file.Name)

	return nil
}
