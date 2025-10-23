package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

type CloneConfig struct {
	GitHubToken string
	GitHubURL   string
	ClonePath   string
	SinceDate   string
	Depth       int
}

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Использование: clone_repo.exe <github_url> <token> <clone_path> [since_date] [depth]")
		fmt.Println("Пример: clone_repo.exe https://github.com/owner/repo.git token123 .\\clone 2024-01-01 50")
		fmt.Println("Пример: clone_repo.exe https://github.com/owner/repo.git token123 .\\clone full")
		os.Exit(1)
	}

	config := CloneConfig{
		GitHubURL:   os.Args[1],
		GitHubToken: os.Args[2],
		ClonePath:   os.Args[3],
	}

	// Опциональные параметры
	if len(os.Args) > 4 && os.Args[4] != "full" {
		config.SinceDate = os.Args[4]
	}
	if len(os.Args) > 5 && os.Args[5] != "full" {
		fmt.Sscanf(os.Args[5], "%d", &config.Depth)
	}

	fmt.Println("🚀 Клонируем репозиторий...")
	if err := cloneRepository(config); err != nil {
		log.Fatal("Ошибка клонирования:", err)
	}

	fmt.Println("✅ Репозиторий успешно клонирован в", config.ClonePath)

	fmt.Printf("📦 Архивируем %s...\n", config.ClonePath)

	archiveName := "archive.zip"
	if err := createArchive(config.ClonePath, archiveName); err != nil {
		log.Fatal("Ошибка архивации:", err)
	}

	fmt.Println("✅ Архив успешно создан:", archiveName)

	projectName := "project_name"
	fmt.Printf("🏗️ Создаем проект %s в GitLab...\n", projectName)

	gitlabURL := "https://localhost"
	token := "glpat-qWVvgvirOKp11KGzWIrlcW86MQp1OjEH.01.0w0crrxmh"
	// description := "test description"

	// projectID, err := createGitLabProject(gitlabURL, token, projectName, description)
	// if err != nil {
	// 	log.Fatal("Ошибка создания проекта:", err)
	// }

	// fmt.Printf("✅ Проект успешно создан. ID: %d\n", projectID)
	projectID := 7
	fmt.Printf("🚚 Мигрируем архив %s в GitLab проект %d...\n", archiveName, projectID)

	// fmt.Sprintf("%d", projectID)
	if err := migrateArchive(gitlabURL, token, fmt.Sprintf("%d", projectID), archiveName); err != nil {
		log.Fatal("Ошибка миграции:", err)
	}

	fmt.Println("✅ Миграция завершена успешно!")
}

func cloneRepository(config CloneConfig) error {
	// Добавляем токен в URL для аутентификации
	authURL := strings.Replace(config.GitHubURL, "https://",
		fmt.Sprintf("https://%s@", config.GitHubToken), 1)

	args := []string{"clone"}

	// Добавляем глубину истории если указана
	if config.Depth > 0 {
		args = append(args, fmt.Sprintf("--depth=%d", config.Depth))
	}

	// Добавляем фильтр по дате если указана
	if config.SinceDate != "" {
		args = append(args, fmt.Sprintf("--shallow-since=%s", config.SinceDate))
	}

	args = append(args, authURL, config.ClonePath)

	cmd := exec.Command("git", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("Выполняем: git %s\n", strings.Join(args, " "))
	return cmd.Run()
}
