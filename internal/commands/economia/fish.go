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
Name:        "fish",
Aliases:     []string{"pescar", "pescado"},
Description: "pesca para ganar cometas",
Category:    "𝐄⵿𝐜𝵿𝐨⵿𝐧𝐨⵿𝐦⵿𝐢⵿𝐚",
Handler:     Fish,
})
}

type fishType struct {
Name   string
Emoji  string
Min    int
Max    int
Weight int
}

var fishTypes = []fishType{
{"Sardina", "🐟", 15, 35, 30},
{"Trucha", "🐠", 30, 60, 25},
{"Salmón", "🍣", 50, 100, 20},
{"Atún", "🐡", 80, 150, 15},
{"Pez Espada", "⚔️", 120, 250, 8},
{"Pez Dorado", "✨", 200, 400, 2},
}

func Fish(evt *events.Message, client *whatsmeow.Client, _ []string) {
jid := evt.Info.Sender.String()

if remaining, ok := store.GlobalEconomy.CheckCooldown(jid, "fish", 900); !ok {
mins := remaining / 60
secs := remaining % 60
sendEcoError(client, evt, fmt.Sprintf("⏳ Ya pescaste recientemente.\nEsperá *%d min %d seg* para volver a pescar.", mins, secs))
return
}

if rand.Intn(100) < 10 {
store.GlobalEconomy.SetCooldown(jid, "fish")
sendEcoError(client, evt, "🎣 Se te rompió la caña de pescar.\nNo atrapaste nada esta vez.")
return
}

totalWeight := 0
for _, f := range fishTypes {
totalWeight += f.Weight
}

randVal := rand.Intn(totalWeight)
cumulative := 0
var chosen fishType

for _, f := range fishTypes {
cumulative += f.Weight
if randVal < cumulative {
chosen = f
break
}
}

amount := rand.Intn(chosen.Max-chosen.Min+1) + chosen.Min
store.GlobalEconomy.AddBalance(jid, amount)
store.GlobalEconomy.SetCooldown(jid, "fish")

text := fmt.Sprintf(`╭ᩙ╌͜╍͜╌͜─͜─͜─͜─፝͜「 `+"`🎣 PESCAR`"+` 」፝͜─͜─͜─͜─͜╌͜━㍃゙⃞
> ㅤ
🫆֟፝𑑓᮫̤٘𑑓̸̷ Atrapaste %s *%s*
🫆֟፝𑑓᮫̤٘𑑓̸̷ 💰 *Ganaste:* %d cometas
╰᳕╌͜╍͜╌͜─͜─͜─͜─͜─͜─⵿֟፝─͜─͜─͜─͜─͜─͜─͜─͜╌͜━㍃゙⃞`, chosen.Emoji, chosen.Name, amount)

client.SendMessage(context.Background(), evt.Info.Chat, &waE2E.Message{
ExtendedTextMessage: &waE2E.ExtendedTextMessage{
Text: proto.String(text),
ContextInfo: &waE2E.ContextInfo{QuotedMessage: evt.Message},
},
})
}
