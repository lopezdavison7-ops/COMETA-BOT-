package main

import (
"bufio"
"context"
"fmt"
"os"
"strings"
"time"

"go.mau.fi/whatsmeow"
"go.mau.fi/whatsmeow/store/sqlstore"
"go.mau.fi/whatsmeow/types"
"go.mau.fi/whatsmeow/types/events"
waLog "go.mau.fi/whatsmeow/util/log"
_ "modernc.org/sqlite"

"cometa-bot/internal/cache"
"cometa-bot/internal/handler"
"cometa-bot/internal/store"
)

func main() {
ctx := context.Background()
dbLog := waLog.Stdout("DB", "ERROR", true)
container, err := sqlstore.New(ctx, "sqlite", "file:cometa.db?_pragma=busy_timeout(10000)&_pragma=journal_mode(WAL)&_pragma=synchronous(NORMAL)&_pragma=foreign_keys(1)", dbLog)
if err != nil {
panic(err)
}
device, err := container.GetFirstDevice(ctx)
if err != nil {
panic(err)
}

device.Platform = "Chrome"
device.BusinessName = ""
device.PushName = "Cometa"

client := whatsmeow.NewClient(device, waLog.Stdout("WA", "INFO", true))
client.EnableAutoReconnect = true

if err := store.InitEconomy("cometa_economy.db"); err != nil {
fmt.Printf("❌ Error iniciando economía: %v\n", err)
} else {
fmt.Println("💾 Base de datos de economía (WAL) lista")
}

qrChan := make(chan string, 1)
pairDone := make(chan struct{})

go func() {
for {
time.Sleep(30 * time.Second)
if client.Store.ID != nil && client.IsConnected() {
_ = client.SendPresence(ctx, types.PresenceAvailable)
time.Sleep(5 * time.Second)
_ = client.SendPresence(ctx, types.PresenceUnavailable)
}
}
}()

go func() {
for {
time.Sleep(60 * time.Second)
if client.Store.ID != nil && !client.IsConnected() {
fmt.Println("⚠️ Watchdog: Desconectado. Reconectando...")
handler.NotifyOwners(client, "Desconexión detectada. Reconectando automáticamente...")
client.Connect()
}
}
}()

client.AddEventHandler(func(evt interface{}) {
switch v := evt.(type) {
case *events.Message:
go handler.Route(v, client)
case *events.QR:
select {
case qrChan <- "ready":
default:
}
case *events.PairSuccess:
fmt.Println("☄️ Cometa vinculado")
_ = client.SendPresence(ctx, types.PresenceAvailable)
close(pairDone)
case *events.Connected:
fmt.Println("☄️ Cometa conectado")
_ = client.SendPresence(ctx, types.PresenceAvailable)
go cacheMenuMedia(client)
case *events.Disconnected:
fmt.Println("⚠️ Desconectado")
case *events.HistorySync:
fmt.Printf("📥 Historial: %d\n", len(v.Data.Conversations))
case *events.AppStateSyncComplete:
fmt.Printf("✅ Sync: %s\n", v.Name)
if v.Name == "regular_low" {
fmt.Println("🚀 Cometa operativo - prueba /menu")
}
case *events.LoggedOut:
fmt.Printf("❌ Logout: %v\n", v.Reason)
handler.NotifyOwners(client, fmt.Sprintf("Sesión cerrada: %v", v.Reason))
os.Exit(1)
case *events.KeepAliveTimeout:
fmt.Println("⏱️ Keepalive timeout")
}
})

if err := client.Connect(); err != nil {
panic(err)
}

if client.Store.ID == nil {
fmt.Println("⏳ Esperando conexión...")
<-qrChan

reader := bufio.NewReader(os.Stdin)
fmt.Print("✨۬ᩙ᜔ Número (sin +): ")
phone, _ := reader.ReadString('\n')
phone = strings.TrimSpace(phone)
phone = strings.TrimPrefix(phone, "+")

code, err := client.PairPhone(ctx, phone, true, whatsmeow.PairClientChrome, "Chrome (Linux)")
if err != nil {
panic(fmt.Sprintf("Error: %v", err))
}
fmt.Printf("🫆֟፝𑑓᮫̤٘𑑓̸̷ Código: %s\n", code)

select {
case <-pairDone:
case <-time.After(2 * time.Minute):
fmt.Println("⚠️ Timeout")
}
}

select {}
}

func cacheMenuMedia(client *whatsmeow.Client) {
if gif, err := os.ReadFile("assets/menu_anime.gif"); err == nil {
cache.MenuGIF = gif
upload, err := client.Upload(context.Background(), gif, whatsmeow.MediaVideo)
if err == nil {
cache.MenuGIFUpload = &upload
cache.MenuType = "gif"
fmt.Println("🎬 GIF del menú cacheado")
return
}
}

if img, err := os.ReadFile("assets/menu_anime.jpg"); err == nil {
cache.MenuImage = img
upload, err := client.Upload(context.Background(), img, whatsmeow.MediaImage)
if err == nil {
cache.MenuUpload = &upload
cache.MenuType = "image"
fmt.Println("📸 Imagen del menú cacheada")
return
}
}

fmt.Println("⚠️ No hay imagen ni GIF en assets/ - menú será solo texto")
}
