package gacha

import (
"context"
"fmt"
"os"

"go.mau.fi/whatsmeow"
"go.mau.fi/whatsmeow/proto/waE2E"
"go.mau.fi/whatsmeow/types/events"
"google.golang.org/protobuf/proto"

"cometa-bot/internal/registry"
"cometa-bot/internal/store"
)

func init() {
registry.Register(&registry.Command{
Name:        "rw",
Aliases:     []string{"roll", "waifu", "husbando"},
Description: "saca un personaje aleatorio del gacha (cooldown 20 min)",
Category:    "𝐆⵿𝐚⵿𝐜⵿𝐡⵿𝐚",
Handler:     RW,
})
}

var rollCache = make(map[string]string)

func GetLastRoll(jid string) (string, bool) {
file, ok := rollCache[jid]
return file, ok
}

func ClearRoll(jid string) {
delete(rollCache, jid)
}

func RW(evt *events.Message, client *whatsmeow.Client, _ []string) {
jid := evt.Info.Sender.String()

if remaining, ok := store.CheckRollCooldown(jid); !ok {
mins := remaining / 60
secs := remaining % 60
client.SendMessage(context.Background(), evt.Info.Chat, &waE2E.Message{
ExtendedTextMessage: &waE2E.ExtendedTextMessage{
Text: proto.String(fmt.Sprintf("╭〔 ⏳ 𝐂𝐎𝐎𝐋𝐃𝐎𝐖𝐍 〕⬣\n\n┃ Ya sacaste una waifu recién, chacho.\n┃ Esperá *%d min %d seg* para volver a tirar.\n┃\n╰━━━━━━━━━━━━━━━━⬣", mins, secs)),
ContextInfo: &waE2E.ContextInfo{QuotedMessage: evt.Message},
},
})
return
}

file, char, ok := store.RollWaifu()
if !ok {
client.SendMessage(context.Background(), evt.Info.Chat, &waE2E.Message{
ExtendedTextMessage: &waE2E.ExtendedTextMessage{
Text: proto.String("❌ No hay personajes disponibles. Usá /genchar para generar más."),
ContextInfo: &waE2E.ContextInfo{QuotedMessage: evt.Message},
},
})
return
}

store.SetRollCooldown(jid)
rollCache[jid] = file

caption := fmt.Sprintf(
"╭〔 🎴 𝐖𝐀𝐈𝐅𝐔 › %s 〕⬣\n┃\n┃ 🎀 *Nombre:* %s\n┃ 📺 *Serie:* %s\n┃ ⚥ *Género:* %s\n┃ 💰 *Valor:* %d cometas\n┃\n╰━━━━━━━━━━━━━━━━⬣\n\n> Usá */claim* para reclamarla ☄️",
char.Series, char.Name, char.Series, char.Gender, char.Value,
)

imgData, err := os.ReadFile(char.ImagePath)
if err != nil {
client.SendMessage(context.Background(), evt.Info.Chat, &waE2E.Message{
ExtendedTextMessage: &waE2E.ExtendedTextMessage{
Text: proto.String(caption),
ContextInfo: &waE2E.ContextInfo{QuotedMessage: evt.Message},
},
})
return
}

upload, err := client.Upload(context.Background(), imgData, whatsmeow.MediaImage)
if err != nil {
client.SendMessage(context.Background(), evt.Info.Chat, &waE2E.Message{
ExtendedTextMessage: &waE2E.ExtendedTextMessage{
Text: proto.String(caption),
ContextInfo: &waE2E.ContextInfo{QuotedMessage: evt.Message},
},
})
return
}

client.SendMessage(context.Background(), evt.Info.Chat, &waE2E.Message{
ImageMessage: &waE2E.ImageMessage{
URL:           &upload.URL,
DirectPath:    &upload.DirectPath,
MediaKey:      upload.MediaKey,
Mimetype:      proto.String("image/jpeg"),
FileEncSHA256: upload.FileEncSHA256,
FileSHA256:    upload.FileSHA256,
FileLength:    proto.Uint64(uint64(len(imgData))),
Caption:       proto.String(caption),
ContextInfo: &waE2E.ContextInfo{
QuotedMessage: evt.Message,
},
},
})
}
