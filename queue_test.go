package queue_test

import (
	"fmt"
	"go-queue"
	"testing"

	"github.com/stretchr/testify/assert"
)

//---------------------------------------------------------------------------------------------------
// Unit Tests
//---------------------------------------------------------------------------------------------------

func TestClose(t *testing.T) {
	//Create Queue
	testQueue := queue.NewQueue(1)
	//Close
	testQueue.Close()
}

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
		testQueue := queue.NewQueue(1)
		defer testQueue.Close()
		//Get Signal
		signal := testQueue.GetSignal()
		//Assert
		assert.NotEqual(t, c.oNotSignal, signal, fmt.Sprintf(cDesc))
	}
}

// func TestDequeue(t *testing.T) {
// 	cases := map[string]struct {
// 		iSize      int
// 		iElemennts []int
// 		oUnderflow bool
// 		oElements  []int
// 	}{
// 		"Normal": {
// 			iSize:      10,
// 			iElemennts: []int{1, 2, 3},
// 			oUnderflow: false,
// 			oElements:  []int{1, 2, 3},
// 		},
// 		"Overflow": {
// 			iSize:      2,
// 			iElemennts: []int{1, 2, 3},
// 			oUnderflow: true,
// 			oElements:  []int{1, 2},
// 		},
// 	}

// 	//Test cases
// 	for cDesc, c := range cases {
// 		//Create Queue
// 		testQueue := queue.NewQueue(c.iSize)
// 		defer testQueue.Close()
// 		//enqueue

// 	}
// }

// func TestDequeueMultiple(t *testing.T) {
// 	cases := map[string]struct{}{}

// 	//Test cases
// 	for cDesc, c := range cases {
// 	}
// }

func TestEnqueue(t *testing.T) {
	cases := map[string]struct {
		iSize     int
		iElements []interface{}
		oOverflow bool
	}{
		"Enqueue_Empty": {
			iSize:     3,
			iElements: []interface{}{},
			oOverflow: false,
		},
		"Enqueue_Nomrmal": {
			iSize:     3,
			iElements: []interface{}{1, 2},
			oOverflow: false,
		},
		"Enqueue_Full": {
			iSize:     3,
			iElements: []interface{}{1, 2, 3},
			oOverflow: false,
		},
		"Enqueue_Over": {
			iSize:     3,
			iElements: []interface{}{1, 2, 3, 4},
			oOverflow: true,
		},
	}

	//Test cases
	for cDesc, c := range cases {
		//Create Queue
		testQueue := queue.NewQueue(c.iSize)
		defer testQueue.Close()
		//Enqueue
		overflow := false
		for _, element := range c.iElements {
			if overflow = testQueue.Enqueue(element); overflow {
				break
			}
		}
		//Assert
		assert.Equal(t, c.oOverflow, overflow, fmt.Sprintf("%s :Overflow", cDesc))
	}
}

