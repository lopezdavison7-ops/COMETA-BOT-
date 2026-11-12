
package economia

import (
	"context"
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
		Name:        "weekly",
		Aliases:     []string{"semanal", "week"},
		Description: "reclama tu recompensa semanal de cometas",
		Category:    "𝐄⵿𝐜𝵿𝐨⵿𝐧𝐨⵿𝐦⵿𝐢⵿𝐚",
		Handler:     Weekly,
	})
}

func Weekly(evt *events.Message, client *whatsmeow.Client, _ []string) {
	if store.GlobalEconomy == nil {
		return
	}

	jid := evt.Info.Sender.String()
	msg, err := store.GlobalEconomy.Weekly(jid)

	var prefix string

	if err != nil {
		if err.Error() == "cooldown" {
			prefix = "⏳ "
		} else {
			prefix = "❌ "
			msg = "Ocurrió un error al procesar tu recompensa."
		}
	} else {
		prefix = "🎁 "
	}

	b := strings.Builder{}

	b.WriteString("╭ᩙ╌͜╍͜╌͜─͜─͜─͜─፝͜「 `WEEKLY` 」፝͜─͜─͜─͜─͜╌͜━㍃゙⃞\n> ㅤ\n")
	b.WriteString("🫆֟፝𑑓᮫̤٘𑑓̸̷ " + prefix + msg + "\n")
	b.WriteString("╰᳕╌͜╍͜╌͜─͜─͜─͜─͜─͜─⵿֟፝─͜─͜─͜─͜─͜─͜─͜─͜╌͜━㍃゙⃞")

	client.SendMessage(context.Background(), evt.Info.Chat, &waE2E.Message{
		ExtendedTextMessage: &waE2E.ExtendedTextMessage{
			Text: proto.String(b.String()),
			ContextInfo: &waE2E.ContextInfo{
				QuotedMessage: evt.Message,
			},
		},
	})
}
