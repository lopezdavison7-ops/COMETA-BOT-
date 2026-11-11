package economia

import (
"context"
"strings"

"go.mau.fi/whatsmeow"
"go.mau.fi/whatsmeow/proto/waE2E"
"go.mau.fi/whatsmeow/types/events"
"google.golang.org/protobuf/proto"

"cometa-bot/internal/registry"
"cometa-bot/internal/store"
)

func getMentionedJID(evt *events.Message) string {
if evt.Message == nil {
return ""
}
if ext := evt.Message.GetExtendedTextMessage(); ext != nil {
if ctx := ext.GetContextInfo(); ctx != nil {
mentions := ctx.GetMentionedJID()
if len(mentions) > 0 {
return mentions[0]
}
}
}
if img := evt.Message.GetImageMessage(); img != nil {
if ctx := img.GetContextInfo(); ctx != nil {
mentions := ctx.GetMentionedJID()
if len(mentions) > 0 {
return mentions[0]
}
}
}
return ""
}

func init() {
registry.Register(&registry.Command{
Name:        "rob",
Aliases:     []string{"robar", "asaltar"},
Description: "intenta robarle cometas a otro usuario",
Category:    "𝐄⵿𝐜𝵿𝐨⵿𝐧𝐨⵿𝐦⵿𝐢⵿𝐚",
Handler:     Rob,
})
}

func Rob(evt *events.Message, client *whatsmeow.Client, _ []string) {
if store.GlobalEconomy == nil {
return
}

victimJID := getMentionedJID(evt)
if victimJID == "" {
client.SendMessage(context.Background(), evt.Info.Chat, &waE2E.Message{
ExtendedTextMessage: &waE2E.ExtendedTextMessage{
Text: proto.String("❌ Debés mencionar a alguien para robarle. Ejemplo: `/rob @usuario`"),
ContextInfo: &waE2E.ContextInfo{QuotedMessage: evt.Message},
},
})
return
}

thiefJID := evt.Info.Sender.String()

if thiefJID == victimJID {
client.SendMessage(context.Background(), evt.Info.Chat, &waE2E.Message{
ExtendedTextMessage: &waE2E.ExtendedTextMessage{
Text: proto.String("❌ No te podés robar a vos mismo, chacho. 😂"),
ContextInfo: &waE2E.ContextInfo{QuotedMessage: evt.Message},
},
})
return
}

msg, err := store.GlobalEconomy.Rob(thiefJID, victimJID)

var prefix string
if err != nil {
if err.Error() == "cooldown" {
prefix = "⏳ "
} else if err.Error() == "poor" {
prefix = "😂 "
} else {
prefix = "👮‍♂️ "
}
} else {
prefix = "🥷 "
}

b := strings.Builder{}
b.WriteString("╭ᩙ╌͜╍͜╌͜─͜─͜─͜─፝͜「 `🥷 ROB` 」፝͜─͜─͜─͜─͜╌͜━㍃゙⃞\n> ㅤ\n")
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
