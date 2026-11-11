package owner

import (
"context"
"fmt"
"os"
"os/exec"
"strings"
"syscall"
"time"

"go.mau.fi/whatsmeow"
"go.mau.fi/whatsmeow/proto/waE2E"
"go.mau.fi/whatsmeow/types/events"
"google.golang.org/protobuf/proto"

"cometa-bot/internal/registry"
)

func init() {
registry.Register(&registry.Command{
Name:        "update",
Aliases:     []string{"actualizar", "fix"},
Description: "actualiza el bot desde GitHub y reinicia",
Category:    "𝐎⵿𝐰⵿𝐧⵿𝐞⵿𝐫",
Handler:     Update,
})
}

func sendUpdateMsg(client *whatsmeow.Client, evt *events.Message, text string) *waE2E.Message {
msg := &waE2E.Message{
ExtendedTextMessage: &waE2E.ExtendedTextMessage{
Text: proto.String(text),
ContextInfo: &waE2E.ContextInfo{
QuotedMessage: evt.Message,
},
},
}
client.SendMessage(context.Background(), evt.Info.Chat, msg)
return msg
}

func Update(evt *events.Message, client *whatsmeow.Client, _ []string) {
if !requireOwner(client, evt) {
return
}

// Mensaje inicial
sendUpdateMsg(client, evt, 
"╭━━〔 🔄 𝐀𝐂𝐓𝐔𝐀𝐋𝐈𝐙𝐀𝐑 〕━━⬣\n"+
"┃\n"+
"┃ 📥 Descargando actualizaciones...\n"+
"┃ ⏳ Espera un momento.\n"+
"┃\n"+
"╰━━━━━━━━━━━━━━━━",
)

// Obtener archivos modificados antes del pull
_, _ = exec.Command("git", "diff", "--name-only", "HEAD").Output()

// Ejecutar git pull
cmd := exec.Command("git", "pull")
cmd.Dir = "."
cmd.Env = os.Environ()

salida, err := cmd.CombinedOutput()
salidaStr := string(salida)

if err != nil {
sendUpdateMsg(client, evt,
"╭━━〔 ❌ 𝐀𝐂𝐓𝐔𝐀𝐋𝐈𝐙𝐀𝐑 〕━━⬣\n"+
"┃\n"+
"┃ ⚠️ Error al descargar actualizaciones.\n"+
"┃\n"+
fmt.Sprintf("┃ 🔍 %s\n", salidaStr)+
"┃\n"+
"╰━━━━━━━━━━━━━━━━⬣",
)
return
}

// Verificar si ya está actualizado
if strings.Contains(salidaStr, "Already up to date") || strings.Contains(salidaStr, "Ya está actualizado") {
sendUpdateMsg(client, evt,
"╭━━〔 ✅ 𝐀𝐂𝐓𝐔𝐀𝐋𝐈𝐙𝐀𝐑 〕━━⬣\n"+
"┃\n"+
"┃ 🎉 El bot ya está actualizado.\n"+
"┃ 📦 No hay cambios nuevos.\n"+
"┃\n"+
"╰━━━━━━━━━━━━━━━━",
)
return
}

// Obtener archivos modificados después del pull
archivosDespues, _ := exec.Command("git", "diff", "--name-only", "HEAD@{1}").Output()
archivosLista := strings.Split(strings.TrimSpace(string(archivosDespues)), "\n")

var cambiosStr string
if len(archivosLista) > 0 && archivosLista[0] != "" {
cambiosStr = "📝 *Archivos modificados:*\n"
for _, a := range archivosLista {
if a != "" {
cambiosStr += fmt.Sprintf("┃ • %s\n", a)
}
}
} else {
cambiosStr = "┃ 📦 Se actualizaron archivos.\n"
}

sendUpdateMsg(client, evt,
"╭━━〔 ✅ 𝐀𝐂𝐓𝐔𝐀𝐋𝐈𝐙𝐀𝐑 〕━━⬣\n"+
"┃\n"+
"┃ 📦 Actualización completada.\n"+
"┃\n"+
cambiosStr+
"┃\n"+
"┃ 🔄 Reiniciando bot en 3 segundos...\n"+
"┃\n"+
"╰━━━━━━━━━━━━━━━━⬣",
)

// Reiniciar el bot después de 3 segundos
time.Sleep(3 * time.Second)

// Reemplazar el proceso actual con una nueva instancia
executable, err := os.Executable()
if err != nil {
fmt.Printf("❌ Error al obtener ejecutable: %v\n", err)
return
}

args := os.Args[1:]
err = syscall.Exec(executable, append([]string{executable}, args...), os.Environ())
if err != nil {
fmt.Printf("❌ Error al reiniciar: %v\n", err)
}
}
