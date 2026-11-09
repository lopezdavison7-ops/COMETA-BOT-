package economia

import (
"context"
"math/rand"
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
Name:        "w",
Aliases:     []string{"work", "trabajar"},
Description: "trabaja para ganar cometas aleatorios",
Category:    "𝐄⵿𝐜𝵿𝐨⵿𝐧𝐨⵿𝐦⵿𝐢⵿𝐚",
Handler:     W,
})
}

func W(evt *events.Message, client *whatsmeow.Client, _ []string) {
if store.GlobalEconomy == nil {
return
}

jid := evt.Info.Sender.String()
amount := rand.Intn(41) + 30
msg, err := store.GlobalEconomy.Work(jid, amount, 360)
if err != nil {
msg = "❌ Error al procesar tu trabajo."
}

b := strings.Builder{}
b.WriteString("╭ᩙ╌͜╍͜╌͜─͜─͜─͜─፝͜「 `💰 WORK` 」፝͜─͜─͜─͜─͜╌͜━㍃゙⃞\n> ㅤ\n")
b.WriteString("🫆֟፝𑑓᮫̤٘𑑓̸̷ " + msg + "\n")
b.WriteString("╰᳕╌͜╍͜╌͜─͜─͜─͜─͜─͜─⵿֟፝─͜─͜─͜─͜─͜─͜─͜─͜╌͜━㍃゙⃞")

client.SendMessage(context.Background(), evt.Info.Chat, &waE2E.Message{
ExtendedTextMessage: &waE2E.ExtendedTextMessage{
Text: proto.String(b.String()),
},
})
}
