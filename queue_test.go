package queue

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	//Errors
	fatalUnderflow string = "%s :Encountered underflow"
	fatalOverflow  string = "%s :Encountered overflow"
	//Formatting
	fmtAssertMsg string = "%s_%s"
)

//assertMsg provides a dirty wrapper for assert messages
func assertMsg(name, desc string) string {
	return fmt.Sprintf(fmtAssertMsg, name, desc)
}

//---------------------------------------------------------------------------------------------------
// Owner
//---------------------------------------------------------------------------------------------------

//TestClose will test close
func TestClose(t *testing.T) {
	//Create Queue
	testQueue := NewQueue(1, false)
	//Close
	testQueue.Close()
}

//TestResize will test the resize
func TestResize(t *testing.T) {
	const name string = "Resize"
	cases := map[string]struct {
		iInitialSize int
		iFinalSize   int
		iElements    []interface{}
		iPriorities  []int
		oInitialSize int
		oFinalSize   int
		oElements    []interface{}
		oPriorities  []int
	}{
		"Valid_Size": {
			iInitialSize: 10,
			iFinalSize:   100,
			iElements:    []interface{}{},
			iPriorities:  []int{},
			oInitialSize: 10,
			oFinalSize:   100,
			oElements:    nil,
			oPriorities:  nil,
		},
		"Invalid_Size": {
			iInitialSize: 10,
			iFinalSize:   -1,
			iElements:    []interface{}{},
			iPriorities:  []int{},
			oInitialSize: 10,
			oFinalSize:   1,
			oElements:    nil,
			oPriorities:  nil,
		},
		"Valid_Size_With_Content": {
			iInitialSize: 10,
			iFinalSize:   100,
			iElements:    []interface{}{1, 2, 3},
			iPriorities:  []int{1, 2, 3},
			oInitialSize: 10,
			oFinalSize:   100,
			oElements:    []interface{}{3, 2, 1},
			oPriorities:  []int{3, 2, 1},
		},
	}
	//Test cases
	for cDesc, c := range cases {
		//Get the assert message base
		msg := assertMsg(name, cDesc)
		//Create Queue
		testQueue := NewQueue(c.iInitialSize, false)
		defer testQueue.Close()
		//Get Size
		initialSize := testQueue.GetSize()
		//Assert
		assert.Equal(t, c.oInitialSize, initialSize, fmt.Sprintf("%s :initial size", msg))
		//Enqueue
		for index, element := range c.iElements {
			if overflow := testQueue.EnqueuePriority(element, c.iPriorities[index]); overflow {
				t.Fatalf(fatalOverflow, msg)
			}
		}
		//Resize
		elements, priorities := testQueue.Resize(c.iFinalSize)
		//Get Size
		finalSize := testQueue.GetSize()
		//Assert
		assert.Equal(t, c.oFinalSize, finalSize, fmt.Sprintf("%s :final size", msg))
		assert.Equal(t, c.oElements, elements, fmt.Sprintf("%s :elements", msg))
		assert.Equal(t, c.oPriorities, priorities, fmt.Sprintf("%s :priorities", msg))
	}
}

//---------------------------------------------------------------------------------------------------
// Flush
//---------------------------------------------------------------------------------------------------

//TestFlush will test the flush
func TestFlush(t *testing.T) {
	const name string = "Flush"
	cases := map[string]struct {
		iElements []interface{}
		oLength   int
	}{
		"With_Content": {
			iElements: []interface{}{1, 2},
			oLength:   0,
		},
	}

	//Test cases
	for cDesc, c := range cases {
		//Get the assert message base
		msg := assertMsg(name, cDesc)
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
		assert.Equal(t, c.oLength, length, msg)
	}
}

//---------------------------------------------------------------------------------------------------
// Event
//---------------------------------------------------------------------------------------------------

