package main

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"sync"
	"testing"
	"time"
)

func typeTest() {
	var a int
	var b int = 1
	var c = 1

	d := 1
	msg := "hello world"
	fmt.Println(a, b, c, d, msg)
}

func typeTest2() {
	var a int8 = 10
	var c1 byte = 'a'
	var b float32 = 12.2
	ok := false

	fmt.Println(a, c1, string(c1), b, ok)
	if string(c1) == "a" {
		fmt.Println("c1 is a")
	}
}

// 反射获取数据类型
func typeofTest() {
	str1 := "Golang"
	fmt.Println("str1::", reflect.TypeOf(str1)) // string
	fmt.Println("str1-len::", len(str1))

	fmt.Println(str1[0])                                  // char类型
	fmt.Println("Kind::", reflect.TypeOf(str1[0]).Kind()) // uint8

	str2 := "Go语言"
	fmt.Println("str2::", len(str2)) // 8 Go为2，语言占6

	runeArr := []rune(str2)
	fmt.Println("runeArr::", runeArr, len(runeArr))
	fmt.Println("runeArr[2]::", runeArr[2], string(runeArr[2]))
	fmt.Println(reflect.TypeOf(runeArr))

	char1 := 'a'
	fmt.Println("char1::", reflect.TypeOf(char1))

	a := 8
	fmt.Println(reflect.TypeOf(a))

	var b int8 = 1
	fmt.Println(reflect.TypeOf(b))

	var c int16 = 1
	fmt.Println(reflect.TypeOf(c))

	var d int32 = 1
	fmt.Println(reflect.TypeOf(d))

	var e int64 = 1
	fmt.Println(reflect.TypeOf(e))
}

// 数组与切片
func arrayAndSlice() {
	var arr [5]int
	var arrDouble [5][4]int
	fmt.Println(arr, arrDouble)

	// 初始化
	var arr2 = [5]int{1, 2, 3, 4, 5}
	arr3 := [5]int{1, 2, 3, 4, 5}
	fmt.Println(arr2, arr3)

	// 按索引修改数组
	for i := 0; i < len(arr3); i++ {
		arr3[i] += 100
	}
	fmt.Println(arr3)

	// make(类型，长度，容量[可选，不指定则容量和长度相同])
	// 容量为预留空间，比如初始化3的长度，但是预留了5的空间？
	slice1 := make([]float32, 0)
	fmt.Println("slice1::", slice1) // []

	slice2 := make([]float32, 3, 5)       // 3是切片初始长度，5是指定了切片初始容量？
	fmt.Println("slice2::", slice2)       // [0 0 0]
	fmt.Println(len(slice2), cap(slice2)) // 3 5

	// 往slice里添加(预留空间范围内)
	slice2 = append(slice2, 1, 2)
	fmt.Println("slice2-append::", slice2) // [0 0 0]
	fmt.Println(len(slice2), cap(slice2))

	// 超出预留空间范围，则为现有长度容量的两倍？
	slice2 = append(slice2, 3)
	fmt.Println(slice2, len(slice2), cap(slice2))
	slice2 = append(slice2, 4)
	fmt.Println(len(slice2), cap(slice2))

	// 切片取值,左闭右开
	sub1 := slice2[2:2] // []
	sub2 := slice2[1:]  // [0 0 1 2 3 4]
	sub3 := slice2[:3]  // [0 0 0]

	fmt.Println(slice2, sub1, sub2, sub3)

	// // 合并切片，可以看到sub2在前后是有变化的，append非纯函数，会对里面的元素有影响
	// // 安全操作，使用copy操作
	// fmt.Println("sub1::", sub1)
	// fmt.Println("sub2::", sub2) // [0 0 1 2 3 4]
	// combined := append(sub1, sub2...)
	// fmt.Println(combined, sub1, sub2)
	// fmt.Println("sub2::", sub2) // [0 0 0 1 2 3]

	sub4 := slice2[2:2] // []
	sub5 := slice2[1:]  // [0 0 1 2 3 4]
	fmt.Println(slice2, sub4, sub5)

	// 使用copy的方式，不会影响原切片和原数组
	tmp4 := make([]float32, len(sub4))
	copy(tmp4, sub4)
	tmp5 := make([]float32, len(sub5))
	copy(tmp5, sub5)
	combined := append(tmp4, tmp5...)
	fmt.Println("combined::", combined)
	fmt.Println(slice2, sub4, sub5)
}

func testMap() {
	m1 := make(map[string]string)
	fmt.Println("map::", m1)
	m2 := map[string]string{
		"NickName": "Jack",
		"Gender":   "boy",
	}
	fmt.Println(m2, m2["NickName"], m2["Gender"])
	m2["NickName"] = "bob"
	fmt.Println(m2, m2["NickName"], m2["Gender"])
}

func add(num int) {
	num += 1
}
func addOrigin(num *int) {
	*num += 1
}

func pointTest() {
	str := "Golang"
	var p *string = &str
	*p = "Hello"
	fmt.Println(str)

	num := 1
	add(num)
	fmt.Println(num) // 1

	num2 := 1
	addOrigin(&num2)
	fmt.Println(num2)
}

func branchTest() {
	age := 18
	if age < 18 {
		fmt.Println("Kid")
	} else {
		fmt.Println("Adult")
	}

	if age := 18; age < 18 {
		fmt.Println("Kid")
	} else {
		fmt.Println("Adult")
	}
}

