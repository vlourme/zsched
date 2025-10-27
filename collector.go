package zsched

import "context"

type TaskCollectorAction func(c *Collector)

type Collector struct {
	channel chan any
}

func newCollector(bufferSize ...int) *Collector {
	if len(bufferSize) != 0 {
		return &Collector{
			channel: make(chan any, bufferSize[0]),
		}
	}

	return &Collector{
		channel: make(chan any),
	}
}

// Push pushes a value to the collector
func (c *Collector) Push(value any) {
	c.channel <- value
}

// Consume consumes values from the collector
func (c *Collector) Consume(fn func(value any)) {
	c.ConsumeWithCtx(context.Background(), fn)
}

// Consume consumes values from the collector
func (c *Collector) ConsumeWithCtx(ctx context.Context, fn func(value any)) {
	for {
		select {
		case <-ctx.Done():
			return
		case value := <-c.channel:
			fn(value)
		}
	}
}

// Pull pulls a value from the collector
func (c *Collector) Pull() any {
	return <-c.channel
}

// Close closes the collector
func (c *Collector) Close() {
	close(c.channel)
}
