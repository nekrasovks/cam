package main

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func archive() {
	if len(os.Args) < 3 {
		fmt.Println("Использование: archive_repo.exe <source_path> <archive_name>")
		fmt.Println("Пример: archive_repo.exe .\\clone myrepo_backup.zip")
		os.Exit(1)
	}

	sourcePath := os.Args[1]
	archiveName := os.Args[2]

	fmt.Printf("📦 Архивируем %s в %s...\n", sourcePath, archiveName)

	if err := createArchive(sourcePath, archiveName); err != nil {
		log.Fatal("Ошибка архивации:", err)
	}

	fmt.Println("✅ Архив успешно создан:", archiveName)
}

func createArchive(sourcePath, archiveName string) error {
	// Создаем архив
	archiveFile, err := os.Create(archiveName)
	if err != nil {
		return err
	}
	defer archiveFile.Close()

	zipWriter := zip.NewWriter(archiveFile)
	defer zipWriter.Close()

	// Функция для обхода файлов
	return filepath.Walk(sourcePath, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Пропускаем системные файлы Git
		if strings.Contains(filePath, "\\.git\\") || strings.HasSuffix(filePath, "\\.git") {
			return nil
		}

		// Получаем относительный путь
		relPath, err := filepath.Rel(sourcePath, filePath)
		if err != nil {
			return err
		}

		// Создаем заголовок файла в архиве
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// Устанавливаем правильное имя файла
		header.Name = filepath.ToSlash(relPath)

		// Для директорий добавляем слеш
		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate // Сжатие
		}

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		// Если это файл - копируем содержимое
		if !info.IsDir() {
			file, err := os.Open(filePath)
			if err != nil {
				return err
			}
			defer file.Close()

			_, err = io.Copy(writer, file)
			if err != nil {
				return err
			}
		}

		return nil
	})
}
