# Before tests
```
sudo apt-get install libcanberra-gtk-module
sudo apt-get install libcanberra-gtk-module libcanberra-gtk3-module
sudo apt-get install graphviz
```

# Tests
About tools: https://blog.golang.org/pprof
```
go test -bench=. -benchmem -benchtime=10s -cpuprofile=cpu.out -memprofile=mem.out
```

```
go tool pprof ./mem.out
go tool pprof ./cpu.out
```

```
top10
web mallocgc
```
# Tests results

```
go test -bench=. -benchmem -benchtime=10s -cpuprofile=cpu.out -memprofile=mem.out

BenchmarkMap-4
7371068	      1881 ns/op	     585 B/op	       6 allocs/op
BenchmarkRedis-4
29284	    400196 ns/op	    2076 B/op	      34 allocs/op
BenchmarkCacheSyncMap-4
5112874	      2642 ns/op	     552 B/op	      11 allocs/op

```
