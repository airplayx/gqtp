package config

type ConcurrentPool struct {
	Ch chan int
}

func (p *ConcurrentPool) Add() {
	p.Ch <- 1
}

func (p *ConcurrentPool) Done() {
	<-p.Ch
}

func NewPool(number int) *ConcurrentPool {
	pool := &ConcurrentPool{}
	pool.Ch = make(chan int, number)
	return pool
}