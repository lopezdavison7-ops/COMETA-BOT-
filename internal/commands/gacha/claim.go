package gacha

import (
"context"
"fmt"

"go.mau.fi/whatsmeow"
"go.mau.fi/whatsmeow/proto/waE2E"
"go.mau.fi/whatsmeow/types/events"
"google.golang.org/protobuf/proto"

"cometa-bot/internal/registry"
"cometa-bot/internal/store"
)

func init() {
registry.Register(&registry.Command{
Name:        "claim",
Aliases:     []string{"reclamar", "c"},
Description: "reclama el personaje que acabas de sacar con /rw",
Category:    "𝐆⵿𝐚⵿𝐜⵿𝐡⵿𝐚",
Handler:     ClaimCmd,
})
}

func ClaimCmd(evt *events.Message, client *whatsmeow.Client, _ []string) {
jid := evt.Info.Sender.String()
file, ok := GetLastRoll(jid)
if !ok {
client.SendMessage(context.Background(), evt.Info.Chat, &waE2E.Message{
ExtendedTextMessage: &waE2E.ExtendedTextMessage{
Text: proto.String("❌ No tenés ningún personaje para reclamar. Usá */rw* primero."),
ContextInfo: &waE2E.ContextInfo{QuotedMessage: evt.Message},
},
})
return
}

char, err := store.Claim(jid, file)
if err != nil {
client.SendMessage(context.Background(), evt.Info.Chat, &waE2E.Message{
ExtendedTextMessage: &waE2E.ExtendedTextMessage{
Text: proto.String("❌ " + err.Error()),
ContextInfo: &waE2E.ContextInfo{QuotedMessage: evt.Message},
},
})
ClearRoll(jid)
return
}

ClearRoll(jid)

text := fmt.Sprintf(
"╭〔 ✅ 𝐑𝐄𝐂𝐋𝐀𝐌𝐀𝐃𝐎 〕⬣\n┃\n┃ 🎉 *%s* ahora es tuya!\n┃ 📺 *Serie:* %s\n┃ 💰 *Valor:* %d cometas\n┃\n╰━━━━━━━━━━━━━━━━⬣\n\n> Mirá tu inventario con */inv* ☄️",
char.Name, char.Series, char.Value,
)

client.SendMessage(context.Background(), evt.Info.Chat, &waE2E.Message{
ExtendedTextMessage: &waE2E.ExtendedTextMessage{
Text: proto.String(text),
ContextInfo: &waE2E.ContextInfo{QuotedMessage: evt.Message},
},
})
}
