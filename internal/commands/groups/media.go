package groups

import (
"context"

"go.mau.fi/whatsmeow"
"go.mau.fi/whatsmeow/types/events"

"cometa-bot/internal/registry"
"cometa-bot/internal/store"
)

func init() {
registry.Register(&registry.Command{
Name:        "setwelcomeimg",
Aliases:     []string{"setbienvenidaimagen"},
Description: "sube una imagen para la bienvenida",
Category:    "𝐆⵿𝐫⵿𝐮⵿𝐩⵿𝐨⵿𝐬",
Handler:     SetWelcomeImg,
})
registry.Register(&registry.Command{
Name:        "setgoodbyeimg",
Aliases:     []string{"setdespedidaimagen"},
Description: "sube una imagen para la despedida",
Category:    "𝐆⵿𝐫⵿𝐮⵿𝐩⵿𝐨⵿𝐬",
Handler:     SetGoodbyeImg,
})
registry.Register(&registry.Command{
Name:        "delwelcomeimg",
Aliases:     []string{},
Description: "borra la imagen de bienvenida",
Category:    "𝐆⵿𝐫⵿𝐮⵿𝐩⵿𝐨⵿𝐬",
Handler:     DelWelcomeImg,
})
registry.Register(&registry.Command{
Name:        "delgoodbyeimg",
Aliases:     []string{},
Description: "borra la imagen de despedida",
Category:    "𝐆⵿𝐫⵿𝐮⵿𝐩⵿𝐨⵿𝐬",
Handler:     DelGoodbyeImg,
})
}

func downloadImage(client *whatsmeow.Client, evt *events.Message) ([]byte, string) {
if img := evt.Message.GetImageMessage(); img != nil {
data, err := client.Download(context.Background(), img)
if err != nil {
return nil, ""
}
return data, img.GetMimetype()
}
if ext := evt.Message.GetExtendedTextMessage(); ext != nil {
if quoted := ext.GetContextInfo().GetQuotedMessage(); quoted != nil {
if img := quoted.GetImageMessage(); img != nil {
data, err := client.Download(context.Background(), img)
if err != nil {
return nil, ""
}
return data, img.GetMimetype()
}
}
}
return nil, ""
}

func SetWelcomeImg(evt *events.Message, client *whatsmeow.Client, _ []string) {
if !requireAdmin(client, evt) { return }
data, mime := downloadImage(client, evt)
if len(data) == 0 {
sendGroupError(client, evt, "Respondé a una imagen para usarla como bienvenida.")
return
}

// 🚀 PRE-UPLOAD: Subimos la imagen a WhatsApp AHORA para que sea instantáneo después
upload, err := client.Upload(context.Background(), data, whatsmeow.MediaImage)
if err != nil {
sendGroupError(client, evt, "Error al subir la imagen a WhatsApp: "+err.Error())
return
}

jid := evt.Info.Chat.String()
saved := store.SavedMedia{
URL: upload.URL, DirectPath: upload.DirectPath, MediaKey: upload.MediaKey,
FileEncSHA256: upload.FileEncSHA256, FileSHA256: upload.FileSHA256,
FileLength: uint64(len(data)), Mimetype: mime,
}

if err := store.SaveMediaUpload(jid, "welcome", saved); err != nil {
sendGroupError(client, evt, "Error al guardar la configuración.")
return
}
sendGroupSuccess(client, evt, "Imagen de bienvenida configurada y pre-subida. ¡Será instantánea!")
}

func SetGoodbyeImg(evt *events.Message, client *whatsmeow.Client, _ []string) {
if !requireAdmin(client, evt) { return }
data, mime := downloadImage(client, evt)
if len(data) == 0 {
sendGroupError(client, evt, "Respondé a una imagen para usarla como despedida.")
return
}

upload, err := client.Upload(context.Background(), data, whatsmeow.MediaImage)
if err != nil {
sendGroupError(client, evt, "Error al subir la imagen a WhatsApp: "+err.Error())
return
}

jid := evt.Info.Chat.String()
saved := store.SavedMedia{
URL: upload.URL, DirectPath: upload.DirectPath, MediaKey: upload.MediaKey,
FileEncSHA256: upload.FileEncSHA256, FileSHA256: upload.FileSHA256,
FileLength: uint64(len(data)), Mimetype: mime,
}

if err := store.SaveMediaUpload(jid, "goodbye", saved); err != nil {
sendGroupError(client, evt, "Error al guardar la configuración.")
return
}
sendGroupSuccess(client, evt, "Imagen de despedida configurada y pre-subida. ¡Será instantánea!")
}

func DelWelcomeImg(evt *events.Message, client *whatsmeow.Client, _ []string) {
if !requireAdmin(client, evt) { return }
store.DeleteMediaUpload(evt.Info.Chat.String(), "welcome")
sendGroupSuccess(client, evt, "Imagen de bienvenida eliminada.")
}

func DelGoodbyeImg(evt *events.Message, client *whatsmeow.Client, _ []string) {
if !requireAdmin(client, evt) { return }
store.DeleteMediaUpload(evt.Info.Chat.String(), "goodbye")
sendGroupSuccess(client, evt, "Imagen de despedida eliminada.")
}
