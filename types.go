package queue

//---------------------------------------------------------------------------------------------
// Generics
//---------------------------------------------------------------------------------------------

const (
	//DefaultSize is the smallest possible queue size
	DefaultSize int = 1
	//DefaultPriority defines the priority when none is given
	DefaultPriority int = 0
)

//---------------------------------------------------------------------------------------------
// Element Container
//---------------------------------------------------------------------------------------------

//container is a single Element Container
//This provides a heapable container for the priority queue to use heap
type container struct {
	element  interface{}
	priority int
}
