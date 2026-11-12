package economia

import (
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"
)

// W conserva el nombre de la función anterior para compatibilidad interna.
// El comando se registra una sola vez en work.go.
func W(evt *events.Message, client *whatsmeow.Client, args []string) {
	Work(evt, client, args)
}