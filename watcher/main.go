package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

/*
1. 不存在配置文件，提示自动生成配置文件-watcher.config.json
2. 根据配置文件，根据include和exclude来确定要监听哪些文件，实现监听
3. 监听变化，执行配置的要执行的方法
4. TODO://提供go常量和方法，供配置，如sleep，ip地址
5. 可指定配置目录
6. 5s内相同的命令仅执行一次
**/

// 当前命令运行的目录
var currentDir string

func getDefaultConfig() string {
	defaultIncludes := currentDir
	defaultExcludes := []string{
		currentDir + "/node_modules",
		currentDir + "/dist",
		currentDir + "/.git",
	}
	// defaultExcludesStr := strings.Join(defaultExcludes, ",")
	defaultExcludesStr := fmt.Sprintf(`"%s"`, strings.Join(defaultExcludes, `","`))
	defaultConfig := fmt.Sprintf(`
	{
		"watchers":[
			{
				"include":[
					"%s"
				],
				"exclude":[
					%s
				],
				"cmds":[
					"npm start"
				]
			}
		]
	}
	`, defaultIncludes, defaultExcludesStr)
	return defaultConfig
}

type Config struct {
	Watchers []Watcher `json:"watchers"`
}

type Watcher struct {
	Include    []string `json:"include"`
	Exclude    []string `json:"exclude"`
	Cmds       []string `json:"cmds"`
	Extensions []string `json:"extensions"`
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || !os.IsNotExist(err)
}

func getCurrentDir() string {
	workingDir, err := os.Getwd()
	if err != nil {
		fmt.Println("获取当前工作目录失败:", err)
		return ""
	}
	return workingDir
	// fmt.Println("命令执行的工作目录:", workingDir)

	// // 获取当前可执行文件的路径
	// exePath, err := os.Executable()
	// if err != nil {
	// 	return ""
	// }
	// exeDir := filepath.Dir(exePath)
	// return exeDir
}

func getConfig(configFileName string) string {
	configFilePath := currentDir + "/" + configFileName
	// 查询当前是否存在配置文件
	isConfigExist := fileExists(configFilePath)
	if !isConfigExist {
		fmt.Println("配置文件不存在，自动生成", configFilePath)
		file, err := os.Create(configFilePath)
		if err != nil {
			fmt.Println("创建配置文件失败", err)
			return ""
		}
		defer file.Close()

		defaultConfigStr := getDefaultConfig()
		_, err = file.WriteString(defaultConfigStr)
		if err != nil {
			fmt.Println("写入默认配置失败：", err)
			return ""
		}
	}

	content, err := os.ReadFile(configFilePath)
	if err != nil {
		fmt.Println("读取配置文件失败", err)
	}
	return string(content)
}

// 判断文件是否在排除列表中
func isExcluded(file string, excludes []string) bool {
	for _, exclude := range excludes {
		if isFromDir(file, exclude) {
			return true
		}
	}
	return false
}

// 检查文件是否匹配特定的后缀
func isValidExtension(file string, extensions []string) bool {
	for _, ext := range extensions {
		if strings.HasSuffix(file, ext) {
			return true
		}
	}
	return false
}

// 检查文件是否在特定目录下
func isFromDir(file string, dir string) bool {
	relPath, err := filepath.Rel(dir, file)
	if err != nil {
		return false
	}
	// 如果相对路径是以 ".." 开头，说明 file 不在 dir 中
	return !strings.HasPrefix(relPath, "..")
}

// Debounce 是防抖函数封装
func Debounce(fn func(), delay time.Duration) func() {
	var timer *time.Timer
	var mu sync.Mutex

	return func() {
		mu.Lock()
		defer mu.Unlock()

		// 如果 timer 已经存在，重置它
		if timer != nil {
			timer.Stop()
		}

		// 新建一个定时器，等到 delay 时间后执行目标函数
		timer = time.AfterFunc(delay, func() {
			mu.Lock()
			defer mu.Unlock()
			fn() // 延迟时间结束后调用目标函数
		})
	}
}

// 执行命令
var lastRunTime time.Time

