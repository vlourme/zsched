package main

import (
	"os"
	"time"

	"github.com/vlourme/zsched/pkg/ctx"
	"github.com/vlourme/zsched/pkg/engine"
	"github.com/vlourme/zsched/pkg/hooks"
	"github.com/vlourme/zsched/pkg/task"
)

var helloTask = task.NewTask(
	"hello",
	func(ctx *ctx.C) error {
		uc := ctx.UserContext().(*UserCtx)
		time.Sleep(5000 * time.Millisecond)

		ctx.Infoln("Hello " + ctx.Get("name").String() + " from " + uc.Name + "!")

		return nil
	},
	task.WithConcurrency(10),
	task.WithMaxRetries(3),
	// task.WithSchedule("* * * * * *", map[string]any{"name": "John"}),
	// task.WithSchedule("*/20 * * * * *", map[string]any{"name": "Mike"}),
	// task.WithSchedule("*/30 * * * * *", map[string]any{"name": "Carl"}),
)

type UserCtx struct {
	Name string
}

func main() {
	userCtx := UserCtx{
		Name: "John",
	}

	engine, err := engine.NewEngine(
		engine.WithRabbitMQBroker(os.Getenv("RABBITMQ_URL")),
		engine.WithQuestDBStorage(os.Getenv("QUESTDB_URL")),
		engine.WithAPI(":8080"),
		engine.WithUserContext(&userCtx),
		engine.WithHooks(hooks.NewPrometheusHook()),
	)

	if err != nil {
		panic(err)
	}

	engine.Register(helloTask)
	engine.Start()
}
