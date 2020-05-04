package queue_test

//---------------------------------------------------------------------------------------------------
// Unit Tests
//---------------------------------------------------------------------------------------------------

// Close()
// GetSignal() (signal <-chan struct{})
// Dequeue() (element interface{}, underflow bool)
// DequeueMultiple(count int) (elements []interface{}, underflow bool)
// Enqueue(element interface{}) (overflow bool)
// EnqueueMultiple(elements []interface{}) (overflow bool)
// Flush() (elements []interface{})
// Peek() (elements []interface{})
// PeekHead() (element interface{}, underflow bool)
// PeekTail() (element interface{}, underflow bool)
// GetSize() (size int)
// GetLength() (len int)
