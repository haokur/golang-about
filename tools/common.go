package tools

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// 获取用户的路径
func GetUserHomePath() string {
	dirPath, err := os.UserHomeDir()
	if err != nil {
		return "~"
	}
	return dirPath
}

// 获取git根目录
func GetGitRootDir() (string, error) {
	// 使用 'git rev-parse --show-toplevel' 获取Git根目录
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get git root directory: %w", err)
	}

	// 去除输出中的换行符和空白
	gitRootDir := strings.TrimSpace(string(out))
	return gitRootDir, nil
}
