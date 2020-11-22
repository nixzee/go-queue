package queue

import (
	"sort"
	"sync"
)

//TODO:
// * Add more testing
// * Clean up
// * Documentation

//---------------------------------------------------------------------------------------------------
// Interface
//---------------------------------------------------------------------------------------------------

//Owner provides high level methods of a queue
type Owner interface {
	//Close provides cleanup
	Close()
	//Resize will flush and resize the queue
	Resize(size int) (elements []interface{}, priorities []int)
}

//Flush provides methods of flushing the queue
type Flush interface {
	//Flush will flush the queue of all elements and return what was in it
	Flush() (elements []interface{}, priorities []int)
}

//Event provides methods of events of the queue
type Event interface {
	//GetSignal returns a signal channel that can be monitored if something is enqueue (event driven)
	GetSignal() (signal <-chan struct{})
}

//Info provides methods of getting info about a queue
type Info interface {
	//GetSize will return the size (max size of the queue)
	GetSize() (size int)
	//GetLength will return the current length of the queue
	GetLength() (len int)
}

//Peek provides methodes to peek at the queue
type Peek interface {
	//Peek allows for peeking at all elements in queue
	Peek() (elements []interface{}, empty bool)
	//PeekHead allows for peek at last element
	PeekHead() (element interface{}, empty bool)
	//PeekTail allows for peek at first element
	PeekTail() (element interface{}, empty bool)
}

//PeekPriority provides methodes to peek at the queue (with priority)
type PeekPriority interface {
	//Peek allows for peeking at all elements and priorities in queue
	PeekPriority() (elements []interface{}, priorities []int, empty bool)
	//PeekHead allows for peek at last element and priority
	PeekHeadPriority() (element interface{}, priority int, empty bool)
	//PeekTail allows for peek at first element and priority
	PeekTailPriority() (element interface{}, priority int, empty bool)
}

//Dequeue provides methods to dequeue
type Dequeue interface {
	//Dequeue will dequeue a single (last) element
	Dequeue() (element interface{}, underflow bool)
}

//DequeuePriority provides methods to dequeue with and get the priority
type DequeuePriority interface {
	//Dequeue will dequeue a single (last) element and priority
	DequeuePriority() (element interface{}, priority int, underflow bool)
}

//Enqueue provides methods to enqueue
type Enqueue interface {
	//Enqueue will enqueue a single element
	Enqueue(element interface{}) (overflow bool)
}

//EnqueuePriority provides methods to enqueue with priority
type EnqueuePriority interface {
	//Enqueue will enqueue a single element with priority
	EnqueuePriority(element interface{}, priority int) (overflow bool)
}

//---------------------------------------------------------------------------------------------------
// Implementation
//---------------------------------------------------------------------------------------------------

//Ensure the implementation
var _ Owner = &queue{}
var _ Flush = &queue{}
var _ Event = &queue{}
var _ Info = &queue{}
var _ Peek = &queue{}
var _ PeekPriority = &queue{}
var _ Dequeue = &queue{}
var _ DequeuePriority = &queue{}
var _ Enqueue = &queue{}
var _ EnqueuePriority = &queue{}

//HAHAHAAAHAHA (The interface is named Interface...)
var _ sort.Interface = &queue{}

//NewQueue returns a new queue
func NewQueue(size int, polling bool) interface {
	Owner
	Flush
	Event
	Info
	Peek
	PeekPriority
	Dequeue
	DequeuePriority
	Enqueue
	EnqueuePriority
} {
	//Check if size is valid
	if size <= 0 {
		size = DefaultSize
	}
	//Signaling...check if polling
	var signal chan struct{}
	if !polling {
		signal = make(chan struct{}, size)
	}
	//Create the containers
	containers := make([]*container, 0)
	//Create the queue
	return &queue{
		size:       size,
		signal:     signal,
		polling:    polling,
		containers: containers,
	}
}

//queue provides a pointer implementation of Queue
type queue struct {
	sync.Mutex
	containers []*container  //containers
	size       int           //the max size of the queue
	signal     chan struct{} //signal to notify that element has been enqueued
	polling    bool          //Don't use signal if polling
}

//---------------------------------------------------------------------------------------------------
// Owner Implementation
//---------------------------------------------------------------------------------------------------

