# go-anywhere
code, docs and ideas about go

## Concurrency

goroutine 是由 go runtime 管理的轻量级线程。goroutines 都运行在相同的地址空间，因此获取共享内存的时候必须进行同步。`sync` 包中提供了一些有用的同步原语如 `Mutex` 。

### Channel

`channel` 是有类型的管道，可以使用 `<-` 操作符在 `channel` 上发送和接收值，`<-` 箭头上的方向指明了数据流动的方向。例如下面的代码把对一个切片中数字求和的任务，平均分配给两个goroutine。

```go
func sum(s []int, c chan int) {
	sum := 0
	for _, v := range s {
		sum += v
	}
	c <- sum // send sum to c
}

func main() {
	s := []int{7, 2, 8, -9, 4, 0}

	c := make(chan int)
	go sum(s[:len(s)/2], c)
	go sum(s[len(s)/2:], c)
	x, y := <-c, <-c // receive from c

	fmt.Println(x, y, x+y)
}
```

对于无缓冲的 `channel` , 发送方和接受方在对方没有准备好之前都处于阻塞状态，类似于java中的同步队列 `SynchronousQueue` 。所以下面的两段代码都能打印出 "hello, world" 这条消息。

```go
var done = make(chan bool) 
var msg string
func aGoroutine() {
    msg = "hello, world"
    done <- true
}
func main() {
    go aGoroutine()
    <-done
    println(msg)
}
```
```go
var done = make(chan bool)
var msg string

func aGoroutine() {
    msg = "hello, world"
    <-done
}
func main() {
    go aGoroutine()
    done <- true
    println(msg)
}
```

 对于带有缓存的 `channel` , 当缓存满时发送方将阻塞，当缓存空时接收方将阻塞，类似于java中的有界阻塞队列`ArrayBlockingQueue`和` LinkedBlockingQueue`。例如下面这段代码发送者将一直阻塞，从而导致死锁。

 ```go
 func main() {
	ch := make(chan int, 2)
	ch <- 1
	ch <- 2
	ch <- 3
	fmt.Println(<-ch)
	fmt.Println(<-ch)
}
```

`channle` 使得goroutines可以不显式锁或条件变量的情况下进行同步。可以使用 `range` 循环来接受channel上的值，直到发送方主动 close channel时才退出循环。

```go
func fibonacci(n int, c chan int) {
	x, y := 0, 1
	for i := 0; i < n; i++ {
		c <- x
		x, y = y, x+y
	}
	close(c)
}

func main() {
	c := make(chan int, 10)
	go fibonacci(cap(c), c)
	for i := range c {
		fmt.Println(i)
	}
}
```

注意: 只有发送方才应该关闭通道，而不是接收方。给一个已经关闭的 channle发送值会引起panic。另外，channel于file不同，通常并不需要关闭channel。只有当接收方必须被告知没有更多的值时才需要关闭，例如终止一个range循环。


#### Select

`select` 语句和 c语言里面 `switch-case` 用法很像， 可以让一个goroutine对等待在多个 channel 上。`select` 将一直阻塞, 直到其中的一个 `case` 可以运行才去执行对应的 `case`。 如果有多个 `case` 都准备好了，它会随机选择一个 `case` 执行。如某个case可以发送退出循环的消息或计时消息或超时消息等。

```go
func fibonacci(c, quit chan int) {
	x, y := 0, 1
	for {
		select {
		case c <- x:
			x, y = y, x+y
		case <-quit:
			fmt.Println("quit")
			return
		}
	}
}

func main() {
	c := make(chan int)
	quit := make(chan int)
	go func() {
		for i := 0; i < 10; i++ {
			fmt.Println(<-c)
		}
		quit <- 0
	}()
	fibonacci(c, quit)
}
```

```go
func main() {
	tick := time.Tick(100 * time.Millisecond)
	boom := time.After(500 * time.Millisecond)
	for {
		select {
		case <-tick:
			fmt.Println("tick.")
		case <-boom:
			fmt.Println("BOOM!")
			return
		default:
			fmt.Println("    .")
			time.Sleep(50 * time.Millisecond)
		}
	}
}
```

