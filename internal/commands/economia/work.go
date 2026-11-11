package economia

import (
"context"
"fmt"
"math/rand"

"go.mau.fi/whatsmeow"
"go.mau.fi/whatsmeow/proto/waE2E"
"go.mau.fi/whatsmeow/types/events"
"google.golang.org/protobuf/proto"

"cometa-bot/internal/registry"
"cometa-bot/internal/store"
)

func init() {
registry.Register(&registry.Command{
Name:        "w",
Aliases:     []string{"work", "trabajar"},
Description: "trabaja para ganar cometas aleatorios",
Category:    "𝐄⵿𝐜𝵿𝐨⵿𝐧𝐨⵿𝐦⵿𝐢⵿𝐚",
Handler:     Work,
})
}

func Work(evt *events.Message, client *whatsmeow.Client, _ []string) {
if store.GlobalEconomy == nil {
sendEcoError(client, evt, "La economía no está inicializada.")
return
}

jid := evt.Info.Sender.String()
amount := rand.Intn(150) + 50

msg, err := store.GlobalEconomy.Work(jid, amount, 600)
if err != nil {
sendEcoError(client, evt, msg)
return
}

client.SendMessage(context.Background(), evt.Info.Chat, &waE2E.Message{
ExtendedTextMessage: &waE2E.ExtendedTextMessage{
Text: proto.String(fmt.Sprintf("╭ᩙ╌͜╍͜╌͜─͜─͜─͜─፝͜「 `💼 TRABAJO` 」፝͜─͜─͜─͜─͜╌͜━㍃゙⃞\n> ㅤ\n🫆֟፝𑑓᮫̤٘𑑓̸̷ %s\n╰᳕╌͜╍͜╌͜─͜─͜─͜─͜─͜─⵿֟፝─͜─͜─͜─͜─͜─͜─͜─͜╌͜━㍃゙⃞", msg)),
ContextInfo: &waE2E.ContextInfo{
QuotedMessage: evt.Message,
},
},
})
}
