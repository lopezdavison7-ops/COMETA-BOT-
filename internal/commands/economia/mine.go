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
Name:        "mine",
Aliases:     []string{"minar"},
Description: "mina recursos para ganar cometas",
Category:    "𝐄⵿𝐜𝵿𝐨⵿𝐧𝐨⵿𝐦⵿𝐢⵿𝐚",
Handler:     Mine,
})
}

type mineResource struct {
Name   string
Emoji  string
Min    int
Max    int
Weight int
}

var mineResources = []mineResource{
{"Piedra", "🪨", 10, 30, 30},
{"Carbón", "⚫", 20, 50, 25},
{"Hierro", "⚙️", 40, 80, 20},
{"Oro", "🥇", 80, 150, 15},
{"Diamante", "💎", 150, 300, 8},
{"Esmeralda", "💚", 200, 400, 2},
}

func Mine(evt *events.Message, client *whatsmeow.Client, _ []string) {
jid := evt.Info.Sender.String()

if remaining, ok := store.GlobalEconomy.CheckCooldown(jid, "mine", 900); !ok {
mins := remaining / 60
secs := remaining % 60
sendEcoError(client, evt, fmt.Sprintf("⏳ Ya minaste recientemente.\nEsperá *%d min %d seg* para volver a minar.", mins, secs))
return
}

totalWeight := 0
for _, r := range mineResources {
totalWeight += r.Weight
}

randVal := rand.Intn(totalWeight)
cumulative := 0
var chosen mineResource

for _, r := range mineResources {
cumulative += r.Weight
if randVal < cumulative {
chosen = r
break
}
}

amount := rand.Intn(chosen.Max-chosen.Min+1) + chosen.Min
store.GlobalEconomy.AddBalance(jid, amount)
store.GlobalEconomy.SetCooldown(jid, "mine")

text := fmt.Sprintf(`╭ᩙ╌͜╍͜╌͜─͜─͜─͜─፝͜「 `+"`⛏️ MINAR`"+` 」፝͜─͜─͜─͜─͜╌͜━㍃゙⃞
> ㅤ
🫆֟፝𑑓᮫̤٘𑑓̸̷ Encontraste %s *%s*
🫆֟፝𑑓᮫̤٘𑑓̸̷ 💰 *Ganaste:* %d cometas
╰᳕╌͜╍͜╌͜─͜─͜─͜─͜─͜─⵿֟፝─͜─͜─͜─͜─͜─͜─͜─͜╌͜━㍃゙⃞`, chosen.Emoji, chosen.Name, amount)

client.SendMessage(context.Background(), evt.Info.Chat, &waE2E.Message{
ExtendedTextMessage: &waE2E.ExtendedTextMessage{
Text: proto.String(text),
ContextInfo: &waE2E.ContextInfo{QuotedMessage: evt.Message},
},
})
}
