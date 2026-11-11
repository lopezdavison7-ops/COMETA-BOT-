package system

import (
"context"
"fmt"
"strings"
"time"

"go.mau.fi/whatsmeow"
"go.mau.fi/whatsmeow/proto/waE2E"
"go.mau.fi/whatsmeow/types/events"
"google.golang.org/protobuf/proto"

"cometa-bot/internal/cache"
"cometa-bot/internal/registry"
)

var startTime = time.Now()

// Datos de tu canal
const canalJID = "120363429914158508@newsletter"
const canalName = "COMETA-BOT "

var catAliases = map[string]string{
"grupos":   "𝐆⵿𝐫⵿⵿𝐩⵿⵿𝐬",
"economia": "𝐄⵿𝵿𝐨⵿𝐨⵿𝐦𝐢⵿𝐚",
"gacha":    "𝐆⵿𝐚⵿⵿𝐡⵿𝐚",
"sistema":  "𝐒⵿𝐢𝐬⵿⵿𝐞𝐦⵿𝐚",
"owner":    "𝐎⵿𝐰⵿𝐧𝐞⵿𝐫",
"pokemon":  "𝐏⵿𝐨⵿⵿𝐞𝐦⵿𝐨𝐧",
"perfil":   "𝐏⵿⵿𝐫𝐟⵿𝐢𝐥",
"juegos":   "𝐉⵿⵿𝐞⵿⵿𝐨⵿𝐬",
"sticker":  "𝐒⵿𝐭𝐢⵿𝐜𝐤⵿𝐞𝐫",
"ia":       "𝐈⵿𝐧⵿𝐭⵿𝐞⵿𝐥𝐢⵿⵿𝐞⵿⵿𝐜⵿⵿𝐚",
}

func init() {
registry.Register(&registry.Command{
Name:        "menu",
Aliases:     []string{"help", "ayuda", "comandos"},
Description: "muestra el menú de comandos (usá /menu <categoría> para filtrar)",
Category:    "𝐒⵿𝐢⵿𝐬⵿𝐭⵿𝐞⵿𝐦⵿𝐚",
Handler:     Menu,
})
}

func Menu(evt *events.Message, client *whatsmeow.Client, args []string) {
uptime := time.Since(startTime).Round(time.Second)
totalCommands := registry.Count()
categories := registry.Categories()

var b strings.Builder

// Header aesthetic
b.WriteString(fmt.Sprintf(`✨۬ᩙ᜔ Hola %s, Ⴆιɳʋҽɳιԃσ α ɱι ɱҽɳυ ԃҽ ƈσɱαɳԃσʂ   𐔌՞ ܸ.ˬ.ܸ ՞

✨᜔۬ ☄️ Ꮯ𝗈𝗆𝖾𝗍𝖺    ̶꯭͞ᵒᵗ 🍒ި᜔⃞ᩙ᜔

┎╌͜╍͜╌͜─͜─͜─͜─͜─፝͜─͜─֟፝͜͠⃞─͜─፝͜─͜─͜─፝͜─͜─͜╌͜━㍃⃞
> ㅤ
> 𖤐⭒๋࣭  `+"`𝗜⵿⵿𝗙𝗢 ⵿𝗕⵿𝗢⵿𝗧`"+` ᩘ⃞゙ᩧ ๋࣭⭒𖤐
> ㅤ
> ፞․̸̷๐  ᴄʀᴇᴀᴅᴏʀ ::  Luis González
> ❐፞․̸̷̤๐  ᴠᴇ𝘴ɪᴏ́ɴ  :: 3.0.0-Go
> ❐፞․̸̷๐  ᴄᴏᴍɴᴅᴏ𝘴 :: %d
> ❐፞․̷̤๐  ᴄᴀᴛᴇɢᴏɪ́ᴀ𝘴 :: %d
> ❐፞․̷̤๐  ᴘᴛɪᴍᴇ :: %s
> ❐፞․̷̤๐  ғᴇᴄʜᴀ :: %s
͜╍͜╌͜─͜─͜─͜─͜─͜─͜─⵿፝─͜─͜──͜─͜─͜─͜─͜╌͜━゙⃞

`, evt.Info.PushName, totalCommands, len(categories), uptime, time.Now().Format("02/01/2006 15:04")))

// Filtrar por categoría
filterCat := ""
if len(args) > 0 {
arg := strings.ToLower(args[0])
if alias, ok := catAliases[arg]; ok {
filterCat = alias
} else {
for _, cat := range categories {
if strings.Contains(strings.ToLower(cat), arg) {
filterCat = cat
break
}
}
}

if filterCat == "" {
b.WriteString("❌ Categoría no encontrada. Prueba con: grupos, economia, gacha, perfil, pokemon, etc.\n")
sendMenuPart(client, evt, b.String())
return
}
}

// Cuerpo del menú
for _, cat := range categories {
if filterCat != "" && cat != filterCat {
continue
}

cmds := registry.GetByCategory(cat)
if len(cmds) == 0 {
continue
}

b.WriteString(fmt.Sprintf("╭ᩙ╌͜╍͜╌͜─͜─͜─͜─፝͜「 `%s` 」፝͜─͜─͜─͜─͜╌͜━㍃゙\n> ㅤ\n", cat))
for _, cmd := range cmds {
b.WriteString(fmt.Sprintf("🫆֟፝𑑓᮫̤٘𑑓̸̷ /%s\n   ↳ %s\n\n", cmd.Name, cmd.Description))
}
b.WriteString("╰᳕╌͜͜╌͜─͜──͜─͜─͜─֟፝─͜──͜─͜─͜──͜─͜╌━㍃゙⃞\n\n")
}

sendMenuPart(client, evt, b.String())
}

