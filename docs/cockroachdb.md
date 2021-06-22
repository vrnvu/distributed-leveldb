# CockrachDB on RockDB

If, on a final exam at a database class, you asked students whether to build a database on a log-structured merge tree (LSM) or a BTree-based storage engine, 90% of your students would probably respond that the decision hinges on your workload. "LSMs are for write-heavy workloads and BTrees are for read-heavy workloads", the conscientious ones would write. If you surveyed most NewSQL (or distributed SQL) databases today, most of them are built on top of an LSM, namely, RocksDB. You might thus conclude that this is because modern applications have shifted to more write-heavy workloads. You would be incorrect.

The main motivation behind RocksDB adoption has nothing to do with its choice of LSM data structure. In our case, the compelling drivers were its rich feature set which turns out to be necessary for a complex product like a distributed database. 

## What is a storage engine?

A storage engine's job is to write things to disk on a single node. For many databases — which may be single node databases themselves — this constitutes a large part of the engineering effort, and so is pretty intricately tied to the engineering effort of the database building itself.

However, in a distributed setting, the distribution, replication, and transaction coordination parts are an added engineering complexity.

- **So what exactly is a storage engine?** A simple first answer, considering the big asks for any database - Atomicity, Consistency, Isolation, and Durability, or ACID - is that the storage engine's responsibility is Atomicity and Durability. This frees up the higher layers of the database to focus on the distributed coordination required to get strong Isolation guarantees like serializability, and providing the Consistency primitives needed to ensure data integrity.
  
Beyond atomicity and durability though, a storage engine has another big job: performance. The performance characteristics of a storage engine determine to a large extent the ceiling on the performance of the entire database. For example: almost every storage engine has a write-ahead log for quickly making a write durable, but which will later get moved to the main indexed data structure. This is a performance optimization that is pretty vital, but also tricky to get right, and involves delicate performance tradeoffs, as committing writes to the write-ahead-log while maintaining high concurrency is tricky.

Storage engines also have to provide a defined isolation model (e.g. that your reads will reflect writes that are still in the write-ahead log and will be as of some single point-in-time snapshot), while supporting many concurrent operations. While that isolation model might be simpler than the overall isolation provided by the database (for instance, RocksDB provides snapshot isolation, while CockroachDB does extra bookkeeping to provide serializability), maintaining it at high performance is still complex.


## A RocksDB primer

RocksDB is a single-node key-value storage engine. The design is based on log-structured merge trees (LSMs).

In RocksDB, keys and values are stored as sorted strings in files called SSTables. These SSTables are arranged in several levels. Within a single level, SSTables are non-overlapping: one SSTable might contain keys covering the range [a,b), the next [b,d), and so on. The key-space does overlap between levels: if you have two levels, the first might have two SSTables (covering the ranges above), but the second level might have a single SSTable over the keyspace [a,e). Looking for the key aardvark requires looking in two SSTables: the [a,b) SSTable in Level 1, and the [a,e) SSTable in Level 2. Each SSTable is internally sorted (as in the name), so lookups within an SSTable take log(n) time. SSTables store their keys in blocks, and have an internal index, so even though a single SSTable may be very large (gigabytes in size), only the index and the relevant block needs to be loaded into memory.

The levels are structured roughly so that each level is in total 10x as large as the level above it. New keys arrive at the highest layer, and as that level gets larger and larger and hits a threshold, some SSTables at that level get compacted into fewer (but larger) SSTables one level lower. The precise details of when to compact (and how) greatly affect performance;

