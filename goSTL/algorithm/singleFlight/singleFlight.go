package singleFlight

import "sync"

//呼叫请求结构体
type call struct {
	wg  sync.WaitGroup //可重入锁
	val interface{}    //请求结果
	err error          //错误反馈
}

type Group struct {
	mu sync.Mutex // protects m
	m  map[string]*call//所有请求
}

//防止击穿缓存,对同一个key进行请求时需要分别进行,利用可重入锁实现
func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*call)
	}
	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		c.wg.Wait()         // 如果请求正在进行中，则等待
		return c.val, c.err // 请求结束，返回结果
	}
	c := new(call)
	c.wg.Add(1)  // 发起请求前加锁
	g.m[key] = c // 添加到 g.m，表明 key 已经有对应的请求在处理
	g.mu.Unlock()

	c.val, c.err = fn() // 调用 fn，发起请求
	c.wg.Done()         // 请求结束

	g.mu.Lock()
	delete(g.m, key) // 更新 g.m
	g.mu.Unlock()

	return c.val, c.err // 返回结果
}

func (g *Group) DoChan(key string, fn func() (interface{}, error)) (ch chan interface{}) {
	ch = make(chan interface{})
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*call)
	}
	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		c.wg.Wait() // 如果请求正在进行中，则等待
		ch <- c.val
		return ch
	}
	c := new(call)
	c.wg.Add(1)  // 发起请求前加锁
	g.m[key] = c // 添加到 g.m，表明 key 已经有对应的请求在处理
	g.mu.Unlock()

	c.val, c.err = fn() // 调用 fn，发起请求
	c.wg.Done()         // 请求结束
	ch <- c.val

	g.mu.Lock()
	delete(g.m, key) // 更新 g.m
	g.mu.Unlock()

	return ch
}
