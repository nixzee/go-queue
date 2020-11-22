package queue

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	fatalUnderflow string = "Unit Test %s encountered underflow"
	fatalOverflow  string = "Unit Test %s encountered overflow"
)

//TestClose will test close
func TestClose(t *testing.T) {
	//Create Queue
	testQueue := NewQueue(1, false)
	//Close
	testQueue.Close()
}

//TestGetSignal will test getting the signal
func TestGetSignal(t *testing.T) {
	cases := map[string]struct {
		oNotSignal chan struct{}
	}{
		"GetSignal": {
			oNotSignal: nil,
		},
	}

	//Test cases
	for cDesc, c := range cases {
		//Create Queue
		testQueue := NewQueue(1, false)
		defer testQueue.Close()
		//Get Signal
		signal := testQueue.GetSignal()
		//Assert
		assert.NotEqual(t, c.oNotSignal, signal, fmt.Sprintf(cDesc))
	}
}

//TestFlush will test the flush
func TestFlush(t *testing.T) {
	cases := map[string]struct {
		iElements []interface{}
		oLength   int
	}{
		"Flush": {
			iElements: []interface{}{1, 2},
			oLength:   0,
		},
	}

	//Test cases
	for cDesc, c := range cases {
		//Create Queue
		testQueue := NewQueue(len(c.iElements), true)
		defer testQueue.Close()
		//Enqueue
		for _, element := range c.iElements {
			testQueue.Enqueue(element)
		}
		// testQueue.EnqueueMultiple(c.iElements)
		//Flush
		testQueue.Flush()
		//GetIndex
		length := testQueue.GetLength()
		//Assert
		assert.Equal(t, c.oLength, length, fmt.Sprintf(cDesc))
	}
}

//TestDequeueConcurrent will test dequeue using go routines (concurrency)
func TestDequeueConcurrent(t *testing.T) {
	cases := map[string]struct {
		iSize     int
		iElements []interface{}
		oElements []interface{}
	}{
		"Dequeue_Single": {
			iSize:     10,
			iElements: []interface{}{1},
			oElements: []interface{}{1},
		},
		"Dequeue_Multiple": {
			iSize:     10,
			iElements: []interface{}{1, 2, 3},
			oElements: []interface{}{1, 2, 3},
		},
	}

	//Test cases
	for cDesc, c := range cases {
		var wg sync.WaitGroup
		stop := make(chan struct{})
		var elements []interface{}
		//Create Queue
		testQueue := NewQueue(c.iSize, false)
		defer testQueue.Close()
		//dequeue routine
		wg.Add(1)
		go func(t *testing.T) {
			defer wg.Done()
			signal := testQueue.GetSignal()
			for {
				select {
				case <-stop:
					return
				case <-signal:
					//Dequeue
					element, underflow := testQueue.Dequeue()
					//Check underflow
					if underflow {
						t.Fatalf(fatalUnderflow, "Dequeue_Concurrent")
					}
					//elements
					elements = append(elements, element)
				}
			}
		}(t)
		time.Sleep(10 * time.Millisecond)
		//enqueue routine
		wg.Add(1)
		go func(t *testing.T) {
			defer wg.Done()
			//enqueue
			for _, element := range c.iElements {
				if overflow := testQueue.Enqueue(element); overflow {
					t.Fatalf(fatalOverflow, "Dequeue_Concurrent")
				}
				time.Sleep(100 * time.Millisecond)
			}
			//stop the dequeue
			close(stop)
		}(t)
		//Wait
		wg.Wait()
		//Assert
		assert.Equal(t, c.oElements, elements, fmt.Sprintf("%s :Elements", cDesc))
	}
}

