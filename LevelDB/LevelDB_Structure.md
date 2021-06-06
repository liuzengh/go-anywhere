
- SSTable: Sorted Strings Table (borrowed from google) is a file of key/value string pairs, sorted by keys. SStables are generally immutable by design.

Sorted Strings Table (SSTable) is a file format used by Apache Cassandra, Scylla, and other NoSQL databases when memtables are flushed to durable storage from memory. Scylla has always tried to maintain compatibility with Apache Cassandra, and file formats are no exception. SSTable is saved as a persistent, ordered, immutable set of files on disk. They are created by a memtable flush and are deleted by a compaction.

log files : *.log 
`memtable`
sorted table: 

测试文件： db_test