//TestGetSignal will test getting the signal
func TestGetSignal(t *testing.T) {
	const name string = "GetSignal"
	cases := map[string]struct {
		oSignal chan struct{}
	}{
		"Valid": {
			oSignal: nil,
		},
	}

	//Test cases
	for cDesc, c := range cases {
		//Get the assert message base
		msg := assertMsg(name, cDesc)
		//Create Queue
		testQueue := NewQueue(1, false)
		defer testQueue.Close()
		//Get Signal
		signal := testQueue.GetSignal()
		//Assert
		assert.NotEqual(t, c.oSignal, signal, msg)
	}
}

//---------------------------------------------------------------------------------------------------
// Info
//---------------------------------------------------------------------------------------------------

//TestGetSize will test getting the size
func TestGetSize(t *testing.T) {
	const name string = "GetSize"
	cases := map[string]struct {
		iSize int
		oSize int
	}{
		"Valid": {
			iSize: 10,
			oSize: 10,
		},
		"Invalid": {
			iSize: -1,
			oSize: 1,
		},
	}

	//Test cases
	for cDesc, c := range cases {
		//Get the assert message base
		msg := assertMsg(name, cDesc)
		//Create Queue
		testQueue := NewQueue(c.iSize, false)
		defer testQueue.Close()
		//Get Size
		size := testQueue.GetSize()
		//Assert
		assert.Equal(t, c.oSize, size, msg)
	}
}

//TestGetLength will test get length
func TestGetLength(t *testing.T) {
	const name string = "GetLength"
	cases := map[string]struct {
		iSize     int
		iElements []interface{}
		oLength   int
	}{
		"Empty": {
			iSize:     10,
			iElements: []interface{}{},
			oLength:   0,
		},
		"Single_Element": {
			iSize:     10,
			iElements: []interface{}{1},
			oLength:   1,
		},
		"Multiple_Element": {
			iSize:     10,
			iElements: []interface{}{1, 2, 3},
			oLength:   3,
		},
	}

	//Test cases
	for cDesc, c := range cases {
		//Get the assert message base
		msg := assertMsg(name, cDesc)
		//Create Queue
		testQueue := NewQueue(c.iSize, true)
		defer testQueue.Close()
		//Enqueue
		for _, element := range c.iElements {
			if overflow := testQueue.Enqueue(element); overflow {
				t.Fatalf(fatalOverflow, msg)
			}
		}
		//Get Length
		length := testQueue.GetLength()
		//Assert
		assert.Equal(t, c.oLength, length, fmt.Sprintf("%s :Length", msg))
	}
}

//---------------------------------------------------------------------------------------------------
// Peek
//---------------------------------------------------------------------------------------------------

//TestPeek will test the peek
const name string = "Peek"

func TestPeek(t *testing.T) {
	cases := map[string]struct {
		iElements []interface{}
		oElements []interface{}
		oEmpty    bool
	}{
		"Empty": {
			iElements: []interface{}{},
			oElements: nil,
			oEmpty:    true,
		},
		"Single_Element": {
			iElements: []interface{}{1},
			oElements: []interface{}{1},
			oEmpty:    false,
		},
		"Multiple_Element": {
			iElements: []interface{}{1, 2, 3},
			oElements: []interface{}{1, 2, 3},
			oEmpty:    false,
		},
	}

	//Test cases
	for cDesc, c := range cases {
		//Get the assert message base
		msg := assertMsg(name, cDesc)
		//Create Queue
		testQueue := NewQueue(len(c.iElements), true)
		defer testQueue.Close()
		//Enqueue
		for _, element := range c.iElements {
			if overflow := testQueue.Enqueue(element); overflow {
				t.Fatalf(fatalOverflow, msg)
			}
		}
		//Peak
		elements, empty := testQueue.Peek()
		//Assert
		assert.Equal(t, c.oElements, elements, fmt.Sprintf("%s :Elements", msg))
		assert.Equal(t, c.oEmpty, empty, fmt.Sprintf("%s :Empty", msg))
	}
}

