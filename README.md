gochannels
==========

Simple program to test Go (golang) channels and system performance.


###Usage
```
Usage of ./channels:
  -D=false: debug output
  -b=false: enable blocking queue behavior (unbuffered channel)
  -d=0: microsecond delay after worker processes job before getting another one
  -i=100: number of jobs to run
  -p=1: number of logical CPUs to use (0 means use ALL)
  -q=10: number of jobs to hold in the queue
  -v=false: detailed output
  -w=10: number of workers threads
```
