package pokemon

import (
"context"
"encoding/json"
"fmt"
"io"
"math/rand"
"net/http"
"strings"
"time"

"go.mau.fi/whatsmeow"
"go.mau.fi/whatsmeow/proto/waE2E"
"go.mau.fi/whatsmeow/types/events"
"google.golang.org/protobuf/proto"

"cometa-bot/internal/registry"
"cometa-bot/internal/store"
)

var httpClient = &http.Client{Timeout: 10 * time.Second}

func init() {
registry.Register(&registry.Command{
Name:        "poke",
Aliases:     []string{"pokemon", "catch", "atrapar"},
Description: "atrapa un pokemon aleatorio (3 pokéballs al día)",
Category:    "𝐏⵿𝐨⵿𝐤⵿𝐞⵿𝐦⵿𝐨⵿𝐧",
Handler:     Poke,
})
}

type PokeAPIResponse struct {
Name    string `json:"name"`
ID      int    `json:"id"`
Sprites struct {
FrontDefault string `json:"front_default"`
FrontShiny   string `json:"front_shiny"`
Other        struct {
Home struct {
FrontDefault string `json:"front_default"`
FrontShiny   string `json:"front_shiny"`
} `json:"home"`
} `json:"other"`
} `json:"sprites"`
Types []struct {
Type struct {
Name string `json:"name"`
} `json:"type"`
} `json:"types"`
}

func titleCase(s string) string {
if s == "" {
return ""
}
return strings.ToUpper(s[:1]) + s[1:]
}

func sendText(client *whatsmeow.Client, evt *events.Message, text string) {
client.SendMessage(context.Background(), evt.Info.Chat, &waE2E.Message{
ExtendedTextMessage: &waE2E.ExtendedTextMessage{
Text: proto.String(text),
ContextInfo: &waE2E.ContextInfo{QuotedMessage: evt.Message},
},
})
}

func Poke(evt *events.Message, client *whatsmeow.Client, _ []string) {
jid := evt.Info.Sender.String()

_, canUse := store.CheckBalls(jid)
if !canUse {
sendText(client, evt, fmt.Sprintf("╭〔 🎒 𝐒𝐈𝐍 𝐏𝐎𝐊𝐄𝐁𝐀𝐋𝐋𝐒 〕⬣\n\n┃ Ya usaste tus *3 Pokéballs* de hoy.\n┃ Vuelve mañana para recargar.\n┃\n╰━━━━━━━━━━━━━━━━⬣"))
return
}

pokeID := rand.Intn(1025) + 1
url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%d", pokeID)

resp, err := httpClient.Get(url)
if err != nil || resp.StatusCode != 200 {
sendText(client, evt, "❌ Error al conectar con PokeAPI. Intenta de nuevo.")
return
}
defer resp.Body.Close()

var data PokeAPIResponse
if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
sendText(client, evt, "❌ Error al procesar los datos del Pokémon.")
return
}

isShiny := rand.Intn(100) < 10
imgURL := data.Sprites.Other.Home.FrontDefault

if isShiny && data.Sprites.Other.Home.FrontShiny != "" {
imgURL = data.Sprites.Other.Home.FrontShiny
}
if imgURL == "" {
imgURL = data.Sprites.FrontDefault
if isShiny && data.Sprites.FrontShiny != "" {
imgURL = data.Sprites.FrontShiny
}
}
if imgURL == "" {
imgURL = fmt.Sprintf("https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/other/home/%d.png", pokeID)
}

imgResp, err := httpClient.Get(imgURL)
if err != nil || imgResp.StatusCode != 200 {
sendText(client, evt, "❌ Error al descargar la imagen del Pokémon.")
return
}
defer imgResp.Body.Close()

imgData, err := io.ReadAll(imgResp.Body)
if err != nil || len(imgData) == 0 {
sendText(client, evt, "❌ La imagen está corrupta o vacía.")
return
}

var types []string
for _, t := range data.Types {
types = append(types, titleCase(t.Type.Name))
}
typeStr := strings.Join(types, " / ")
if typeStr == "" {
typeStr = "Desconocido"
}

remaining := store.UseBall(jid)
store.AddPokemon(jid, store.Pokemon{
ID:      data.ID,
Name:    titleCase(data.Name),
Image:   imgURL,
IsShiny: isShiny,
Date:    time.Now().Format("02/01"),
})

shinyText := ""
if isShiny {
shinyText = "\n┃ ✨ *¡ES SHINY!* ✨"
}

caption := fmt.Sprintf("╭〔 🎯 𝐏𝐎𝐊𝐄𝐁𝐀𝐋𝐋 〕⬣\n┃\n┃ ¡Atrapaste a *%s*!\n┃ 🔢 *Nº:* %03d\n┃ ⚡ *Tipo:* %s%s\n┃\n┃ 🎒 *Pokéballs restantes:* %d/3\n┃\n╰━━━━━━━━━━━━━━━━⬣\n\n> Usa */pokeinv* para ver tu colección ☄️",
titleCase(data.Name), data.ID, typeStr, shinyText, remaining)

upload, err := client.Upload(context.Background(), imgData, whatsmeow.MediaImage)
if err != nil {
sendText(client, evt, caption+"\n\n⚠️ (No se pudo adjuntar la imagen)")
return
}

client.SendMessage(context.Background(), evt.Info.Chat, &waE2E.Message{
ImageMessage: &waE2E.ImageMessage{
URL:           &upload.URL,
DirectPath:    &upload.DirectPath,
MediaKey:      upload.MediaKey,
Mimetype:      proto.String("image/png"),
FileEncSHA256: upload.FileEncSHA256,
FileSHA256:    upload.FileSHA256,
FileLength:    proto.Uint64(uint64(len(imgData))),
Caption:       proto.String(caption),
ContextInfo: &waE2E.ContextInfo{
QuotedMessage: evt.Message,
},
},
})
}
