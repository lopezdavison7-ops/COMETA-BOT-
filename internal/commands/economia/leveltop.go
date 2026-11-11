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
Name:        "leveltop",
Aliases:     []string{"rangotop", "toptop"},
Description: "muestra el ranking de niveles",
Category:    "𝐄⵿𝐜𝵿𝐨⵿𝐧𝐨⵿𝐦⵿𝐢⵿𝐚",
Handler:     LevelTop,
})
}

func LevelTop(evt *events.Message, client *whatsmeow.Client, _ []string) {
if store.GlobalEconomy == nil {
return
}

tops := store.GlobalEconomy.GetLevelTop(10)
if len(tops) == 0 {
client.SendMessage(context.Background(), evt.Info.Chat, &waE2E.Message{
ExtendedTextMessage: &waE2E.ExtendedTextMessage{
Text: proto.String("Aún no hay usuarios con nivel registrado."),
ContextInfo: &waE2E.ContextInfo{QuotedMessage: evt.Message},
},
})
return
}

b := strings.Builder{}
b.WriteString("╭ᩙ╌͜╍͜╌͜─͜─͜─͜─፝͜「 `🏆 LEVEL TOP` 」፝͜─͜─͜─͜─͜╌͜━㍃゙⃞\n> ㅤ\n")

var mentions []string
for i, u := range tops {
jid, _ := types.ParseJID(u.JID)
b.WriteString(fmt.Sprintf("🫆֟፝𑑓᮫̤٘𑑓̸̷ *%d.* @%s (Nivel %d)\n", i+1, jid.User, u.Level))
mentions = append(mentions, u.JID) // Agregamos a TODOS a la lista de menciones
}
b.WriteString("╰᳕╌͜╍͜╌͜─͜─͜─͜─͜─͜─⵿֟፝─͜─͜─͜─͜─͜─͜─͜─͜╌͜━㍃゙⃞")

client.SendMessage(context.Background(), evt.Info.Chat, &waE2E.Message{
ExtendedTextMessage: &waE2E.ExtendedTextMessage{
Text: proto.String(b.String()),
ContextInfo: &waE2E.ContextInfo{
QuotedMessage: evt.Message,
MentionedJID:  mentions, // ¡Ahora sí envía todas las menciones correctamente!
},
},
})
}
