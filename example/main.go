package main

import (
	"fmt"
	"os"
	"time"

	"github.com/vlourme/zsched"
	"github.com/vlourme/zsched/pkg/hooks"
)

type UserCtx struct {
	Name string
}

var helloTask = zsched.NewTask(
	"hello",
	func(ctx *zsched.Context[*UserCtx]) error {
		time.Sleep(10000 * time.Millisecond)

		ctx.Infoln("Executed", ctx.GetStr("name"))

		ctx.Push(1)

		return nil
	},
	zsched.WithCollector(func(c *zsched.Collector, userCtx any) {
		total := 0
		c.Consume(func(value any) {
			total += value.(int)
			fmt.Println("Value:", value, "Total:", total)
		})
	}),
	zsched.WithConcurrency(10),
	zsched.WithMaxRetries(3),
	zsched.WithSchedule("* * * * * *", map[string]any{"name": "John"}),
)

func main() {
	userCtx := UserCtx{
		Name: "John",
	}

	engine, err := zsched.NewBuilder(&userCtx).
		WithRabbitMQBroker(os.Getenv("RABBITMQ_URL")).
		WithQuestDBStorage(os.Getenv("QUESTDB_URL")).
		WithHooks(&hooks.PrometheusHook{}, &hooks.TaskLoggerHook{}).
		WithAPI(":8080").
		Build()
	if err != nil {
		panic(err)
	}

	engine.Register(helloTask)
	engine.Start()
}
