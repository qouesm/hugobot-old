package main

import "github.com/qouesm/hugobot/commands"

func exportCommands() []commands.Command {
	return []commands.Command{
		// commands.ClassClear,
		// commands.Ping,
		// commands.Quietping,
		commands.Role,
		// commands.Q,
	}
}
