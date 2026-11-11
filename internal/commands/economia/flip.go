package economia

import (
"fmt"
"math/rand"
"strconv"
"strings"

"go.mau.fi/whatsmeow"
"go.mau.fi/whatsmeow/types/events"

"cometa-bot/internal/registry"
"cometa-bot/internal/store"
)

func init() {
registry.Register(&registry.Command{
Name:        "flip",
Aliases:     []string{"apostar", "caraocruz", "coinflip"},
Description: "apuesta cometas a cara o cruz (50/50)",
Category:    "𝐄⵿𝐜𝵿𝐨⵿𝐧𝐨⵿𝐦⵿𝐢⵿𝐚",
Handler:     Flip,
})
}

func Flip(evt *events.Message, client *whatsmeow.Client, args []string) {
if len(args) < 2 {
sendEcoError(client, evt, "Usá: /flip <cantidad> <cara|cruz>")
return
}

amount, err := strconv.Atoi(args[0])
if err != nil || amount <= 0 {
sendEcoError(client, evt, "La cantidad debe ser un número mayor a 0.")
return
}

choice := strings.ToLower(args[1])
if choice != "cara" && choice != "cruz" {
sendEcoError(client, evt, "Debés elegir *cara* o *cruz*.")
return
}

jid := evt.Info.Sender.String()
balance := store.GlobalEconomy.GetBalance(jid)

if amount > balance {
sendEcoError(client, evt, "No tenés suficientes cometas en la billetera.")
return
}

result := "cara"
if rand.Intn(2) == 1 {
result = "cruz"
}

if choice == result {
store.GlobalEconomy.AddBalance(jid, amount)
sendEcoSuccess(client, evt, fmt.Sprintf("🎉 ¡Ganaste! Salió *%s*. Ganaste *%d* cometas.", result, amount))
} else {
store.GlobalEconomy.AddBalance(jid, -amount)
sendEcoError(client, evt, fmt.Sprintf("💀 ¡Perdiste! Salió *%s*. Perdiste *%d* cometas.", result, amount))
}
}
