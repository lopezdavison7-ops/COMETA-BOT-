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
Name:        "inv",
Aliases:     []string{"inventario", "bag", "mochila"},
Description: "muestra tu inventario de personajes",
Category:    "𝐆⵿𝐚⵿𝐜⵿𝐡⵿𝐚",
Handler:     Inv,
})
}

func Inv(evt *events.Message, client *whatsmeow.Client, _ []string) {
target := evt.Info.Sender.String()
var mentions []string

if m := getMentionedJID(evt); m != "" {
target = m
mentions = append(mentions, m)
}

chars := store.GetUserInventory(target)

if len(chars) == 0 {
client.SendMessage(context.Background(), evt.Info.Chat, &waE2E.Message{
ExtendedTextMessage: &waE2E.ExtendedTextMessage{
Text: proto.String("╭〔 🎴 𝐈𝐍𝐕𝐄𝐍𝐓𝐀𝐑𝐈𝐎 〕⬣\n┃\n┃ Aún no hay personajes aquí.\n┃ Usá */rw* para sacar tu primera waifu!\n┃\n╰━━━━━━━━━━━━━━━━⬣"),
ContextInfo: &waE2E.ContextInfo{
QuotedMessage: evt.Message,
MentionedJID:  mentions,
},
},
})
return
}

b := strings.Builder{}
b.WriteString("╭〔 🎴 𝐈𝐍𝐕𝐄𝐍𝐓𝐀𝐑𝐈𝐎 〕⬣\n┃\n")

total := 0
maxShow := 15
for i, c := range chars {
total += c.Value
if i >= maxShow {
break
}
gender := "♀️"
if c.Gender == "Masculino" {
gender = "♂️"
}
b.WriteString(fmt.Sprintf("┃ *%d.* %s %s\n┃    📺 %s | 💰 %d cometas\n", i+1, gender, c.Name, c.Series, c.Value))
}

if len(chars) > maxShow {
b.WriteString(fmt.Sprintf("┃ ...y %d personajes más\n", len(chars)-maxShow))
}

b.WriteString(fmt.Sprintf("┃\n┃ 👥 Total: %d personajes\n┃ 💎 Valor total: %d cometas\n", len(chars), total))
b.WriteString("╰━━━━━━━━━━━━━━━━⬣")

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
