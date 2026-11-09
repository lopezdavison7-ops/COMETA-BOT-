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
Name:        "bal",
Aliases:     []string{"balance", "saldo", "dinero"},
Description: "muestra tus cometas o los de un usuario mencionado",
Category:    "𝐄⵿𝐜𝵿𝐨⵿𝐧𝐨⵿𝐦⵿𝐢⵿𝐚",
Handler:     Bal,
})
}

// getMentions extrae los JIDs mencionados del ContextInfo del mensaje
func getMentions(evt *events.Message) []types.JID {
if evt.Message == nil {
return nil
}

var ctx *waE2E.ContextInfo
if ext := evt.Message.GetExtendedTextMessage(); ext != nil {
ctx = ext.GetContextInfo()
} else if img := evt.Message.GetImageMessage(); img != nil {
ctx = img.GetContextInfo()
} else if vid := evt.Message.GetVideoMessage(); vid != nil {
ctx = vid.GetContextInfo()
} else if doc := evt.Message.GetDocumentMessage(); doc != nil {
ctx = doc.GetContextInfo()
}

if ctx != nil && len(ctx.MentionedJID) > 0 {
var jids []types.JID
for _, jidStr := range ctx.MentionedJID {
jid, err := types.ParseJID(jidStr)
if err == nil {
jids = append(jids, jid)
}
}
return jids
}
return nil
}

func Bal(evt *events.Message, client *whatsmeow.Client, _ []string) {
if store.GlobalEconomy == nil {
return
}

targetJID := evt.Info.Sender
mentions := getMentions(evt)

// Si hay menciones, usamos la primera como objetivo
if len(mentions) > 0 {
targetJID = mentions[0]
}

balance, err := store.GlobalEconomy.GetBalance(targetJID.String())
if err != nil {
balance = 0
}

userName := targetJID.User
if evt.Info.PushName != "" && targetJID == evt.Info.Sender {
userName = evt.Info.PushName
}

b := strings.Builder{}
b.WriteString("╭ᩙ╌͜╍͜╌͜─͜─͜─͜─፝͜「 `💎 BALANCE` 」፝͜─͜─͜─͜─͜╌͜━㍃゙⃞\n> ㅤ\n")
b.WriteString(fmt.Sprintf("🫆֟፝𑑓᮫̤٘𑑓̸̷ Usuario: *%s*\n", userName))
b.WriteString(fmt.Sprintf("🫆֟፝𑑓᮫̤٘𑑓̸̷ Cometas: *%d* ☄️\n", balance))
b.WriteString("╰᳕╌͜╍͜╌͜─͜─͜─͜─͜─͜─⵿֟፝─͜─͜─͜─͜─͜─͜─͜─͜╌͜━㍃゙⃞")

client.SendMessage(context.Background(), evt.Info.Chat, &waE2E.Message{
ExtendedTextMessage: &waE2E.ExtendedTextMessage{
Text: proto.String(b.String()),
},
})
}