//Close is provides cleanup
func (q *queue) Close() {
	q.Lock()
	defer q.Unlock()
	//Cleanup
	q.signal, q.containers = nil, nil
	q.size = 0
	return
}

//Resize will flush and resize the queue
func (q *queue) Resize(size int) (elements []interface{}, priorities []int) {
	q.Lock()
	defer q.Unlock()
	//Get the elements and priorities
	for _, container := range q.containers {
		elements = append(elements, container.element)
		priorities = append(priorities, container.priority)
	}
	//Reset the containers
	// q.containers = make([]container, q.size)
	q.containers = make([]*container, 0)
	//Check if size is valid
	if size <= 0 {
		size = DefaultSize
	}
	q.size = size
	return
}

//---------------------------------------------------------------------------------------------------
// Flush Implementation
//---------------------------------------------------------------------------------------------------

//Flush will flush the queue of all elements and return what was in it
func (q *queue) Flush() (elements []interface{}, priorities []int) {
	q.Lock()
	defer q.Unlock()
	//Get the elements and priorities
	for index, container := range q.containers {
		elements = append(elements, container.element)
		priorities = append(priorities, container.priority)
		q.containers[index] = nil
	}
	//Reset the containers
	// q.containers = make([]container, q.size)
	q.containers = make([]*container, 0)
	return
}

//---------------------------------------------------------------------------------------------------
// Event Implementation
//---------------------------------------------------------------------------------------------------

//GetSignal returns a signal channel that can be monitored if something is enqueue (event driven)
func (q *queue) GetSignal() (signal <-chan struct{}) {
	q.Lock()
	defer q.Unlock()
	signal = q.signal
	return
}

//---------------------------------------------------------------------------------------------------
// Info Implementation
//---------------------------------------------------------------------------------------------------

//GetSize will return the size (max size of the queue)
func (q *queue) GetSize() (size int) {
	q.Lock()
	defer q.Unlock()
	size = q.size
	return
}

//GetLength will return the current length of the queue
func (q *queue) GetLength() (len int) {
	q.Lock()
	defer q.Unlock()
	len = q.Len()
	return
}

//---------------------------------------------------------------------------------------------------
// Peek Implementation
//---------------------------------------------------------------------------------------------------

//Peek allows for peeking at all elements in queue
func (q *queue) Peek() (elements []interface{}, empty bool) {
	q.Lock()
	defer q.Unlock()
	len := q.Len()
	//Check if empty
	if empty = q.checkIfEmpty(); empty {
		return
	}
	//Copy
	containers := make([]*container, len)
	copy(containers, q.containers)
	//Populate the elements
	elements = make([]interface{}, 0)
	for _, container := range containers {
		elements = append(elements, container.element)
	}
	return
}

//PeekHead allows for peek at last element
func (q *queue) PeekHead() (element interface{}, empty bool) {
	q.Lock()
	defer q.Unlock()
	//Check if empty
	if empty = q.checkIfEmpty(); empty {
		return
	}
	//Get first element
	element = q.containers[0].element
	return
}

//PeekTail allows for peek at first element
func (q *queue) PeekTail() (element interface{}, empty bool) {
	q.Lock()
	defer q.Unlock()
	//Check if empty
	if empty = q.checkIfEmpty(); empty {
		return
	}
	//Get last element
	element = q.containers[q.Len()-1].element
	return
}

//---------------------------------------------------------------------------------------------------
// Peek Priority Implementation
//---------------------------------------------------------------------------------------------------

//Peek allows for peeking at all elements in queue
func (q *queue) PeekPriority() (elements []interface{}, priorities []int, empty bool) {
	q.Lock()
	defer q.Unlock()
	//Check if empty
	if empty = q.checkIfEmpty(); empty {
		return
	}
	len := q.Len()
	//Copy
	containers := make([]*container, len)
	copy(containers, q.containers)
	//Populate the elements
	elements = make([]interface{}, 0)
	priorities = make([]int, 0)
	for _, container := range containers {
		elements = append(elements, container.element)
		priorities = append(priorities, container.priority)
	}
	return
}

