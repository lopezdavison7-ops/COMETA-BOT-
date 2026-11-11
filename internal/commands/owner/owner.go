package owner

import (
"context"
"fmt"
"strings"

"go.mau.fi/whatsmeow"
"go.mau.fi/whatsmeow/proto/waE2E"
"go.mau.fi/whatsmeow/types"
"go.mau.fi/whatsmeow/types/events"
"google.golang.org/protobuf/proto"

"cometa-bot/internal/cache"
"cometa-bot/internal/registry"
"cometa-bot/internal/store"
)

func sendOwnerMsg(client *whatsmeow.Client, evt *events.Message, text string) {
client.SendMessage(context.Background(), evt.Info.Chat, &waE2E.Message{
ExtendedTextMessage: &waE2E.ExtendedTextMessage{
Text: proto.String(text),
ContextInfo: &waE2E.ContextInfo{QuotedMessage: evt.Message},
},
})
}

func requireOwner(client *whatsmeow.Client, evt *events.Message) bool {
if !store.IsOwner(evt.Info.Sender.String()) {
sendOwnerMsg(client, evt, "❌ Solo el owner puede usar este comando.")
return false
}
return true
}

func init() {
registry.Register(&registry.Command{
Name:        "setbanner",
Aliases:     []string{"setmenu", "banner"},
Description: "cambia la imagen/gif/video del menú",
Category:    "𝐎⵿𝐰⵿𝐧⵿𝐞⵿𝐫",
Handler:     SetBanner,
})
registry.Register(&registry.Command{
Name:        "setpp",
Aliases:     []string{"setppbot"},
Description: "cambia la foto de perfil del bot",
Category:    "𝐎⵿𝐰⵿𝐧⵿𝐞⵿𝐫",
Handler:     SetPP,
})
registry.Register(&registry.Command{
Name:        "setname",
Aliases:     []string{"setbotname"},
Description: "cambia el nombre del bot",
Category:    "𝐎⵿𝐰⵿𝐧⵿𝐞⵿𝐫",
Handler:     SetName,
})
}

func SetBanner(evt *events.Message, client *whatsmeow.Client, _ []string) {
if !requireOwner(client, evt) {
return
}

var mediaData []byte
var mime string
var mediaType string

if img := evt.Message.GetImageMessage(); img != nil {
mediaType = "image"
mime = img.GetMimetype()
var err error
mediaData, err = client.Download(context.Background(), img)
if err != nil {
sendOwnerMsg(client, evt, "❌ Error al descargar imagen.")
return
}
} else if vid := evt.Message.GetVideoMessage(); vid != nil {
mediaType = "video"
mime = vid.GetMimetype()
var err error
mediaData, err = client.Download(context.Background(), vid)
if err != nil {
sendOwnerMsg(client, evt, "❌ Error al descargar video.")
return
}
} else if sticker := evt.Message.GetStickerMessage(); sticker != nil {
mediaType = "image"
mime = "image/webp"
var err error
mediaData, err = client.Download(context.Background(), sticker)
if err != nil {
sendOwnerMsg(client, evt, "❌ Error al descargar sticker.")
return
}
} else if evt.Message.GetExtendedTextMessage() != nil && evt.Message.GetExtendedTextMessage().GetContextInfo() != nil {
quoted := evt.Message.GetExtendedTextMessage().GetContextInfo().GetQuotedMessage()
if quoted == nil {
sendOwnerMsg(client, evt, "❌ Enviá o respondé a una imagen, gif, video o sticker.")
return
}
if img := quoted.GetImageMessage(); img != nil {
mediaType = "image"
mime = img.GetMimetype()
var err error
mediaData, err = client.Download(context.Background(), img)
if err != nil {
sendOwnerMsg(client, evt, "❌ Error al descargar.")
return
}
} else if vid := quoted.GetVideoMessage(); vid != nil {
mediaType = "video"
mime = vid.GetMimetype()
var err error
mediaData, err = client.Download(context.Background(), vid)
if err != nil {
sendOwnerMsg(client, evt, "❌ Error al descargar.")
return
}
} else if sticker := quoted.GetStickerMessage(); sticker != nil {
mediaType = "image"
mime = "image/webp"
var err error
mediaData, err = client.Download(context.Background(), sticker)
if err != nil {
sendOwnerMsg(client, evt, "❌ Error al descargar.")
return
}
}
}

if len(mediaData) == 0 {
sendOwnerMsg(client, evt, "❌ Enviá o respondé a una imagen, gif, video o sticker.")
return
}

store.SetBanner(mediaData, mime, mediaType)

cache.MenuImage = mediaData
if mediaType == "image" {
upload, err := client.Upload(context.Background(), mediaData, whatsmeow.MediaImage)
if err == nil {
cache.MenuUpload = &upload
}
}

typeName := mediaType
if strings.Contains(mime, "gif") || strings.Contains(mime, "mp4") {
typeName = "gif/video"
}

sendOwnerMsg(client, evt, fmt.Sprintf("✅ Banner actualizado.\n📎 Tipo: %s\n📏 Tamaño: %d bytes", typeName, len(mediaData)))
}

func SetPP(evt *events.Message, client *whatsmeow.Client, _ []string) {
if !requireOwner(client, evt) {
return
}

var mediaData []byte
if img := evt.Message.GetImageMessage(); img != nil {
var err error
mediaData, err = client.Download(context.Background(), img)
if err != nil {
sendOwnerMsg(client, evt, "❌ Error al descargar.")
return
}
} else if evt.Message.GetExtendedTextMessage() != nil && evt.Message.GetExtendedTextMessage().GetContextInfo() != nil {
quoted := evt.Message.GetExtendedTextMessage().GetContextInfo().GetQuotedMessage()
if quoted != nil && quoted.GetImageMessage() != nil {
var err error
mediaData, err = client.Download(context.Background(), quoted.GetImageMessage())
if err != nil {
sendOwnerMsg(client, evt, "❌ Error al descargar.")
return
}
}
}

if len(mediaData) == 0 {
sendOwnerMsg(client, evt, "❌ Enviá o respondé a una imagen.")
return
}

_, err := client.SetGroupPhoto(context.Background(), *client.Store.ID, mediaData)
if err != nil {
sendOwnerMsg(client, evt, "❌ Error al cambiar foto: "+err.Error())
return
}
sendOwnerMsg(client, evt, "✅ Foto de perfil actualizada.")
}

func SetName(evt *events.Message, client *whatsmeow.Client, args []string) {
if !requireOwner(client, evt) {
return
}
if len(args) == 0 {
sendOwnerMsg(client, evt, "❌ Usá: /setname <nombre>")
return
}
name := strings.Join(args, " ")
if len(name) > 25 {
name = name[:25]
}
store.SetBotName(name)
client.SendMessage(context.Background(), evt.Info.Chat, &waE2E.Message{
ExtendedTextMessage: &waE2E.ExtendedTextMessage{
Text: proto.String("✅ Nombre del bot actualizado a: *" + name + "*"),
ContextInfo: &waE2E.ContextInfo{QuotedMessage: evt.Message},
},
})
}

var _ = types.JID{}
