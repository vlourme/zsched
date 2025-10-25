# zsched

A lightweight, opinionated mix of a queue system and task orchestrator built in and for Go, using QuestDB and AMQP broker. Designed for simplicity and performance.

![Web UI](media/screenshot.png)

## ✨ Features

- **Queue System**: Built on LavinMQ (AMQP 0.9.1) for reliable message delivery
- **Cron Scheduling**: Second-precision cron expressions for precise task timing
- **Retry Logic**: Configurable retry mechanisms for handling task failures
- **Concurrency Control**: Fine-grained control over task execution concurrency
- **Persistent Storage**: QuestDB integration for storing tasks and execution logs
- **REST API**: Complete HTTP API for task dispatch and log retrieval
- **Web Dashboard**: Clean, lightweight UI for task management and monitoring
- **Hooks**: Execute actions before and after task executions, including **Prometheus** metrics

## 🚀 Quick Start

### Installation

```bash
go get github.com/vlourme/zsched
```

### Basic Usage

1. **Define a task**:

```go
var helloTask = task.NewTask(
    "hello-world", // Queue name
    func(ctx *ctx.C) error {
        ctx.Info("Hello " + state.Get("name").String())

        if ctx.Get("start_child").Bool() {
            ctx.Execute(map[string]any{"name": "child"}, ctx.State)
        }

        return nil
    },
    task.WithConcurrency(10),
    task.WithMaxRetries(3),
    task.WithSchedule("0 * * * * *", map[string]any{"name": "John"}), // Every second
)
```

2. **Create and configure the engine**:

```go
func main() {
    engine, err := engine.NewEngine(
        engine.WithRabbitMQBroker("amqp://guest:guest@localhost:5672/"),
        engine.WithQuestDBStorage("postgres://admin:quest@localhost:8812/qdb?sslmode=disable"),
        engine.WithAPI(":8080"),
    )
    if err != nil {
        panic(err)
    }

    engine.Register(helloTask)
    engine.Start()
}
```

Check out the [example](example/main.go) for a complete example.

3. **Run your application**:

```bash
go run main.go
```

4. **Dispatch tasks via API**:

```bash
curl -X POST http://localhost:8080/tasks/hello-world \
  -H "Content-Type: application/json" \
  -d '{"name": "John"}'
```

### Docker support

You will have to dockerize your tasks and the engine, we advice to follow the Dockerfile (`example/Dockerfile`) and docker-compose.yml files to get started.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
