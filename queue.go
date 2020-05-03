package queue

import "sync"

//
type Queue interface {
	Close()
	GetSignal()
	Dequeue()
	DequeueMultiple()
	Enqueue()
	EnqueueMultiple()
	Flush()
	Peek()
	PeekHead()
	PeekTail()
}

//
type queue struct {
	sync.Mutex               //
	data       []interface{} //
	size       int           //
	index      int           //
	singal     chan struct{} //
}

// //NewQueue
// func NewQueue(size int) interface {
// 	Queue
// } {
// 	//Check if size is valid
// 	if size <= 0 {
// 		size = DefaultSize
// 	}
// 	return &queue{
// 		size: size,
// 	}
// }

func (q *queue) Close() {
	q.Lock()
	defer q.Unlock()
	//TODO
	return
}

func (q *queue) GetSignal() (signal <-chan struct{}) {
	q.Lock()
	defer q.Unlock()
	//signal
	signal = q.singal
	return
}

func (q *queue) Dequeue() (element interface{}, underflow bool) {
	q.Lock()
	defer q.Unlock()
	//dequeue
	underflow, q.len, q.data, element = dequeue(q.len, q.size, q.data)
	//TODO signal
	return
}

func (q *queue) DequeueMultiple(count int) (elements []interface{}, underflow bool) {
	q.Lock()
	defer q.Unlock()
	//dequeue based on count
	for i := 0; i < count; i++ {
		var element interface{}
		//dequeue
		underflow, q.len, q.data, element = dequeue(q.len, q.size, q.data)
		//check for underflow
		if underflow {
			return
		}
		elements = append(elements, element)
		//TODO signal
	}
	return
}

func (q *queue) Enqueue(element interface{}) (overflow bool) {
	q.Lock()
	defer q.Unlock()
	//TODO
	return
}

func (q *queue) EnqueueMultiple(elements []interface{}) (overflow bool) {
	q.Lock()
	defer q.Unlock()
	//TODO
	return
}

func (q *queue) Flush() (elements []interface{}) {
	q.Lock()
	defer q.Unlock()
	//TODO
	return
}

func (q *queue) Peek() (elements []interface{}) {
	q.Lock()
	defer q.Unlock()
	//TODO
	return
}

func (q *queue) PeekHead() (element interface{}) {
	q.Lock()
	defer q.Unlock()
	//TODO
	return
}

func (q *queue) PeekTail() (element interface{}) {
	q.Lock()
	defer q.Unlock()
	//TODO
	return
}
