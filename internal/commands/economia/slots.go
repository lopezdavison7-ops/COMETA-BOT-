package economia

import (
"context"
"fmt"
"math/rand"
"strconv"

"go.mau.fi/whatsmeow"
"go.mau.fi/whatsmeow/proto/waE2E"
"go.mau.fi/whatsmeow/types/events"
"google.golang.org/protobuf/proto"

"cometa-bot/internal/registry"
"cometa-bot/internal/store"
)

func init() {
registry.Register(&registry.Command{
Name:        "slots",
Aliases:     []string{"tragamonedas", "slot"},
Description: "juega al tragamonedas",
Category:    "𝐄⵿𝐜𝵿𝐨⵿𝐧𝐨⵿𝐦⵿𝐢⵿𝐚",
Handler:     Slots,
})
}

var slotSymbols = []string{"🍒", "🍋", "🍊", "🍇", "⭐", "💎", "7️⃣"}

func Slots(evt *events.Message, client *whatsmeow.Client, args []string) {
if len(args) == 0 {
sendEcoError(client, evt, "Usá: /slots <cantidad>")
return
}

amount, err := strconv.Atoi(args[0])
if err != nil || amount <= 0 {
sendEcoError(client, evt, "Cantidad inválida.")
return
}

jid := evt.Info.Sender.String()
balance := store.GlobalEconomy.GetBalance(jid)

if amount > balance {
sendEcoError(client, evt, "No tenés suficientes cometas.")
return
}

if amount < 10 {
sendEcoError(client, evt, "Mínimo 10 cometas.")
return
}

store.GlobalEconomy.AddBalance(jid, -amount)

bonus := 0
if store.HasItem(jid, "luck_charm") {
bonus = 10
store.UseItem(jid, "luck_charm")
}

r1 := slotSymbols[rand.Intn(len(slotSymbols))]
r2 := slotSymbols[rand.Intn(len(slotSymbols))]
r3 := slotSymbols[rand.Intn(len(slotSymbols))]

if bonus > 0 && rand.Intn(100) < bonus {
r3 = r1
}

winnings := 0
multiplier := 0

if r1 == r2 && r2 == r3 {
switch r1 {
case "7️⃣":
multiplier = 10
case "💎":
multiplier = 7
case "⭐":
multiplier = 5
default:
multiplier = 3
}
} else if r1 == r2 || r2 == r3 || r1 == r3 {
multiplier = 1
}

winnings = amount * multiplier

result := "💀 *PERDISTE*"
if winnings > 0 {
store.GlobalEconomy.AddBalance(jid, winnings)
if multiplier >= 3 {
result = fmt.Sprintf("🎉 *¡JACKPOT!* x%d", multiplier)
} else {
result = fmt.Sprintf("✨ *Ganaste* x%d", multiplier)
}
}

text := fmt.Sprintf(`╭ᩙ╌͜╍͜╌͜─͜─͜─͜─፝͜「 `+"`🎰 SLOTS`"+` 」፝͜─͜─͜─͜─͜╌͜━㍃゙⃞
> ㅤ
🫆֟፝𑑓᮫̤٘𑑓̸̷ *Apuesta:* %d cometas
🫆֟፝𑑓᮫̤٘𑑓̸̷ *Resultado:* %s %s %s
🫆֟፝𑑓᮫̤٘𑑓̸̷ %s
%s
╰᳕╌͜╍͜╌͜─͜─͜─͜─͜─͜─⵿֟፝─͜─͜─͜─͜─͜─͜─͜─͜╌͜━㍃゙⃞`,
amount, r1, r2, r3, result,
func() string {
if winnings > 0 {
return fmt.Sprintf("💰 *Ganaste:* %d cometas", winnings)
}
return ""
}())

client.SendMessage(context.Background(), evt.Info.Chat, &waE2E.Message{
ExtendedTextMessage: &waE2E.ExtendedTextMessage{
Text: proto.String(text),
ContextInfo: &waE2E.ContextInfo{QuotedMessage: evt.Message},
},
})
}