func runCmds(cmds []string) {
	currentTime := time.Now()
	if currentTime.Sub(lastRunTime) < 5*time.Second {
		// fmt.Println("频繁操作，被忽略")
		return
	}
	lastRunTime = currentTime

	// 暂停1000ms，等待所有文件保存成功，再运行
	// time.Sleep(1000 * time.Millisecond)

	for _, cmd := range cmds {
		parts := strings.Fields(cmd)
		head := parts[0]
		// 剩余的作为参数
		args := parts[1:]

		// 如果是cd到一个目录
		if head == "cd" {
			// 切换到目标目录
			targetDir := args[0]
			if err := os.Chdir(targetDir); err != nil {
				fmt.Printf("切换到目录 %s 失败: %v\n", targetDir, err)
			}
			continue
		}

		fmt.Printf("[执行命令]: %s %s\n", head, strings.Join(args, " "))
		command := exec.Command(head, args...)
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr
		err := command.Run()
		if err != nil {
			fmt.Printf("命令执行失败: %s, 错误: %s\n", cmd, err)
		}
	}
}

// 监听文件变化
func watchFiles(watcher *fsnotify.Watcher, w *Watcher) {
	include := w.Include
	exclude := w.Exclude
	cmds := w.Cmds
	extensions := w.Extensions

	for _, path := range include {
		// 递归添加目录，跳过排除路径
		err := filepath.Walk(path, func(file string, fi os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			// 如果是目录且不在排除列表中，则添加监听
			if fi.IsDir() {
				if !isExcluded(file, exclude) {
					if err := watcher.Add(file); err != nil {
						return err
					}
					fmt.Printf("[监听目录]: %s\n", file)
				} else {
					// fmt.Printf("跳过排除的目录: %s\n", file)
				}
			}
			return nil
		})
		if err != nil {
			log.Fatalf("添加监听失败: %s, 错误: %v", path, err)
		}
	}

	// 处理文件变化事件
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				// 检查事件来自哪个文件夹
				for _, dir := range include {
					if isFromDir(event.Name, dir) {
						// fmt.Printf("文件事件来自: %s, 文件: %s\n", dir, event.Name)

						// 过滤掉 exclude 列表中的文件
						if isExcluded(event.Name, exclude) {
							// fmt.Printf("文件被忽略: %s\n", event.Name)
							continue
						}

						// 如果是新建文件夹，则重启监听
						if event.Op&fsnotify.Create == fsnotify.Create {
							// 使用 os.Stat 来检查是否是目录
							fi, err := os.Stat(event.Name)
							if err == nil && fi.IsDir() {
								// fmt.Printf("新文件夹创建: %s，加入监听\n", event.Name)
								if err := watcher.Add(event.Name); err != nil {
									fmt.Printf("无法添加新文件夹到监听: %s\n", err)
								}
								continue
							}
						}

						// 如果限定了后缀，不符合的后缀直接忽略
						if len(extensions) != 0 {
							if !isValidExtension(event.Name, extensions) {
								// fmt.Printf("文件后缀不符合: %s\n", event.Name)
								continue
							}
						}

						// 处理文件变化
						if event.Op&fsnotify.Create == fsnotify.Create {
							fmt.Printf("[文件创建]: %s\n", event.Name)
							runCmds(cmds)
						} else if event.Op&fsnotify.Write == fsnotify.Write {
							fmt.Printf("[文件修改]: %s\n", event.Name)
							runCmds(cmds)
						} else if event.Op&fsnotify.Remove == fsnotify.Remove {
							fmt.Printf("[文件删除]: %s\n", event.Name)
							runCmds(cmds)
						} else if event.Op&fsnotify.Rename == fsnotify.Rename {
							fmt.Printf("[文件重命名]: %s\n", event.Name)
							runCmds(cmds)
						} else {
							runCmds(cmds)
						}
					}
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				fmt.Println("监听错误:", err)
			}
		}
	}()

	// 初始化启动执行
	// 先进入当前运行的目录
	targetDir := currentDir
	fmt.Println("进入目录：", targetDir)
	if err := os.Chdir(targetDir); err != nil {
		fmt.Printf("切换到目录 %s 失败: %v\n", targetDir, err)
	}
	runCmds(cmds)
}

func main() {
	// 定义命令行参数 --config
	configFile := flag.String("config", "watcher.config.json", "配置文件的路径")
	flag.Parse() // 解析命令行参数

	// 获取当前命令所在目录
	currentDir = getCurrentDir()
	// 获取对应的配置文件
	jsonStr := getConfig(*configFile)

	// 解析配置
	var config Config
	err := json.Unmarshal([]byte(jsonStr), &config)
	if err != nil {
		fmt.Println("解析 JSON 失败:", err)
		return
	}
	watchers := config.Watchers

	// 创建 fsnotify 监听器
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println("创建监听器失败:", err)
		return
	}
	defer watcher.Close()

	// 处理每个 watcher 的监听
	for _, w := range watchers {
		watchFiles(watcher, &w)
	}

	// 阻止主协程退出
	done := make(chan bool)
	<-done
}
