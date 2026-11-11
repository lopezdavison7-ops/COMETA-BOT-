package ai

import (
"context"
"strings"

"go.mau.fi/whatsmeow"
"go.mau.fi/whatsmeow/proto/waE2E"
"go.mau.fi/whatsmeow/types/events"
"google.golang.org/protobuf/proto"

ai_core "cometa-bot/internal/ai"
"cometa-bot/internal/registry"
)

func init() {
registry.Register(&registry.Command{
Name:        "ia",
Aliases:     []string{"chatgpt", "bot", "pregunta"},
Description: "pregunta algo a la inteligencia artificial",
Category:    "𝐈⵿𝐧⵿𝐭⵿𝐞⵿𝐥⵿𝐢⵿𝐠⵿𝐞⵿𝐧⵿𝐜⵿𝐢⵿𝐚",
Handler:     IA,
})
}

func IA(evt *events.Message, client *whatsmeow.Client, args []string) {
prompt := strings.Join(args, " ")
if prompt == "" {
client.SendMessage(context.Background(), evt.Info.Chat, &waE2E.Message{
ExtendedTextMessage: &waE2E.ExtendedTextMessage{
Text: proto.String("❌ Escribe algo para preguntar. Ejemplo: `/ia como hacer vigorón`"),
},
})
return
}

thinkingMsg := &waE2E.Message{
ExtendedTextMessage: &waE2E.ExtendedTextMessage{
Text: proto.String("⏳ Pensando... ☄️"),
ContextInfo: &waE2E.ContextInfo{
QuotedMessage: evt.Message,
},
},
}

msgResp, err := client.SendMessage(context.Background(), evt.Info.Chat, thinkingMsg)
if err != nil {
return
}

res := ai_core.Ask(prompt)

finalText := "╭ᩙ╌͜╍͜╌͜─͜─͜─͜─፝͜「 `🤖 COMETA IA` 」፝͜─͜─͜─͜─͜╌͜━㍃゙⃞\n> ㅤ\n🫆֟፝𑑓᮫̤٘𑑓̸̷ " + res + "\n╰᳕╌͜╍͜╌͜─͜─͜─͜─͜─͜─⵿֟፝─͜─͜─͜─͜─͜─͜─͜─͜╌͜━㍃゙⃞"

editMsg := client.BuildEdit(evt.Info.Chat, msgResp.ID, &waE2E.Message{
ExtendedTextMessage: &waE2E.ExtendedTextMessage{
Text: proto.String(finalText),
ContextInfo: &waE2E.ContextInfo{
QuotedMessage: evt.Message,
},
},
})

client.SendMessage(context.Background(), evt.Info.Chat, editMsg)
}
