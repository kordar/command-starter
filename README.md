# command-starter

[`github.com/kordar/command`](../../packages/command) 的轻量启动器，负责命令注册、flag 绑定与启动编排。

## 子模块

| 目录 | 说明 |
|------|------|
| [fx](./fx) | Uber Fx 集成，通过依赖注入聚合命令并自动执行 |

## 快速开始（独立模式）

### 安装

```bash
go get github.com/kordar/command-starter
```

### 使用

```go
package main

import (
    starter "github.com/kordar/command-starter"
    "github.com/kordar/command"
)

type Hello struct{}

func (Hello) Name() string                    { return "hello" }
func (Hello) Execute(args ...interface{})      { println("hello world") }
func (Hello) GetArgs() []interface{}           { return nil }

func main() {
    starter.AddCli(&Hello{})
    starter.StartWithName("f")
}
```

```bash
go run . -f hello
```

### API

```go
// 注册命令
starter.AddCli(cmd1, cmd2, cmd3)

// 完整配置
starter.Start(starter.Config{
    Name:      "f",
    Usage:     "执行的命令名称",
    Value:     "",
    FlagParse: true,
})

// 便捷方式
starter.StartWithName("f")
```

### Config

| 字段 | 默认值 | 说明 |
|------|--------|------|
| `Name` | `fname` | 命令行参数名 |
| `Usage` | `function name` | 参数说明 |
| `Value` | — | 默认值（为空时不执行） |
| `FlagParse` | `false` | 是否自动 `flag.Parse()` |

---

## Fx 集成

详见 [fx/README.md](./fx/README.md)（待补充）。

```go
import starter "github.com/kordar/command-starter/fx/v2"

// 通过 gocfg 配置驱动加载
obs := starter.StarterModule("command")
app := fx.New(obs.Load(yourConfigData)...)
```
