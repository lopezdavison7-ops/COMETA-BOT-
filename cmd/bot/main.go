package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	_ "modernc.org/sqlite"

	"cometa-bot/internal/channel"
	"cometa-bot/internal/ui"
)

var botStartTime = time.Now()

func main() {
	dbLog := waLog.Stdout("Database", "ERROR", true)
	container, err := sqlstore.New("sqlite", "file:cometa.db?_pragma=foreign_keys(1)", dbLog)
	if err != nil {
		panic(err)
	}

	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		panic(err)
	}

	clientLog := waLog.Stdout("Client", "INFO", true)
	client := whatsmeow.NewClient(deviceStore, clientLog)

	client.AddEventHandler(func(evt interface{}) {
		switch v := evt.(type) {
		case *events.Message:
			go handleMessage(v, client)
		case *events.Connected:
			channel.StartAutoPromo(client, 4*time.Hour)
		}
	})

	if client.Store.ID == nil {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("✨۬ᩙ᜔ Ingresa tu número (sin +, ej: 5057777777): ")
		phone, _ := reader.ReadString('\n')
		phone = strings.TrimSpace(phone)
		phone = strings.TrimPrefix(phone, "+")

		code, err := client.PairPhone(phone, true, whatsmeow.PairClientChrome, "Cometa Core")
		if err != nil {
			panic(err)
		}
		fmt.Printf("🫆֟፝𑑓᮫̤٘𑑓̸̷ Código de vinculación: %s\n", code)
		fmt.Println("> Ingresa este código en WhatsApp > Dispositivos vinculados > Vincular con número")
	}

	err = client.Connect()
	if err != nil {
		panic(err)
	}

	select {}
}

func handleMessage(evt *events.Message, client *whatsmeow.Client) {
	text := evt.Message.GetConversation()
	if text != ".menu" && text != ".help" {
		return
	}

	cats := []ui.Category{
		{"𝐀⵿𝐧𝐢⵿𝐦⵿𝐞", []string{".anime", ".animesearch", ".buscaranime", ".ani"}},
		{"𝐁⵿𝐮𝐬⵿𝐪𝐮⵿𝐞𝐝⵿𝐚", []string{".peli", ".movie", ".film", ".pelicula"}},
		{"𝐂⵿𝐚⵿𝐧𝐯⵿𝐚", []string{".xnxx", ".xnxxcard", ".xcard", ".bratanime", ".brat", ".animebrat"}},
		{"𝐃⵿𝐞𝐬⵿𝐜𝐚⵿𝐫𝐠𝐚⵿𝐬", []string{".play",  ".musica", ".reproductor", ".song", ".play2", ".playvideo", ".video", ".applemusic", ".amusic", ".am", ".applemusicdl", ".instagram", ".ig", ".igdl", ".instadl", ".insta", ".mediafire", ".mf", ".mfdl", ".mediafiredl", ".fb", ".facebook", ".fbdl", ".facebookdl", ".letra", ".lyrics", ".letras", ".karaoke", ".reproductor", ".pinterest", ".pin", ".pinterestimg", ".soundcloud", ".sc", ".sound", ".spotify", ".sp", ".spoti", ".spotifydl", ".tt", ".tiktok", ".tk", ".ttv", ".tiktokdl", ".tiktoksearch", ".xvideos", ".xv", ".xvsearch", ".xvdl"}},
		{"𝐅⵿𝐮⵿𝐧", []string{".adoptar", ".adopt", ".hijo", ".hija", ".adoptados", ".misadoptados", ".listaadoptados", ".rw", ".recompensa", ".mascota", ".pet", ".tamagotchi", ".mascotas", ".ship", ".pareja", ".amor", ".compatibilidad", ".shippear", ".8ball", ".bola", ".magicball", ".adivinanza", ".acertijo", ".adivinanzacancelar", ".cancelaradivinanza", ".akinador", ".aki", ".genio", ".adivina", ".chiste", ".chistes", ".dado", ".dados", ".dino", ".dinosaurio", ".trex", ".chrome", ".doxeo", ".dox", ".fakedox", ".info", ".encuesta", ".poll", ".fake", ".fakem", ".fakemessage", ".falso", ".prediccion", ".predecir", ".futuro", ".oraculo", ".lid", ".lo.inocente.de", ".men", ".memoria", ".memory", ".juego", ".pelea", ".batalla", ".combate", ".fight", ".duelo", ".ppt", ".piedrapapeltijera", ".prehot", ".preguntahot", ".hot", ".runner", ".correr", ".esquivar", ".lane", ".snake", ".viborita", ".serpiente", ".tetris", ".tetriscancelar", ".cancelar.tetris", ".top", ".toplocos", ".ranking", ".trivia", ".preguntados", ".triviacancelar", ".cancelartrivia", ".ttt", ".tresenraya", ".teteti", ".tttcancelar", ".cancelarttt", ".verdad", ".pregunta.verdad"}},
		{"𝐄⵿𝐜𝐨⵿𝐧𝐨⵿𝐦𝐢⵿𝐚", []string{".rw", ".recompensa", ".top", ".toplocos", ".ranking"}},
	}

	img, caption, err := ui.RenderMenu("Luis González", "3.0.0-Go", cats, botStartTime)
	if err != nil {
		client.SendMessage(context.Background(), evt.Info.Chat, &waE2E.Message{
			ExtendedTextMessage: &waE2E.ExtendedTextMessage{
				Text: &caption,
			},
		})
		return
	}

	uploaded, err := client.Upload(context.Background(), img, whatsmeow.MediaImage)
	if err != nil {
		return
	}

	client.SendMessage(context.Background(), evt.Info.Chat, &waE2E.Message{
		ImageMessage: &waE2E.ImageMessage{
			Caption:       &caption,
			Mimetype:      proto.String("image/jpeg"),
			JpegThumbnail: uploaded.Content[:32],
			Url:           &uploaded.URL,
			DirectPath:    &uploaded.DirectPath,
			MediaKey:      uploaded.MediaKey,
			FileEncSha256: uploaded.FileEncSHA256,
			FileSha256:    uploaded.FileSHA256,
			FileLength:    proto.Uint64(uint64(len(img))),
		},
	})
}

func init() {
	fmt.Println("☄️ Ꮯ᥆ոᥱ𝗍ᥲ Bot Iniciado")
}