func switchTest() {
	type Gender int8
	const (
		MALE   Gender = 1
		FEMALE Gender = 2
	)
	gender := MALE

	// 不需手动break，需要继续往下传递使用 fallthrough
	switch gender {
	case FEMALE:
		fmt.Println("female")
	case MALE:
		fmt.Println("male")
		fallthrough
	default:
		fmt.Println("default here")
	}
}

func forTest() {
	sum := 0
	for i := 0; i < 10; i++ {
		if sum > 20 {
			break
		}
		sum += i
	}
	fmt.Println(sum)
}

// 对数组，slice，map的遍历
func loopTest() {
	arr := []int{1, 2, 3, 4, 5}
	for i, num := range arr {
		fmt.Println(i, num)
	}

	fmt.Println(arr[0])

	obj := map[string]string{
		"Name":   "jack",
		"Gender": "male",
	}
	for key, value := range obj {
		fmt.Println(key, value)
	}
	fmt.Println(obj["Name"]) // 不能使用obj.Name
}

func count(a, b int) int {
	return a + b
}

func funcTest() {
	result := count(1, 2)
	fmt.Println(result)
}

// 自定义错误(意料之中的错误)
func hello(name string) error {
	if len(name) == 0 {
		return errors.New("error:name is null")
	}
	fmt.Println("Hello", name)
	return nil
}

func errorTest() {
	_, err := os.Open("filename.txt")
	if err != nil {
		fmt.Println("error::", err)
	}

	if err := hello(""); err != nil {
		fmt.Println("执行出错：", err)
	}
}

// 类似try-catch操作
func get(index int) (ret int) {
	// 处理错误，保证即使报错，后续代码也能执行
	// 该函数内出现了错误，则会触发panic，控制权交给defer
	// defer中，使用recover，使程序恢复正常
	defer func() {
		r := recover() // 这里将程序恢复
		fmt.Println(r) // 0
		if r != nil {
			fmt.Println("Error::", r)
		}
		ret = -1 // 设置错误时的返回值
		// if r := recover(); r != nil {
		// 	fmt.Println("some error happened::", r)
		// }
	}()

	arr := [3]int{2, 3, 4}
	return arr[index]
}
func tryCatchTest() {
	fmt.Println("get(1)::", get(1))
	fmt.Println("get(5)::", get(5))
	fmt.Println("continue")
}

// 结构体=》类似js中的class
type Student struct {
	name string
	age  int
}

// 在Student上挂载方法
func (student *Student) greet() string {
	return fmt.Sprintf("hello, i am %s", student.name)
}

func structTest() {
	// 实例化方式一
	stu := &Student{
		name: "Jack",
	}
	msg := stu.greet()
	fmt.Println("msg::", msg)

	// 实例化方式二
	stu2 := new(Student)
	stu2.name = "Mick"
	fmt.Println(stu2.greet())
}

// 接口
type IAnima interface {
	eat() string
}
type Dog struct {
	name string
	age  int
}

func (dog *Dog) eat() string {
	return fmt.Sprintf("%s is eating", dog.name)
}

// struct定义基本的属性，可以将多个方法使用func (instance *struct) 来关联一个struct
// 在实例化时，var dog IAnima，使用 IAnima 来限制和检验Dog实现了 IAnima 的接口
func interfaceTest() {
	var dog IAnima = &Dog{
		name: "white",
		age:  10,
	}
	fmt.Println(dog.eat())

	// 空接口的使用
	m := make(map[string]interface{})
	m["name"] = "Tom"
	m["age"] = 18
	m["scores"] = [3]int{98, 99, 199}
	fmt.Println(m)
}

// 并发编程之-sync
// 各协程间不需要通信
var wg sync.WaitGroup

func download(url string) {
	fmt.Println("start to download", url)
	time.Sleep(time.Second)
	// wg减去一个计数
	wg.Done()
}

func syncTest() {
	for i := 0; i < 3; i++ {
		// 为wg添加一个技术
		wg.Add(1)
		// 启动新的协程并发执行
		go download("a.com/" + string(i+'0'))
	}
	// 等到所有都执行结束
	wg.Wait()
	fmt.Println("all done!")
}

// 并发编程之-channel
// 可以在协程之间传递消息，阻塞等待并发协程返回信息
var ch = make(chan string, 10) // 创建大小为10的缓冲信道
func download2(url string) {
	fmt.Println("start to download", url)
	time.Sleep(time.Second)
	ch <- url + "-下载的结果" // 将 url 发送给信道
}

func channelTest() {
	for i := 0; i < 3; i++ {
		go download2("a.com/" + string(i+'0'))
	}
	for i := 0; i < 3; i++ {
		// 从信道中拿出msg，即可以是执行的结果
		msg := <-ch
		fmt.Println("finished", msg)
	}
	fmt.Println("All Done")
}

// 单元测试
func plus(num1, num2 int) int {
	return num1 + num2
}
func testPlus(t *testing.T) {
	if ans := plus(1, 2); ans != 3 {
		t.Error("add(1,2) should be equal to 3")
	}
}

func main() {
	typeTest()
	typeTest2()
	typeofTest()
	arrayAndSlice()
	testMap()
	pointTest()
	branchTest()
	switchTest()
	forTest()
	loopTest()
	funcTest()
	errorTest()
	tryCatchTest()
	structTest()
	interfaceTest()
	syncTest()
	channelTest()
}