应用1：判断等价二叉树
```go
package main

import (
	"fmt"
	"golang.org/x/tour/tree"
//	"reflect"
)

// Walk walks the tree t sending all values
// from the tree to the channel ch.
func Walk(t *tree.Tree, ch chan int) {
	if t == nil {
		return
	}

	Walk(t.Left, ch)
	ch <- t.Value
	Walk(t.Right, ch)
}

// Same determines whether the trees
// t1 and t2 contain the same values.
func Same(t1, t2 *tree.Tree) bool {
	ch1 := make(chan int, 1)
	ch2 := make(chan int, 1)
	go Walk(t1, ch1)
	go Walk(t2, ch2)
	// result1 := []int{}
	// result2 := make([]int, 10)
	for i := 0; i < 10; i++ {
		// result1 = append(result1, <-ch1)
		// result2[i] = <-ch2
		x, y := <-ch1, <-ch2
		if x != y {
			return false
		}
		
	}
	return true
	// return reflect.DeepEqual(result1, result2)
}

func main() {
	/*
		ch := make(chan int)
		go Walk(tree.New(1), ch)
		for i := 0; i < 10; i++ {
			c := <-ch
			fmt.Println(c)
		}
	*/
	/*
		for c := range ch {
			fmt.Println(c)
		}
	*/

	fmt.Println(Same(tree.New(1), tree.New(2)))
	fmt.Println(Same(tree.New(1), tree.New(1)))
}
```
```go
package tree // import "golang.org/x/tour/tree"

import (
	"fmt"
	"math/rand"
)

// A Tree is a binary tree with integer values.
type Tree struct {
	Left  *Tree
	Value int
	Right *Tree
}

// New returns a new, random binary tree holding the values k, 2k, ..., 10k.
func New(k int) *Tree {
	var t *Tree
	for _, v := range rand.Perm(10) {
		t = insert(t, (1+v)*k)
	}
	return t
}

func insert(t *Tree, v int) *Tree {
	if t == nil {
		return &Tree{nil, v, nil}
	}
	if v < t.Value {
		t.Left = insert(t.Left, v)
	} else {
		t.Right = insert(t.Right, v)
	}
	return t
}

func (t *Tree) String() string {
	if t == nil {
		return "()"
	}
	s := ""
	if t.Left != nil {
		s += t.Left.String() + " "
	}
	s += fmt.Sprint(t.Value)
	if t.Right != nil {
		s += " " + t.Right.String()
	}
	return "(" + s + ")"
}
```
应用2： 并发网络爬虫

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

type SafeCache struct {
	mu        sync.Mutex
	container map[string]bool
}

func (safe_cache *SafeCache) Set(key string) {
	safe_cache.mu.Lock()
	safe_cache.container[key] = true
	safe_cache.mu.Unlock()
}

func (safe_cache *SafeCache) Get(key string) bool {
	safe_cache.mu.Lock()
	defer safe_cache.mu.Unlock()
	return safe_cache.container[key]
}

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher, cache SafeCache) {
	// TODO: Fetch URLs in parallel.
	// TODO: Don't fetch the same URL twice.
	// This implementation doesn't do either:
	if depth <= 0 {
		return
	}
	/*
	cache.mu.Lock()
	_, exist := cache.container[url]
	cache.mu.Unlock()
	*/
	if cache.Get(url) {
		fmt.Println("in cache")
		return
	}
	

	body, urls, err := fetcher.Fetch(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	/*
	cache.mu.Lock()
	cache.container[url] = true
	cache.mu.Unlock()
	*/
	cache.Set(url)

	fmt.Printf("found: %s %q\n", url, body)
	for _, u := range urls {
		go Crawl(u, depth-1, fetcher, cache)
	}
	time.Sleep(time.Second)
	return
}

func main() {
	cache := SafeCache{container: make(map[string]bool)}
	Crawl("https://golang.org/", 4, fetcher, cache)
}

// fakeFetcher is Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

func (f fakeFetcher) Fetch(url string) (string, []string, error) {
	if res, ok := f[url]; ok {
		return res.body, res.urls, nil
	}
	return "", nil, fmt.Errorf("not found: %s", url)
}

// fetcher is a populated fakeFetcher.
var fetcher = fakeFetcher{
	"https://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"https://golang.org/pkg/",
			"https://golang.org/cmd/",
		},
	},
	"https://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"https://golang.org/",
			"https://golang.org/cmd/",
			"https://golang.org/pkg/fmt/",
			"https://golang.org/pkg/os/",
		},
	},
	"https://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
	"https://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
}
```
## RPC

### 消息传递

远程过程调用(Remote Procedure Call, RPC)是一种请求响应协议，是分布式系统的关键组件之一。下图展示了RPC的消息传递过程。
```
Client             Server
	request--->
         <---response
