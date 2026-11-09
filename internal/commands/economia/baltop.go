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
Name:        "baltop",
Aliases:     []string{"top", "leaderboard", "ricos"},
Description: "muestra el ranking de los usuarios con más cometas",
Category:    "𝐄⵿𝐜𝵿𝐨⵿𝐧𝐨⵿𝐦⵿𝐢⵿𝐚",
Handler:     BalTop,
})
}

func BalTop(evt *events.Message, client *whatsmeow.Client, _ []string) {
if store.GlobalEconomy == nil {
return
}

tops, err := store.GlobalEconomy.GetTopUsers(10)
if err != nil || len(tops) == 0 {
client.SendMessage(context.Background(), evt.Info.Chat, &waE2E.Message{
ExtendedTextMessage: &waE2E.ExtendedTextMessage{
Text: proto.String("❌ No hay datos en el ranking aún."),
},
})
return
}

b := strings.Builder{}
b.WriteString("╭ᩙ╌͜╍͜╌͜─͜─͜─͜─፝͜「 `🏆 TOP COMETAS` 」፝͜─͜─͜─͜─͜╌͜━㍃゙⃞\n> ㅤ\n")

medals := []string{"🥇", "🥈", "🥉", "4️⃣", "5️⃣", "6️⃣", "7️⃣", "8️⃣", "9️⃣", "🔟"}
mentions := []string{}

for i, top := range tops {
medal := "🔹"
if i < len(medals) {
medal = medals[i]
}
user := strings.Split(top.JID, "@")[0]
b.WriteString(fmt.Sprintf("%s *%d.* @%s → *%d* ☄️\n", medal, top.Rank, user, top.Balance))
mentions = append(mentions, top.JID)
}

b.WriteString("╰᳕╌͜╍͜╌͜─͜─͜─͜─͜─͜─⵿֟፝─͜─͜─͜─͜─͜─͜─͜─͜╌͜━㍃゙⃞")

client.SendMessage(context.Background(), evt.Info.Chat, &waE2E.Message{
ExtendedTextMessage: &waE2E.ExtendedTextMessage{
Text: proto.String(b.String()),
ContextInfo: &waE2E.ContextInfo{
MentionedJID: mentions,
},
},
})
}
