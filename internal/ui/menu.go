package ui

import (
"fmt"
"os"
"strings"
"time"

"go.mau.fi/whatsmeow/types"
)

var StartTime = time.Now()

type Command struct {
Name string
Desc string
}

type Category struct {
Name     string
Commands []Command
}

func RenderMenu(creator, version string, cats []Category, start time.Time, userName string, userJID types.JID) (string, []types.JID) {
totalCmds := 0
for _, c := range cats {
totalCmds += len(c.Commands)
}

var saludo string
mentions := []types.JID{}

if userName != "" {
saludo = fmt.Sprintf("✨۬ᩙ᜔ Hola %s, Ⴆιҽɳʋҽɳιԃσ α ɱι ɱҽɳυ ԃҽ ƈσɱαɳԃσʂ   𐔌՞ ܸ.ˬ.ܸ ՞𐦯", userName)
} else {
saludo = fmt.Sprintf("✨۬ᩙ᜔ @%s Ⴆιҽɳʋҽɳιԃσ α ɱι ɱҽɳυ ԃҽ ƈσɱαɳԃσʂ   𐔌՞ ܸ.ˬ.ܸ ՞𐦯", userJID.User)
mentions = append(mentions, userJID)
}

b := strings.Builder{}
b.Grow(700 + (totalCmds * 60))
b.WriteString(saludo + "\n\n")
b.WriteString("✨۬ᩙ᜔ ☄️ Ꮯ𝗈𝗆𝖾𝗍𝖺    ̶͞ᵇ꯭ᵒᵗ 🍒ި᜔⃞🌸ᩙ᜔\n\n")
b.WriteString("┎╌͜╍͜╌͜─፝͜─͜─͜─͜─፝͜─͜─֟፝͜͠⃞─͜─፝͜─͜─͜─፝͜─͜─͜╌͜━㍃゙⃞\n> ㅤ\n")
b.WriteString("> 𖤐⭒๋࣭ ⭑ `𝗜⵿𝗡⵿𝗙⵿𝗢 ⵿𝗕⵿𝗢⵿𝗧` ᩘ📊⃞゙ᩧ ⭑๋࣭⭒𖤐\n> ㅤ\n")
b.WriteString("> ❐፞․̸̷̤๐  ᴄʀᴇᴀᴅᴏʀ ::  " + creator + "\n")
b.WriteString("> ❐፞․̸̷̤๐  ᴠᴇʀ𝘴ɪᴏ́ɴ  :: " + version + "\n")
b.WriteString(fmt.Sprintf("> ❐፞․̸̷̤๐  ᴄᴏᴍᴀɴᴅᴏ𝘴 :: %d\n", totalCmds))
b.WriteString(fmt.Sprintf("> ❐፞․̸̷̤๐  ᴄᴀᴛᴇɢᴏʀɪ́ᴀ𝘴 :: %d\n", len(cats)))
b.WriteString("> ❐፞․̸̷̤๐  ᴜᴘᴛɪᴍᴇ :: " + time.Since(start).Round(time.Second).String() + "\n")
b.WriteString("> ❐፞․̸̷̤๐  ғᴇᴄʜᴀ :: " + time.Now().Format("02/01/2006 15:04") + "\n")
b.WriteString("┖͜╍͜╌͜─͜─͜─͜─͜─͜─͜─⵿֟፝─͜─͜─͜─͜─͜─͜─͜─͜╌͜━㍃゙⃞\n\n")
for _, cat := range cats {
b.WriteString("╭ᩙ╌͜╍͜╌͜─͜─͜─͜─፝͜「 `" + cat.Name + "` 」፝͜─͜─͜─͜─͜╌͜━㍃゙⃞\n> ㅤ\n")
for _, cmd := range cat.Commands {
b.WriteString("🫆֟፝𑑓᮫̤٘𑑓̸̷ " + cmd.Name + "\n")
b.WriteString("   ↳ " + cmd.Desc + "\n\n")
}
b.WriteString("╰᳕╌͜╍͜╌͜─͜─͜─͜─͜─͜─⵿֟፝─͜─͜─͜─͜─͜─͜─͜─͜╌͜━㍃゙⃞\n\n")
}
return b.String(), mentions
}

func LoadMenuImage() ([]byte, error) {
return os.ReadFile("assets/menu_anime.jpg")
}
