package system

import (
"context"
"strings"

"go.mau.fi/whatsmeow"
"go.mau.fi/whatsmeow/proto/waE2E"
"go.mau.fi/whatsmeow/types/events"
"google.golang.org/protobuf/proto"

"cometa-bot/internal/registry"
)

func init() {
registry.Register(&registry.Command{
Name:        "canal",
Aliases:     []string{"channel", "soporte"},
Description: "muestra el enlace al canal oficial de Cometa",
Category:    "𝐒⵿𝐢⵿𝐬⵿𝐭⵿𝐞⵿𝐦⵿𝐚",
Handler:     Canal,
})
}

func Canal(evt *events.Message, client *whatsmeow.Client, _ []string) {
b := strings.Builder{}
b.WriteString("╭ᩙ╌͜╍͜╌͜─͜─͜─͜─፝͜「 `📢 CANAL OFICIAL` 」፝͜─͜─͜─͜─͜╌͜━㍃゙⃞\n> ㅤ\n")
b.WriteString("🫆֟፝𑑓᮫̤٘𑑓̸̷ ¡Únete para novedades, actualizaciones\n")
b.WriteString("🫆֟፝𑑓᮫̤٘𑑓̸̷ y soporte de *Cometa Bot*!\n> ㅤ\n")
b.WriteString("🔗 https://whatsapp.com/channel/0029Vb8eeKGG3R3kwBcZdp2Q\n")
b.WriteString("╰᳕╌͜╍͜╌͜─͜─͜─͜─͜─͜─⵿֟፝─͜─͜─͜─͜─͜─͜─͜─͜╌͜━㍃゙⃞")

client.SendMessage(context.Background(), evt.Info.Chat, &waE2E.Message{
ExtendedTextMessage: &waE2E.ExtendedTextMessage{
Text: proto.String(b.String()),
},
})
}