func TestEnqueueMultiple(t *testing.T) {
	cases := map[string]struct {
		iSize     int
		iElements []interface{}
		oOverflow bool
	}{
		"EnqueueMultiple_Empty": {
			iSize:     3,
			iElements: []interface{}{},
			oOverflow: false,
		},
		"EnqueueMultiple_Nomrmal": {
			iSize:     3,
			iElements: []interface{}{1, 2},
			oOverflow: false,
		},
		"EnqueueMultiple_Full": {
			iSize:     3,
			iElements: []interface{}{1, 2, 3},
			oOverflow: false,
		},
		"EnqueueMultiple_Over": {
			iSize:     3,
			iElements: []interface{}{1, 2, 3, 4},
			oOverflow: true,
		},
	}

	//Test cases
	for cDesc, c := range cases {
		//Create Queue
		testQueue := queue.NewQueue(c.iSize)
		defer testQueue.Close()
		//Enqueue
		overflow := false
		if overflow = testQueue.EnqueueMultiple(c.iElements); overflow {
			break
		}
		//Assert
		assert.Equal(t, c.oOverflow, overflow, fmt.Sprintf("%s :Overflow", cDesc))
	}
}

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
		testQueue := queue.NewQueue(len(c.iElements))
		defer testQueue.Close()
		//Enqueue
		testQueue.EnqueueMultiple(c.iElements)
		//Flush
		testQueue.Flush()
		//GetIndex
		length := testQueue.GetLength()
		//Assert
		assert.Equal(t, c.oLength, length, fmt.Sprintf(cDesc))
	}
}

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
		"Peek_SingleElement": {
			iElements: []interface{}{1},
			oElements: []interface{}{1},
			oEmpty:    false,
		},
		"Peek_MultipleElement": {
			iElements: []interface{}{1, 2, 3},
			oElements: []interface{}{1, 2, 3},
			oEmpty:    false,
		},
	}

	//Test cases
	for cDesc, c := range cases {
		//Create Queue
		testQueue := queue.NewQueue(len(c.iElements))
		defer testQueue.Close()
		//Enqueue
		for _, element := range c.iElements {
			if overflow := testQueue.Enqueue(element); overflow {
				//TODO: Fatal
			}
		}
		//Peak
		elements, empty := testQueue.Peek()
		//Assert
		assert.Equal(t, c.oElements, elements, fmt.Sprintf("%s :Elements", cDesc))
		assert.Equal(t, c.oEmpty, empty, fmt.Sprintf("%s :Empty", cDesc))
	}
}

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
		"PeekHead_SingleElement": {
			iElements: []interface{}{1},
			oElement:  1,
			oEmpty:    false,
		},
		"PeekHead_MultipleElement": {
			iElements: []interface{}{1, 2, 3},
			oElement:  1,
			oEmpty:    false,
		},
	}

	//Test cases
	for cDesc, c := range cases {
		//Create Queue
		testQueue := queue.NewQueue(len(c.iElements))
		defer testQueue.Close()
		//Enqueue
		for _, element := range c.iElements {
			if overflow := testQueue.Enqueue(element); overflow {
				//TODO: Fatal
			}
		}
		//PeakHead
		element, empty := testQueue.PeekHead()
		//Assert
		assert.Equal(t, c.oElement, element, fmt.Sprintf("%s :Element", cDesc))
		assert.Equal(t, c.oEmpty, empty, fmt.Sprintf("%s :Empty", cDesc))
	}
}

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
		"PeekTail_SingleElement": {
			iElements: []interface{}{1},
			oElement:  1,
			oEmpty:    false,
		},
		"PeekTail_MultipleElement": {
			iElements: []interface{}{1, 2, 3},
			oElement:  3,
			oEmpty:    false,
		},
	}

	//Test cases
	for cDesc, c := range cases {
		//Create Queue
		testQueue := queue.NewQueue(len(c.iElements))
		defer testQueue.Close()
		//Enqueue
		for _, element := range c.iElements {
			if overflow := testQueue.Enqueue(element); overflow {
				//TODO: Fatal
			}
		}
		//PeakTail
		element, empty := testQueue.PeekTail()
		//Assert
		assert.Equal(t, c.oElement, element, fmt.Sprintf("%s :Element", cDesc))
		assert.Equal(t, c.oEmpty, empty, fmt.Sprintf("%s :Empty", cDesc))
	}
}

func TestGetSize(t *testing.T) {
	cases := map[string]struct {
		iSize int
		oSize int
	}{
		"GetSize": {
			iSize: 10,
			oSize: 10,
		},
	}

	//Test cases
	for cDesc, c := range cases {
		//Create Queue
		testQueue := queue.NewQueue(c.iSize)
		defer testQueue.Close()
		//Get Size
		size := testQueue.GetSize()
		//Assert
		assert.Equal(t, c.oSize, size, fmt.Sprintf(cDesc))
	}
}

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
		"GetLength_SingleElement": {
			iSize:     10,
			iElements: []interface{}{1},
			oLength:   1,
		},
		"GetLength_MultipleElement": {
			iSize:     10,
			iElements: []interface{}{1, 2, 3},
			oLength:   3,
		},
	}

	//Test cases
	for cDesc, c := range cases {
		//Create Queue
		testQueue := queue.NewQueue(c.iSize)
		defer testQueue.Close()
		//Enqueue
		for _, element := range c.iElements {
			if overflow := testQueue.Enqueue(element); overflow {
				//TODO: Fatal
			}
		}
		//Get Length
		length := testQueue.GetLength()
		//Assert
		assert.Equal(t, c.oLength, length, fmt.Sprintf("%s :Length", cDesc))
	}
}
