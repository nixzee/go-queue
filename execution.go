package queue

//triggerSignal
func triggerSignal(signal chan struct{}) {
	//TODO: Non-block maybe?
	signal <- struct{}{}
}

//shift
func shift(dataIn []interface{}) (dataOut []interface{}) {
	l := len(dataIn)
	dataOut = append(dataIn[l-1:], dataIn[:l-1]...)
	return
}

//enqueue
func enqueue(lenIn, size int, dataIn []interface{}, element interface{}) (overflow bool, lenOut int, dataOut []interface{}) {
	//Check if queue is full (overflow)
	if lenIn >= size {
		overflow = true
		lenOut, dataOut = lenIn, dataIn
		return
	}
	//shift data
	dataOut = shift(dataIn)
	//insert element (use zero)
	dataOut[0] = element
	//increment len
	lenOut = lenIn + 1
	return
}

//dequeue
func dequeue(lenIn, size int, dataIn []interface{}) (underflow bool, lenOut int, dataOut []interface{}, element interface{}) {
	//Check if queue is empty (underflow)
	if lenIn == 0 {
		underflow = true
		lenOut, dataOut = lenIn, dataIn
		return
	}
	//Get the last value in
	element = dataIn[lenIn-1]
	//clear element from data
	dataOut[lenIn-1] = nil
	//decrement len
	lenOut = lenIn - 1
	return
}

//checkifEmpty
func checkifEmpty() (empty bool) {
	return
}

//checkifFull
func checkifFull() (full bool) {
	return
}