```
- 请求过程：客户端发起RPC，向已知的远程服务器发送请求消息，远程服务器使用客户端提供的参数执行指定的程序。
- 响应过程：远程服务器向客户机发送响应，应用程序继续其进程。当服务器正在处理调用时，客户端被阻塞(它等待直到服务器完成处理后才恢复执行)，除非客户端向服务器发送一个异步请求，例如XMLHttpRequest。

RPC调用流程：

1. 客户端调用客户端stub（client stub）。这个调用是本地调用，并将调用参数push到栈（stack）中。
2. 客户端stub（client stub）将这些参数包装，并通过系统调用发送到服务端机器。打包的过程叫 序列化(marshalling)。（常见方式：XML、JSON、二进制编码）
3. 客户端本地操作系统发送信息至服务器。（可通过自定义TCP协议或HTTP传输）
4. 服务器系统将信息传送至服务端stub（server stub）。
5. 服务端stub（server stub）解析信息。该过程叫 反序列化(unmarshalling)。
6. 服务端stub（server stub）调用程序，并通过类似的方式返回给客户端。

RPC也可视为是进程间通信(Inter-process Communication, IPC)的一种形式，因为不同的进程有不同的地址空间。

- 在同一台主机：不同的进程有不同的虚拟地址空间，即使物理地址空间是相同的
- 在不同的主机上，物理地址空间是不同的



### go中的net/rpc 包

go中的rpc包提供对象的导出方法的远程访问（通过网络或其他I/O连接）。服务器注册一个对象，使其作为一个具有对象类型名称的服务对远端的客户端可见。注册后，客户端可以远程访问导出的对象方法。服务器可以注册不同类型的多个对象(服务)，但不允许注册相同类型的多个对象。

```
client app       handler fns
stub fns         dispatcher
RPC lib          RPC lib
net ------------ net
```

一般而言对象中的方法，需要满足如下的3条标准才可用于远程访问：

1. 方法和方法的类型是  exported.
2. 方法有两个参数, 参数的类型都是 exported (or builtin); 方法的第二个参数是指针类型
3. 方法的返回值类型是 error

通常来说，需要可以用于远程调用的方法具有如下形式：

```go
func (t *T) MethodName(argType T1, replyType *T2) error
```

- 第一个参数代表调用者提供的参数
- 第二个参数代表返回给调用者的参数
- 方法的返回值，如果非nil，将被作为字符串回传，在客户端看来就和errors.New创建的一样。如果返回了错误，返回给调用者的参数将不会被发送给

服务端：单个连接调用 `ServerConn`, 一般而言调用 `Accept` 创建网络监听器；对于 Http监听器，调用 `HandleHttp` 和 `http.Serve`。

客户端：想要使用(服务端提供的)服务的客户端需要先和服务端建立连接，并且在该连接上产生一个 `NewClient`。 `Dial` ( `DialHttp` )函数在一个原始的网络连接（一个Http连接）上执行了以上必要两个步骤，该方法返回的 `Client` 对象包含两个成员方法：`Call` 和 `Go`。

- `Call`方法: 同步调用，等待远程调用完成
- `Go`方法：异步调用，使用Call结构体中的成员变量Done（类型为channel）发出完成信号。

```go
func (client *Client) Go(serviceMethod string, args interface{}, reply interface{}, done chan *Call) *Call {
	call := new(Call)
	call.ServiceMethod = serviceMethod
	call.Args = args
	call.Reply = reply
	if done == nil {
		done = make(chan *Call, 10) // buffered.
	} else {
		if cap(done) == 0 {
			log.Panic("rpc: done channel is unbuffered")
		}
	}
	call.Done = done
	client.send(call)
	return call
}
func (client *Client) Call(serviceMethod string, args interface{}, reply interface{}) error {
	call := <-client.Go(serviceMethod, args, reply, make(chan *Call, 1)).Done
	return call.Error
}
```

可以看到上述两个方法都参数都需要：指定要调用的服务和方法、一个包含参数的指针和一个接收结果参数的指针。而同步调用方法`Call`实际上是在异步调用方法`Go`上做了一层包装，接受到Call结构体中的成员变量Done（类型为channel）发出完成信号才返回。

编码解码器：除非设置了显式的编解码器，否则使用包 `encoding/gob` 来传输数据。

实现上的一些细节： 

1. 能将哪些将数据格式化为数据包（序列化）？
Go的RPC库可以传递strings、arrays、objects、maps和&c。Go通过复制指向的数据传递指针。但是不能传递 channels 和 functions
2. 客户端如何知道要与哪个服务器计算机进行通信（绑定）?
对于Go的RPC，服务器名/端口是Dial的一个参数。对于大系统有某种名称或配置服务器。

### RPC中的故障处理

RPC是如何处理丢包、网络中断、服务器慢、服务器崩溃等故障(failures)的？对于客户端RPC库来说，当故障发生时， 客户端永远不会看到来自服务器的响应。客户端不知道服务器是否看到了请求! 可能服务器从未看到请求；可能服务器执行了调用程序，在发送应答之前崩溃；可能服务器执行了调用程序，发送应答的时候网络挂掉了。

最简单的故障处理方案-“尽最大努力”。 Call()等待响应一段时间；如果没有到达，重新发送请求；多做几次；然后放弃并返回一个错误。该方案对只读操作, 如果重复则什么也不做的操作，例如DB检查记录是否已经被插入，是合适的。

更好的RPC方案-“最多一次”： 服务器RPC代码检测重复的请求， 返回先前的回复，而不是重新运行处理程序。客户端包括唯一的ID (XID)，每个请求使用相同的XID进行重新发送来检测重复请求。如果最多一次服务器崩溃并重新启动怎么办?如果内存中有最多一次的重复信息，服务器会忘记并在重新启动后接受重复的请求，也许它应该把重复的信息写到磁盘上，也许replica server也应该复制重复的信息。

其他方案-"正好一次": 无限重试加上重复检测加上容错服务

go中RPC是采用的是一种简单的“至多一次”策略：

1. 打开TCP连接
2. 写请求到TCP连接RPC， 从来没有重新去请求，所以服务器不会看到重复的请求。
3. RPC代码返回一个错误，如果没有得到一个答复。
   - 也许是TCP超时
   - 也许是服务器没看到请求
   - 也许是服务器处理了请求但是在回复返回前，服务器或网络崩溃了



## Leetcode

### Array
lc1. Two Sum  

访问权限： golang中根据首字母的大小写来确定可以访问的权限。无论是方法名、常量、变量名还是结构体的名称，如果首字母大写，则可以被其他的包访问；如果首字母小写，则只能在本包中使用。

### List
lc2. Add Two Numbers

内存分配：作为一种限制情况，如果一个 composite literal 根本不包含字段，它将为该类型创建一个零值。表达式 `new(File)` 和 `&File{}`是等价的。

### String

lc3. Longest Substring Without Repeating Characters

函数重载： go语言math包里面定义了 `min/max` 函数，但只有 `float64` 类型的，为了保持语言的简洁性性， `golang` 是不支持重载，所以并没有整数类型的 `min/max` 。

字符串： 在Golang中 string 底层是用byte字节数组存储的，并且是不可以修改的。Go语言中byte和rune实质上就是uint8和int32类型。byte用来强调数据是raw data，而不是数字；而rune用来表示Unicode的code point, 是含有语义的字符。

>在 Go 语言中，一个string类型的值既可以被拆分为一个包含多个字符的序列，也可以被拆分为一个包含多个字节的序列。前者可以由一个以rune为元素类型的切片来表示，而后者则可以由一个以byte为元素类型的切片代表。 rune是 Go 语言特有的一个基本数据类型，它的一个值就代表一个字符，即：一个 Unicode 字符。比如，’G’、’o’、’爱’、’好’、’者’代表的就都是一个 Unicode 字符。

lc6. ZigZag Conversion

make 一个 slice：切片由指向数据的pointer, 切片长度以及容量三个属性来表示。
```go
// 创建一个长度为10，容量为100的切片，切片中的元素被初始化int类型的默认值0
make([]int, 10, 100)
```

byte类型的默认值：byte是数值类型，等同于uint8, 其zero value是 0

二维切片：切片是变长的，所以每个维度上的切片长度可以不同，这一点和 `vector<vector<T>>` 很像。

lc8. String to Integer (atoi)

有限状态机: 开始状态为1，终止状态为2，下图中没有表示出来的是：在状态4的时候，当超过数值取值范围时，会截断，然后转向终止状态。

```
┌──┐
│ ┌┴──┐        '+'or'-' ┌───┐
│ │ 1 ├───────┬────────►| 3 ├───┐
│ └─▲─┤       │         └┬──┘   |
└───┘ │       │digit     │digit │
 ' '  │       │        ┌─▼─┐    │
      │       └────────► 4 ├─┐  │
      │other           └┬─▲┘ │  │
      │                 | └──┘  |
    ┌─▼─┐               | digit | 
    │ 2 ◄───────────────┘       │
    └─▲─┘ other                 │
      └─────────────────────────┘
        other
