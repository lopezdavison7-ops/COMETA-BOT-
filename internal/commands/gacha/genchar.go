package gacha

import (
"context"
"encoding/json"
"fmt"
"io"
"net/http"

"go.mau.fi/whatsmeow"
"go.mau.fi/whatsmeow/proto/waE2E"
"go.mau.fi/whatsmeow/types/events"
"google.golang.org/protobuf/proto"

"cometa-bot/internal/registry"
"cometa-bot/internal/store"
)

func init() {
registry.Register(&registry.Command{
Name:        "genchar",
Aliases:     []string{"generate", "addchar"},
Description: "(Owner) Genera personajes desde konachan para el gacha",
Category:    "𝐆⵿𝐚⵿𝐜⵿𝐡⵿𝐚",
Handler:     Genchar,
})
}

type KonachanPost struct {
ID      int    `json:"id"`
Tags    string `json:"tags"`
FileURL string `json:"file_url"`
}

func Genchar(evt *events.Message, client *whatsmeow.Client, args []string) {
if !store.IsOwner(evt.Info.Sender.String()) {
sendGachaError(client, evt, "Solo el owner puede usar este comando.")
return
}

resp, err := http.Get("https://konachan.com/post.json?limit=1&tags=rating:safe")
if err != nil {
sendGachaError(client, evt, "Error al conectar con la API.")
return
}
defer resp.Body.Close()

var posts []KonachanPost
if err := json.NewDecoder(resp.Body).Decode(&posts); err != nil || len(posts) == 0 {
sendGachaError(client, evt, "No se pudo obtener el personaje.")
return
}

post := posts[0]
imgResp, err := http.Get(post.FileURL)
if err != nil {
sendGachaError(client, evt, "Error al descargar la imagen.")
return
}
defer imgResp.Body.Close()

imgData, err := io.ReadAll(imgResp.Body)
if err != nil || len(imgData) == 0 {
sendGachaError(client, evt, "Imagen corrupta o vacía.")
return
}

upload, err := client.Upload(context.Background(), imgData, whatsmeow.MediaImage)
if err != nil {
sendGachaError(client, evt, "Error al subir la imagen a WhatsApp.")
return
}

caption := fmt.Sprintf("✅ Personaje generado y agregado al gacha!\n🆔 ID: %d\n🏷️ Tags: %s", post.ID, post.Tags)

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
ContextInfo:   &waE2E.ContextInfo{QuotedMessage: evt.Message},
},
})
}

func sendGachaError(client *whatsmeow.Client, evt *events.Message, text string) {
client.SendMessage(context.Background(), evt.Info.Chat, &waE2E.Message{
ExtendedTextMessage: &waE2E.ExtendedTextMessage{
Text:        proto.String("❌ " + text),
ContextInfo: &waE2E.ContextInfo{QuotedMessage: evt.Message},
},
})
}
