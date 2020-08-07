# Hitler
Hitler 是 golang 实现的, 基于Snowflake算法的唯一 ID 生成器。
Hitler 以组件形式工作在应用项目中, 支持自定义 workerId 位数和初始化策略, 从而适用于docker等虚拟化环境下实例自动重启、漂移等场景。 
在实现上, Hitler 通过借用未来时间来解决 sequence 天然存在的并发限制; 
采用 RingBuffer 来缓存已生成的 UID, 并行化 UID 的生产和消费, 同时对 CacheLine 补齐，避免了由 RingBuffer 带来的硬件级
「伪共享」问题. 最终单机 QPS 可达 -(待测试) 万。
---- 参考BaiDu的 [UidGenerator]

## example
```go
	rb := NewRBuffer(&RBufferConfig{Size: 1000, MachineCap: 22, DbUrl: "root:Dachang1234!@(127.0.0.1:3306)/id_gen?charset=utf8mb4"})
	for i := 0; i < 20000; i++ {
		time.Sleep(time.Millisecond)
		log.Println(rb.GetID())
	}
```
## Snowflake 算法
```go

+=========+==================+=======================+================+
|   bit   |     delta second |  work node id         |    sequence    |
+=========+==================+=======================+================+
|     1   |          28      |       22              |        13      |
|=========+==================+=======================+================+

```
Snowflake 算法描述：指定机器 & 同一时刻 & 某一并发序列，是唯一的。据此可生成一个 64 bits 的唯一 ID（long）。默认采用上图字节分配方式：

* sign(1bit)
固定 1bit 符号标识，即生成的 UID 为正数。

* delta seconds (28 bits)
当前时间，相对于时间基点"2016-05-20"的增量值，单位：秒，最多可支持约 8.7 年

* worker id (22 bits)
机器 id，最多可支持约 420w 次机器启动。内置实现为在启动时由数据库分配，默认分配策略为用后即弃，后续可提供复用策略。

* sequence (13 bits)
每秒下的并发序列，13 bits 可支持每秒 8192 个并发。

## Double Ring Buffer
RingBuffer 环形数组，数组每个元素成为一个 slot。RingBuffer 容量，默认为 Snowflake 算法中 sequence 最大值，且为 2^N。可通过boostPower配置进行扩容，以提高 RingBuffer 读写吞吐量。

Tail 指针、Cursor 指针用于环形数组上读写 slot：

* Tail 指针
表示 Producer 生产的最大序号(此序号从 0 开始，持续递增)。Tail 不能超过 Cursor，即生产者不能覆盖未消费的 slot。当 Tail 已赶上 curosr，此时可通过rejectedPutBufferHandler指定 PutRejectPolicy

* Cursor 指针
表示 Consumer 消费到的最小序号(序号序列与 Producer 序列相同)。Cursor 不能超过 Tail，即不能消费未生产的 slot。当 Cursor 已赶上 tail，此时可通过rejectedTakeBufferHandler指定 TakeRejectPolicy
![RingBuffer](doc/ringbuffer.png)
CachedUidGenerator 采用了双 RingBuffer，Uid-RingBuffer 用于存储 Uid、Flag-RingBuffer 用于存储 Uid 状态(是否可填充、是否可消费)

由于数组元素在内存中是连续分配的，可最大程度利用 CPU cache 以提升性能。但同时会带来「伪共享」FalseSharing 问题，为此在 Tail、Cursor 指针、Flag-RingBuffer 中采用了 CacheLine
补齐方式。

![FalseSharing](doc/cacheline_padding.png) 
#### RingBuffer 填充时机 ####
* 初始化预填充  
  RingBuffer 初始化时，预先填充满整个 RingBuffer.
  
* 周期填充  
  通过 Schedule 线程，定时补全空闲 slots。可通过```scheduleInterval```配置，以应用定时填充功能，并指定 Schedule 时间间隔 
  
Quick Start
------------
```go
go get -u github.com/Callmedachang/hitler
```

### 关于 UID 比特分配的建议
对于并发数要求不高、期望长期使用的应用, 可增加```timeCap```位数, 减少```sequenceCap```位数. 例如节点采取用完即弃的 WorkerIdAssigner 策略, 重启频率为 12 次/天,
那么配置成```{"machineCap":23,"timeCap":31,"sequenceCap":9}```时, 可支持 28 个节点以整体并发量 14400 UID/s 的速度持续运行 68 年.

对于节点重启频率频繁、期望长期使用的应用, 可增加```machineCap```和```timeCap```位数, 减少```sequenceCap```位数. 例如节点采取用完即弃的 WorkerIdAssigner 策略, 重启频率为 24*12 次/天,
那么配置成```{"machineCap":27,"timeCap":30,"sequenceCap":6}```时, 可支持 37 个节点以整体并发量 2400 UID/s 的速度持续运行 34 年.

### 性能测试
之后补充~