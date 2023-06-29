package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/zelalemmk/key-val-store/cmd/web/handlers"
)

var storagePath = "/tmp/keyVal"

const (
	fileName = "keyVal.json"
	filePerm = 0644
	dirPerm  = 0755
)

func LoadData() {
	if err := createDirIfNotExist(storagePath); err != nil {
		log.Fatalf("Failed to create storage directory: %v", err)
	}

	filePath := filepath.Join(storagePath, fileName)

	file, err := os.OpenFile(filePath, os.O_RDONLY|os.O_CREATE, filePerm)
	if err != nil {
		log.Fatalf("Failed to open data file: %v", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&handlers.KeyVal); err != nil {
		log.Printf("Failed to decode data file: %v", err)
		return
	}

	log.Println("Data loaded from disk storage file")
}

func SaveData() {
	if err := createDirIfNotExist(storagePath); err != nil {
		log.Fatalf("Failed to create storage directory: %v", err)
	}

	filePath := filepath.Join(storagePath, fileName)

	file, err := os.Create(filePath)
	if err != nil {
		log.Fatalf("Failed to create data file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(handlers.KeyVal); err != nil {
		log.Fatalf("Failed to encode data: %v", err)
	}

	log.Println("Data saved to disk storage file")
}

func GenerateData(numEntries int) error {
	for i := 0; i < numEntries; i++ {
		key := fmt.Sprintf("key%d", i)
		value := fmt.Sprintf("value%d", i)
		handlers.KeyVal[key] = value
	}

	SaveData()

	return nil
}

func createDirIfNotExist(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(path, dirPerm); err != nil {
			return err
		}
	}

	return nil
}