//TestPeekHead will test peek head
func TestPeekHead(t *testing.T) {
	const name string = "PeekHead"
	cases := map[string]struct {
		iElements []interface{}
		oElement  interface{}
		oEmpty    bool
	}{
		"Empty": {
			iElements: []interface{}{},
			oElement:  nil,
			oEmpty:    true,
		},
		"Single_Element": {
			iElements: []interface{}{1},
			oElement:  1,
			oEmpty:    false,
		},
		"Multiple_Element": {
			iElements: []interface{}{1, 2, 3},
			oElement:  1,
			oEmpty:    false,
		},
	}

	//Test cases
	for cDesc, c := range cases {
		//Get the assert message base
		msg := assertMsg(name, cDesc)
		//Create Queue
		testQueue := NewQueue(len(c.iElements), true)
		defer testQueue.Close()
		//Enqueue
		for _, element := range c.iElements {
			if overflow := testQueue.Enqueue(element); overflow {
				t.Fatalf(fatalOverflow, msg)
			}
		}
		//PeakHead
		element, empty := testQueue.PeekHead()
		//Assert
		assert.Equal(t, c.oElement, element, fmt.Sprintf("%s :Element", msg))
		assert.Equal(t, c.oEmpty, empty, fmt.Sprintf("%s :Empty", msg))
	}
}

//TestPeekTail will test peek tail
func TestPeekTail(t *testing.T) {
	const name string = "PeekTail"
	cases := map[string]struct {
		iElements []interface{}
		oElement  interface{}
		oEmpty    bool
	}{
		"Empty": {
			iElements: []interface{}{},
			oElement:  nil,
			oEmpty:    true,
		},
		"Single_Element": {
			iElements: []interface{}{1},
			oElement:  1,
			oEmpty:    false,
		},
		"Multiple_Element": {
			iElements: []interface{}{1, 2, 3},
			oElement:  3,
			oEmpty:    false,
		},
	}

	//Test cases
	for cDesc, c := range cases {
		//Get the assert message base
		msg := assertMsg(name, cDesc)
		//Create Queue
		testQueue := NewQueue(len(c.iElements), true)
		defer testQueue.Close()
		//Enqueue
		for _, element := range c.iElements {
			if overflow := testQueue.Enqueue(element); overflow {
				t.Fatalf(fatalOverflow, msg)
			}
		}
		//PeakTail
		element, empty := testQueue.PeekTail()
		//Assert
		assert.Equal(t, c.oElement, element, fmt.Sprintf("%s :Element", msg))
		assert.Equal(t, c.oEmpty, empty, fmt.Sprintf("%s :Empty", msg))
	}
}

//---------------------------------------------------------------------------------------------------
// Peek Priority
//---------------------------------------------------------------------------------------------------

//TestPeek will test the peek with priority
func TestPeekPriority(t *testing.T) {
	const name string = "PeekPriority"
	cases := map[string]struct {
		iElements   []interface{}
		iPriorities []int
		oElements   []interface{}
		oPriorities []int
		oEmpty      bool
	}{
		"Empty": {
			iElements:   []interface{}{},
			iPriorities: []int{},
			oElements:   nil,
			oPriorities: nil,
			oEmpty:      true,
		},
		"Single_Element": {
			iElements:   []interface{}{1},
			iPriorities: []int{0},
			oElements:   []interface{}{1},
			oPriorities: []int{0},
			oEmpty:      false,
		},
		"Multiple_Element": {
			iElements:   []interface{}{1, 2, 3},
			iPriorities: []int{0, 0, 0},
			oElements:   []interface{}{1, 2, 3},
			oPriorities: []int{0, 0, 0},
			oEmpty:      false,
		},
	}

	//Test cases
	for cDesc, c := range cases {
		//Get the assert message base
		msg := assertMsg(name, cDesc)
		//Create Queue
		testQueue := NewQueue(len(c.iElements), true)
		defer testQueue.Close()
		//Enqueue
		for index, element := range c.iElements {
			if overflow := testQueue.EnqueuePriority(element, c.iPriorities[index]); overflow {
				t.Fatalf(fatalOverflow, msg)
			}
		}
		//Peak
		elements, priorities, empty := testQueue.PeekPriority()
		//Assert
		assert.Equal(t, c.oElements, elements, fmt.Sprintf("%s :Elements", msg))
		assert.Equal(t, c.oPriorities, priorities, fmt.Sprintf("%s :Priorities", msg))
		assert.Equal(t, c.oEmpty, empty, fmt.Sprintf("%s :Empty", msg))
	}
}