```

lc14. Longest Common Prefix

### Binary Search

lc4. Median of Two Sorted Arrays

内置 `append` 函数: 该函数的参数是可变的，可以将0个或多个值 `x` 附加到类型为 `S` 的 `s` (必须是切片类型)后面，并返回结果切片(也是类型为 `S` )。

```go
append(s S, x ...T) S  // T is the element type of S
s0 := []int{0, 0}
s1 := append(s0, 2)                // append a single element     s1 == []int{0, 0, 2}
s2 := append(s1, 3, 5, 7)          // append multiple elements    s2 == []int{0, 0, 2, 3, 5, 7}
s3 := append(s2, s0...)            // append a slice              s3 == []int{0, 0, 2, 3, 5, 7, 0, 0}
s4 := append(s3[3:6], s3[2:]...)   // append overlapping slice    s4 == []int{3, 5, 7, 2, 3, 5, 7, 0, 0}
```
思路：找到一个合适划分 `(splitA, splitB)` , 使得划分后得到的集合 `left_part` 和 `right_part` 满足： `max(left_part) <= min(right_part), |left_part| - |right_part| <= 1`
```
          left_part          |         right_part
    A[0], A[1], ..., A[i-1]  |  A[i], A[i+1], ..., A[m-1]
    B[0], B[1], ..., B[j-1]  |  B[j], B[j+1], ..., B[n-1]
