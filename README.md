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
...之后补充吧