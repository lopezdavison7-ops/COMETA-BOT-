package system

import (
"context"
"fmt"

"go.mau.fi/whatsmeow"
"go.mau.fi/whatsmeow/proto/waE2E"
"go.mau.fi/whatsmeow/types/events"
"google.golang.org/protobuf/proto"

"cometa-bot/internal/cache"
"cometa-bot/internal/registry"
"cometa-bot/internal/ui"
)

const canalURL = "https://whatsapp.com/channel/0029Vb8eeKGG3R3kwBcZdp2Q"

func init() {
registry.Register(&registry.Command{
Name:        "menu",
Aliases:     []string{"m", "help", "ayuda", "commands"},
Description: "muestra este menú de comandos",
Category:    "𝐒⵿𝐢⵿𝐬⵿𝐭⵿𝐞⵿𝐦⵿𝐚",
Handler:     Menu,
})
}

func Menu(evt *events.Message, client *whatsmeow.Client, _ []string) {
uniqueCmds := registry.GetAllUnique()
catsMap := make(map[string][]ui.Command)
for _, cmd := range uniqueCmds {
catsMap[cmd.Category] = append(catsMap[cmd.Category], ui.Command{
Name: "/" + cmd.Name,
Desc: cmd.Description,
})
}

cats := make([]ui.Category, 0, len(catsMap))
for catName, cmds := range catsMap {
cats = append(cats, ui.Category{Name: catName, Commands: cmds})
}

caption, mentions := ui.RenderMenu("Luis González", "3.0.0-Go", cats, ui.StartTime, evt.Info.PushName, evt.Info.Sender)

mentionedJID := make([]string, 0, len(mentions))
for _, jid := range mentions {
mentionedJID = append(mentionedJID, jid.String())
}

if cache.MenuUpload == nil || cache.MenuImage == nil {
client.SendMessage(context.Background(), evt.Info.Chat, &waE2E.Message{
ExtendedTextMessage: &waE2E.ExtendedTextMessage{
Text:        proto.String(caption + "\n\n📢 " + canalURL),
ContextInfo: &waE2E.ContextInfo{MentionedJID: mentionedJID},
},
})
return
}

thumbnail := cache.MenuImage
if len(cache.MenuImage) > 32 {
thumbnail = cache.MenuImage[:32]
}

msg := &waE2E.Message{
ImageMessage: &waE2E.ImageMessage{
Caption:       proto.String(caption),
Mimetype:      proto.String("image/jpeg"),
JPEGThumbnail: thumbnail,
URL:           &cache.MenuUpload.URL,
DirectPath:    &cache.MenuUpload.DirectPath,
MediaKey:      cache.MenuUpload.MediaKey,
FileEncSHA256: cache.MenuUpload.FileEncSHA256,
FileSHA256:    cache.MenuUpload.FileSHA256,
FileLength:    proto.Uint64(uint64(len(cache.MenuImage))),
ContextInfo: &waE2E.ContextInfo{
MentionedJID: mentionedJID,
ExternalAdReply: &waE2E.ContextInfo_ExternalAdReplyInfo{
Title:     proto.String("☄️ Cometa Bot - Canal Oficial"),
Body:      proto.String("Únete para novedades y soporte"),
MediaType: waE2E.ContextInfo_ExternalAdReplyInfo_IMAGE.Enum(),
SourceURL: proto.String(canalURL),
Thumbnail: thumbnail,
},
},
},
}

if _, err := client.SendMessage(context.Background(), evt.Info.Chat, msg); err != nil {
fmt.Printf("❌ Error: %v\n", err)
}
}