//PeekHead allows for peek at last element
func (q *queue) PeekHeadPriority() (element interface{}, priority int, empty bool) {
	q.Lock()
	defer q.Unlock()
	//Check if empty
	if empty = q.checkIfEmpty(); empty {
		return
	}
	//Get first element
	container := q.containers[0]
	element = container.element
	priority = container.priority
	return
}

//PeekTail allows for peek at first element
func (q *queue) PeekTailPriority() (element interface{}, priority int, empty bool) {
	q.Lock()
	defer q.Unlock()
	//Check if empty
	if empty = q.checkIfEmpty(); empty {
		return
	}
	//Check if empty
	empty = q.checkIfEmpty()
	//Get last element
	container := q.containers[q.Len()-1]
	element = container.element
	priority = container.priority
	return
}

//---------------------------------------------------------------------------------------------------
// Dequeue Implementation
//---------------------------------------------------------------------------------------------------

//Dequeue will dequeue a single (last) element
func (q *queue) Dequeue() (element interface{}, underflow bool) {
	q.Lock()
	defer q.Unlock()
	//Dequeue
	underflow, element, _ = q.dequeue()
	return
}

//---------------------------------------------------------------------------------------------------
// Dequeue Priority Implementation
//---------------------------------------------------------------------------------------------------

//Dequeue will dequeue a single (last) element
func (q *queue) DequeuePriority() (element interface{}, priority int, underflow bool) {
	q.Lock()
	defer q.Unlock()
	//Dequeue
	underflow, element, priority = q.dequeue()
	return
}

//---------------------------------------------------------------------------------------------------
// Enqueue Implementation
//---------------------------------------------------------------------------------------------------

//Enqueue will enqueue a single element
func (q *queue) Enqueue(element interface{}) (overflow bool) {
	q.Lock()
	defer q.Unlock()
	//Enqueue
	overflow = q.enqueue(element, DefaultPriority)
	//Trigger signal
	q.triggerSignal()
	return
}

//---------------------------------------------------------------------------------------------------
// Enqueue Priority Implementation
//---------------------------------------------------------------------------------------------------

//Enqueue will enqueue a single element
func (q *queue) EnqueuePriority(element interface{}, priority int) (overflow bool) {
	q.Lock()
	defer q.Unlock()
	//Enqueue
	overflow = q.enqueue(element, priority)
	//Trigger signal
	q.triggerSignal()
	return
}

//---------------------------------------------------------------------------------------------------
// Hidden
//---------------------------------------------------------------------------------------------------

//triggerSignal will send the signal that element(s) have been enqueued (non-blocking)
func (q *queue) triggerSignal() {
	//check polling
	if q.polling {
		return
	}
	//Fire
	select {
	case q.signal <- struct{}{}:
	default:
	}
}

//checkIfEmpty will check if the queue is empty
func (q *queue) checkIfEmpty() (empty bool) {
	empty = q.Len() <= 0
	return
}

//checkIfFull will check if the queue is full
func (q *queue) checkIfFull() (full bool) {
	full = q.Len() >= q.size
	return
}

//enqueue performs the enqueue logic
func (q *queue) enqueue(element interface{}, priority int) (overflow bool) {
	//Check if queue is full (overflow)
	if q.checkIfFull() {
		overflow = true
		return
	}
	//Push
	// heap.Push(q, container{element: element, priority: priority})
	q.containers = append(q.containers, &container{element: element, priority: priority})
	sort.Sort(q)
	return
}

//dequeue performs the dequeue logic
func (q *queue) dequeue() (underflow bool, element interface{}, priority int) {
	//Check if queue is empty (underflow)
	if q.checkIfEmpty() {
		underflow = true
		return
	}
	//Pop
	container := q.containers[0]
	element = container.element
	priority = container.priority
	q.containers[0] = nil //Come garbage collect
	q.containers = q.containers[1:]
	return
}

//---------------------------------------------------------------------------------------------------
// Sort
//---------------------------------------------------------------------------------------------------

//Len implements Length
func (q *queue) Len() int {
	return len(q.containers)
}

//Less implements Length
//Note: This is technically backwards to make it a "max"
func (q *queue) Less(i, j int) bool {
	return q.containers[i].priority > q.containers[j].priority
}

//Swap implements Swap
func (q *queue) Swap(i, j int) {
	q.containers[i], q.containers[j] = q.containers[j], q.containers[i]
}
