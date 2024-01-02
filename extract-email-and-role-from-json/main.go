package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Define the structure of the JSON data
type Data map[string][]Result

type Result struct {
	UserProfileId            int    `json:"UserProfileId"`
	FirstName                string `json:"FirstName"`
	LastName                 string `json:"LastName"`
	DisplayName              string `json:"DisplayName"`
	Email                    string `json:"Email"`
	Role                     string `json:"Role"`
	RoleId                   int    `json:"RoleId"`
	CreatedDate              string `json:"CreatedDate"`
	Organization             string `json:"Organization"`
	RequestStatus            int    `json:"RequestStatus"`
	RequestStatusDisplayText string `json:"RequestStatusDisplayText"`
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run main.go input.json")
		os.Exit(1)
	}

	inputFilePath := os.Args[1]

	// Open the JSON file
	jsonFile, err := os.Open(inputFilePath)
	if err != nil {
		fmt.Println("Error opening JSON file:", err)
		os.Exit(1)
	}
	defer jsonFile.Close()

	// Decode the JSON file into a Data structure
	var jsonData Data
	jsonDecoder := json.NewDecoder(jsonFile)
	err = jsonDecoder.Decode(&jsonData)
	if err != nil {
		fmt.Println("Error decoding JSON data:", err)
		os.Exit(1)
	}

	// Check for duplicates and extract unique emails and roles
	uniqueEmailRoleMap := make(map[string]string)
	for _, results := range jsonData {
		for _, result := range results {
			emailRole := fmt.Sprintf("Email: %s, Role: %s", result.Email, result.Role)
			uniqueEmailRoleMap[emailRole] = emailRole
		}
	}

	// Create the output directory if it doesn't exist
	outputDir := "output"
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		os.Mkdir(outputDir, 0755)
	}

	// Save unique email-role pairs to "output.txt"
	outputFilePath := filepath.Join(outputDir, "output.txt")
	var lines []string
	for _, uniqueEmailRole := range uniqueEmailRoleMap {
		lines = append(lines, uniqueEmailRole)
	}

	if err := os.WriteFile(outputFilePath, []byte(strings.Join(lines, "\n")), 0644); err != nil {
		fmt.Println("Error writing to output file:", err)
		os.Exit(1)
	}

	fmt.Println("Unique Emails and Roles extracted and saved to:", outputFilePath)
}
