package profile

import (
"context"
"fmt"
"io"
"net/http"
"strconv"
"strings"
"time"

"go.mau.fi/whatsmeow"
"go.mau.fi/whatsmeow/proto/waE2E"
"go.mau.fi/whatsmeow/types"
"go.mau.fi/whatsmeow/types/events"
"google.golang.org/protobuf/proto"

"cometa-bot/internal/registry"
"cometa-bot/internal/store"
)

func init() {
registry.Register(&registry.Command{
Name:        "perfil",
Aliases:     []string{"profile", "me", "yo"},
Description: "muestra tu perfil con foto y economía",
Category:    "𝐏⵿𝐞⵿𝐫⵿𝐟⵿𝐢⵿𝐥",
Handler:     Perfil,
})
registry.Register(&registry.Command{
Name:        "setdesc",
Aliases:     []string{"setdescripcion"},
Description: "cambia tu descripción",
Category:    "𝐏⵿𝐞⵿𝐫⵿𝐟⵿𝐢⵿𝐥",
Handler:     SetDesc,
})
registry.Register(&registry.Command{
Name:        "setedad",
Aliases:     []string{"setage"},
Description: "cambia tu edad",
Category:    "𝐏⵿𝐞⵿𝐫⵿𝐟⵿𝐢⵿𝐥",
Handler:     SetEdad,
})
registry.Register(&registry.Command{
Name:        "setgenero",
Aliases:     []string{"setgender"},
Description: "cambia tu género",
Category:    "𝐏⵿𝐞⵿𝐫⵿𝐟⵿𝐢⵿𝐥",
Handler:     SetGenero,
})
registry.Register(&registry.Command{
Name:        "setcancion",
Aliases:     []string{"setsong"},
Description: "cambia tu canción favorita",
Category:    "𝐏⵿𝐞⵿𝐫⵿𝐟⵿𝐢⵿𝐥",
Handler:     SetCancion,
})
registry.Register(&registry.Command{
Name:        "setcrush",
Aliases:     []string{"setlike"},
Description: "cambia la persona que te gusta",
Category:    "𝐏⵿𝐞⵿𝐫⵿𝐟⵿𝐢⵿𝐥",
Handler:     SetCrush,
})
}

func sendText(client *whatsmeow.Client, evt *events.Message, text string) {
client.SendMessage(context.Background(), evt.Info.Chat, &waE2E.Message{
ExtendedTextMessage: &waE2E.ExtendedTextMessage{
Text: proto.String(text),
ContextInfo: &waE2E.ContextInfo{QuotedMessage: evt.Message},
},
})
}

func getTargetJID(evt *events.Message, args []string) types.JID {
ext := evt.Message.GetExtendedTextMessage()
if ext != nil {
ctx := ext.GetContextInfo()
if ctx != nil {
mentioned := ctx.GetMentionedJID()
if len(mentioned) > 0 {
jid, _ := types.ParseJID(mentioned[0])
return jid
}
if p := ctx.GetParticipant(); p != "" {
jid, _ := types.ParseJID(p)
return jid
}
}
}
if len(args) > 0 {
jid, _ := types.ParseJID(args[0])
return jid
}
return evt.Info.Sender
}

