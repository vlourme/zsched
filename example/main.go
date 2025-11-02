package main

import (
	"errors"
	"math/rand"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/vlourme/zsched"
	"github.com/vlourme/zsched/pkg/hooks"
)

type UserCtx struct {
	Name string
}

var helloTask = zsched.NewTask(
	"hello",
	func(ctx *zsched.Context[*UserCtx]) error {
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)

		ctx.Infoln("Hello", ctx.GetStr("name"))

		if rand.Intn(100) < 30 {
			return errors.New("random error")
		}

		return nil
	},
	zsched.WithConcurrency(10),
	zsched.WithMaxRetries(3),
	// zsched.WithSchedule("* * * * * *", map[string]any{"name": "John"}),
	// zsched.WithDefaultParameters(map[string]any{
	// 	"name": "World",
	// }),
)

var dispatchTask = zsched.NewTask(
	"dispatch",
	func(ctx *zsched.Context[*UserCtx]) error {
		for range ctx.GetInt("count", 10) {
			helloTask.Execute(map[string]any{"name": "World"})
		}
		return nil
	},
	zsched.WithDefaultParameters(map[string]any{
		"count": 10,
	}),
)

func main() {
	godotenv.Load()

	userCtx := UserCtx{
		Name: "John",
	}

	engine, err := zsched.NewBuilder(&userCtx).
		WithRabbitMQBroker(os.Getenv("RABBITMQ_URL")).
		WithTimescaleDBStorage(os.Getenv("POSTGRES_URL")).
		WithHooks(&hooks.PrometheusHook{}, &hooks.TaskLoggerHook{}).
		WithAPI(":8080").
		Build()
	if err != nil {
		panic(err)
	}

	engine.Register(helloTask, dispatchTask)
	engine.Start()
}
