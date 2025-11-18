package main

import (
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
		time.Sleep(time.Duration(rand.Intn(3000)) * time.Millisecond)

		ctx.Infoln("Hello", ctx.GetStr("name"))

		return nil
	},
	zsched.WithConcurrency(10),
	zsched.WithMaxRetries(3),
	zsched.WithTags("stockx"),
	// zsched.WithDefaultParameters(map[string]any{
	// 	"name": "World",
	// }),
)

var dispatchTask = zsched.NewTask(
	"dispatch",
	func(ctx *zsched.Context[*UserCtx]) error {
		time.Sleep(5 * time.Second)
		for range ctx.GetInt("count", 10) {
			helloTask.Execute(map[string]any{"name": "World"})
		}
		return nil
	},
	zsched.WithMaxRetries(0),
	zsched.WithTags("goat"),
	zsched.WithDefaultParameters(map[string]any{
		"count": 10,
	}),
	zsched.WithSchedule("* * 0 * * *", map[string]any{"name": "John", "count": 10, "message": "Hello"}),
	zsched.WithSchedule("* * 1 * * *", map[string]any{"name": "John", "count": 10, "message": "Hello"}),
	zsched.WithSchedule("* * 2 * * *", map[string]any{"name": "John", "count": 10, "message": "Hello"}),
	zsched.WithSchedule("* * 3 * * *", map[string]any{"name": "John", "count": 10, "message": "Hello"}),
	zsched.WithSchedule("* * 4 * * *", map[string]any{"name": "John", "count": 10, "message": "Hello"}),
	zsched.WithSchedule("* * 5 * * *", map[string]any{"name": "John", "count": 10, "message": "Hello"}),
	zsched.WithSchedule("* * 6 * * *", map[string]any{"name": "John", "count": 10, "message": "Hello"}),
)

func main() {
	godotenv.Load()

	userCtx := UserCtx{
		Name: "John",
	}

	engine, err := zsched.NewBuilder(&userCtx).
		WithRabbitMQBroker(os.Getenv("RABBITMQ_URL")).
		WithTimescaleDBStorage(os.Getenv("POSTGRES_URL")).
		WithHooks(&hooks.PrometheusHook{}).
		WithAPI(":8080").
		Build()
	if err != nil {
		panic(err)
	}

	engine.Register(helloTask, dispatchTask)
	engine.Start()
}
