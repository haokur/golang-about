package oss

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/manifoldco/promptui"
)

// window 杀进程
// unix 杀进程
func killProcessByPortUnix(port string) {
	cmd := exec.Command("lsof", "-t", "-i", fmt.Sprintf(":%s", port))
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error finding process:%v\n", err)
	}
	pidStr := strings.TrimSpace(string(output))
	if pidStr == "" {
		fmt.Println("No process found using the port.", port)
	}
	fmt.Printf("Process found on port %s with PID:%s\n", port, pidStr)

	for _, pid := range strings.Split(pidStr, "\n") {
		// 杀死进程
		killCmd := exec.Command("kill", pid)
		if err := killCmd.Run(); err != nil {
			fmt.Printf("Error killing process:%v\n", err)
			continue
		}
	}
	fmt.Println("Process killed successfully.")
}

func UserSelect() {
	// 定义选项
	options := []string{"Option 1", "Option 2", "Option 3"}

	// 创建选择提示
	prompt := promptui.Select{
		Label: "Please choose an option",
		Items: options,
	}

	// 获取用户选择
	_, result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	// 输出选择结果
	fmt.Printf("You chose: %s\n", result)
}

func KillProcessByPort(port string) {
	switch runtime.GOOS {
	case "windows":
		fmt.Println("windows的杀进程")
	case "darwin", "linux":
		fmt.Println("linux杀进程")
		UserSelect()
	default:
		fmt.Println("Unsupported operating system")
	}
}

func KillProcess(pid int) {
	process, err := os.FindProcess(pid)
	fmt.Println(process)
	if err != nil {
		fmt.Println("FindProcess error:", err)
	}
	process.Kill()
	state, err := process.Wait()
	if err != nil {
		fmt.Println("Wait error", err)
	}
	fmt.Println(state)
}

func OsTest() {
	// fmt.Println("os test")

	// 获取信息类
	execPath, _ := os.Executable()
	fmt.Println("获取当前可执行文件的路径：", execPath)

	execDir, _ := os.Getwd()
	fmt.Println("运行命令的所在文件目录：", execDir)

	userHomeDir, _ := os.UserHomeDir()
	fmt.Println("用户的HOME目录：", userHomeDir)

	fmt.Println("用户的id？", os.Getuid(), os.Geteuid())

	fmt.Println("当前程序的pid：", os.Getpid())

	// KillProcess(os.Getpid())
	// killProcessByPortUnix("8000")
	KillProcessByPort("8000")
}
