package gacha

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
Name:        "sell",
Aliases:     []string{"vender"},
Description: "vende un personaje de tu inventario por cometas",
Category:    "𝐆⵿𝐚⵿𝐜⵿𝐡⵿𝐚",
Handler:     SellCmd,
})
}

func SellCmd(evt *events.Message, client *whatsmeow.Client, args []string) {
if len(args) == 0 {
client.SendMessage(context.Background(), evt.Info.Chat, &waE2E.Message{
ExtendedTextMessage: &waE2E.ExtendedTextMessage{
Text: proto.String("❌ Usá: `/sell <nombre>`\nEjemplo: `/sell asuna`"),
ContextInfo: &waE2E.ContextInfo{QuotedMessage: evt.Message},
},
})
return
}

jid := evt.Info.Sender.String()
query := strings.Join(args, " ")

char, err := store.SellCharacter(jid, query)
if err != nil {
client.SendMessage(context.Background(), evt.Info.Chat, &waE2E.Message{
ExtendedTextMessage: &waE2E.ExtendedTextMessage{
Text: proto.String("❌ " + err.Error()),
ContextInfo: &waE2E.ContextInfo{QuotedMessage: evt.Message},
},
})
return
}

// Sumar cometas al balance del usuario
if store.GlobalEconomy != nil {
store.GlobalEconomy.AddBalance(jid, char.Value)
}

text := fmt.Sprintf(
"╭〔 💰 𝐕𝐄𝐍𝐃𝐈𝐃𝐎 〕⬣\n┃\n┃ Vendiste a *%s* por %d cometas\n┃ 📺 *Serie:* %s\n┃ 💳 Se sumaron a tu balance\n┃\n╰━━━━━━━━━━━━━━━━⬣",
char.Name, char.Value, char.Series,
)

client.SendMessage(context.Background(), evt.Info.Chat, &waE2E.Message{
ExtendedTextMessage: &waE2E.ExtendedTextMessage{
Text: proto.String(text),
ContextInfo: &waE2E.ContextInfo{QuotedMessage: evt.Message},
},
})
}
