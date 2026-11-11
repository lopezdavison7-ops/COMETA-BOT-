package sticker

import (
"context"
"fmt"
"strings"

"go.mau.fi/whatsmeow"
"go.mau.fi/whatsmeow/proto/waE2E"
"go.mau.fi/whatsmeow/types/events"
"google.golang.org/protobuf/proto"

"cometa-bot/internal/registry"
"cometa-bot/internal/store"
)

func init() {
registry.Register(&registry.Command{
Name:        "setmeta",
Aliases:     []string{"marcasticker", "meta", "setpack"},
Description: "cambia tu marca de agua personalizada para stickers",
Category:    "𝐒⵿𝐭⵿𝐢⵿𝐜⵿𝐤⵿𝐞⵿𝐫",
Handler:     SetMeta,
})
}

func SetMeta(evt *events.Message, client *whatsmeow.Client, args []string) {
ctx := context.Background()
senderJID := evt.Info.Sender.String()

// Leer el texto crudo del mensaje y quitar solo el comando
raw := ""
if c := evt.Message.GetConversation(); c != "" {
raw = c
} else if e := evt.Message.GetExtendedTextMessage(); e != nil {
raw = e.GetText()
}

fields := strings.Fields(raw)
text := strings.TrimSpace(strings.Join(args, " "))
if len(fields) > 1 {
text = strings.TrimSpace(strings.Join(fields[1:], " "))
}

fmt.Printf("🔍 /setmeta recibió: %q\n", text)

if text == "" {
helpText := "╭━━〔 ⚙️ 𝐒𝐄𝐓𝐌𝐄𝐓𝐀 〕━━⬣\n┃\n┃ ❌ Uso:\n┃ › /setmeta Autor\n┃ › /setmeta Pack | Autor\n┃\n┃ 💡 Ejemplos:\n┃ › /setmeta Cometa\n┃ › /setmeta Cometa Bot | @DevLuis\n┃\n╰━━━━━━━━━━━━━━━━⬣"

client.SendMessage(ctx, evt.Info.Chat, &waE2E.Message{
ExtendedTextMessage: &waE2E.ExtendedTextMessage{
Text: proto.String(helpText),
ContextInfo: &waE2E.ContextInfo{
QuotedMessage: evt.Message,
},
},
})
return
}

if store.HasStickerMeta(senderJID) {
client.SendMessage(ctx, evt.Info.Chat, &waE2E.Message{
ExtendedTextMessage: &waE2E.ExtendedTextMessage{
Text: proto.String("⚠️ Ya tienes una marca establecida.\nUsa /delmeta para eliminarla primero."),
ContextInfo: &waE2E.ContextInfo{
QuotedMessage: evt.Message,
},
},
})
return
}

var packname, author string
if strings.Contains(text, "|") {
parts := strings.SplitN(text, "|", 2)
packname = strings.TrimSpace(parts[0])
author = strings.TrimSpace(parts[1])
if packname == "" {
packname = text
}
if author == "" {
author = text
}
} else {
packname = text
author = text
}

store.SetStickerMeta(senderJID, packname, author)

successText := fmt.Sprintf("╭━━〔 ✅ 𝐌𝐀𝐑𝐂𝐀 𝐀𝐂𝐓𝐔𝐀𝐋𝐈𝐙𝐀𝐃𝐀 〕━━⬣\n┃\n┃ 📦 *Pack:* %s\n┃ ✍️ *Autor:* %s\n┃\n┃ Ahora tus stickers llevarán\n┃ esta marca de agua personalizada.\n┃\n╰━━━━━━━━━━━━━━━━⬣", packname, author)

client.SendMessage(ctx, evt.Info.Chat, &waE2E.Message{
ExtendedTextMessage: &waE2E.ExtendedTextMessage{
Text: proto.String(successText),
ContextInfo: &waE2E.ContextInfo{
QuotedMessage: evt.Message,
},
},
})
}
