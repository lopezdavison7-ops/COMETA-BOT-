package registry

import (
"sort"
"sync"

"go.mau.fi/whatsmeow"
"go.mau.fi/whatsmeow/types/events"
)

type Command struct {
Name        string
Aliases     []string
Description string
Category    string
Handler     func(evt *events.Message, client *whatsmeow.Client, args []string)
}

var (
commands = make(map[string]*Command)
aliases  = make(map[string]string)
mu       sync.RWMutex
)

func Register(cmd *Command) {
mu.Lock()
defer mu.Unlock()
commands[cmd.Name] = cmd
for _, alias := range cmd.Aliases {
aliases[alias] = cmd.Name
}
}

func Get(name string) *Command {
mu.RLock()
defer mu.RUnlock()
if cmd, ok := commands[name]; ok {
return cmd
}
if realName, ok := aliases[name]; ok {
return commands[realName]
}
return nil
}

func Count() int {
mu.RLock()
defer mu.RUnlock()
return len(commands)
}

func Categories() []string {
mu.RLock()
defer mu.RUnlock()

catSet := make(map[string]bool)
for _, cmd := range commands {
if cmd.Category != "" {
catSet[cmd.Category] = true
}
}

var cats []string
for cat := range catSet {
cats = append(cats, cat)
}
sort.Strings(cats)
return cats
}

func GetByCategory(category string) []*Command {
mu.RLock()
defer mu.RUnlock()

var cmds []*Command
for _, cmd := range commands {
if cmd.Category == category {
cmds = append(cmds, cmd)
}
}

sort.Slice(cmds, func(i, j int) bool {
return cmds[i].Name < cmds[j].Name
})

return cmds
}

func AllCommands() []*Command {
mu.RLock()
defer mu.RUnlock()

var cmds []*Command
for _, cmd := range commands {
cmds = append(cmds, cmd)
}

sort.Slice(cmds, func(i, j int) bool {
return cmds[i].Name < cmds[j].Name
})

return cmds
}
