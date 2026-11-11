package economia

import (
"context"
"fmt"
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
Name:        "dep",
Aliases:     []string{"depositar", "deposit"},
Description: "deposita cometas en el banco (a salvo de robos)",
Category:    "𝐄⵿𝐜𝵿𝐨⵿𝐧𝐨⵿𝐦⵿𝐢⵿𝐚",
Handler:     Dep,
})
registry.Register(&registry.Command{
Name:        "with",
Aliases:     []string{"retirar", "withdraw"},
Description: "retira cometas del banco a tu billetera",
Category:    "𝐄⵿𝐜𝵿𝐨⵿𝐧𝐨⵿𝐦⵿𝐢⵿𝐚",
Handler:     With,
})
}

func Dep(evt *events.Message, client *whatsmeow.Client, args []string) {
if len(args) == 0 {
sendEcoError(client, evt, "Usá: /dep <cantidad> o /dep all")
return
}
jid := evt.Info.Sender.String()
balance := store.GlobalEconomy.GetBalance(jid)
amount := parseAmount(args[0], balance)

if amount <= 0 {
sendEcoError(client, evt, "Monto inválido.")
return
}
if amount > balance {
sendEcoError(client, evt, "No tenés tantos cometas en la billetera.")
return
}

store.GlobalEconomy.AddBalance(jid, -amount)
store.GlobalEconomy.AddBank(jid, amount)

sendEcoSuccess(client, evt, fmt.Sprintf("🏦 Depositaste *%d* cometas en el banco.", amount))
}

func With(evt *events.Message, client *whatsmeow.Client, args []string) {
if len(args) == 0 {
sendEcoError(client, evt, "Usá: /with <cantidad> o /with all")
return
}
jid := evt.Info.Sender.String()
bank := store.GlobalEconomy.GetBank(jid)
amount := parseAmount(args[0], bank)

if amount <= 0 {
sendEcoError(client, evt, "Monto inválido.")
return
}
if amount > bank {
sendEcoError(client, evt, "No tenés tantos cometas en el banco.")
return
}

store.GlobalEconomy.AddBank(jid, -amount)
store.GlobalEconomy.AddBalance(jid, amount)

sendEcoSuccess(client, evt, fmt.Sprintf("💵 Retiraste *%d* cometas del banco.", amount))
}

func parseAmount(arg string, max int) int {
if strings.ToLower(arg) == "all" || arg == "todo" {
return max
}
amount, err := strconv.Atoi(arg)
if err != nil {
return 0
}
return amount
}

func sendEcoError(client *whatsmeow.Client, evt *events.Message, text string) {
client.SendMessage(context.Background(), evt.Info.Chat, &waE2E.Message{
ExtendedTextMessage: &waE2E.ExtendedTextMessage{
Text: proto.String("❌ " + text),
ContextInfo: &waE2E.ContextInfo{QuotedMessage: evt.Message},
},
})
}

func sendEcoSuccess(client *whatsmeow.Client, evt *events.Message, text string) {
client.SendMessage(context.Background(), evt.Info.Chat, &waE2E.Message{
ExtendedTextMessage: &waE2E.ExtendedTextMessage{
Text: proto.String("✅ " + text),
ContextInfo: &waE2E.ContextInfo{QuotedMessage: evt.Message},
},
})
}