```

### Dynamic Programming

lc5. Longest Palindromic Substring

map的key类型：只有可以比较的类型才能作为map的key, 

> As mentioned earlier, map keys may be of any type that is comparable. The language spec defines this precisely, but in short, comparable types are boolean, numeric, string, pointer, channel, and interface types, and structs or arrays that contain only those types. Notably absent from the list are slices, maps, and functions; these types cannot be compared using ==, and may not be used as map keys.


二维切片：golang中的array和slice都是一维的，如果要创建二维slice需要预先定义或动态分配

```go
type LinesOfText [][]byte     // A slice of byte slices.
text := LinesOfText{
	[]byte("Now is the time"),
	[]byte("for all good gophers"),
	[]byte("to bring some fun to the party."),
}
// Allocate the top-level slice.
picture := make([][]uint8, YSize) // One row per unit of y.
// Loop over the rows, allocating the slice for each row.
for i := range picture {
	picture[i] = make([]uint8, XSize)
}
```

lc10. Regular Expression Matching

以下用python写得伪代码很好的应用了带备忘录的动态规划方法。
```python
def isMatch(self, s: str, p: str):
    text, pattern = s, p
    memo = {}
    def dp(i, j):
        if (i, j) not in memo:
            if j == len(pattern):
                ans = i == len(text)
            else:
                first_match = i < len(text) and pattern[j] in {text[i], '.'}
                if j+1 < len(pattern) and pattern[j+1] == '*':
                    ans = dp(i, j+2) or first_match and dp(i+1, j)
                else:
                    ans = first_match and dp(i+1, j+1)
            memo[i, j] = ans
        return memo[i, j]
    return dp(0, 0)
