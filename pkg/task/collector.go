package task

import "github.com/vlourme/zsched/pkg/ctx"

// TaskCollectorAction is a function that performs the task collector action
type TaskCollectorAction func(collector ctx.Collector)

// Collector is an interface to interact with a collector, an object that
// collects data from the task and returns it outside its scope, useful
// for preparing batch operations.
type Collector struct {
	channel chan any
}

// NewCollector creates a new collector
// The buffer size is the size of the channel, if not provided, the channel is unbuffered.
func NewCollector(bufferSize ...int) *Collector {
	var channel chan any
	if len(bufferSize) > 0 {
		channel = make(chan any, bufferSize[0])
	} else {
		channel = make(chan any)
	}

	return &Collector{
		channel: channel,
	}
}

// Push pushes a value to the collector
func (c *Collector) Push(value any) {
	c.channel <- value
}

// Pull pulls a value from the collector
func (c *Collector) Pull() any {
	return <-c.channel
}

// Channel returns the channel of the collector
func (c *Collector) Channel() <-chan any {
	return c.channel
}
