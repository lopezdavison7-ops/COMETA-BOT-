package economia

import (
"context"
"fmt"
"math/rand"
"strconv"
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
Name:        "roulette",
Aliases:     []string{"ruleta", "rt"},
Description: "apuesta a color o número en la ruleta",
Category:    "𝐄⵿𝐜𝵿𝐨⵿𝐧𝐨⵿𝐦⵿𝐢⵿𝐚",
Handler:     Roulette,
})
}

func Roulette(evt *events.Message, client *whatsmeow.Client, args []string) {
if len(args) < 2 {
sendEcoError(client, evt, "Usá: /roulette <cantidad> <rojo|negro|verde|número>")
return
}

amount, err := strconv.Atoi(args[0])
if err != nil || amount <= 0 {
sendEcoError(client, evt, "Cantidad inválida.")
return
}

bet := strings.ToLower(args[1])

jid := evt.Info.Sender.String()
balance := store.GlobalEconomy.GetBalance(jid)

if amount > balance {
sendEcoError(client, evt, "No tenés suficientes cometas.")
return
}

store.GlobalEconomy.AddBalance(jid, -amount)

result := rand.Intn(37)
color := "verde"
if result == 0 {
color = "verde"
} else if result%2 == 0 {
color = "negro"
} else {
color = "rojo"
}

winnings := 0
betType := ""

if bet == "rojo" || bet == "r" {
betType = "Rojo"
if color == "rojo" {
winnings = amount * 2
}
} else if bet == "negro" || bet == "n" {
betType = "Negro"
if color == "negro" {
winnings = amount * 2
}
} else if bet == "verde" || bet == "v" {
betType = "Verde"
if color == "verde" {
winnings = amount * 14
}
} else {
num, err := strconv.Atoi(bet)
if err == nil && num >= 0 && num <= 36 {
betType = fmt.Sprintf("Número %d", num)
if num == result {
winnings = amount * 36
}
} else {
sendEcoError(client, evt, "Apuesta inválida. Usá: rojo, negro, verde o un número (0-36)")
store.GlobalEconomy.AddBalance(jid, amount)
return
}
}

if winnings > 0 {
store.GlobalEconomy.AddBalance(jid, winnings)
}

resultStr := fmt.Sprintf("%d (%s)", result, strings.Title(color))

text := fmt.Sprintf(`╭ᩙ╌͜╍͜╌͜─͜─͜─͜─፝͜「 `+"`🎡 RULETA`"+` 」፝͜─͜─͜─͜─͜╌͜━㍃゙⃞
> ㅤ
🫆֟፝𑑓᮫̤٘𑑓̸̷ *Apuesta:* %d cometas a *%s*
🫆֟፝𑑓᮫̤٘𑑓̸̷ *Resultado:* %s
%s
╰᳕╌͜╍͜╌͜─͜─͜─͜─͜─͜─⵿֟፝─͜─͜─͜─͜─͜─͜─͜─͜╌͜━㍃゙⃞`,
amount, betType, resultStr,
func() string {
if winnings > 0 {
return fmt.Sprintf("🎉 *¡Ganaste!* %d cometas", winnings)
}
return "💀 *Perdiste*"
}())

client.SendMessage(context.Background(), evt.Info.Chat, &waE2E.Message{
ExtendedTextMessage: &waE2E.ExtendedTextMessage{
Text: proto.String(text),
ContextInfo: &waE2E.ContextInfo{QuotedMessage: evt.Message},
},
})
}
