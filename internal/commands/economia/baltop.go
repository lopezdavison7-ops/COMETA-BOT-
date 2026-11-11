package economia

import (
"context"
"fmt"
"strings"

"go.mau.fi/whatsmeow"
"go.mau.fi/whatsmeow/proto/waE2E"
"go.mau.fi/whatsmeow/types"
"go.mau.fi/whatsmeow/types/events"
"google.golang.org/protobuf/proto"

"cometa-bot/internal/registry"
"cometa-bot/internal/store"
)

func init() {
registry.Register(&registry.Command{
Name:        "baltop",
Aliases:     []string{"top", "ranking", "topsaldo"},
Description: "muestra el ranking de los más ricos",
Category:    "𝐄⵿𝐜𝵿𝐨⵿𝐧𝐨⵿𝐦⵿𝐢⵿𝐚",
Handler:     BalTop,
})
}

func BalTop(evt *events.Message, client *whatsmeow.Client, _ []string) {
if store.GlobalEconomy == nil {
return
}

tops := store.GlobalEconomy.GetTop(10)
if len(tops) == 0 {
client.SendMessage(context.Background(), evt.Info.Chat, &waE2E.Message{
ExtendedTextMessage: &waE2E.ExtendedTextMessage{
Text: proto.String("Aún no hay usuarios con saldo."),
ContextInfo: &waE2E.ContextInfo{QuotedMessage: evt.Message},
},
})
return
}

b := strings.Builder{}
b.WriteString("╭ᩙ╌͜╍͜╌͜─͜─͜─͜─፝͜「 `🏆 BAL TOP` 」፝͜─͜─͜─͜─͜╌͜━㍃゙⃞\n> ㅤ\n")

var mentions []string
for i, u := range tops {
jid, _ := types.ParseJID(u.JID)
b.WriteString(fmt.Sprintf("🫆֟፝𑑓᮫̤٘𑑓̸̷ *%d.* @%s (%d cometas - Nvl %d)\n", i+1, jid.User, u.Balance, u.Level))
mentions = append(mentions, u.JID)
}
b.WriteString("╰᳕╌͜╍͜╌͜─͜─͜─͜─͜─͜─⵿֟፝─͜─͜─͜─͜─͜─͜─͜─͜╌͜━㍃゙⃞")

client.SendMessage(context.Background(), evt.Info.Chat, &waE2E.Message{
ExtendedTextMessage: &waE2E.ExtendedTextMessage{
Text: proto.String(b.String()),
ContextInfo: &waE2E.ContextInfo{
QuotedMessage: evt.Message,
MentionedJID:  mentions,
},
},
})
}
