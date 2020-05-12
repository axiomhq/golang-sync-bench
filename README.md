Benchmarking Go Synchronization Primitives
====

This repository contains some simple benchmarks of some different
synchronization primitives in Golang.

Run with
```
$ go test -bench .
```

The benchmarks are implemented as a loop:
```.go
for {
    takeLock()
    doWork()
    dropLock()
}
```
The work done is purely a simple CPU exercise with no memory allocations.
Different workloads will probably change the outcome by a lot.

Results for Different Go Versions
----
These numbers are from an unreliable developer laptop - so take with a grain of salt.
The numbers changed in a 10% range between different runs, but the overall trends remained the same.

```
go version go1.14.2 linux/amd64
goos: linux
goarch: amd64
BenchmarkSync/Atomic-8         	100000000	        11.2 ns/op
BenchmarkSync/ch0-8            	100000000	        12.7 ns/op
BenchmarkSync/ch1-8            	85027264	        13.8 ns/op
BenchmarkSync/ch10-8           	82793593	        14.2 ns/op
BenchmarkSync/Mutex-8          	57655360	        19.6 ns/op
BenchmarkSync/RWMutex-8        	19018388	        90.8 ns/op
BenchmarkSync/Spamch0-8        	100000000	        13.6 ns/op
BenchmarkSync/Spamch1-8        	72175544	        15.9 ns/op
BenchmarkSync/Spamch10-8       	75115940	        16.2 ns/op
BenchmarkSync/SpamRWMutex-8    	10031592	       105 ns/op
PASS
ok  	_/home/kamstrup/Axiom/hax	13.499s
```

```
go version go1.13.8 linux/amd64
goos: linux
goarch: amd64
BenchmarkSync/Atomic-8         	100000000	        10.2 ns/op
BenchmarkSync/ch0-8            	92577460	        13.7 ns/op
BenchmarkSync/ch1-8            	81593481	        15.1 ns/op
BenchmarkSync/ch10-8           	78057223	        15.0 ns/op
BenchmarkSync/Mutex-8          	55968832	        20.7 ns/op
BenchmarkSync/RWMutex-8        	15184892	        79.9 ns/op
BenchmarkSync/Spamch0-8        	100000000	        14.2 ns/op
BenchmarkSync/Spamch1-8        	63595525	        16.9 ns/op
BenchmarkSync/Spamch10-8       	68215941	        17.2 ns/op
BenchmarkSync/SpamRWMutex-8    	 9844534	       104 ns/op
PASS
ok  	_/home/kamstrup/Axiom/hax	12.932s
```

```
go version go1.12.8 linux/amd64
goos: linux
goarch: amd64
BenchmarkSync/Atomic-8         	50000000	        37.5 ns/op
BenchmarkSync/ch0-8            	200000000	         8.21 ns/op
BenchmarkSync/ch1-8            	100000000	        13.0 ns/op
BenchmarkSync/ch10-8           	100000000	        14.2 ns/op
BenchmarkSync/Mutex-8          	50000000	        26.2 ns/op
BenchmarkSync/RWMutex-8        	50000000	        96.3 ns/op
BenchmarkSync/Spamch0-8        	100000000	        13.9 ns/op
BenchmarkSync/Spamch1-8        	100000000	        16.0 ns/op
BenchmarkSync/Spamch10-8       	100000000	        16.7 ns/op
BenchmarkSync/SpamRWMutex-8    	10000000	       142 ns/op
PASS
ok  	_/home/kamstrup/Axiom/hax	21.329s
```

Read the Code
----
Please read the code before drawing your own conclusions,
it has been highly tailored for what we needed.
There are some interesting caveats in the comments as well! 