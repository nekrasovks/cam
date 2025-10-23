package main

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

func migrateArch() {
	if len(os.Args) < 5 {
		fmt.Println("Использование: migrate_to_gitlab.exe <gitlab_url> <token> <project_id> <archive_path>")
		fmt.Println("Пример: migrate_to_gitlab.exe http://localhost:8080 glpat-token123 15 myrepo_backup.zip")
		os.Exit(1)
	}

	gitlabURL := os.Args[1]
	token := os.Args[2]
	projectID := os.Args[3]
	archivePath := os.Args[4]

	fmt.Printf("🚚 Мигрируем архив %s в GitLab проект %s...\n", archivePath, projectID)

	if err := migrateArchive(gitlabURL, token, projectID, archivePath); err != nil {
		log.Fatal("Ошибка миграции:", err)
	}

	fmt.Println("✅ Миграция завершена успешно!")
}

type CreateCommitBody struct {
	Branch        string `json:"branch"`         //"main",
	Content       string `json:"content"`        // "some content"
	CommitMessage string `json:"commit_message"` // "commit_message": "create a new file"
	Encoding      string `json:"encoding"`
}

func migrateArchive(gitlabURL, token, projectID, archivePath string) error {
	// Открываем архив
	fileContent, err := os.ReadFile(archivePath)
	if err != nil {
		return err
	}

	// Отправляем в GitLab
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	url := fmt.Sprintf("%s/api/v4/projects/%s/repository/files/%s", gitlabURL, projectID, archivePath)

	body := CreateCommitBody{
		Branch:        "main",
		CommitMessage: "copy from git",
		Content:       base64.StdEncoding.EncodeToString(fileContent),
		Encoding:      "base64",
	}

	rawBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(rawBody))
	if err != nil {
		return err
	}

	req.Header.Set("PRIVATE-TOKEN", token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("GitLab API вернул ошибку: %s", resp.Status)
	}

	return nil
}
