# go-anywhere
code, docs and ideas about go


## 

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