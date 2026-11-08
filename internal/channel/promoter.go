package channel

import (
	"context"
	"time"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
)

const TargetChannel = "0029Vb8eeKGG3R3kwBcZdp2Q@newsletter"

func StartAutoPromo(client *whatsmeow.Client, interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			sendPromo(client)
		}
	}()
}

func sendPromo(client *whatsmeow.Client) {
	jid, err := types.ParseJID(TargetChannel)
	if err != nil {
		return
	}

	msg := &waE2E.Message{
		ExtendedTextMessage: &waE2E.ExtendedTextMessage{
			Text: proto.String("✨۬ᩙ᜔ ☄️ Ꮯ᥆ոᥱ𝗍ᥲ    ̶͞ᵇ꯭ᵒᵗ 🍒ິ᜔⃞🌸ᩙ᜔\n\n🫆֟፝𑑓᮫̤٘𑑓̸̷ 𝗨𝗻𝗲𝘁𝗲 𝗮𝗹 𝗰𝗮𝗻𝗮𝗹 𝗼𝗳𝗶𝗰𝗶𝗮𝗹:\nhttps://whatsapp.com/channel/0029Vb8eeKGG3R3kwBcZdp2Q"),
			ContextInfo: &waE2E.ContextInfo{
				ActionLink: &waE2E.ActionLink{
					Url:        proto.String("https://whatsapp.com/channel/0029Vb8eeKGG3R3kwBcZdp2Q"),
					ButtonTitle: proto.String("Ver Canal"),
				},
			},
		},
	}

	client.SendMessage(context.Background(), jid, msg)
}