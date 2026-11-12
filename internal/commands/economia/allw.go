
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
		Name:        "allw",
		Aliases:     []string{"workall", "trabajartodo", "todo"},
		Description: "ejecuta todas tus actividades económicas disponibles",
		Category:    "𝐄⵿𝐜𝵿𝐨⵿𝐧𝐨⵿𝐦⵿𝐢⵿𝐚",
		Handler:     Allw,
	})
}

func Allw(evt *events.Message, client *whatsmeow.Client, _ []string) {
	if store.GlobalEconomy == nil {
		sendEcoError(client, evt, "La economía no está inicializada.")
		return
	}

	jid := evt.Info.Sender.String()
	result, err := store.GlobalEconomy.AllWork(jid)

	prefix := "🛡️ "

	if err != nil {
		if err.Error() == "cooldown" {
			prefix = "⏳ "
		} else {
			prefix = "❌ "
			result = "Ocurrió un error al procesar ALL WORK."
		}
	}

	var b strings.Builder

	b.WriteString("╭ᩙ╌͜╍͜╌͜─͜─͜─͜─፝͜「 `ALL WORK` 」፝͜─͜─͜─͜─͜╌͜━㍃゙⃞\n> ㅤ\n")
	b.WriteString(prefix + result + "\n")
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