//TestDequeuePriorityConcurrent will test dequeue priority using go routines (concurrency)
func TestDequeuePriorityConcurrent(t *testing.T) {
	cases := map[string]struct {
		iSize       int
		iElements   []interface{}
		iPriorities []int
		oElements   []interface{}
		oPriorities []int
	}{
		"DequeuePriority_Single": {
			iSize:       10,
			iElements:   []interface{}{1},
			iPriorities: []int{10},
			oElements:   []interface{}{1},
			oPriorities: []int{10},
		},
		"DequeuePriority_Multiple": {
			iSize:       10,
			iElements:   []interface{}{1, 2, 3},
			iPriorities: []int{10, 0, 100},
			oElements:   []interface{}{1, 2, 3},
			oPriorities: []int{10, 0, 100},
		},
	}

	//Test cases
	for cDesc, c := range cases {
		var wg sync.WaitGroup
		stop := make(chan struct{})
		var elements []interface{}
		var priorities []int
		//Create Queue
		testQueue := NewQueue(c.iSize, false)
		defer testQueue.Close()
		//dequeue routine
		wg.Add(1)
		go func(t *testing.T) {
			defer wg.Done()
			signal := testQueue.GetSignal()
			for {
				select {
				case <-stop:
					return
				case <-signal:
					//Dequeue
					element, priority, underflow := testQueue.DequeuePriority()
					//Check underflow
					if underflow {
						t.Fatalf(fatalUnderflow, "DequeuePriority_Concurrent")
					}
					//elements
					elements = append(elements, element)
					priorities = append(priorities, priority)
				}
			}
		}(t)
		time.Sleep(10 * time.Millisecond)
		//enqueue routine
		wg.Add(1)
		go func(t *testing.T) {
			defer wg.Done()
			//enqueue
			for index, element := range c.iElements {
				if overflow := testQueue.EnqueuePriority(element, c.iPriorities[index]); overflow {
					t.Fatalf(fatalOverflow, "DequeuePriority_Concurrent")
				}
				time.Sleep(100 * time.Millisecond)
			}
			//stop the dequeue
			close(stop)
		}(t)
		//Wait
		wg.Wait()
		//Assert
		assert.Equal(t, c.oElements, elements, fmt.Sprintf("%s :Elements", cDesc))
		assert.Equal(t, c.oPriorities, priorities, fmt.Sprintf("%s :Priorities", cDesc))
	}
}

//TestDequeueSerialized will test dequeue serialized
func TestDequeueSerialized(t *testing.T) {
	cases := map[string]struct {
		iSize     int
		iElements []interface{}
		oElements []interface{}
	}{
		"Dequeue_Multiple": {
			iSize:     10,
			iElements: []interface{}{1, 2, 3, 4, 5, 6},
			oElements: []interface{}{1, 2, 3, 4, 5, 6},
		},
	}

	//Test cases
	for cDesc, c := range cases {
		var elements []interface{}
		//Create Queue
		testQueue := NewQueue(c.iSize, true)
		defer testQueue.Close()
		//enqueue
		for _, element := range c.iElements {
			if overflow := testQueue.Enqueue(element); overflow {
				t.Fatalf(fatalOverflow, "Dequeue_Serialized")
			}
		}
		//dequeue
		for range c.iElements {
			element, underflow := testQueue.Dequeue()
			//Check underflow
			if underflow {
				t.Fatalf(fatalUnderflow, "Dequeue_Serialized")
			}
			//Append
			elements = append(elements, element)
		}
		//Assert
		assert.Equal(t, c.oElements, elements, fmt.Sprintf("%s :Elements", cDesc))
		//Delete
		// fmt.Println(elements)
		// fmt.Println(priorities)
	}
}

//TestDequeuePrioritySerialized will test dequeue priority serialized
func TestDequeuePrioritySerialized(t *testing.T) {
	cases := map[string]struct {
		iSize       int
		iElements   []interface{}
		iPriorities []int
		oElements   []interface{}
		oPriorities []int
	}{
		"DequeuePriority_Multiple_Random": {
			iSize:       10,
			iElements:   []interface{}{1, 2, 3},
			iPriorities: []int{10, 0, 100},
			oElements:   []interface{}{3, 1, 2},
			oPriorities: []int{100, 10, 0},
		},
	}

	//Test cases
	for cDesc, c := range cases {
		var elements []interface{}
		var priorities []int
		//Create Queue
		testQueue := NewQueue(c.iSize, true)
		defer testQueue.Close()
		//enqueue
		for index, element := range c.iElements {
			if overflow := testQueue.EnqueuePriority(element, c.iPriorities[index]); overflow {
				t.Fatalf(fatalOverflow, "DequeuePriority_Serialized")
			}
		}
		//dequeue
		for range c.iElements {
			element, priority, underflow := testQueue.DequeuePriority()
			//Check underflow
			if underflow {
				t.Fatalf(fatalUnderflow, "DequeuePriority_Serialized")
			}
			//Append
			elements = append(elements, element)
			priorities = append(priorities, priority)
		}
		//Assert
		assert.Equal(t, c.oElements, elements, fmt.Sprintf("%s :Elements", cDesc))
		assert.Equal(t, c.oPriorities, priorities, fmt.Sprintf("%s :Priorities", cDesc))
	}
}