```

go中闭包中使用递归：可以用先声明再定义的方式
```go
var recur func()
recur = func(){
    recur()
}
```
### Math

lc7. Reverse Integer

数值类型的范围： `math` 包含了数值的范围， 如： `math.MaxInt32 = 2**31 - 1, math.MinInt32 = -2**31`

lc9. Palindrome Number

字符串转换：strconv包中实现了字符串和基本类型的转换，如最为常见的数值转换：

```go
i, err := strconv.Atoi("-42")
s := strconv.Itoa(-42)
```
lc12. Integer to Roman

预先定义符号表

```go
vals := [...]int {1000, 900, 500, 400, 100, 90, 50, 40, 10, 9, 5, 4, 1}
syms := [...]string{"M", "CM", "D", "CD", "C", "XC", "L", "XL", "X", "IX", "V", "IV", "I"}
```

lc13. Roman to Integer

### Sort

lc451. Sort Characters By Frequency

Go中的接口提供了一种方法来指定对象的行为:如果有什么东西可以做到这一点，那么它可以在这里使用。一个类型可以实现多个接口。例如，如果一个集合实现了 `sort.Interface` 接口，那么它可以通过package sort中的方法进行排序。排序接口，包含`Len()`， `Less(i, j int) bool`和 `Swap(i, j int)`

### Two Pointers

lc11. Container With Most Water

`min(A[left], A[right]) * (right - left)`

1. 在初始时，左右指针分别指向数组的左右两端
2. 求出当前双指针对应的容器的容量
3. 对应数字较小的那个指针以后不可能作为容器的边界了，将其丢弃，并移动对应的指针。

### 单调栈

lc42. Trapping Rain Water

给定 n 个非负整数表示每个宽度为 1 的柱子的高度图，计算按此排列的柱子，下雨之后能接多少雨水。

维护一个单调栈，单调栈存储的是下标，满足从栈底到栈顶的下标对应的数组 `height` 中的元素递减。

从左到右遍历数组，遍历到下标 i 时，如果栈内至少有两个元素，记栈顶元素为 `top`，`top` 的下面一个元素是 `left`，则一定有 `height[left] >= height[top]`。如果`height[left] > height[top]` ，则得到一个可以接雨水的区域，该区域的宽度是 `i−left−1`，高度是 `min(height[left], height[i]) - height[top]`，根据宽度和高度即可计算得到该区域能接的雨水量。

参考资料： https://oi-wiki.org/ds/monotonous-stack/

```go
func trap(height []int) int {
    s, top := []int{}, -1
    var result int
    for index, item := range height {
        if top == -1  || item <= height[s[top]]{
            top++
            s = append(s, index)
        } else {
            for -1 < top  && height[s[top]] < item {
                top--
                if top == -1 {
                    break
                }
                distance := index - s[top] - 1 
                delta := MinInt(item, height[s[top]]) - height[s[top+1]] 
                result += distance * delta 
                
            }
            s = s[0:top+1]
            s = append(s, index)
            top++
        }
    }
    return result
}

func MinInt(x, y int) int{
    if x < y {
        return x
    }
    return y
}
```
