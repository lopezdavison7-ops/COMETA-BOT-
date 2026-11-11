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
Name:        "slut",
Aliases:     []string{},
Description: "gana cometas de forma... creativa",
Category:    "𝐄⵿𝐜𝵿𝐨⵿𝐧𝐨⵿𝐦⵿𝐢⵿𝐚",
Handler:     Slut,
})
}

type slutOutcome struct {
Text   string
Min    int
Max    int
Weight int
}

var slutOutcomes = []slutOutcome{
{"Un cliente generoso te dio propina", 50, 150, 30},
{"Trabajaste toda la noche", 100, 250, 25},
{"Un VIP te pagó extra", 200, 400, 20},
{"Ganaste un concurso de baile", 300, 600, 15},
{"Te contrataron para un evento exclusivo", 500, 1000, 8},
{"La policía te multó", -100, -50, 2},
}

func Slut(evt *events.Message, client *whatsmeow.Client, _ []string) {
jid := evt.Info.Sender.String()

if remaining, ok := store.GlobalEconomy.CheckCooldown(jid, "slut", 1200); !ok {
mins := remaining / 60
secs := remaining % 60
sendEcoError(client, evt, fmt.Sprintf("⏳ Ya trabajaste recientemente.\nEsperá *%d min %d seg* para volver a trabajar.", mins, secs))
return
}

totalWeight := 0
for _, o := range slutOutcomes {
totalWeight += o.Weight
}

randVal := rand.Intn(totalWeight)
cumulative := 0
var chosen slutOutcome

for _, o := range slutOutcomes {
cumulative += o.Weight
if randVal < cumulative {
chosen = o
break
}
}

amount := rand.Intn(chosen.Max-chosen.Min+1) + chosen.Min

if amount < 0 {
balance := store.GlobalEconomy.GetBalance(jid)
if balance < -amount {
amount = -balance
}
store.GlobalEconomy.AddBalance(jid, amount)
} else {
store.GlobalEconomy.AddBalance(jid, amount)
}

store.GlobalEconomy.SetCooldown(jid, "slut")

var text string
if amount >= 0 {
text = fmt.Sprintf(`╭ᩙ╌͜╍͜╌͜─͜─͜─͜─፝͜「 `+"`💋 SLUT`"+` 」፝͜─͜─͜─͜─͜╌͜━㍃゙⃞
> ㅤ
🫆֟፝𑑓᮫̤٘𑑓̸̷ %s
🫆֟፝𑑓᮫̤٘𑑓̸̷ 💰 *Ganaste:* %d cometas
╰᳕╌͜╍͜╌͜─͜─͜─͜─͜─͜─⵿֟፝─͜─͜─͜─͜─͜─͜─͜─͜╌͜━㍃゙⃞`, chosen.Text, amount)
} else {
text = fmt.Sprintf(`╭ᩙ╌͜╍͜╌͜─͜─͜─͜─፝͜「 `+"`💋 SLUT`"+` 」፝͜─͜─͜─͜─͜╌͜━㍃゙⃞
> ㅤ
🫆֟፝𑑓᮫̤٘𑑓̸̷ %s
🫆֟፝𑑓᮫̤٘𑑓̸̷ 💸 *Perdiste:* %d cometas
╰᳕╌͜╍͜╌͜─͜─͜─͜─͜─͜─⵿֟፝─͜─͜─͜─͜─͜─͜─͜─͜╌͜━㍃゙⃞`, chosen.Text, -amount)
}

client.SendMessage(context.Background(), evt.Info.Chat, &waE2E.Message{
ExtendedTextMessage: &waE2E.ExtendedTextMessage{
Text: proto.String(text),
ContextInfo: &waE2E.ContextInfo{QuotedMessage: evt.Message},
},
})
}
