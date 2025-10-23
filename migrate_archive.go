package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
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

func migrateArchive(gitlabURL, token, projectID, archivePath string) error {
	// Открываем архив
	file, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Создаем multipart форму
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filepath.Base(archivePath))
	if err != nil {
		return err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return err
	}

	err = writer.Close()
	if err != nil {
		return err
	}

	// Отправляем в GitLab
	url := fmt.Sprintf("%s/api/v4/projects/%s/import", gitlabURL, projectID)
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return err
	}

	req.Header.Set("PRIVATE-TOKEN", token)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("GitLab API вернул ошибку: %s", resp.Status)
	}

	return nil
}
