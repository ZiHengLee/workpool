# workpool
Golang 协程池

go test -bench=. -benchmem -run=none

对比 原生协程和20个groutine的协程池

运行结果
BenchmarkNoGroutineWork-8                      1        27077517237 ns/op       11201624 B/op      26303 allocs/op sum:1000000000
BenchmarkSimpleGoroutineSetTimes-8             1        24806976018 ns/op       787622776 B/op   1823431 allocs/op sum:1815067704
BenchmarkPoolPutSetTimes-8                     1        28313076354 ns/op           5248 B/op         34 allocs/op sum:2999996278
BenchmarkPoolTimeLifeSetTimes-8                1        29985336842 ns/op           3624 B/op         30 allocs/op sum:4000000000

结论
使用协程池可以减少内存分配空间和内存分配次数
但是性能比原生协程有所减少
可以稍微提高协程池大小，性能可以与原生协程差不多甚至略有提升
