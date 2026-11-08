package ui

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type Category struct {
	Name     string
	Commands []string
}

func RenderMenu(creator, version string, cats []Category, start time.Time) ([]byte, string, error) {
	totalCmds := 0
	for _, c := range cats {
		totalCmds += len(c.Commands)
	}

	b := strings.Builder{}
	b.Grow(700 + (totalCmds * 25))

	b.WriteString("✨۬ᩙ᜔ @user Ⴆιҽɳʋҽɳιԃσ α ɱι ɱҽɳυ ԃҽ ƈσɱαɳԃσʂ   𐔌՞ ܸ.ˬ.ܸ ՞𐦯\n\n")
	b.WriteString("✨۬ᩙ᜔ ☄️ Ꮯ᥆ոᥱ𝗍ᥲ    ̶͞ᵇ꯭ᵒᵗ 🍒ິ᜔⃞🌸ᩙ᜔\n\n")
	b.WriteString("┎╌͜╍͜╌͜─፝͜─͜─͜─͜─፝͜─͜─֟፝͜͠⃞─͜─፝͜─͜─͜─፝͜─͜─͜╌͜━㍃゙⃞\n> ㅤ\n")
	b.WriteString("> 𖤐⭒๋࣭ ⭑ `𝗜⵿𝗡⵿𝗙⵿𝗢 ⵿𝗕⵿𝗢⵿𝗧` ᩘ📊⃞゙ᩧ ⭑๋࣭⭒𖤐\n> ㅤ\n")
	b.WriteString("> ❐፞․̸̷̤๐  ᴄʀᴇᴀᴅᴏʀ ::  ")
	b.WriteString(creator)
	b.WriteString("\n> ❐፞․̸̷̤๐  ᴠᴇʀ𝘴ɪᴏ́ɴ  :: ")
	b.WriteString(version)
	b.WriteString("\n> ❐፞․̸̷̤๐  ᴄᴏᴍᴀɴᴅᴏ𝘴 :: ")
	b.WriteString(fmt.Sprintf("%d", totalCmds))
	b.WriteString("\n> ❐፞․̸̷̤๐  ᴄᴀᴛᴇɢᴏʀɪ́ᴀ𝘴 :: ")
	b.WriteString(fmt.Sprintf("%d", len(cats)))
	b.WriteString("\n> ❐፞․̸̷̤๐  ᴜᴘᴛɪᴍᴇ :: ")
	b.WriteString(time.Since(start).Round(time.Second).String())
	b.WriteString("\n> ❐፞․̸̷̤๐  ғᴇᴄʜᴀ :: ")
	b.WriteString(time.Now().Format("02/01/2006 15:04"))
	b.WriteString("\n┖͜╍͜╌͜─͜─͜─͜─͜─͜─͜─⵿֟፝─͜─͜─͜─͜─͜─͜─͜─͜╌͜━㍃゙⃞\n\n")

	for _, cat := range cats {
		b.WriteString("╭ᩙ╌͜╍͜╌͜─͜─͜─͜─፝͜「 `")
		b.WriteString(cat.Name)
		b.WriteString("` 」፝͜─͜─͜─͜─͜╌͜━㍃゙⃞\n> ㅤ\n🫆֟፝𑑓᮫̤٘𑑓̸̷ ")
		b.WriteString(strings.Join(cat.Commands, " "))
		b.WriteString("\n╰᳕╌͜╍͜╌͜─͜─͜─͜─͜─͜─⵿֟፝─͜─͜─͜─͜─͜─͜─͜─͜╌͜━㍃゙⃞\n\n")
	}

	img, err := os.ReadFile("assets/menu_anime.jpg")
	if err != nil {
		return nil, b.String(), err
	}

	return img, b.String(), nil
}