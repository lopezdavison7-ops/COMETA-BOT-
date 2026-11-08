package handler

import (
	"context"
	"strings"
	"time"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"

	"cometa-bot/internal/ui"
)

type CommandFunc func(evt *events.Message, client *whatsmeow.Client, args []string)

var commands = map[string]CommandFunc{}

func Register(name string, fn CommandFunc) {
	commands[name] = fn
}

func init() {
	Register("ping", cmdPing)
	Register("menu", cmdMenu)
	Register("help", cmdMenu)
}

func Route(evt *events.Message, client *whatsmeow.Client) {
	text := strings.TrimSpace(evt.Message.GetConversation())
	if !strings.HasPrefix(text, ".") {
		return
	}

	parts := strings.Fields(text)
	cmd := strings.ToLower(strings.TrimPrefix(parts[0], "."))
	args := parts[1:]

	if fn, ok := commands[cmd]; ok {
		go fn(evt, client, args)
	}
}

func cmdPing(evt *events.Message, client *whatsmeow.Client, _ []string) {
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

	client.SendMessageResponse(context.Background(), evt.Info.Chat, msgID.ID, &waE2E.Message{
		ExtendedTextMessage: &waE2E.ExtendedTextMessage{
			Text: proto.String(b.String()),
		},
	})
}

func cmdMenu(evt *events.Message, client *whatsmeow.Client, _ []string) {
	cats := []ui.Category{
		{"𝐒⵿𝐢⵿𝐬⵿𝐭⵿𝐞⵿𝐦⵿𝐚", []string{".ping", ".menu", ".help"}},
	}

	img, caption, err := ui.RenderMenu("Luis González", "3.0.0-Go", cats, time.Time{})
	if err != nil {
		client.SendMessage(context.Background(), evt.Info.Chat, &waE2E.Message{
			ExtendedTextMessage: &waE2E.ExtendedTextMessage{Text: &caption},
		})
		return
	}

	uploaded, err := client.Upload(context.Background(), img, whatsmeow.MediaImage)
	if err != nil {
		return
	}

	client.SendMessage(context.Background(), evt.Info.Chat, &waE2E.Message{
		ImageMessage: &waE2E.ImageMessage{
			Caption:       &caption,
			Mimetype:      proto.String("image/jpeg"),
			JpegThumbnail: uploaded.Content[:32],
			Url:           &uploaded.URL,
			DirectPath:    &uploaded.DirectPath,
			MediaKey:      uploaded.MediaKey,
			FileEncSha256: uploaded.FileEncSHA256,
			FileSha256:    uploaded.FileSHA256,
			FileLength:    proto.Uint64(uint64(len(img))),
		},
	})
}