package groups

import (
"context"
"fmt"

"go.mau.fi/whatsmeow"
"go.mau.fi/whatsmeow/proto/waE2E"
"go.mau.fi/whatsmeow/types"
"go.mau.fi/whatsmeow/types/events"
"google.golang.org/protobuf/proto"

"cometa-bot/internal/store"
)

func isGroup(jid types.JID) bool {
return jid.Server == "g.us"
}

func isAdmin(client *whatsmeow.Client, chat, sender types.JID) bool {
group, err := client.GetGroupInfo(context.Background(), chat)
if err != nil {
return false
}
senderStr := sender.String()
for _, p := range group.Participants {
if p.JID.String() == senderStr {
return p.IsAdmin || p.IsSuperAdmin
}
}
return false
}

func isOwner(sender types.JID) bool {
return store.IsOwner(sender.String())
}

func sendGroupMsg(client *whatsmeow.Client, chat types.JID, text string, quoted *waE2E.Message) {
client.SendMessage(context.Background(), chat, &waE2E.Message{
ExtendedTextMessage: &waE2E.ExtendedTextMessage{
Text: proto.String(text),
ContextInfo: &waE2E.ContextInfo{
QuotedMessage: quoted,
},
},
})
}

func sendGroupError(client *whatsmeow.Client, evt *events.Message, text string) {
sendGroupMsg(client, evt.Info.Chat, "❌ "+text, evt.Message)
}

func sendGroupSuccess(client *whatsmeow.Client, evt *events.Message, text string) {
sendGroupMsg(client, evt.Info.Chat, "✅ "+text, evt.Message)
}

func requireGroup(client *whatsmeow.Client, evt *events.Message) bool {
if !isGroup(evt.Info.Chat) {
sendGroupError(client, evt, "Este comando solo funciona en grupos.")
return false
}
return true
}

func requireAdmin(client *whatsmeow.Client, evt *events.Message) bool {
if !requireGroup(client, evt) {
return false
}
if !isAdmin(client, evt.Info.Chat, evt.Info.Sender) && !isOwner(evt.Info.Sender) {
sendGroupError(client, evt, "Necesitás ser admin para usar este comando.")
return false
}
return true
}

func requireOwner(client *whatsmeow.Client, evt *events.Message) bool {
if !isOwner(evt.Info.Sender) {
sendGroupError(client, evt, "Solo el owner puede usar este comando.")
return false
}
return true
}

func parseMention(text string) (types.JID, error) {
if len(text) < 3 {
return types.JID{}, fmt.Errorf("mención inválida")
}
num := text[1 : len(text)-1]
return types.NewJID(num, types.DefaultUserServer), nil
}
