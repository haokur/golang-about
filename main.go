package main

import (
	"fmt"
	"mainmodule/cmd"
	"mainmodule/tools"
)

func testCmd() {
	// input答案
	projectName, _ := cmd.Input("请输入项目名", "my-project")

	// 多选select
	stylelintOptions := []string{"eslint", "typescript", "prettier"}
	stylelint, _, _ := cmd.Check("请选择你喜欢的选项", &stylelintOptions, true)

	// 单选radio
	webFrameOptions := []string{"Vue3", "React", "Angular"}
	webFrame, _ := cmd.Radio("请选择你的框架", &webFrameOptions)

	// yes or no => confirm
	installAuto, _ := cmd.Confirm("是否自动执行npm install安装依赖", true)

	fmt.Println(map[string]interface{}{
		"projectName": projectName,
		"lintRules":   stylelint,
		"webFrame":    webFrame,
		"installAuto": installAuto,
	})
}

// 测试杀进程
func testKill() {
	// tools.KillProcess(&[]string{"5173", "node"})
	tools.KillProcess(&[]string{"5173", "obsidian"})
	// tools.KillProcess(&[]string{"5173"})
	// port := "5173"
	// // // lsof -c nginx | awk '{print $2, $1}'
	// // // cmd := exec.Command("lsof", "-i", fmt.Sprintf(":%s", port))
	// // cmd := exec.Command("lsof", "-i", fmt.Sprintf(":%s", port), "|", "awk", "'{print $2, $1}'")
	// // output, err := cmd.Output()
	// // if err != nil {
	// // 	fmt.Println(err)
	// // }

	// // Step 1: Execute lsof command and capture its output
	// lsofCmd := exec.Command("lsof", "-i", fmt.Sprintf(":%s", port))

	// // Get the output from lsof
	// lsofOutput, err := lsofCmd.Output()
	// if err != nil {
	// 	log.Fatalf("Error executing lsof command: %v", err)
	// }

	// lsofOutputStr := string(lsofOutput)
	// fmt.Println(lsofOutputStr)
	// lsofOutputRows := strings.Split(lsofOutputStr, "\n")
	// labelRow := lsofOutputRows[0]
	// labelArr := strings.Fields(labelRow)
	// fmt.Println(labelArr)

	// valueRows := lsofOutputRows[1:]

	// lsofProcessList := []map[string]string{}
	// for _, v := range valueRows {
	// 	valueRowArr := strings.Fields(v)
	// 	if len(valueRowArr) != 0 {
	// 		rowMap := map[string]string{}
	// 		for index, rowMapKey := range labelArr {
	// 			rowMap[rowMapKey] = valueRowArr[index]
	// 		}
	// 		lsofProcessList = append(lsofProcessList, rowMap)
	// 	}

	// }

	// fmt.Println(lsofProcessList)
	// for k, v := range lsofProcessList {
	// 	fmt.Println(k, v, v["COMMAND"])
	// }
	// // Step 2: Execute awk command and pass the lsof output to it
	// awkCmd := exec.Command("awk", "{print $2, $1}")

	// // Pass lsof output to awk input (via stdin)
	// awkCmd.Stdin = strings.NewReader(string(lsofOutput))

	// // Capture the output from awk
	// awkOutput, err := awkCmd.Output()
	// if err != nil {
	// 	log.Fatalf("Error executing awk command: %v", err)
	// }

	// outputString := string(awkOutput)
	// outputArray := strings.Split(strings.TrimSpace(outputString), "\n")

	// // Print the result array
	// for i, line := range outputArray {
	// 	fmt.Printf("Line %d: %s\n", i+1, line)
	// }

	// fmt.Println(string(output))
}

func main() {
	testKill()
	// 多选select
	// stylelintOptions := []string{"eslint", "typescript", "prettier"}
	// stylelint, _, _ := cmd.Check("请选择你喜欢的选项", &stylelintOptions, true)
	// fmt.Println(stylelint)
}
