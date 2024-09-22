package main

import (
	"fmt"
	"mainmodule/cmd"
)

func main() {
	// input答案
	projectName, _ := cmd.Input("请输入项目名", "my-project")

	// 多选select
	stylelintOptions := []string{"eslint", "typescript", "prettier"}
	stylelint, _ := cmd.Check("请选择你喜欢的选项", &stylelintOptions, true)

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
