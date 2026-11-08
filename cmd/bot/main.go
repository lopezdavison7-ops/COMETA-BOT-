package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	_ "modernc.org/sqlite"

	"cometa-bot/internal/handler"
)

func main() {
	dbLog := waLog.Stdout("DB", "ERROR", true)
	container, err := sqlstore.New("sqlite", "file:cometa.db?_pragma=foreign_keys(1)", dbLog)
	if err != nil {
		panic(err)
	}

	device, err := container.GetFirstDevice()
	if err != nil {
		panic(err)
	}

	client := whatsmeow.NewClient(device, waLog.Stdout("WA", "INFO", true))

	client.AddEventHandler(func(evt interface{}) {
		if msg, ok := evt.(*events.Message); ok {
			handler.Route(msg, client)
		}
	})

	if client.Store.ID == nil {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("✨۬ᩙ᜔ Número (sin +, ej 5057777777): ")
		phone, _ := reader.ReadString('\n')
		phone = strings.TrimSpace(phone)
		phone = strings.TrimPrefix(phone, "+")

		code, err := client.PairPhone(phone, true, whatsmeow.PairClientChrome, "Cometa Core")
		if err != nil {
			panic(err)
		}
		fmt.Printf("🫆֟፝𑑓᮫̤٘𑑓̸̷ Código: %s\n", code)
		fmt.Println("> WhatsApp > Dispositivos > Vincular con número")
	}

	if err := client.Connect(); err != nil {
		panic(err)
	}

	fmt.Println("☄️ Ꮯ᥆ոᥱ𝗍ᥲ conectado")
	select {}
}