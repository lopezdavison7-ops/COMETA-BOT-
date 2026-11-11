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

var catAliases = map[string]string{
"grupos":   "𝐆⵿𝐫⵿𝐮⵿𝐩⵿𝐨⵿𝐬",
"economia": "𝐄⵿𝐜𝐨⵿𝐧𝐨⵿𝐦⵿𝐢⵿𝐚",
"gacha":    "𝐆⵿𝐚⵿𝐜⵿𝐡⵿𝐚",
"sistema":  "𝐒⵿𝐢⵿𝐬⵿𝐭⵿𝐞⵿𝐦⵿𝐚",
"owner":    "𝐎⵿𝐰⵿𝐧⵿𝐞⵿𝐫",
"pokemon":  "𝐏⵿𝐨⵿𝐤⵿𝐞⵿𝐦⵿𝐨⵿𝐧",
"perfil":   "𝐏⵿𝐞⵿𝐫⵿𝐟⵿𝐢⵿𝐥",
"juegos":   "𝐉⵿𝐮⵿𝐞⵿𝐠⵿𝐨⵿𝐬",
"sticker":  "𝐒⵿𝐭⵿𝐢⵿𝐜⵿𝐤⵿𝐞⵿𝐫",
"ia":       "𝐈⵿𝐧⵿𝐭⵿𝐞⵿𝐥⵿𝐢⵿𝐠⵿𝐞⵿𝐧⵿𝐜⵿𝐢⵿𝐚",
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
b.WriteString(fmt.Sprintf(`✨۬ᩙ᜔ Hola %s, Ⴆιҽɳʋҽɳιԃσ α ɱι ɱҽɳυ ԃҽ ƈσɱαɳԃσʂ   𐔌՞ ܸ.ˬ.ܸ ՞𐦯

✨۬ᩙ᜔ ☄️ Ꮯ𝗈𝗆𝖾𝗍𝖺    ̶͞ᵇ꯭ᵒᵗ 🍒ި᜔⃞🌸ᩙ᜔

┎╌͜╍͜╌͜─፝͜─͜─͜─͜─፝͜─͜─֟፝͜͠⃞─͜─፝͜─͜─͜─፝͜─͜─͜╌͜━㍃゙⃞
> ㅤ
> 𖤐⭒๋࣭ ⭑ `+"`𝗜⵿𝗡⵿𝗙⵿𝗢 ⵿𝗕⵿𝗢⵿𝗧`"+` ᩘ📊⃞゙ᩧ ⭑๋࣭⭒𖤐
> ㅤ
> ❐፞․̸̷̤๐  ᴄʀᴇᴀᴅᴏʀ :: 🎀devluis🍒
> ❐፞․̸̷̤๐  ᴠᴇʀ𝘴ɪᴏ́ɴ  :: 2.0-Go
> ❐፞․̸̷̤๐  ᴄᴏᴍᴀɴᴅᴏ𝘴 :: %d
> ❐፞․̸̷̤๐  ᴄᴀᴛᴇɢᴏʀɪ́ᴀ𝘴 :: %d
> ❐፞․̸̷̤๐  ᴜᴘᴛɪᴍᴇ :: %s
> ❐፞․̸̷̤๐  ғᴇᴄʜᴀ :: %s
┖͜╍͜╌͜─͜─͜─͜─͜─͜─͜─⵿֟፝─͜─͜─͜─͜─͜─͜─͜─͜╌͜━㍃゙⃞

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

b.WriteString(fmt.Sprintf("╭ᩙ╌͜╍͜╌͜─͜─͜─͜─፝͜「 `%s` 」፝͜─͜─͜─͜─͜╌͜━㍃゙⃞\n> ㅤ\n", cat))
for _, cmd := range cmds {
b.WriteString(fmt.Sprintf("🫆֟፝𑑓᮫̤٘𑑓̸̷ /%s\n   ↳ %s\n\n", cmd.Name, cmd.Description))
}
b.WriteString("╰᳕╌͜╍͜╌͜─͜─͜─͜─͜─͜─⵿֟፝─͜─͜─͜─͜─͜─͜─͜─͜╌͜━㍃゙⃞\n\n")
}



sendMenuPart(client, evt, b.String())
}

func sendMenuPart(client *whatsmeow.Client, evt *events.Message, text string) {
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
ContextInfo: &waE2E.ContextInfo{QuotedMessage: evt.Message},
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
ContextInfo: &waE2E.ContextInfo{QuotedMessage: evt.Message},
},
})
} else {
client.SendMessage(context.Background(), evt.Info.Chat, &waE2E.Message{
ExtendedTextMessage: &waE2E.ExtendedTextMessage{
Text: proto.String(text),
ContextInfo: &waE2E.ContextInfo{QuotedMessage: evt.Message},
},
})
}
}
