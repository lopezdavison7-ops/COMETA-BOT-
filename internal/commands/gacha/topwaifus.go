package gacha

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
Name:        "topwaifus",
Aliases:     []string{"topwaifu", "tw", "mejores"},
Description: "ranking de coleccionistas de waifus",
Category:    "𝐆⵿𝐚⵿𝐜⵿𝐡⵿𝐚",
Handler:     TopWaifus,
})
}

func TopWaifus(evt *events.Message, client *whatsmeow.Client, _ []string) {
tops := store.GetTopCollectors(10)

if len(tops) == 0 {
client.SendMessage(context.Background(), evt.Info.Chat, &waE2E.Message{
ExtendedTextMessage: &waE2E.ExtendedTextMessage{
Text: proto.String("╭〔 👑 𝐓𝐎𝐏 𝐖𝐀𝐈𝐅𝐔𝐒 〕⬣\n┃\n┃ Aún no hay waifus reclamadas.\n┃ Usá */rw* y */claim* para entrar al top!\n┃\n╰━━━━━━━━━━━━━━━━⬣"),
ContextInfo: &waE2E.ContextInfo{QuotedMessage: evt.Message},
},
})
return
}

medals := []string{"🥇", "🥈", "🥉"}

b := strings.Builder{}
b.WriteString("╭〔 👑 𝐓𝐎𝐏 𝐖𝐀𝐈𝐅𝐔𝐒 〕⬣\n┃\n")

var mentions []string
for i, s := range tops {
rank := fmt.Sprintf("*%d.*", i+1)
if i < 3 {
rank = medals[i]
}

ownerTag := "???"
if jid, err := types.ParseJID(s.JID); err == nil {
ownerTag = "@" + jid.User
mentions = append(mentions, s.JID)
}

unidad := "waifus"
if s.Count == 1 {
unidad = "waifu"
}

b.WriteString(fmt.Sprintf("┃ %s %s — %d %s\n┃    💎 Valor: %d cometas\n┃    ⭐ Mejor: %s (%d)\n", rank, ownerTag, s.Count, unidad, s.TotalValue, s.Best.Name, s.Best.Value))
}
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
