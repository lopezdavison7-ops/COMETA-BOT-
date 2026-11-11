package handler

import (
"context"
"fmt"
"strings"
"time"

"go.mau.fi/whatsmeow"
"go.mau.fi/whatsmeow/proto/waE2E"
"go.mau.fi/whatsmeow/types"
"go.mau.fi/whatsmeow/types/events"
"google.golang.org/protobuf/proto"

"cometa-bot/internal/store"
)

func HandleGroupEvent(client *whatsmeow.Client, evt interface{}) {
switch v := evt.(type) {
case *events.GroupInfo:
handleGroupInfo(client, v)
}
}

func handleGroupInfo(client *whatsmeow.Client, evt *events.GroupInfo) {
chatJID := evt.JID
jidStr := chatJID.String()

welcomeOn := store.GlobalEconomy.GetGroupSetting(jidStr, "welcome")
goodbyeOn := store.GlobalEconomy.GetGroupSetting(jidStr, "goodbye")

if welcomeOn != 1 && goodbyeOn != 1 {
return
}

group, err := client.GetGroupInfo(context.Background(), chatJID)
if err != nil {
return
}

for _, participant := range evt.Join {
if welcomeOn == 1 {
sendWelcomeMessage(client, chatJID, group, participant, "add", "welcome")
}
}

for _, participant := range evt.Leave {
if goodbyeOn == 1 {
sendWelcomeMessage(client, chatJID, group, participant, "remove", "goodbye")
}
}
}

func sendWelcomeMessage(client *whatsmeow.Client, chat types.JID, group *types.GroupInfo, user types.JID, action, mediaKey string) {
jidStr := chat.String()

var template string
if action == "add" {
template = store.GlobalEconomy.GetGroupMessage(jidStr, "welcome_msg")
if template == "" {
template = "╭━━━〔 👋 BIENVENIDO 〕━━━⊷\n┃\n┃ ¡Hola @user!\n┃ Bienvenido a @group\n┃\n┃ Fecha: @date\n┃ Miembros: @members\n┃\n━━━━━━━━━━━━━━━━━━━━━━⊷\n\n🎉 ¡Disfrutá tu estancia!"
}
} else {
template = store.GlobalEconomy.GetGroupMessage(jidStr, "goodbye_msg")
if template == "" {
template = "╭━━━〔 👋 ADIOS 〕━━━⊷\n┃\n┃ @user se fue del grupo\n┃ ¡Hasta pronto!\n┃\n━━━━━━━━━━━━━━━━━━━━━━⊷"
}
}

now := time.Now()
dateStr := now.Format("02/01/2006 15:04")

mentionTag := "@" + user.User
mentionJID := user.String()

template = strings.ReplaceAll(template, "@user", mentionTag)
template = strings.ReplaceAll(template, "@name", user.User)
template = strings.ReplaceAll(template, "@group", group.Name)
template = strings.ReplaceAll(template, "@date", dateStr)
template = strings.ReplaceAll(template, "@members", fmt.Sprintf("%d", len(group.Participants)))

mentions := []string{mentionJID}

// 🚀 Cargar imagen pre-subida (Instantáneo)
imgMsg := store.GetMediaUpload(jidStr, mediaKey)

if imgMsg != nil {
imgMsg.Caption = proto.String(template)
imgMsg.ContextInfo = &waE2E.ContextInfo{
MentionedJID: mentions,
}
client.SendMessage(context.Background(), chat, &waE2E.Message{
ImageMessage: imgMsg,
})
fmt.Printf("📨 %s con imagen (pre-upload) enviado para %s en %s\n", action, mentionJID, chat.String())
return
}

// Fallback a solo texto si no hay imagen
client.SendMessage(context.Background(), chat, &waE2E.Message{
ExtendedTextMessage: &waE2E.ExtendedTextMessage{
Text: proto.String(template),
ContextInfo: &waE2E.ContextInfo{
MentionedJID: mentions,
},
},
})

fmt.Printf("📨 %s (texto) enviado para %s en %s\n", action, mentionJID, chat.String())
}