//TestPeek will test the peek
func TestPeek(t *testing.T) {
	cases := map[string]struct {
		iElements []interface{}
		oElements []interface{}
		oEmpty    bool
	}{
		"Peek_Empty": {
			iElements: []interface{}{},
			oElements: nil,
			oEmpty:    true,
		},
		"Peek_Single_Element": {
			iElements: []interface{}{1},
			oElements: []interface{}{1},
			oEmpty:    false,
		},
		"Peek_Multiple_Element": {
			iElements: []interface{}{1, 2, 3},
			oElements: []interface{}{1, 2, 3},
			oEmpty:    false,
		},
	}

	//Test cases
	for cDesc, c := range cases {
		//Create Queue
		testQueue := NewQueue(len(c.iElements), true)
		defer testQueue.Close()
		//Enqueue
		for _, element := range c.iElements {
			if overflow := testQueue.Enqueue(element); overflow {
				t.Fatalf(fatalOverflow, "Peek")
			}
		}
		//Peak
		elements, empty := testQueue.Peek()
		//Assert
		assert.Equal(t, c.oElements, elements, fmt.Sprintf("%s :Elements", cDesc))
		assert.Equal(t, c.oEmpty, empty, fmt.Sprintf("%s :Empty", cDesc))
	}
}

//TestPeek will test the peek with priority
func TestPeekPriority(t *testing.T) {
	cases := map[string]struct {
		iElements   []interface{}
		iPriorities []int
		oElements   []interface{}
		oPriorities []int
		oEmpty      bool
	}{
		"PeekPriority_Empty": {
			iElements:   []interface{}{},
			iPriorities: []int{},
			oElements:   nil,
			oPriorities: nil,
			oEmpty:      true,
		},
		"PeekPriority_Single_Element": {
			iElements:   []interface{}{1},
			iPriorities: []int{0},
			oElements:   []interface{}{1},
			oPriorities: []int{0},
			oEmpty:      false,
		},
		"PeekPriority_Multiple_Element": {
			iElements:   []interface{}{1, 2, 3},
			iPriorities: []int{0, 0, 0},
			oElements:   []interface{}{1, 2, 3},
			oPriorities: []int{0, 0, 0},
			oEmpty:      false,
		},
	}

	//Test cases
	for cDesc, c := range cases {
		//Create Queue
		testQueue := NewQueue(len(c.iElements), true)
		defer testQueue.Close()
		//Enqueue
		for index, element := range c.iElements {
			if overflow := testQueue.EnqueuePriority(element, c.iPriorities[index]); overflow {
				t.Fatalf(fatalOverflow, "PeekPriority")
			}
		}
		//Peak
		elements, priorities, empty := testQueue.PeekPriority()
		//Assert
		assert.Equal(t, c.oElements, elements, fmt.Sprintf("%s :Elements", cDesc))
		assert.Equal(t, c.oPriorities, priorities, fmt.Sprintf("%s :Priorities", cDesc))
		assert.Equal(t, c.oEmpty, empty, fmt.Sprintf("%s :Empty", cDesc))
	}
}

