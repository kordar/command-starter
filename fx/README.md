# command-starter/fx

将 [`github.com/kordar/command`](../../../packages/command) 以 Fx 依赖注入方式挂载到应用中，通过 `fx.ValueGroup` 自动聚合所有 `command.FuncCli` 实现并执行。

## 能力

- 通过 `fx.ValueGroup("command-funccli")` 聚合所有注册的命令
- 支持配置驱动（`StarterModule`）和编程式（`Module`）两种接入方式
- 提供 `ExecuteCommand` 在 Fx 启动后异步执行指定命令
- 提供 `StartFunc` 供外部直接调用启动命令
- 支持 flag 模式：绑定 `flag.StringVar` + 自动 `flag.Parse`

## 注入类型

| 类型 | 说明 |
|------|------|
| `CommandRegistry` | 聚合了所有 `command.FuncCli` 实现的注册表 |
| `ModuleConfig` | 命令名、flag 参数配置 |
| `CmdArgs` | 从外部传入的命令名称与原始参数 |

## 快速开始

```go
package main

import (
    starter "github.com/kordar/command-starter/fx/v2"
    "github.com/kordar/command"
    "go.uber.org/fx"
)

// 实现 command.FuncCli
type Hello struct {
    command.BaseFuncCli
}

func (Hello) Name() string                    { return "hello" }
func (Hello) Execute(args ...interface{})      { println("hello world") }
func (Hello) GetArgs() []interface{}           { return nil }

func main() {
    // 通过 fx.ValueGroup 注入命令
    provider := fx.Annotate(
        func() command.FuncCli { return &Hello{} },
        fx.As(new(command.FuncCli)),
        fx.ResultTags(`group:"command-funccli"`),
    )

    app := fx.New(
        starter.Module(starter.ModuleConfig{
            Name:      "f",
            FlagParse: true,
        }),
        fx.Provide(provider),
    )
    app.Run()
}
```

### 配置驱动（gocfg）

```go
obs := starter.StarterModule("command")
app := fx.New(obs.Load(yourConfigData)...)
```

配置键位于 `[command]` section：

| 键 | 类型 | 默认值 | 说明 |
|----|------|--------|------|
| `name` | string | `fname` | 命令行 flag 名 |
| `usage` | string | `function name` | flag 说明 |
| `value` | string | — | flag 默认值 |
| `flag-parse` | bool | `false` | 是否自动 `flag.Parse()` |

## API

### Module

```go
starter.Module(starter.ModuleConfig{
    Name:      "cmd",
    Usage:     "执行的命令",
    FlagParse: true,
})
```

### ExecuteCommand — Fx 启动后异步执行

```go
fx.Invoke(starter.ExecuteCommand)
```

配合外部传入 `CmdArgs`：

```go
fx.Supply(starter.CmdArgs{
    Name:    "hello",
    RawArgs: os.Args,
})
```

### StartFunc — 直接启动

```go
starter.StartFunc("hello", reg)
```

## 命令注册

命令通过 `fx.ValueGroup("command-funccli")` 注入：

```go
fx.Provide(
    fx.Annotate(
        NewHelloCommand,
        fx.As(new(command.FuncCli)),
        fx.ResultTags(`group:"command-funccli"`),
    ),
    fx.Annotate(
        NewMigrateCommand,
        fx.As(new(command.FuncCli)),
        fx.ResultTags(`group:"command-funccli"`),
    ),
)
```
