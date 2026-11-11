package groups

import (
"strings"

"go.mau.fi/whatsmeow"
"go.mau.fi/whatsmeow/types/events"

"cometa-bot/internal/registry"
"cometa-bot/internal/store"
)

func init() {
registry.Register(&registry.Command{
Name:        "welcome",
Aliases:     []string{"bienvenida"},
Description: "activa/desactiva la bienvenida",
Category:    "𝐆⵿𝐫⵿𝐮⵿𝐩⵿𝐨⵿𝐬",
Handler:     Welcome,
})
registry.Register(&registry.Command{
Name:        "goodbye",
Aliases:     []string{"despedida"},
Description: "activa/desactiva la despedida",
Category:    "𝐆⵿𝐫⵿𝐮⵿𝐩⵿𝐨⵿𝐬",
Handler:     Goodbye,
})
registry.Register(&registry.Command{
Name:        "setwelcome",
Aliases:     []string{"setbienvenida"},
Description: "personaliza el mensaje de bienvenida",
Category:    "𝐆⵿𝐫⵿𝐮⵿𝐩⵿𝐨⵿𝐬",
Handler:     SetWelcome,
})
registry.Register(&registry.Command{
Name:        "setgoodbye",
Aliases:     []string{"setdespedida"},
Description: "personaliza el mensaje de despedida",
Category:    "𝐆⵿𝐫⵿𝐮⵿𝐩⵿𝐨⵿𝐬",
Handler:     SetGoodbye,
})
}

func Welcome(evt *events.Message, client *whatsmeow.Client, args []string) {
if !requireAdmin(client, evt) { return }
jid := evt.Info.Chat.String()
current := store.GlobalEconomy.GetGroupSetting(jid, "welcome")
if len(args) > 0 {
arg := strings.ToLower(args[0])
if arg == "on" || arg == "1" {
store.GlobalEconomy.SetGroupSetting(jid, "welcome", 1)
sendGroupSuccess(client, evt, "Bienvenida activada.")
return
}
if arg == "off" || arg == "0" {
store.GlobalEconomy.SetGroupSetting(jid, "welcome", 0)
sendGroupSuccess(client, evt, "Bienvenida desactivada.")
return
}
}
status := "❌ Desactivada"
if current == 1 { status = "✅ Activada" }
sendGroupMsg(client, evt.Info.Chat, "╭〔 👋 𝐁𝐈𝐄𝐍𝐕𝐄𝐍𝐈𝐃𝐀 〕⬣\n┃\n┃ Estado: "+status+"\n┃\n┃ Uso:\n┃ › /welcome on\n┃ › /welcome off\n┃\n╰━━━━━━━━━━━━━━━━⬣", evt.Message)
}

func Goodbye(evt *events.Message, client *whatsmeow.Client, args []string) {
if !requireAdmin(client, evt) { return }
jid := evt.Info.Chat.String()
current := store.GlobalEconomy.GetGroupSetting(jid, "goodbye")
if len(args) > 0 {
arg := strings.ToLower(args[0])
if arg == "on" || arg == "1" {
store.GlobalEconomy.SetGroupSetting(jid, "goodbye", 1)
sendGroupSuccess(client, evt, "Despedida activada.")
return
}
if arg == "off" || arg == "0" {
store.GlobalEconomy.SetGroupSetting(jid, "goodbye", 0)
sendGroupSuccess(client, evt, "Despedida desactivada.")
return
}
}
status := "❌ Desactivada"
if current == 1 { status = "✅ Activada" }
sendGroupMsg(client, evt.Info.Chat, "╭〔 👋 𝐃𝐄𝐒𝐏𝐄𝐃𝐈𝐃𝐀 〕⬣\n┃\n┃ Estado: "+status+"\n┃\n┃ Uso:\n┃ › /goodbye on\n┃ › /goodbye off\n┃\n╰━━━━━━━━━━━━━━━━⬣", evt.Message)
}

func SetWelcome(evt *events.Message, client *whatsmeow.Client, args []string) {
if !requireAdmin(client, evt) { return }
if len(args) == 0 {
sendGroupError(client, evt, "Usá: /setwelcome <mensaje>\n\nVariables:\n• @user → mención\n• @group → nombre del grupo\n• @name → nombre del usuario")
return
}
jid := evt.Info.Chat.String()
msg := strings.Join(args, " ")
store.GlobalEconomy.SetGroupMessage(jid, "welcome_msg", msg)
sendGroupSuccess(client, evt, "Mensaje de bienvenida actualizado.")
}

func SetGoodbye(evt *events.Message, client *whatsmeow.Client, args []string) {
if !requireAdmin(client, evt) { return }
if len(args) == 0 {
sendGroupError(client, evt, "Usá: /setgoodbye <mensaje>\n\nVariables:\n• @user → mención\n• @group → nombre del grupo\n• @name → nombre del usuario")
return
}
jid := evt.Info.Chat.String()
msg := strings.Join(args, " ")
store.GlobalEconomy.SetGroupMessage(jid, "goodbye_msg", msg)
sendGroupSuccess(client, evt, "Mensaje de despedida actualizado.")
}