//TestPeekHeadPriority will test peek head with priority
func TestPeekHeadPriority(t *testing.T) {
	const name string = "PeekHeadPriority"
	cases := map[string]struct {
		iElements   []interface{}
		iPriorities []int
		oElement    interface{}
		oPriority   int
		oEmpty      bool
	}{
		"Empty": {
			iElements:   []interface{}{},
			iPriorities: []int{},
			oElement:    nil,
			oPriority:   0,
			oEmpty:      true,
		},
		"Single_Element": {
			iElements:   []interface{}{1},
			iPriorities: []int{1},
			oElement:    1,
			oPriority:   1,
			oEmpty:      false,
		},
		"Multiple_Element": {
			iElements:   []interface{}{1, 2, 3},
			iPriorities: []int{1, 1, 1},
			oElement:    1,
			oPriority:   1,
			oEmpty:      false,
		},
	}

	//Test cases
	for cDesc, c := range cases {
		//Get the assert message base
		msg := assertMsg(name, cDesc)
		//Create Queue
		testQueue := NewQueue(len(c.iElements), true)
		defer testQueue.Close()
		//Enqueue
		for index, element := range c.iElements {
			if overflow := testQueue.EnqueuePriority(element, c.iPriorities[index]); overflow {
				t.Fatalf(fatalOverflow, msg)
			}
		}
		//PeakHead
		element, priority, empty := testQueue.PeekHeadPriority()
		//Assert
		assert.Equal(t, c.oElement, element, fmt.Sprintf("%s :Element", msg))
		assert.Equal(t, c.oPriority, priority, fmt.Sprintf("%s :Priority", msg))
		assert.Equal(t, c.oEmpty, empty, fmt.Sprintf("%s :Empty", msg))
	}
}

//TestPeekTailPriority will test peek tail with priority
func TestPeekTailPriority(t *testing.T) {
	const name string = "PeekTailPriority"
	cases := map[string]struct {
		iElements   []interface{}
		iPriorities []int
		oElement    interface{}
		oPriority   int
		oEmpty      bool
	}{
		"Empty": {
			iElements:   []interface{}{},
			iPriorities: []int{},
			oElement:    nil,
			oPriority:   0,
			oEmpty:      true,
		},
		"Single_Element": {
			iElements:   []interface{}{1},
			iPriorities: []int{1},
			oElement:    1,
			oPriority:   1,
			oEmpty:      false,
		},
		"Multiple_Element": {
			iElements:   []interface{}{1, 2, 3},
			iPriorities: []int{1, 1, 1},
			oElement:    3,
			oPriority:   1,
			oEmpty:      false,
		},
	}

	//Test cases
	for cDesc, c := range cases {
		//Get the assert message base
		msg := assertMsg(name, cDesc)
		//Create Queue
		testQueue := NewQueue(len(c.iElements), true)
		defer testQueue.Close()
		//Enqueue
		for index, element := range c.iElements {
			if overflow := testQueue.EnqueuePriority(element, c.iPriorities[index]); overflow {
				t.Fatalf(fatalOverflow, msg)
			}
		}
		//PeakHead
		element, priority, empty := testQueue.PeekTailPriority()
		//Assert
		assert.Equal(t, c.oElement, element, fmt.Sprintf("%s :Element", msg))
		assert.Equal(t, c.oPriority, priority, fmt.Sprintf("%s :Priority", msg))
		assert.Equal(t, c.oEmpty, empty, fmt.Sprintf("%s :Empty", msg))
	}
}

//---------------------------------------------------------------------------------------------------
// Dequeue / Enqueue
//---------------------------------------------------------------------------------------------------

