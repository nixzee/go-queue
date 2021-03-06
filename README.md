# go-queue

[![Test Status](https://github.com/nixzee/go-queue/workflows/Test/badge.svg)](https://github.com/nixzee/go-queue/actions)
[![Coverage Status](https://coveralls.io/repos/github/nixzee/go-queue/badge.svg?branch=master)](https://coveralls.io/github/nixzee/go-queue?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/nixzee/go-queue/actions)](https://github.com/nixzee/go-queue)

`go-queue` is a Finite Synchronous Priority Queue GO package. Let us break that apart:

## Finite

The queue has a finite size. This means if the queue size is set 10 elements it will have room for only 10 elements. Any attempt to enqueue more will cause an overflow flag to be set. Any attempt to dequeue beyond empty will cause an underflow. There are few reasons to use an infinite/lossless queue.

## Synchronous

One of my biggest gripes with other GO queue packages is that they require the developer to poll the dequeue. This means either using dedicated sleep or ticker/timer. You can still do this with this package, but I am providing you with a channel that you can subscribe to. This makes this queue (and your code) event driven.

## Priority

Elements inserted into the queue can be given priority. The higher the number, the higher priority. Elements with higher priority will bumped up in the queue until it reaches the end or finds and element of the same or higher priority. The queue will maintain order like any FIFO would.

## Install

`go get github.com/nixzee/go-queue`
