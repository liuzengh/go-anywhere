
### Level架构

![LevelDBstructure](LevelDBstructure.svg)


当日志文件超过一定大小(默认4MB)时:创建一个全新的memtable和日志文件，并在这里直接进行未来的更新。在后台将进行如下操作：

1. 将前一个memtable的内容写入一个sstable。
2. memtable丢弃。
3. 删除旧的日志文件和旧的memtable。
4. 添加新的sstable到年轻的(level-0)level

### 构建和维护SSTables

排序字符串表(Sorted Strings Table, SSTable)是一个按键排序的键/值字符串对文件。sstable通常是不可变的。在磁盘上，SSTable被保存为一个持久的、有序的、不可变的文件集。它们通过memtable刷新创建，并通过压缩删除。

在磁盘上维护有序结构是可能的（例如B树），但在内存保存则要容易得多。有许多可以使用的众所周知的树形数据结构，例如红黑树或AVL树，使用这些数据结构，可以按任何顺序插入键，并按关键字排序顺序读取记录。而LevelDB使用的是概率数据结构跳表SkipList来维护有序的key-value对。绝大多数操作（读／写）的平均时间复杂度均为 `O(log n)`，有着与平衡树相媲美的操作效率，但是从实现的角度来说简单许多。

现在我们可以使我们的存储引擎工作如下：

1. 写入时，将其添加到内存中的平衡树数据结构（例如，SkipList）。这个内存树有时被称为内存表（memtable）。
2. 当内存表大于某个阈值（通常为几兆字节）时，将其作为SSTable文件写入磁盘。这可以高效地完成，因为树已经维护了按键排序的键值对。新的SSTable文件成为数据库的最新部分。当SSTable被写入磁盘时，写入可以继续到一个新的内存表实例。
3. 为了提供读取请求，首先尝试在内存表中找到关键字，然后在最近的磁盘段中，然后在下一个较旧的段中找到该关键字。
4. 有时会在后台运行合并和压缩过程以组合段文件并丢弃覆盖或删除的值。

这个方案效果很好。它只会遇到一个问题：如果数据库崩溃，则最近的写入（在内存表中，但尚未写入磁盘）将丢失。为了避免这个问题，我们可以在磁盘上保存一个单独的日志，每个写入都会立即被附加到磁盘上，就像在前一节中一样。该日志不是按排序顺序，但这并不重要，因为它的唯一目的是在崩溃后恢复内存表。每当内存表写出到SSTable时，相应的日志都可以被丢弃。



---


log files : *.log 
`memtable`
sorted table: 

测试文件： db_test
