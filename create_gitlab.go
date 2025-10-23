package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type GitLabProject struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Visibility  string `json:"visibility,omitempty"`
}

func createGitlab() {
	if len(os.Args) < 5 {
		fmt.Println("Использование: create_gitlab.exe <gitlab_url> <token> <project_name> [description]")
		fmt.Println("Пример: create_gitlab.exe http://localhost:8080 glpat-token123 myproject \"Мой проект\"")
		os.Exit(1)
	}

	gitlabURL := os.Args[1]
	token := os.Args[2]
	projectName := os.Args[3]
	description := ""
	if len(os.Args) > 4 {
		description = os.Args[4]
	}

	fmt.Printf("🏗️ Создаем проект %s в GitLab...\n", projectName)

	projectID, err := createGitLabProject(gitlabURL, token, projectName, description)
	if err != nil {
		log.Fatal("Ошибка создания проекта:", err)
	}

	fmt.Printf("✅ Проект успешно создан. ID: %d\n", projectID)
}

func createGitLabProject(gitlabURL, token, name, description string) (int, error) {
	project := GitLabProject{
		Name:        name,
		Description: description,
		Visibility:  "private",
	}

	jsonData, err := json.Marshal(project)
	if err != nil {
		return 0, err
	}

	url := fmt.Sprintf("%s/api/v4/projects", gitlabURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return 0, err
	}

	req.Header.Set("PRIVATE-TOKEN", token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return 0, fmt.Errorf("GitLab API вернул ошибку: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}

	// Преобразуем ID в int
	idFloat, ok := result["id"].(float64)
	if !ok {
		return 0, fmt.Errorf("неверный формат ID проекта")
	}

	return int(idFloat), nil
}
