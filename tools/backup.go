package tools

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// 获取项目名称（Git 根目录的名称）
func getProjectName() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get project name: %w", err)
	}
	return filepath.Base(strings.TrimSpace(string(out))), nil
}

// 获取未提交的文件列表
func getUncommittedFiles(sourceDir string) ([]string, error) {
	os.Chdir(sourceDir)
	cmd := exec.Command("git", "status", "--porcelain")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get uncommitted files: %w", err)
	}

	var files []string
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if len(line) > 3 {
			files = append(files, strings.TrimSpace(line[3:])) // 截取文件路径部分
		}
	}
	return files, nil
}

// 复制文件到目标路径，保持层级结构
func copyFile(src, dest string) error {
	// 创建目标文件夹
	if err := os.MkdirAll(filepath.Dir(dest), os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directories: %w", err)
	}

	// 打开源文件
	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer sourceFile.Close()

	// 创建目标文件
	destinationFile, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destinationFile.Close()

	// 复制内容
	if _, err := io.Copy(destinationFile, sourceFile); err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	return nil
}

// 备份未提交的文件
func backupUncommittedFiles(sourceDir string, backupDir string) error {
	files, err := getUncommittedFiles(sourceDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		source := filepath.Join(sourceDir, file)
		dest := filepath.Join(backupDir, file)
		fmt.Printf("Backing up: %s -> %s\n", source, dest)
		if err := copyFile(source, dest); err != nil {
			return fmt.Errorf("failed to backup file %s: %w", source, err)
		}
	}

	return nil
}

func StartBackupUnCommitFiles(sourceDir string, targetDir string) {
	// 获取当前时间戳
	timestamp := time.Now().Format("2006_01_02_150405")

	// 创建备份文件夹
	backupDir := filepath.Join(targetDir, timestamp)
	if err := os.MkdirAll(backupDir, os.ModePerm); err != nil {
		fmt.Println("Error creating backup directory:", err)
		return
	}

	// 备份未提交的文件
	if err := backupUncommittedFiles(sourceDir, backupDir); err != nil {
		fmt.Println("Error backing up files:", err)
		return
	}

	fmt.Println("Backup completed successfully.")
}
