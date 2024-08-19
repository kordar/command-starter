package command_starter

import (
	"flag"
	"github.com/kordar/command"
	"github.com/spf13/cast"
)

var (
	f string
)

type CommandModule struct {
}

func (m CommandModule) Name() string {
	return "command_starter"
}

func (m CommandModule) Load(value interface{}) {
	item := cast.ToStringMapString(value)
	name := "fname"
	if item["name"] != "" {
		name = item["name"]
	}
	usage := "函数名称"
	if item["usage"] != "" {
		usage = item["usage"]
	}
	flag.StringVar(&f, name, item["value"], usage)
	command.StartCmd(name)
}

func (m CommandModule) Close() {
}
