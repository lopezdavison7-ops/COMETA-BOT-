package handler

import (
"fmt"
"strings"

"go.mau.fi/whatsmeow"
"go.mau.fi/whatsmeow/types/events"

"cometa-bot/internal/registry"

_ "cometa-bot/internal/commands/economia"
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

func Route(evt *events.Message, client *whatsmeow.Client) {
text := strings.TrimSpace(getText(evt))
if text == "" || !strings.HasPrefix(text, "/") {
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
