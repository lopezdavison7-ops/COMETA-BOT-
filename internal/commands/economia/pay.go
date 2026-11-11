package economia

import (
"context"
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
Name:        "pay",
Aliases:     []string{"dar", "transferir", "enviar"},
Description: "transfiere cometas a otro usuario",
Category:    "𝐄⵿𝐜𝵿𝐨⵿𝐧𝐨⵿𝐦⵿𝐢⵿𝐚",
Handler:     Pay,
})
}

func Pay(evt *events.Message, client *whatsmeow.Client, args []string) {
if store.GlobalEconomy == nil {
return
}

victimJID := getMentionedJID(evt)
if victimJID == "" || len(args) < 2 {
client.SendMessage(context.Background(), evt.Info.Chat, &waE2E.Message{
ExtendedTextMessage: &waE2E.ExtendedTextMessage{
Text: proto.String("❌ Usá: `/pay @usuario cantidad`. Ejemplo: `/pay @juan 500`"),
ContextInfo: &waE2E.ContextInfo{QuotedMessage: evt.Message},
},
})
return
}

amount, err := strconv.Atoi(args[len(args)-1])
if err != nil || amount <= 0 {
client.SendMessage(context.Background(), evt.Info.Chat, &waE2E.Message{
ExtendedTextMessage: &waE2E.ExtendedTextMessage{
Text: proto.String("❌ La cantidad debe ser un número mayor a 0."),
ContextInfo: &waE2E.ContextInfo{QuotedMessage: evt.Message},
},
})
return
}

senderJID := evt.Info.Sender.String()

if senderJID == victimJID {
client.SendMessage(context.Background(), evt.Info.Chat, &waE2E.Message{
ExtendedTextMessage: &waE2E.ExtendedTextMessage{
Text: proto.String("❌ No te podés enviar dinero a vos mismo. 😂"),
ContextInfo: &waE2E.ContextInfo{QuotedMessage: evt.Message},
},
})
return
}

msg, err := store.GlobalEconomy.Pay(senderJID, victimJID, amount)

var prefix string
if err != nil {
prefix = "❌ "
} else {
prefix = "✅ "
}

b := strings.Builder{}
b.WriteString("╭ᩙ╌͜╍͜╌͜─͜─͜─͜─፝͜「 `💸 PAY` 」፝͜─͜─͜─͜─͜╌͜━㍃゙⃞\n> ㅤ\n")
b.WriteString("🫆֟፝𑑓᮫̤٘𑑓̸̷ " + prefix + msg + "\n")
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
