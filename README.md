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

## JSON vs GOB
```
BenchmarkJSON-4   	  770727	     15493 ns/op
BenchmarkGOB-4   	  145520	     84069 ns/op
```

Table of compare performance ns/op

|               |  JSON  |   GOB  |
|---------------|-------:|-------:|
| Benchmark/4-4 | 15 493 | 84 069 |


## Tests Storages


```
go test -bench=. -benchmem -benchtime=10s -cpuprofile=cpu.out -memprofile=mem.out

BenchmarkMap-4               373377 	     2888 ns/op

BenchmarkMutexMap/1-4  	 3241263	      3660 ns/op	     552 B/op	       6 allocs/op
BenchmarkMutexMap/2-4  	 3252266	      3741 ns/op	     552 B/op	       6 allocs/op
BenchmarkMutexMap/4-4  	 3148250	      3801 ns/op	     552 B/op	       6 allocs/op
BenchmarkMutexMap/8-4  	 3098610	      3909 ns/op	     552 B/op	       6 allocs/op

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

BenchmarkPostgreSQL/1-4      493	   2210216 ns/op
BenchmarkPostgreSQL/2-4 	 494	   2259960 ns/op
BenchmarkPostgreSQL/4-4 	 458	   2378314 ns/op
BenchmarkPostgreSQL/8-4 	 441	   2533566 ns/op
```

Table of compare performance ns/op

|                  | Only map | map with sync.RWMutex | sync.Map | Aerospike | Redis  | PostgreSQL |
|------------------|------:|------:|--------:|--------:|--------:|--------:|
| Benchmark/1-4 | 2 888 | 3 660 | 4 553 | 116 237 | 122 826 | 2 210 216 |
| Benchmark/2-4 | | 3 741 | 4 958 |  94 596 | 104 896 | 2 259 960 |
| Benchmark/4-4 | | 3 801 | 4 824 |  87 570 |  94 707 | 2 378 314 |
| Benchmark/8-4 | | 3 909 | 5 186 |  86 481 |  88 499 | 2 533 566 |

