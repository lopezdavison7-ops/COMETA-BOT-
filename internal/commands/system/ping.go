package system

import (
"context"
"strings"
"time"

"go.mau.fi/whatsmeow"
"go.mau.fi/whatsmeow/proto/waE2E"
"go.mau.fi/whatsmeow/types/events"
"google.golang.org/protobuf/proto"

"cometa-bot/internal/registry"
)

func init() {
registry.Register(&registry.Command{
Name:        "ping",
Aliases:     []string{"p"},
Description: "calcula la lentitud del bot",
Category:    "𝐒⵿𝐢⵿𝐬⵿𝐭⵿𝐞⵿𝐦⵿𝐚",
Handler:     Ping,
})
}

func Ping(evt *events.Message, client *whatsmeow.Client, _ []string) {
start := time.Now()

b := strings.Builder{}
b.Grow(300)
b.WriteString("╭ᩙ╌͜╍͜╌͜─͜─͜─͜─፝͜「 `🏓 PING` 」፝͜─͜─͜─͜─͜╌͜━㍃゙⃞\n")
b.WriteString("> ㅤ\n")
b.WriteString("🫆֟፝𑑓᮫̤٘𑑓̸̷ ᴄᴀʟᴄᴜʟᴀɴᴅᴏ ʟᴀᴛᴇɴᴄɪᴀ...\n")
b.WriteString("╰᳕╌͜╍͜╌͜─͜─͜─͜─͜─͜─⵿֟፝─͜─͜─͜─͜─͜─͜─͜─͜╌͜━㍃゙⃞")

msgID, err := client.SendMessage(context.Background(), evt.Info.Chat, &waE2E.Message{
ExtendedTextMessage: &waE2E.ExtendedTextMessage{
Text: proto.String(b.String()),
},
})
if err != nil {
return
}

latency := time.Since(start)

b.Reset()
b.Grow(300)
b.WriteString("╭ᩙ╌͜╍͜╌͜─͜─͜─͜─፝͜「 `🏓 PONG` 」፝͜─͜─͜─͜─͜╌͜━㍃゙⃞\n")
b.WriteString("> ㅤ\n")
b.WriteString("🫆֟፝𑑓᮫̤٘𑑓̸̷ ʟᴀᴛᴇɴᴄɪᴀ: ")
b.WriteString(latency.Round(time.Millisecond).String())
b.WriteString("\n🫆֟፝𑑓᮫̤٘𑑓̸̷ ᴇsᴛᴀᴅᴏ: ᴀᴄᴛɪᴠᴏ ☄️\n")
b.WriteString("╰᳕╌͜╍͜╌͜─͜─͜─͜─͜─͜─⵿֟፝─͜─͜─͜─͜─͜─͜─͜─͜╌͜━㍃゙⃞")

editedMsg := client.BuildEdit(evt.Info.Chat, msgID.ID, &waE2E.Message{
ExtendedTextMessage: &waE2E.ExtendedTextMessage{
Text: proto.String(b.String()),
},
})

client.SendMessage(context.Background(), evt.Info.Chat, editedMsg)
}
