package common

import (
	"fmt"
	"os"
)

func CreateStaticFolder(folderPath string) {
	// Check if the folder exists
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		// Folder doesn't exist, create it
		err := os.MkdirAll(folderPath, 0755) // 0755 is the permission mode for the directory
		if err != nil {
			fmt.Println("Error creating folder:", err)
			return
		}
	}
}
