package handler

import (
"context"
"strings"

"go.mau.fi/whatsmeow"
"go.mau.fi/whatsmeow/proto/waE2E"
"go.mau.fi/whatsmeow/types"
"google.golang.org/protobuf/proto"

"cometa-bot/internal/store"
)

func SendWelcome(client *whatsmeow.Client, chat types.JID, user types.JID, action string) {
jid := chat.String()
welcomeOn := store.GlobalEconomy.GetGroupSetting(jid, "welcome")
goodbyeOn := store.GlobalEconomy.GetGroupSetting(jid, "goodbye")

var template string
if action == "add" && welcomeOn == 1 {
template = store.GlobalEconomy.GetGroupMessage(jid, "welcome_msg")
if template == "" {
template = "╭〔 👋 𝐁𝐈𝐄𝐍𝐕𝐄𝐍𝐈𝐃𝐎 〕⬣\n┃\n┃ ¡Hola @user!\n┃ Bienvenido a @group\n┃\n┃ Disfrutá tu estancia ☄️\n┃\n╰━━━━━━━━━━━━━━━━⬣"
}
} else if action == "remove" && goodbyeOn == 1 {
template = store.GlobalEconomy.GetGroupMessage(jid, "goodbye_msg")
if template == "" {
template = "╭〔 👋 𝐀𝐃𝐈𝐎𝐒 〕⬣\n┃\n┃ @user se fue del grupo.\n┃ ¡Hasta pronto! 🫡\n┃\n╰━━━━━━━━━━━━━━━━⬣"
}
} else {
return
}

group, err := client.GetGroupInfo(context.Background(), chat)
if err != nil {
return
}

template = strings.ReplaceAll(template, "@user", "@"+user.User)
template = strings.ReplaceAll(template, "@name", user.User)
template = strings.ReplaceAll(template, "@group", group.Name)

client.SendMessage(context.Background(), chat, &waE2E.Message{
ExtendedTextMessage: &waE2E.ExtendedTextMessage{
Text: proto.String(template),
ContextInfo: &waE2E.ContextInfo{
MentionedJID: []string{user.String()},
},
},
})
}
