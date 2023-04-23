package pool

import (
	"fmt"
	"sync"
)

//写一个Job,暴露一个回调方法
type cb func()
type Job struct {
	id int
	f  cb
}

func NewJob(i int, f1 cb) *Job {
	return &Job{
		id: i,
		f:  f1,
	}
}
func (this *Job) Do() {
	//	fmt.Printf("i am job id=%d\n", this.id)
	this.f()
}

type Pool struct {
	EntryPoint chan Job
	JobQueue   chan Job
	Num        int
	wg         *sync.WaitGroup
}

//new一个线程池对象
func NewPool(n int) *Pool {
	return &Pool{
		EntryPoint: make(chan Job),
		JobQueue:   make(chan Job),
		Num:        n,
		wg:         new(sync.WaitGroup),
	}
}

//写一个线程池，对job的外部 channel
//写一个线程池，对内的job队伍 channel
//写一个work干活的方法去调用 job的方法
func (this *Pool) work(id int) {
	defer func() {
		this.wg.Done()
		fmt.Printf("i am work id=%d,exit\n", id)
	}()
	for job := range this.JobQueue {
		job.Do()
		fmt.Printf("i am work id=%d\n", id)
	}
}

//写一个并发数
//写一个Hold住的run方法
func (this *Pool) Run() {
	//开启work goroutine
	for i := 0; i < this.Num; i++ {
		this.wg.Add(1)
		go this.work(i)
	}
loop:
	for {
		select {
		case job, ok := <-this.EntryPoint:
			if ok {
				this.JobQueue <- job
			} else {
				close(this.JobQueue)
				break loop
			}
		default:
		}
	}
	this.wg.Wait()
}
func (this *Pool) Close() {
	close(this.EntryPoint)
}