//TestPeekHead will test peek head
func TestPeekHead(t *testing.T) {
	cases := map[string]struct {
		iElements []interface{}
		oElement  interface{}
		oEmpty    bool
	}{
		"PeekHead_Empty": {
			iElements: []interface{}{},
			oElement:  nil,
			oEmpty:    true,
		},
		"PeekHead_Single_Element": {
			iElements: []interface{}{1},
			oElement:  1,
			oEmpty:    false,
		},
		"PeekHead_Multiple_Element": {
			iElements: []interface{}{1, 2, 3},
			oElement:  1,
			oEmpty:    false,
		},
	}

	//Test cases
	for cDesc, c := range cases {
		//Create Queue
		testQueue := NewQueue(len(c.iElements), true)
		defer testQueue.Close()
		//Enqueue
		for _, element := range c.iElements {
			if overflow := testQueue.Enqueue(element); overflow {
				t.Fatalf(fatalOverflow, "PeekHead")
			}
		}
		//PeakHead
		element, empty := testQueue.PeekHead()
		//Assert
		assert.Equal(t, c.oElement, element, fmt.Sprintf("%s :Element", cDesc))
		assert.Equal(t, c.oEmpty, empty, fmt.Sprintf("%s :Empty", cDesc))
	}
}

//TestPeekHeadPriority will test peek head with priority
func TestPeekHeadPriority(t *testing.T) {
	cases := map[string]struct {
		iElements   []interface{}
		iPriorities []int
		oElement    interface{}
		oPriority   int
		oEmpty      bool
	}{
		"PeekHeadPriority_Empty": {
			iElements:   []interface{}{},
			iPriorities: []int{},
			oElement:    nil,
			oPriority:   0,
			oEmpty:      true,
		},
		"PeekHeadPriority_Single_Element": {
			iElements:   []interface{}{1},
			iPriorities: []int{1},
			oElement:    1,
			oPriority:   1,
			oEmpty:      false,
		},
		"PeekHeadPriority_Multiple_Element": {
			iElements:   []interface{}{1, 2, 3},
			iPriorities: []int{1, 1, 1},
			oElement:    1,
			oPriority:   1,
			oEmpty:      false,
		},
	}

	//Test cases
	for cDesc, c := range cases {
		//Create Queue
		testQueue := NewQueue(len(c.iElements), true)
		defer testQueue.Close()
		//Enqueue
		for index, element := range c.iElements {
			if overflow := testQueue.EnqueuePriority(element, c.iPriorities[index]); overflow {
				t.Fatalf(fatalOverflow, "PeekHeadPriority")
			}
		}
		//PeakHead
		element, priority, empty := testQueue.PeekHeadPriority()
		//Assert
		assert.Equal(t, c.oElement, element, fmt.Sprintf("%s :Element", cDesc))
		assert.Equal(t, c.oPriority, priority, fmt.Sprintf("%s :Priority", cDesc))
		assert.Equal(t, c.oEmpty, empty, fmt.Sprintf("%s :Empty", cDesc))
	}
}

//TestPeekTail will test peek tail
func TestPeekTail(t *testing.T) {
	cases := map[string]struct {
		iElements []interface{}
		oElement  interface{}
		oEmpty    bool
	}{
		"PeekTail_Empty": {
			iElements: []interface{}{},
			oElement:  nil,
			oEmpty:    true,
		},
		"PeekTail_Single_Element": {
			iElements: []interface{}{1},
			oElement:  1,
			oEmpty:    false,
		},
		"PeekTail_Multiple_Element": {
			iElements: []interface{}{1, 2, 3},
			oElement:  3,
			oEmpty:    false,
		},
	}

	//Test cases
	for cDesc, c := range cases {
		//Create Queue
		testQueue := NewQueue(len(c.iElements), true)
		defer testQueue.Close()
		//Enqueue
		for _, element := range c.iElements {
			if overflow := testQueue.Enqueue(element); overflow {
				t.Fatalf(fatalOverflow, "PeekTail")
			}
		}
		//PeakTail
		element, empty := testQueue.PeekTail()
		//Assert
		assert.Equal(t, c.oElement, element, fmt.Sprintf("%s :Element", cDesc))
		assert.Equal(t, c.oEmpty, empty, fmt.Sprintf("%s :Empty", cDesc))
	}
}

