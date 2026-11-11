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
Name:        "rank",
Aliases:     []string{"nivel", "lvl", "xp"},
Description: "muestra tu nivel y experiencia actual",
Category:    "𝐄⵿𝐜𝵿𝐨⵿𝐧𝐨⵿𝐦⵿𝐢⵿𝐚",
Handler:     Rank,
})
}

func Rank(evt *events.Message, client *whatsmeow.Client, _ []string) {
if store.GlobalEconomy == nil { return }

jid := evt.Info.Sender.String()
level, xp := store.GlobalEconomy.GetProfile(jid)
xpNeeded := level * 100
progress := xp
if progress > xpNeeded { progress = xpNeeded }

stars := ""
for i := 0; i < level; i++ { stars += "⭐" }
if stars == "" { stars = "🌱" }

b := strings.Builder{}
b.WriteString("╭ᩙ╌͜╍͜╌͜─͜─͜─͜─፝͜「 `📊 RANK` 」፝͜─͜─͜─͜─͜╌͜━㍃゙⃞\n> ㅤ\n")
b.WriteString("🫆֟፝𑑓᮫̤٘𑑓̸̷ *Nivel:* " + stars + " (" + strconv.Itoa(level) + ")\n")
b.WriteString("🫆֟፝𑑓᮫̤٘𑑓̸̷ *XP:* " + strconv.Itoa(progress) + " / " + strconv.Itoa(xpNeeded) + "\n")
b.WriteString("╰᳕╌͜╍͜╌͜─͜─͜─͜─͜─͜─⵿֟፝─͜─͜─͜─͜─͜─͜─͜─͜╌͜━㍃゙⃞")

client.SendMessage(context.Background(), evt.Info.Chat, &waE2E.Message{
ExtendedTextMessage: &waE2E.ExtendedTextMessage{
Text: proto.String(b.String()),
ContextInfo: &waE2E.ContextInfo{QuotedMessage: evt.Message},
},
})
}
