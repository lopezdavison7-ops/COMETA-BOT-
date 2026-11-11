package handler

import (
"context"
"fmt"
"strings"

"go.mau.fi/whatsmeow"
"go.mau.fi/whatsmeow/types"
"go.mau.fi/whatsmeow/types/events"

"cometa-bot/internal/registry"
"cometa-bot/internal/store"

_ "cometa-bot/internal/commands/ai"
_ "cometa-bot/internal/commands/economia"
_ "cometa-bot/internal/commands/fun"
_ "cometa-bot/internal/commands/gacha"
_ "cometa-bot/internal/commands/groups"
_ "cometa-bot/internal/commands/owner"
_ "cometa-bot/internal/commands/pokemon"
_ "cometa-bot/internal/commands/profile"
_ "cometa-bot/internal/commands/sticker"
_ "cometa-bot/internal/commands/system"
)

func getText(evt *events.Message) string {
if evt.Message == nil {
return ""
}
if text := evt.Message.GetConversation(); text != "" {
return text
}
if ext := evt.Message.GetExtendedTextMessage(); ext != nil {
return ext.GetText()
}
if img := evt.Message.GetImageMessage(); img != nil {
return img.GetCaption()
}
if vid := evt.Message.GetVideoMessage(); vid != nil {
return vid.GetCaption()
}
if doc := evt.Message.GetDocumentMessage(); doc != nil {
return doc.GetCaption()
}
return ""
}

func registerLIDsFromGroup(client *whatsmeow.Client, chatJID types.JID) {
if chatJID.Server != "g.us" {
return
}
group, err := client.GetGroupInfo(context.Background(), chatJID)
if err != nil {
return
}
for _, p := range group.Participants {
lidStr := p.LID.String()
phoneStr := p.PhoneNumber.User
if lidStr != "" && phoneStr != "" {
store.RegisterLID(lidStr, phoneStr)
fmt.Printf("🔍 Registrado LID: %s -> %s\n", lidStr, phoneStr)
}
}
}

func Route(evt *events.Message, client *whatsmeow.Client) {

if evt.Info.Chat.Server == "g.us" {
go registerLIDsFromGroup(client, evt.Info.Chat)
}

text := strings.TrimSpace(getText(evt))
if text == "" {
return
}

// 🔧 ARREGLO: Eliminar espacios después del /
// Si el texto empieza con "/" seguido de espacios, los quitamos
if strings.HasPrefix(text, "/") {
// Buscar el primer carácter que no sea espacio después del /
trimmed := strings.TrimLeft(text[1:], " ")
text = "/" + trimmed
}

if !strings.HasPrefix(text, "/") {
return
}

parts := strings.Fields(text)
if len(parts) == 0 {
return
}

cmd := strings.ToLower(strings.TrimPrefix(parts[0], "/"))
if cmd == "" {
return
}

args := parts[1:]

if command := registry.Get(cmd); command != nil {
go func() {
defer func() {
if r := recover(); r != nil {
NotifyOwners(client, fmt.Sprintf("Pánico en /%s\nError: %v\nUser: %s", cmd, r, evt.Info.Sender.String()))
}
}()
command.Handler(evt, client, args)
}()
}
}
