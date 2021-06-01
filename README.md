# go-anywhere
code, docs and ideas about go


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

### Math

lc7. Reverse Integer

数值类型的范围： `math` 包含了数值的范围， 如： `math.MaxInt32 = 2**31 - 1, math.MinInt32 = -2**31`

lc9. Palindrome Number

字符串转换：strconv包中实现了字符串和基本类型的转换，如最为常见的数值转换：

```go
i, err := strconv.Atoi("-42")
s := strconv.Itoa(-42)
```

### Sort

lc451. Sort Characters By Frequency

Go中的接口提供了一种方法来指定对象的行为:如果有什么东西可以做到这一点，那么它可以在这里使用。一个类型可以实现多个接口。例如，如果一个集合实现了 `sort.Interface` 接口，那么它可以通过package sort中的方法进行排序。排序接口，包含`Len()`， `Less(i, j int) bool`和 `Swap(i, j int)`