func Perfil(evt *events.Message, client *whatsmeow.Client, args []string) {
target := getTargetJID(evt, args)
jid := target.String()

profile := store.GetProfile(jid)

var balance, bank, level, xp int
if store.GlobalEconomy != nil {
balance = store.GlobalEconomy.GetBalance(jid)
bank = store.GlobalEconomy.GetBank(jid)
level, xp = store.GlobalEconomy.GetProfile(jid)
}

elapsed := time.Now().Unix() - profile.FirstSeen
days := int(elapsed / 86400)
hours := int((elapsed % 86400) / 3600)
timeStr := fmt.Sprintf("%d días, %d horas", days, hours)

var b strings.Builder
b.WriteString("╭━━━〔 👤 PERFIL 〕━━━⊷\n")
b.WriteString("┃\n")

if profile.Age > 0 {
b.WriteString(fmt.Sprintf("┃ 🎂 Edad: %d años\n", profile.Age))
}
if profile.Gender != "" {
b.WriteString(fmt.Sprintf("┃ ⚧ Género: %s\n", profile.Gender))
}
if profile.Desc != "" {
b.WriteString(fmt.Sprintf("┃ 📝 Desc: %s\n", profile.Desc))
}
if profile.FavSong != "" {
b.WriteString(fmt.Sprintf("┃ 🎵 Canción: %s\n", profile.FavSong))
}
if profile.Crush != "" {
b.WriteString(fmt.Sprintf("┃ 💖 Crush: %s\n", profile.Crush))
}

b.WriteString(fmt.Sprintf("┃ 💰 Billetera: %d\n", balance))
b.WriteString(fmt.Sprintf("┃ 🏦 Banco: %d\n", bank))
b.WriteString(fmt.Sprintf("┃ 📊 Nivel: %d (%d/%d XP)\n", level, xp, level*100))
b.WriteString(fmt.Sprintf("┃ ⏱️ Tiempo: %s\n", timeStr))
b.WriteString("┃\n")
b.WriteString("╰━━━━━━━━━━━━━━━━━━━━━━⊷")

caption := b.String()

// 🔍 INVESTIGACIÓN WHATSMEOW: Obtener foto de perfil real
// El método es GetProfilePictureInfo, que nos da la URL directa del CDN de WhatsApp
picInfo, err := client.GetProfilePictureInfo(context.Background(), target, nil)

if err == nil && picInfo != nil && picInfo.URL != "" {
// 1. Descargar la imagen desde la URL
resp, err := http.Get(picInfo.URL)
if err == nil {
defer resp.Body.Close()
imgData, err := io.ReadAll(resp.Body)
if err == nil && len(imgData) > 0 {
// 2. Subir la imagen a los servidores de WhatsApp
upload, err := client.Upload(context.Background(), imgData, whatsmeow.MediaImage)
if err == nil {
// 3. Enviar como imagen con el perfil de caption
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
MentionedJID:  []string{jid},
QuotedMessage: evt.Message,
},
},
})
return
}
}
}
}

// Fallback: Si no tiene foto o falla la descarga, manda solo texto
client.SendMessage(context.Background(), evt.Info.Chat, &waE2E.Message{
ExtendedTextMessage: &waE2E.ExtendedTextMessage{
Text: proto.String(caption),
ContextInfo: &waE2E.ContextInfo{
MentionedJID:  []string{jid},
QuotedMessage: evt.Message,
},
},
})
}

func SetDesc(evt *events.Message, client *whatsmeow.Client, args []string) {
if len(args) == 0 {
sendText(client, evt, "Usá: /setdesc <tu descripción>")
return
}
jid := evt.Info.Sender.String()
desc := strings.Join(args, " ")
if len(desc) > 100 {
desc = desc[:100]
}
store.SetProfileField(jid, "desc", desc)
sendText(client, evt, "✅ Descripción actualizada: "+desc)
}

func SetEdad(evt *events.Message, client *whatsmeow.Client, args []string) {
if len(args) == 0 {
sendText(client, evt, "Usá: /setedad <número>")
return
}
age, err := strconv.Atoi(args[0])
if err != nil || age < 1 || age > 120 {
sendText(client, evt, "Edad inválida. Debe ser entre 1 y 120.")
return
}
jid := evt.Info.Sender.String()
store.SetProfileField(jid, "age", fmt.Sprintf("%d", age))
sendText(client, evt, fmt.Sprintf("✅ Edad actualizada: %d años", age))
}

func SetGenero(evt *events.Message, client *whatsmeow.Client, args []string) {
if len(args) == 0 {
sendText(client, evt, "Usá: /setgenero <M/F/No binario/Otro>")
return
}
jid := evt.Info.Sender.String()
gender := strings.Join(args, " ")
store.SetProfileField(jid, "gender", gender)
sendText(client, evt, "✅ Género actualizado: "+gender)
}

func SetCancion(evt *events.Message, client *whatsmeow.Client, args []string) {
if len(args) == 0 {
sendText(client, evt, "Usá: /setcancion <nombre de la canción>")
return
}
jid := evt.Info.Sender.String()
song := strings.Join(args, " ")
store.SetProfileField(jid, "fav_song", song)
sendText(client, evt, "✅ Canción favorita actualizada: "+song)
}

func SetCrush(evt *events.Message, client *whatsmeow.Client, args []string) {
if len(args) == 0 {
sendText(client, evt, "Usá: /setcrush <nombre de la persona>")
return
}
jid := evt.Info.Sender.String()
crush := strings.Join(args, " ")
store.SetProfileField(jid, "crush", crush)
sendText(client, evt, "✅ Crush actualizado: "+crush)
}