// buildContextInfo crea el ContextInfo con metadatos de canal inyectados
func buildContextInfo(evt *events.Message) *waE2E.ContextInfo {
return &waE2E.ContextInfo{
ForwardingScore: proto.Uint32(999),
IsForwarded:     proto.Bool(true),
QuotedMessage:   evt.Message,
// 🔮 Inyectar metadatos de newsletter para el botón "Ver canal"
ForwardedNewsletterMessageInfo: &waE2E.ContextInfo_ForwardedNewsletterMessageInfo{
NewsletterJID:     proto.String(canalJID),
ServerMessageID:   proto.Int32(1),
NewsletterName:    proto.String(canalName),
},
}
}

func sendMenuPart(client *whatsmeow.Client, evt *events.Message, text string) {
ctx := buildContextInfo(evt)

if cache.MenuType == "gif" && cache.MenuGIFUpload != nil {
client.SendMessage(context.Background(), evt.Info.Chat, &waE2E.Message{
VideoMessage: &waE2E.VideoMessage{
URL:           &cache.MenuGIFUpload.URL,
DirectPath:    &cache.MenuGIFUpload.DirectPath,
MediaKey:      cache.MenuGIFUpload.MediaKey,
Mimetype:      proto.String("video/mp4"),
FileEncSHA256: cache.MenuGIFUpload.FileEncSHA256,
FileSHA256:    cache.MenuGIFUpload.FileSHA256,
FileLength:    proto.Uint64(cache.MenuGIFUpload.FileLength),
Caption:       proto.String(text),
GifPlayback:   proto.Bool(true),
ContextInfo:   ctx,
},
})
} else if cache.MenuType == "image" && cache.MenuUpload != nil {
client.SendMessage(context.Background(), evt.Info.Chat, &waE2E.Message{
ImageMessage: &waE2E.ImageMessage{
URL:           &cache.MenuUpload.URL,
DirectPath:    &cache.MenuUpload.DirectPath,
MediaKey:      cache.MenuUpload.MediaKey,
Mimetype:      proto.String("image/jpeg"),
FileEncSHA256: cache.MenuUpload.FileEncSHA256,
FileSHA256:    cache.MenuUpload.FileSHA256,
FileLength:    proto.Uint64(cache.MenuUpload.FileLength),
Caption:       proto.String(text),
ContextInfo:   ctx,
},
})
} else {
client.SendMessage(context.Background(), evt.Info.Chat, &waE2E.Message{
ExtendedTextMessage: &waE2E.ExtendedTextMessage{
Text:        proto.String(text),
ContextInfo: ctx,
},
})
}
}
