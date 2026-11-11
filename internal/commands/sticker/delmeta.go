package sticker

import (
"context"

"go.mau.fi/whatsmeow"
"go.mau.fi/whatsmeow/proto/waE2E"
"go.mau.fi/whatsmeow/types/events"
"google.golang.org/protobuf/proto"

"cometa-bot/internal/registry"
"cometa-bot/internal/store"
)

func init() {
registry.Register(&registry.Command{
Name:        "delmeta",
Aliases:     []string{"delmarcasticker", "delsticker", "resetmeta"},
Description: "elimina tu marca de agua personalizada para stickers",
Category:    "𝐒⵿𝐭⵿𝐢⵿𝐜⵿𝐤⵿𝐞⵿𝐫",
Handler:     DelMeta,
})
}

func DelMeta(evt *events.Message, client *whatsmeow.Client, args []string) {
ctx := context.Background()
senderJID := evt.Info.Sender.String()

if !store.HasStickerMeta(senderJID) {
client.SendMessage(ctx, evt.Info.Chat, &waE2E.Message{
ExtendedTextMessage: &waE2E.ExtendedTextMessage{
Text: proto.String("ℹ️ No tienes una marca personalizada establecida."),
ContextInfo: &waE2E.ContextInfo{
QuotedMessage: evt.Message,
},
},
})
return
}

store.DeleteStickerMeta(senderJID)

client.SendMessage(ctx, evt.Info.Chat, &waE2E.Message{
ExtendedTextMessage: &waE2E.ExtendedTextMessage{
Text: proto.String("✅ Marca de stickers eliminada.\nAhora usarás la marca por defecto."),
ContextInfo: &waE2E.ContextInfo{
QuotedMessage: evt.Message,
},
},
})
}
