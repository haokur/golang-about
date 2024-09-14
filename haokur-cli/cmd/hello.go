package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var helloCmd = &cobra.Command{
	Use:   "hello",
	Short: "Say hello",
	Run: func(cmd *cobra.Command, args []string) {
		// out, err := exec.Command("cd ~/haokur/code/Github/koa-server && npm start").Output()
		// if err != nil {
		// 	fmt.Printf("Error: %v\n", err)
		// 	return
		// }
		// // 打印命令输出
		// fmt.Printf("Output:\n%s\n", out)

		dir := "/Users/haokur/code/Github/koa-server"
		// 进入指定的目录
		if err := os.Chdir(dir); err != nil {
			fmt.Printf("Failed to change directory to %s: %v\n", dir, err)
			return
		}

		// 确定当前工作目录
		wd, err := os.Getwd()
		if err != nil {
			fmt.Printf("Failed to get working directory: %v\n", err)
			return
		}
		fmt.Printf("Working in directory: %s\n", wd)

		// 创建系统命令
		cmdToRun := exec.Command("npm", "start")

		// 将命令的标准输出和标准错误连接到当前进程
		cmdToRun.Stdout = os.Stdout
		cmdToRun.Stderr = os.Stderr

		if err := cmdToRun.Run(); err != nil {
			fmt.Printf("Command execution failed: %v\n", err)
		}

		// 启动命令
		if err := cmdToRun.Start(); err != nil {
			fmt.Printf("Failed to start command: %v\n", err)
			return
		}

		// 等待命令执行完成（或被中断）
		if err := cmdToRun.Wait(); err != nil {
			fmt.Printf("Command finished with error: %v\n", err)
		}

		// 执行系统命令
		// out, err := exec.Command(args[0], args[1:]...).CombinedOutput()
		// out, err := exec.Command("npm", "start").CombinedOutput()
		// out, err := exec.Command("ls").CombinedOutput()
		// if err != nil {
		// 	fmt.Printf("Command execution failed: %v\n", err)
		// 	return
		// }

		// 打印命令输出
		// fmt.Printf("Output:\n%s\n", out)
	},
}

func init() {
	rootCmd.AddCommand(helloCmd)
}
