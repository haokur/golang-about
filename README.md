### 使用本地模块，例如根目录的 `main.go` 使用 `helper` 下的方法

1. 根目录下初始化主模块

```sh
# mainmodule 自定义
go mod init mainmodule
```

2. 创建子模块及初始化

```sh
go mod init mainmodule/helper
```

3. 子模块的包名需要是定义的 helper

- helper/calc.go

```golang
package helper

// Reduce 相减
func Reduce(num1, num2 int) int {
	return num1 - num2
}
```

- helper/go.mod

```text
module mainmodule/helper

go 1.19
```

4. 回到根目录，mainmodule 模块使用 helper 模块

- main.go 的代码

```golang
package main

import (
	"fmt"
	"mainmodule/helper"
)

func main() {
	result := helper.Reduce(2, 1)
	fmt.Println("result::", result)
}
```

> 注意这里的 package 还是 main 而不是 mainmodule

- 根模块定义加载子模块

```sh
go mod edit -replace mainmodule/helper=./helper
go mod tidy
```

执行完上面命令，根目录下的 go.mod 文件将变更成下面这样

```text
module mainmodule

go 1.19

replace mainmodule/helper => ./helper

require mainmodule/helper v0.0.0-00010101000000-000000000000
```

注意，先在 main.go 中使用 helper，然后再执行 `go mod tidy` 命令

### 使用 nodemon 和 air 实现代码更新自动刷新重启应用

#### 一、nodemon

1. 安装 nodemon

```sh
npm install nodemon --save-dev
```

2. 新建 nodemon 配置

- basic-nodemon.json

```json
{
  "watch": ["basic/*.go"],
  "quiet": true,
  "ext": "go",
  "exec": "sh -c 'cd basic && go run main.go'"
}
```

3. package.json 中添加命令

```json
{
  "scripts": {
    "start": "nodemon --config basic-nodemon.json"
  },
  "devDependencies": {
    "nodemon": "^3.1.4"
  }
}
```

4. 运行命令

```sh
npm start
```

5. 更改 basic/main.go 的代码，自动运行 go run main.go

#### 二、air

> 在开发 gin 应用时，使用 nodemon 重新运行，因为服务型应用一直挂着占着端口，重新启动时需要先结束之前启动的进程，略不适配。air 和 gin 可以很好地配合使用

1. 全局安装 air

```sh
go install github.com/cosmtrek/air@latest
```

或者直接安装到项目目录下使用：

```sh
curl -sSfL https://raw.githubusercontent.com/air-verse/air/master/install.sh | sh -s
```

2. 在 gin 服务目录下，初始化配置文件

```sh
cd gin-server && air init
```

3. 在 package.json 中 `scripts` 添加命令

```json
{
  "scripts": {
    "gin": "cd gin-serve && air"
  }
}
```

4. 项目根目录下运行

```sh
npm run gin
```

### 使用 go 监听文件夹变化，变化后执行对应命令

```sh
# 打包watcher，根目录下运行
npm run build:watcher

# 将watcher复制到$PATH路径目录下面
# 打开命令行执行
watcher --config <绝对路径的配置文件>

# 或让其自动生成
watcher
```

自动生成的配置文件如下：

```json
{
  "watchers": [
    {
      "include": ["/Users/xxxxx/golang-about/watcher"],
      "exclude": [
        "/Users/xxxxx/golang-about/watcher/node_modules",
        "/Users/xxxxx/golang-about/watcher/dist",
        "/Users/xxxxx/golang-about/watcher/.git"
      ],
      "cmds": ["npm start"]
    }
  ]
}
```

配置文件说明：

- watchers，监听对象，可同时存在多个
  - include，要监听的目录或文件
  - exclude，忽略监听的目录或文件
  - cmds，要运行的命令，可多个
  - extensions，字符串数组，只监听的后缀名如：[".js",".ts"]

### 基于 bubbletea 开发命令行交互

- 封装代码在 /cmd 目录下

```go
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
```
