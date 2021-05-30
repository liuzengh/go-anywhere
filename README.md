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