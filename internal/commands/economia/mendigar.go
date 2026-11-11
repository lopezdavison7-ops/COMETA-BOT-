package economia

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
Name:        "mendigar",
Aliases:     []string{"beg", "pedir"},
Description: "pide cometas a los transeúntes (cooldown 5 min)",
Category:    "𝐄⵿𝐜𝵿𝐨⵿𝐧𝐨⵿𝐦⵿𝐢⵿𝐚",
Handler:     Mendigar,
})
}

func Mendigar(evt *events.Message, client *whatsmeow.Client, _ []string) {
if store.GlobalEconomy == nil {
return
}

jid := evt.Info.Sender.String()
msg, err := store.GlobalEconomy.Beg(jid)

var prefix string
if err != nil {
if err.Error() == "cooldown" {
prefix = "⏳ "
} else {
prefix = "❌ "
}
} else {
prefix = "🥺 "
}

client.SendMessage(context.Background(), evt.Info.Chat, &waE2E.Message{
ExtendedTextMessage: &waE2E.ExtendedTextMessage{
Text: proto.String(fmt.Sprintf("╭ᩙ╌͜╍͜╌͜─͜─͜─͜─፝͜「 `🥺 MENDIGAR` 」፝͜─͜─͜─͜─͜╌͜━㍃゙⃞\n> ㅤ\n🫆֟፝𑑓᮫̤٘𑑓̸̷ %s%s\n╰᳕╌͜╍͜╌͜─͜─͜─͜─͜─͜─⵿֟፝─͜─͜─͜─͜─͜─͜─͜─͜╌͜━㍃゙⃞", prefix, msg)),
ContextInfo: &waE2E.ContextInfo{QuotedMessage: evt.Message},
},
})
}
