# go-queue

`go-queue` is a Finite Synchronous Queue GO package. Lets break that apart:

## Finite

The queue has a finite size. This means if the queue size is set to lets say 10 elements, it will have room for 10 elements. Any attempt to enqueue more will cause an overflow flag to be set. Any attempt to dequeue beyond empty will cause an underflow.

## Synchronous

One of my biggest gripes with other GO queue packages is that they require the developer to poll the dequeue. This means either using dedicated sleep or ticker/timer. You can still do this with this package but I am providing you with a channel that you can subscribe to. This makes this queue (and your code) event driven. 

## Install