//TestDequeueConcurrent will test dequeue using go routines (concurrency)
func TestDequeueConcurrent(t *testing.T) {
	const name string = "DequeueConcurrent"
	cases := map[string]struct {
		iSize     int
		iElements []interface{}
		oElements []interface{}
	}{
		"Single": {
			iSize:     10,
			iElements: []interface{}{1},
			oElements: []interface{}{1},
		},
		"Multiple": {
			iSize:     10,
			iElements: []interface{}{1, 2, 3},
			oElements: []interface{}{1, 2, 3},
		},
	}

	//Test cases
	for cDesc, c := range cases {
		//Get the assert message base
		msg := assertMsg(name, cDesc)
		var wg sync.WaitGroup
		stop := make(chan struct{})
		var elements []interface{}
		//Create Queue
		testQueue := NewQueue(c.iSize, false)
		defer testQueue.Close()
		//dequeue routine
		underflow, overflow := false, false
		wg.Add(1)
		go func(t *testing.T) {
			defer wg.Done()
			signal := testQueue.GetSignal()
			for {
				select {
				case <-stop:
					return
				case <-signal:
					var element interface{}
					//Dequeue
					if element, underflow = testQueue.Dequeue(); underflow {
						return
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
				if overflow = testQueue.Enqueue(element); overflow {
					return
				}
				time.Sleep(100 * time.Millisecond)
			}
			//stop the dequeue
			close(stop)
		}(t)
		//Wait
		wg.Wait()
		//Check for underflow and overflow
		if underflow {
			t.Fatalf(fatalUnderflow, msg)
		}
		if overflow {
			t.Fatalf(fatalOverflow, msg)
		}
		//Assert
		assert.Equal(t, c.oElements, elements, fmt.Sprintf("%s :Elements", msg))
	}
}

//TestDequeueSerialized will test dequeue serialized
func TestDequeueSerialized(t *testing.T) {
	const name string = "DequeueSerialized"
	cases := map[string]struct {
		iSize     int
		iElements []interface{}
		oElements []interface{}
	}{
		"Multiple": {
			iSize:     10,
			iElements: []interface{}{1, 2, 3, 4, 5, 6},
			oElements: []interface{}{1, 2, 3, 4, 5, 6},
		},
	}

	//Test cases
	for cDesc, c := range cases {
		//Get the assert message base
		msg := assertMsg(name, cDesc)
		var elements []interface{}
		//Create Queue
		testQueue := NewQueue(c.iSize, true)
		defer testQueue.Close()
		//enqueue
		for _, element := range c.iElements {
			if overflow := testQueue.Enqueue(element); overflow {
				t.Fatalf(fatalOverflow, msg)
			}
		}
		//dequeue
		for range c.iElements {
			element, underflow := testQueue.Dequeue()
			//Check underflow
			if underflow {
				t.Fatalf(fatalUnderflow, msg)
			}
			//Append
			elements = append(elements, element)
		}
		//Assert
		assert.Equal(t, c.oElements, elements, fmt.Sprintf("%s :Elements", msg))
	}
}

func TestUnderflow(t *testing.T) {
	const name string = "Underflow"
	cases := map[string]struct {
		iSize     int
		iPolling  bool
		oUndeflow bool
	}{
		"Valid": {
			1,
			true,
			true,
		},
	}

	//Test cases
	for cDesc, c := range cases {
		//Get the assert message base
		msg := assertMsg(name, cDesc)
		//Create Queue
		testQueue := NewQueue(c.iSize, c.iPolling)
		defer testQueue.Close()
		//Dequeue
		var underflow bool
		for i := 0; i < c.iSize; i++ {
			_, underflow = testQueue.Dequeue()
		}

		//Assert
		assert.Equal(t, underflow, c.oUndeflow, msg)
	}
}

func TestOverflow(t *testing.T) {
	const name string = "Overflow"
	cases := map[string]struct {
		iSize     int
		iPolling  bool
		oOverflow bool
	}{
		"Valid": {
			1,
			true,
			true,
		},
	}

	//Test cases
	for cDesc, c := range cases {
		//Get the assert message base
		msg := assertMsg(name, cDesc)
		//Create Queue
		testQueue := NewQueue(c.iSize, c.iPolling)
		defer testQueue.Close()
		//Enqueue
		var overflow bool
		for i := 0; i < c.iSize+1; i++ { //Plus one is important...
			overflow = testQueue.Enqueue(struct{}{})
		}
		//Assert
		assert.Equal(t, overflow, c.oOverflow, msg)
	}
}

//---------------------------------------------------------------------------------------------------
// Dequeue / Enqueue Priority
//---------------------------------------------------------------------------------------------------

//TestDequeuePriorityConcurrent will test dequeue priority using go routines (concurrency)
func TestDequeuePriorityConcurrent(t *testing.T) {
	const name string = "DequeuePriorityConcurrent"
	cases := map[string]struct {
		iSize       int
		iElements   []interface{}
		iPriorities []int
		oElements   []interface{}
		oPriorities []int
	}{
		"Single": {
			iSize:       10,
			iElements:   []interface{}{1},
			iPriorities: []int{10},
			oElements:   []interface{}{1},
			oPriorities: []int{10},
		},
		"Multiple": {
			iSize:       10,
			iElements:   []interface{}{1, 2, 3},
			iPriorities: []int{10, 0, 100},
			oElements:   []interface{}{1, 2, 3},
			oPriorities: []int{10, 0, 100},
		},
	}

	//Test cases
	for cDesc, c := range cases {
		//Get the assert message base
		msg := assertMsg(name, cDesc)
		var wg sync.WaitGroup
		stop := make(chan struct{})
		var elements []interface{}
		var priorities []int
		//Create Queue
		testQueue := NewQueue(c.iSize, false)
		defer testQueue.Close()
		//dequeue routine
		underflow, overflow := false, false
		wg.Add(1)
		go func(t *testing.T) {
			defer wg.Done()
			signal := testQueue.GetSignal()
			for {
				select {
				case <-stop:
					return
				case <-signal:
					var element interface{}
					var priority int
					//Dequeue
					if element, priority, underflow = testQueue.DequeuePriority(); underflow {
						return
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
				if overflow = testQueue.EnqueuePriority(element, c.iPriorities[index]); overflow {
					return
				}
				time.Sleep(100 * time.Millisecond)
			}
			//stop the dequeue
			close(stop)
		}(t)
		//Wait
		wg.Wait()
		//Check for underflow and overflow
		if underflow {
			t.Fatalf(fatalUnderflow, msg)
		}
		if overflow {
			t.Fatalf(fatalOverflow, msg)
		}
		//Assert
		assert.Equal(t, c.oElements, elements, fmt.Sprintf("%s :Elements", msg))
		assert.Equal(t, c.oPriorities, priorities, fmt.Sprintf("%s :Priorities", msg))
	}
}

//TestDequeuePrioritySerialized will test dequeue priority serialized
func TestDequeuePrioritySerialized(t *testing.T) {
	const name string = "DequeuePrioritySerialized"
	cases := map[string]struct {
		iSize       int
		iElements   []interface{}
		iPriorities []int
		oElements   []interface{}
		oPriorities []int
	}{
		"Multiple_Random": {
			iSize:       10,
			iElements:   []interface{}{1, 2, 3},
			iPriorities: []int{10, 0, 100},
			oElements:   []interface{}{3, 1, 2},
			oPriorities: []int{100, 10, 0},
		},
	}

	//Test cases
	for cDesc, c := range cases {
		//Get the assert message base
		msg := assertMsg(name, cDesc)
		var elements []interface{}
		var priorities []int
		//Create Queue
		testQueue := NewQueue(c.iSize, true)
		defer testQueue.Close()
		//enqueue
		for index, element := range c.iElements {
			if overflow := testQueue.EnqueuePriority(element, c.iPriorities[index]); overflow {
				t.Fatalf(fatalOverflow, msg)
			}
		}
		//dequeue
		for range c.iElements {
			element, priority, underflow := testQueue.DequeuePriority()
			//Check underflow
			if underflow {
				t.Fatalf(fatalUnderflow, msg)
			}
			//Append
			elements = append(elements, element)
			priorities = append(priorities, priority)
		}
		//Assert
		assert.Equal(t, c.oElements, elements, fmt.Sprintf("%s :Elements", msg))
		assert.Equal(t, c.oPriorities, priorities, fmt.Sprintf("%s :Priorities", msg))
	}
}
