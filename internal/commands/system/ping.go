package system

import (
"context"
"fmt"
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
Aliases:     []string{"p", "pong"},
Description: "verifica la latencia del bot",
Category:    "𝐒⵿𝐢⵿𝐬⵿𝐭⵿𝐞⵿𝐦⵿𝐚",
Handler:     Ping,
})
}

func Ping(evt *events.Message, client *whatsmeow.Client, _ []string) {
latency := time.Now().UnixMilli() - evt.Info.Timestamp.UnixMilli()
if latency < 0 {
latency = 0
}

text := fmt.Sprintf("╭ᩙ╌͜╍͜╌͜─͜─͜─͜─፝͜「 `🏓 PONG` 」፝͜─͜─͜─͜─͜╌͜━㍃゙⃞\n> ㅤ\n🫆֟፝𑑓᮫̤٘𑑓̸̷ ʟᴀᴛᴇɴᴄɪᴀ: %dms\n🫆֟፝𑑓᮫̤٘𑑓̸̷ ᴇsᴛᴀᴅᴏ: ᴀᴄᴛɪᴠᴏ ☄️\n╰᳕╌͜╍͜╌͜─͜─͜─͜─͜─͜─⵿֟፝─͜─͜─͜─͜─͜─͜─͜─͜╌͜━㍃゙⃞", latency)

client.SendMessage(context.Background(), evt.Info.Chat, &waE2E.Message{
ExtendedTextMessage: &waE2E.ExtendedTextMessage{
Text: proto.String(text),
ContextInfo: &waE2E.ContextInfo{
QuotedMessage: evt.Message,
},
},
})
}
