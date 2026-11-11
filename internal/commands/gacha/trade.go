package gacha

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
Name:        "trade",
Aliases:     []string{"cambiar", "intercambiar"},
Description: "intercambia un personaje con otro usuario",
Category:    "𝐆⵿𝐚⵿𝐜⵿𝐡⵿𝐚",
Handler:     TradeCmd,
})
}

func TradeCmd(evt *events.Message, client *whatsmeow.Client, args []string) {
if len(args) < 3 {
client.SendMessage(context.Background(), evt.Info.Chat, &waE2E.Message{
ExtendedTextMessage: &waE2E.ExtendedTextMessage{
Text: proto.String("❌ Usá: `/trade @usuario mi_personaje su_personaje`\nEjemplo: `/trade @juan asuna mikasa`"),
ContextInfo: &waE2E.ContextInfo{QuotedMessage: evt.Message},
},
})
return
}

// Obtener el JID mencionado
mentionedJID := getMentionedJID(evt)
if mentionedJID == "" {
client.SendMessage(context.Background(), evt.Info.Chat, &waE2E.Message{
ExtendedTextMessage: &waE2E.ExtendedTextMessage{
Text: proto.String("❌ Debés mencionar al usuario con el que querés intercambiar."),
ContextInfo: &waE2E.ContextInfo{QuotedMessage: evt.Message},
},
})
return
}

senderJID := evt.Info.Sender.String()
if senderJID == mentionedJID {
client.SendMessage(context.Background(), evt.Info.Chat, &waE2E.Message{
ExtendedTextMessage: &waE2E.ExtendedTextMessage{
Text: proto.String("❌ No podés intercambiar con vos mismo, chacho 😂"),
ContextInfo: &waE2E.ContextInfo{QuotedMessage: evt.Message},
},
})
return
}

// Los últimos dos argumentos son los nombres
senderQuery := strings.ToLower(strings.Join(args[1:len(args)-1], " "))
receiverQuery := strings.ToLower(args[len(args)-1])

sentChar, receivedChar, err := store.TradeCharacter(senderJID, mentionedJID, senderQuery, receiverQuery)
if err != nil {
client.SendMessage(context.Background(), evt.Info.Chat, &waE2E.Message{
ExtendedTextMessage: &waE2E.ExtendedTextMessage{
Text: proto.String("❌ " + err.Error()),
ContextInfo: &waE2E.ContextInfo{QuotedMessage: evt.Message},
},
})
return
}

text := fmt.Sprintf(
"╭〔 🔄 𝐈𝐍𝐓𝐄𝐑𝐂𝐀𝐌𝐁𝐈𝐎 〕⬣\n┃\n┃ ✅ Intercambio exitoso!\n┃\n┃ Vos diste: *%s*\n┃ (@%s)\n┃\n┃ Vos recibiste: *%s*\n┃ (@%s)\n┃\n╰━━━━━━━━━━━━━━━━⬣",
sentChar.Name, sentChar.Series, receivedChar.Name, receivedChar.Series,
)

client.SendMessage(context.Background(), evt.Info.Chat, &waE2E.Message{
ExtendedTextMessage: &waE2E.ExtendedTextMessage{
Text: proto.String(text),
ContextInfo: &waE2E.ContextInfo{QuotedMessage: evt.Message},
},
})
}

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
return ""
}
