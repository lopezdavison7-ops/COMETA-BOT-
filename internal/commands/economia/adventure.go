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
Name:        "adventure",
Aliases:     []string{"aventura", "adv"},
Description: "explora y gana cometas",
Category:    "𝐄⵿𝐜𝵿𝐨⵿𝐧𝐨⵿𝐦⵿𝐢⵿𝐚",
Handler:     Adventure,
})
}

type adventureOutcome struct {
Text   string
Min    int
Max    int
Weight int
}

var adventureOutcomes = []adventureOutcome{
{"Encontraste un cofre olvidado", 100, 250, 25},
{"Derrotaste a un bandido", 150, 300, 20},
{"Rescataste a un mercader", 200, 400, 15},
{"Descubriste una cueva secreta", 250, 500, 12},
{"Completaste una misión épica", 400, 800, 8},
{"Encontraste un tesoro legendario", 600, 1200, 5},
{"Te perdiste y no ganaste nada", 0, 0, 15},
}

func Adventure(evt *events.Message, client *whatsmeow.Client, _ []string) {
jid := evt.Info.Sender.String()

if remaining, ok := store.GlobalEconomy.CheckCooldown(jid, "adventure", 1200); !ok {
mins := remaining / 60
secs := remaining % 60
sendEcoError(client, evt, fmt.Sprintf("⏳ Ya exploraste recientemente.\nEsperá *%d min %d seg* para otra aventura.", mins, secs))
return
}

totalWeight := 0
for _, o := range adventureOutcomes {
totalWeight += o.Weight
}

randVal := rand.Intn(totalWeight)
cumulative := 0
var chosen adventureOutcome

for _, o := range adventureOutcomes {
cumulative += o.Weight
if randVal < cumulative {
chosen = o
break
}
}

amount := 0
if chosen.Max > 0 {
amount = rand.Intn(chosen.Max-chosen.Min+1) + chosen.Min
store.GlobalEconomy.AddBalance(jid, amount)
}

store.GlobalEconomy.SetCooldown(jid, "adventure")

var text string
if amount > 0 {
text = fmt.Sprintf(`╭ᩙ╌͜╍͜╌͜─͜─͜─͜─፝͜「 `+"`🗺️ AVENTURA`"+` 」፝͜─͜─͜─͜─͜╌͜━㍃゙⃞
> ㅤ
🫆֟፝𑑓᮫̤٘𑑓̸̷ %s
🫆֟፝𑑓᮫̤٘𑑓̸̷ 💰 *Ganaste:* %d cometas
╰᳕╌͜╍͜╌͜─͜─͜─͜─͜─͜─⵿֟፝─͜─͜─͜─͜─͜─͜─͜─͜╌͜━㍃゙⃞`, chosen.Text, amount)
} else {
text = fmt.Sprintf(`╭ᩙ╌͜╍͜╌͜─͜─͜─͜─፝͜「 `+"`🗺️ AVENTURA`"+` 」፝͜─͜─͜─͜─͜╌͜━㍃゙⃞
> ㅤ
🫆֟፝𑑓᮫̤٘𑑓̸̷ %s
🫆֟፝𑑓᮫̤٘𑑓̸̷ 💀 *No ganaste nada*
╰᳕╌͜╍͜╌͜─͜─͜─͜─͜─͜─⵿֟፝─͜─͜─͜─͜─͜─͜─͜─͜╌͜━㍃゙⃞`, chosen.Text)
}

client.SendMessage(context.Background(), evt.Info.Chat, &waE2E.Message{
ExtendedTextMessage: &waE2E.ExtendedTextMessage{
Text: proto.String(text),
ContextInfo: &waE2E.ContextInfo{QuotedMessage: evt.Message},
},
})
}