//TestPeekTailPriority will test peek tail with priority
func TestPeekTailPriority(t *testing.T) {
	cases := map[string]struct {
		iElements   []interface{}
		iPriorities []int
		oElement    interface{}
		oPriority   int
		oEmpty      bool
	}{
		"PeekTailPriority_Empty": {
			iElements:   []interface{}{},
			iPriorities: []int{},
			oElement:    nil,
			oPriority:   0,
			oEmpty:      true,
		},
		"PeekTailPriority_Single_Element": {
			iElements:   []interface{}{1},
			iPriorities: []int{1},
			oElement:    1,
			oPriority:   1,
			oEmpty:      false,
		},
		"PeekTailPriority_Multiple_Element": {
			iElements:   []interface{}{1, 2, 3},
			iPriorities: []int{1, 1, 1},
			oElement:    3,
			oPriority:   1,
			oEmpty:      false,
		},
	}

	//Test cases
	for cDesc, c := range cases {
		//Create Queue
		testQueue := NewQueue(len(c.iElements), true)
		defer testQueue.Close()
		//Enqueue
		for index, element := range c.iElements {
			if overflow := testQueue.EnqueuePriority(element, c.iPriorities[index]); overflow {
				t.Fatalf(fatalOverflow, "PeekTailPriority")
			}
		}
		//PeakHead
		element, priority, empty := testQueue.PeekTailPriority()
		//Assert
		assert.Equal(t, c.oElement, element, fmt.Sprintf("%s :Element", cDesc))
		assert.Equal(t, c.oPriority, priority, fmt.Sprintf("%s :Priority", cDesc))
		assert.Equal(t, c.oEmpty, empty, fmt.Sprintf("%s :Empty", cDesc))
	}
}

//TestGetSize will test getting the size
func TestGetSize(t *testing.T) {
	cases := map[string]struct {
		iSize int
		oSize int
	}{
		"GetSize_Valid": {
			iSize: 10,
			oSize: 10,
		},
		"GetSize_Invalid": {
			iSize: -1,
			oSize: 1,
		},
	}

	//Test cases
	for cDesc, c := range cases {
		//Create Queue
		testQueue := NewQueue(c.iSize, false)
		defer testQueue.Close()
		//Get Size
		size := testQueue.GetSize()
		//Assert
		assert.Equal(t, c.oSize, size, fmt.Sprintf(cDesc))
	}
}

//TestGetLength will test get length
func TestGetLength(t *testing.T) {
	cases := map[string]struct {
		iSize     int
		iElements []interface{}
		oLength   int
	}{
		"GetLength_Empty": {
			iSize:     10,
			iElements: []interface{}{},
			oLength:   0,
		},
		"GetLength_Single_Element": {
			iSize:     10,
			iElements: []interface{}{1},
			oLength:   1,
		},
		"GetLength_Multiple_Element": {
			iSize:     10,
			iElements: []interface{}{1, 2, 3},
			oLength:   3,
		},
	}

	//Test cases
	for cDesc, c := range cases {
		//Create Queue
		testQueue := NewQueue(c.iSize, true)
		defer testQueue.Close()
		//Enqueue
		for _, element := range c.iElements {
			if overflow := testQueue.Enqueue(element); overflow {
				t.Fatalf(fatalOverflow, "GetLength")
			}
		}
		//Get Length
		length := testQueue.GetLength()
		//Assert
		assert.Equal(t, c.oLength, length, fmt.Sprintf("%s :Length", cDesc))
	}
}

//TestResize will test the resize
func TestResize(t *testing.T) {
	cases := map[string]struct {
		iInitialSize int
		iFinalSize   int
		oInitialSize int
		oFinalSize   int
	}{
		"Resize_Valid": {
			iInitialSize: 10,
			iFinalSize:   100,
			oInitialSize: 10,
			oFinalSize:   100,
		},
	}
	//Test cases
	for cDesc, c := range cases {
		//Create Queue
		testQueue := NewQueue(c.iInitialSize, false)
		defer testQueue.Close()
		//Get Size
		initialSize := testQueue.GetSize()
		//Assert
		assert.Equal(t, c.oInitialSize, initialSize, fmt.Sprintf(cDesc))
		//Resize
		testQueue.Resize(c.iFinalSize)
		//Get Size
		finalSize := testQueue.GetSize()
		//Assert
		assert.Equal(t, c.oFinalSize, finalSize, fmt.Sprintf(cDesc))
	}
}
