package command_starter

import (
	"flag"
	"github.com/kordar/command"
	"github.com/spf13/cast"
)

var (
	f string
)

func AddCli(cli ...command.FuncCli) {
	for _, fCli := range cli {
		command.SetCmd(fCli)
	}
}

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
	usage := "function name"
	if item["usage"] != "" {
		usage = item["usage"]
	}
	val := cast.ToString(item["value"])
	flag.StringVar(&f, name, val, usage)
	flagParse := cast.ToBool(item["flagParse"])
	if flagParse {
		flag.Parse()
	}
	command.StartCmd(f)
}

func (m CommandModule) Close() {
}
