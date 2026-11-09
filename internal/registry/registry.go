package registry

import (
"go.mau.fi/whatsmeow"
"go.mau.fi/whatsmeow/types/events"
)

type CommandFunc func(evt *events.Message, client *whatsmeow.Client, args []string)

type Command struct {
Name        string
Aliases     []string
Description string
Category    string
Handler     CommandFunc
}

var commands = map[string]*Command{}

func Register(cmd *Command) {
commands[cmd.Name] = cmd
for _, alias := range cmd.Aliases {
commands[alias] = cmd
}
}

func Get(name string) *Command {
return commands[name]
}

func GetAll() map[string]*Command {
return commands
}

func GetAllUnique() []*Command {
seen := map[string]bool{}
result := []*Command{}
for _, cmd := range commands {
if !seen[cmd.Name] {
seen[cmd.Name] = true
result = append(result, cmd)
}
}
return result
}
