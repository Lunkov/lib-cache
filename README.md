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

Intel® Core™ i5 CPU 760 @ 2.80GHz × 4 (https://www.cpubenchmark.net/cpu.php?cpu=Intel+Core+i5-760S+%40+2.53GHz&id=782)

```
go test -bench=. -benchmem -benchtime=10s -cpuprofile=cpu.out -memprofile=mem.out

BenchmarkMap/1-4       	 3241263	      3660 ns/op	     552 B/op	       6 allocs/op
BenchmarkMap/2-4       	 3252266	      3741 ns/op	     552 B/op	       6 allocs/op
BenchmarkMap/4-4       	 3148250	      3801 ns/op	     552 B/op	       6 allocs/op
BenchmarkMap/8-4       	 3098610	      3909 ns/op	     552 B/op	       6 allocs/op

BenchmarkSyncMap/1-4   	 2677140	      4553 ns/op	     616 B/op	       9 allocs/op
BenchmarkSyncMap/2-4   	 2502535	      4958 ns/op	     616 B/op	       9 allocs/op
BenchmarkSyncMap/4-4   	 2506807	      4824 ns/op	     616 B/op	       9 allocs/op
BenchmarkSyncMap/8-4   	 2471282	      5186 ns/op	     616 B/op	       9 allocs/op

BenchmarkAerospike/1-4 	  100814	    116237 ns/op	    3973 B/op	     111 allocs/op
BenchmarkAerospike/2-4 	  131563	     94596 ns/op	    3972 B/op	     111 allocs/op
BenchmarkAerospike/4-4 	  137200	     87570 ns/op	    3971 B/op	     111 allocs/op
BenchmarkAerospike/8-4 	  138630	     86481 ns/op	    3971 B/op	     111 allocs/op

BenchmarkRedis/1-4     	  105042	    122826 ns/op	    2236 B/op	      33 allocs/op
BenchmarkRedis/2-4     	  111357	    104896 ns/op	    2233 B/op	      33 allocs/op
BenchmarkRedis/4-4     	  131476	     94707 ns/op	    2233 B/op	      33 allocs/op
BenchmarkRedis/8-4     	  131264	     88499 ns/op	    2233 B/op	      33 allocs/op
```

Table of compare performance ns/op

|                  | map with sync.RWMutex | sync.Map | Aerospike | Redis  |
|------------------|:-----|:----:|-------:|-------:|
| BenchmarkMap/1-4 | 3660 | 4553 | 116237 | 122826 |
| BenchmarkMap/2-4 | 3741 | 4958 |  94596 | 104896 |
| BenchmarkMap/4-4 | 3801 | 4824 |  87570 |  94707 |
| BenchmarkMap/8-4 | 3909 | 5186 |  86481 |  88499 |

