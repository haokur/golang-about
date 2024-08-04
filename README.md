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