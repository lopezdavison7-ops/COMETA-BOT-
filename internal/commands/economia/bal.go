package economia

import (
"context"
"fmt"
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
Name:        "bal",
Aliases:     []string{"balance", "saldo", "dinero"},
Description: "muestra tu saldo de cometas",
Category:    "𝐄⵿𝐜𝵿𝐨⵿𝐧𝐨⵿𝐦⵿𝐢⵿𝐚",
Handler:     Bal,
})
}

func Bal(evt *events.Message, client *whatsmeow.Client, _ []string) {
if store.GlobalEconomy == nil {
return
}

jid := evt.Info.Sender.String()
balance := store.GlobalEconomy.GetBalance(jid)
bank := store.GlobalEconomy.GetBank(jid)
level, xp := store.GlobalEconomy.GetProfile(jid)
needXP := level * 100

b := strings.Builder{}
b.WriteString("╭ᩙ╌͜╍͜╌͜─͜─͜─͜─፝͜「 `💰 BALANCE` 」፝͜─͜─͜─͜─͜╌͜━㍃゙⃞\n> ㅤ\n")
b.WriteString(fmt.Sprintf("🫆֟፝𑑓᮫̤٘𑑓̸̷ *Billetera:* %d cometas\n", balance))
b.WriteString(fmt.Sprintf("🏦֟፝𑑓᮫̤٘𑑓̸̷ *Banco:* %d cometas\n", bank))
b.WriteString(fmt.Sprintf("📊֟፝𑑓᮫̤٘𑑓̸̷ *Nivel:* %d (%d/%d XP)\n", level, xp, needXP))
b.WriteString("╰᳕╌͜╍͜╌͜─͜─͜─͜─͜─͜─⵿֟፝─͜─͜─͜─͜─͜─͜─͜─͜╌͜━㍃゙⃞")

client.SendMessage(context.Background(), evt.Info.Chat, &waE2E.Message{
ExtendedTextMessage: &waE2E.ExtendedTextMessage{
Text: proto.String(b.String()),
ContextInfo: &waE2E.ContextInfo{
QuotedMessage: evt.Message,
},
},
})
}
