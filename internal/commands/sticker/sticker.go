package sticker

import (
"context"
"fmt"
"strings"

"go.mau.fi/whatsmeow"
"go.mau.fi/whatsmeow/proto/waE2E"
"go.mau.fi/whatsmeow/types"
"go.mau.fi/whatsmeow/types/events"
"google.golang.org/protobuf/proto"

"cometa-bot/internal/registry"
"cometa-bot/internal/store"
"cometa-bot/internal/utils"
)

func init() {
registry.Register(&registry.Command{
Name:        "sticker",
Aliases:     []string{"s", "stik", "stick"},
Description: "crea stickers desde imagen, video o sticker",
Category:    "𝐒⵿𝐭⵿𝐢⵿𝐜⵿𝐤⵿𝐞⵿𝐫",
Handler:     Sticker,
})
}

func Sticker(evt *events.Message, client *whatsmeow.Client, args []string) {
ctx := context.Background()
senderJID := evt.Info.Sender.String()

packname, author := store.GetStickerMeta(senderJID)

if len(args) > 0 {
tempText := strings.Join(args, " ")
if strings.Contains(tempText, "|") {
parts := strings.SplitN(tempText, "|", 2)
if parts[0] != "" {
packname = strings.TrimSpace(parts[0])
}
if len(parts) > 1 && parts[1] != "" {
author = strings.TrimSpace(parts[1])
}
} else {
packname = tempText
author = tempText
}
}

var targetMsg *events.Message
var targetType string

if evt.Message.GetImageMessage() != nil || evt.Message.GetVideoMessage() != nil || evt.Message.GetStickerMessage() != nil {
targetMsg = evt
if evt.Message.GetImageMessage() != nil {
targetType = "image"
} else if evt.Message.GetVideoMessage() != nil {
targetType = "video"
} else {
targetType = "sticker"
}
} else if ctxInfo := evt.Message.GetExtendedTextMessage().GetContextInfo(); ctxInfo != nil && ctxInfo.GetQuotedMessage() != nil {
quoted := ctxInfo.GetQuotedMessage()
if quoted.GetImageMessage() != nil || quoted.GetVideoMessage() != nil || quoted.GetStickerMessage() != nil {
targetMsg = &events.Message{
Info: types.MessageInfo{
MessageSource: types.MessageSource{
Chat: evt.Info.Chat,
},
},
Message: quoted,
}
if quoted.GetImageMessage() != nil {
targetType = "image"
} else if quoted.GetVideoMessage() != nil {
targetType = "video"
} else {
targetType = "sticker"
}
}
}

if targetMsg == nil {
helpText := `╭━━〔 🖼️ 𝐒𝐓𝐈𝐂𝐊𝐄𝐑 〕━━⬣
┃
┃  Envía o responde una imagen,
┃  video (máx 15s) o sticker.
┃
┃  Ejemplos:
┃ › /s ➔ sticker normal
┃ › /s MiMarca ➔ marca temporal
┃ › /s Marca | Autor ➔ personalizada
┃
╰━━━━━━━━━━━━━━━━⬣`

client.SendMessage(ctx, evt.Info.Chat, &waE2E.Message{
ExtendedTextMessage: &waE2E.ExtendedTextMessage{
Text: proto.String(helpText),
ContextInfo: &waE2E.ContextInfo{
QuotedMessage: evt.Message,
},
},
})
return
}

var mediaBuffer []byte
var err error

switch targetType {
case "image":
img := targetMsg.Message.GetImageMessage()
mediaBuffer, err = client.Download(ctx, img)
case "video":
vid := targetMsg.Message.GetVideoMessage()
mediaBuffer, err = client.Download(ctx, vid)
case "sticker":
stk := targetMsg.Message.GetStickerMessage()
mediaBuffer, err = client.Download(ctx, stk)
}

if err != nil {
client.SendMessage(ctx, evt.Info.Chat, &waE2E.Message{
ExtendedTextMessage: &waE2E.ExtendedTextMessage{
Text: proto.String("❌ Error al descargar el media."),
ContextInfo: &waE2E.ContextInfo{
QuotedMessage: evt.Message,
},
},
})
return
}

var webpBuffer []byte
switch targetType {
case "image":
webpBuffer, err = utils.ConvertImageToWebp(mediaBuffer)
case "video":
webpBuffer, err = utils.ConvertVideoToWebp(mediaBuffer)
case "sticker":
webpBuffer, err = utils.CleanSticker(mediaBuffer)
}

if err != nil {
client.SendMessage(ctx, evt.Info.Chat, &waE2E.Message{
ExtendedTextMessage: &waE2E.ExtendedTextMessage{
Text: proto.String(fmt.Sprintf("❌ %v", err)),
ContextInfo: &waE2E.ContextInfo{
QuotedMessage: evt.Message,
},
},
})
return
}

stickerFinal, err := utils.AddExif(webpBuffer, packname, author)
if err != nil {
stickerFinal = webpBuffer
}

upload, err := client.Upload(ctx, stickerFinal, whatsmeow.MediaImage)
if err != nil {
client.SendMessage(ctx, evt.Info.Chat, &waE2E.Message{
ExtendedTextMessage: &waE2E.ExtendedTextMessage{
Text: proto.String("❌ Error al subir el sticker."),
ContextInfo: &waE2E.ContextInfo{
QuotedMessage: evt.Message,
},
},
})
return
}

stickerMsg := &waE2E.Message{
StickerMessage: &waE2E.StickerMessage{
URL:           &upload.URL,
DirectPath:    &upload.DirectPath,
MediaKey:      upload.MediaKey,
FileEncSHA256: upload.FileEncSHA256,
FileSHA256:    upload.FileSHA256,
FileLength:    proto.Uint64(uint64(len(stickerFinal))),
Mimetype:      proto.String("image/webp"),
IsAnimated:    proto.Bool(targetType == "video"),
ContextInfo: &waE2E.ContextInfo{
QuotedMessage: evt.Message,
},
},
}

client.SendMessage(ctx, evt.Info.Chat, stickerMsg)
}
