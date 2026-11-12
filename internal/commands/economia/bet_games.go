package economia

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"

	"cometa-bot/internal/registry"
	"cometa-bot/internal/store"
)

const bettingCategory = "𝐄⵿𝐜𝵿𝐨⵿𝐧𝐨⵿𝐦⵿𝐢⵿𝐚"

func init() {
	commands := []*registry.Command{
		{Name: "dice", Aliases: []string{"dado"}, Description: "apuesta al resultado de un dado", Category: bettingCategory, Handler: BetDice},
		{Name: "highlow", Aliases: []string{"altobajo", "mayormenor"}, Description: "adivina si saldrá alto o bajo", Category: bettingCategory, Handler: BetHighLow},
		{Name: "rpsbet", Aliases: []string{"piedrapapel"}, Description: "apuesta en piedra, papel o tijera", Category: bettingCategory, Handler: BetRPS},
		{Name: "blackjack", Aliases: []string{"bj"}, Description: "juega una ronda automática de blackjack", Category: bettingCategory, Handler: BetBlackjack},
		{Name: "race", Aliases: []string{"carrera"}, Description: "apuesta por un caballo de carreras", Category: bettingCategory, Handler: BetRace},
		{Name: "plinko", Aliases: []string{"caida"}, Description: "lanza una ficha por Plinko", Category: bettingCategory, Handler: BetPlinko},
		{Name: "limbo", Aliases: nil, Description: "elige un multiplicador y prueba tu suerte", Category: bettingCategory, Handler: BetLimbo},
		{Name: "numberbet", Aliases: []string{"apuestanumero"}, Description: "adivina un número del 0 al 9", Category: bettingCategory, Handler: BetNumber},
		{Name: "crash", Aliases: []string{"choque"}, Description: "elige cuándo retirar antes del crash", Category: bettingCategory, Handler: BetCrash},
		{Name: "evenodd", Aliases: []string{"parimpar"}, Description: "apuesta a un número par o impar", Category: bettingCategory, Handler: BetEvenOdd},
		{Name: "minesbet", Aliases: []string{"minasapuesta"}, Description: "elige una casilla y evita las minas", Category: bettingCategory, Handler: BetMines},
		{Name: "jackpot", Aliases: []string{"premio"}, Description: "busca el premio mayor", Category: bettingCategory, Handler: BetJackpot},
		{Name: "wheel", Aliases: []string{"ruedafortuna"}, Description: "gira una rueda de premios", Category: bettingCategory, Handler: BetWheel},
		{Name: "baccarat", Aliases: []string{"bacara"}, Description: "apuesta a jugador, banca o empate", Category: bettingCategory, Handler: BetBaccarat},
		{Name: "war", Aliases: []string{"guerra"}, Description: "compara tu carta contra la banca", Category: bettingCategory, Handler: BetWar},
		{Name: "bethelp", Aliases: []string{"juegos"}, Description: "muestra los juegos de apuestas", Category: bettingCategory, Handler: BetHelp},
	}

	for _, command := range commands {
		registry.Register(command)
	}
}

func sendBetText(client *whatsmeow.Client, evt *events.Message, message string) {
	if client == nil || evt == nil {
		return
	}

	client.SendMessage(context.Background(), evt.Info.Chat, &waE2E.Message{
		ExtendedTextMessage: &waE2E.ExtendedTextMessage{
			Text:        proto.String(message),
			ContextInfo: &waE2E.ContextInfo{QuotedMessage: evt.Message},
		},
	})
}

func startBet(evt *events.Message, client *whatsmeow.Client, args []string, required int, usage string) (int, bool) {
	if store.GlobalEconomy == nil {
		sendBetText(client, evt, "⚠️ La economía no está inicializada.")
		return 0, false
	}

	if len(args) < required {
		sendBetText(client, evt, "🎲 Uso: "+usage)
		return 0, false
	}

	amount, err := strconv.Atoi(args[0])
	if err != nil || amount < 10 || amount > store.MaxBetAmount {
		sendBetText(client, evt, "⚠️ La apuesta debe ser de 10 a 10.000.000 cometas.")
		return 0, false
	}

	jid := evt.Info.Sender.String()
	if store.GlobalEconomy.GetBalance(jid) < amount {
		sendBetText(client, evt, "💸 No tienes suficientes cometas en la billetera.")
		return 0, false
	}

	return amount, true
}

func settleBet(evt *events.Message, client *whatsmeow.Client, game string, amount, payout int, details string) {
	if payout < 0 {
		payout = 0
	}

	err := store.GlobalEconomy.SettleBet(evt.Info.Sender.String(), amount, payout)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrInsufficientBetBalance):
			sendBetText(client, evt, "💸 Tu saldo cambió y ya no alcanza para cubrir esta apuesta. No se realizó ningún cobro.")
		case errors.Is(err, store.ErrInvalidBetAmount):
			sendBetText(client, evt, "⚠️ La cantidad de la apuesta no es válida.")
		default:
			sendBetText(client, evt, "⚠️ No se pudo procesar la apuesta: "+err.Error())
		}
		return
	}

	net := payout - amount
	status := "💀 *Perdiste.*"

	if net > 0 {
		status = "🎉 *Ganaste.*"
	} else if net == 0 {
		status = "➖ *Apuesta devuelta.*"
	}

	message := fmt.Sprintf(
		"╭〔 🎲 %s 〕\n> ㅤ\n%s\n🪙 *Apuesta:* %s cometas\n💰 *Retorno:* %s cometas\n📊 *Balance neto:* %s cometas\n%s\n╰━━━━━━━━━━━━━━",
		strings.ToUpper(game),
		details,
		formatComets(amount),
		formatComets(payout),
		signedComets(net),
		status,
	)

	sendBetText(client, evt, message)
}

func formatComets(value int) string {
	return strconv.Itoa(value)
}

func signedComets(value int) string {
	if value > 0 {
		return "+" + formatComets(value)
	}
	return formatComets(value)
}

func parseMultiplier(value string, minimum, maximum float64) (float64, bool) {
	value = strings.ReplaceAll(strings.TrimSpace(value), ",", ".")

	multiplier, err := strconv.ParseFloat(value, 64)
	if err != nil || math.IsNaN(multiplier) || math.IsInf(multiplier, 0) {
		return 0, false
	}

	return multiplier, multiplier >= minimum && multiplier <= maximum
}

func roundedPayout(amount int, multiplier float64) int {
	return int(math.Round(float64(amount) * multiplier))
}

// ============================================================
// DADO
// /dice <cantidad> <1-6>
// ============================================================

func BetDice(evt *events.Message, client *whatsmeow.Client, args []string) {
	amount, ok := startBet(evt, client, args, 2, "/dice <cantidad> <1-6>")
	if !ok {
		return
	}

	guess, err := strconv.Atoi(args[1])
	if err != nil || guess < 1 || guess > 6 {
		sendBetText(client, evt, "Elige un número del 1 al 6. Ejemplo: /dice 100 4")
		return
	}

	result := rand.Intn(6) + 1
	payout := 0

	if guess == result {
		payout = amount * 5
	}

	settleBet(evt, client, "DADO", amount, payout,
		fmt.Sprintf("🎯 Elegiste *%d*. Salió *%d*.", guess, result))
}

// ============================================================
// ALTO O BAJO
// /highlow <cantidad> <alto|bajo>
// ============================================================

func BetHighLow(evt *events.Message, client *whatsmeow.Client, args []string) {
	amount, ok := startBet(evt, client, args, 2, "/highlow <cantidad> <alto|bajo>")
	if !ok {
		return
	}

	choice := strings.ToLower(args[1])
	if choice != "alto" && choice != "bajo" && choice != "a" && choice != "b" {
		sendBetText(client, evt, "Debes elegir *alto* o *bajo*.")
		return
	}

	result := rand.Intn(100) + 1
	isHigh := result >= 51

	won := (isHigh && (choice == "alto" || choice == "a")) ||
		(!isHigh && (choice == "bajo" || choice == "b"))

	payout := 0
	if won {
		payout = roundedPayout(amount, 1.9)
	}

	label := "BAJO"
	if isHigh {
		label = "ALTO"
	}

	settleBet(evt, client, "ALTO O BAJO", amount, payout,
		fmt.Sprintf("🎯 Elegiste *%s*. Salió *%d (%s)*.", strings.ToUpper(choice), result, label))
}

// ============================================================
// PIEDRA, PAPEL O TIJERA
// /rpsbet <cantidad> <piedra|papel|tijera>
// ============================================================

func BetRPS(evt *events.Message, client *whatsmeow.Client, args []string) {
	amount, ok := startBet(evt, client, args, 2, "/rpsbet <cantidad> <piedra|papel|tijera>")
	if !ok {
		return
	}

	choice := strings.ToLower(args[1])
	valid := map[string]bool{
		"piedra": true,
		"papel":  true,
		"tijera": true,
	}

	if !valid[choice] {
		sendBetText(client, evt, "Elige piedra, papel o tijera.")
		return
	}

	options := []string{"piedra", "papel", "tijera"}
	result := options[rand.Intn(len(options))]
	payout := 0
	status := "perdiste"

	if choice == result {
		payout = amount
		status = "empate"
	} else if (choice == "piedra" && result == "tijera") ||
		(choice == "papel" && result == "piedra") ||
		(choice == "tijera" && result == "papel") {
		payout = roundedPayout(amount, 1.9)
		status = "ganaste"
	}

	settleBet(evt, client, "PIEDRA PAPEL TIJERA", amount, payout,
		fmt.Sprintf("🫵 Tú: *%s*\n🤖 Bot: *%s*\nResultado: *%s*.", choice, result, status))
}

// ============================================================
// BLACKJACK
// /blackjack <cantidad>
// ============================================================

type blackjackCard struct {
	name  string
	value int
}

func drawBlackjackCard() blackjackCard {
	value := rand.Intn(13) + 1

	switch value {
	case 1:
		return blackjackCard{name: "A", value: 11}
	case 11:
		return blackjackCard{name: "J", value: 10}
	case 12:
		return blackjackCard{name: "Q", value: 10}
	case 13:
		return blackjackCard{name: "K", value: 10}
	default:
		return blackjackCard{name: strconv.Itoa(value), value: value}
	}
}

func blackjackScore(cards []blackjackCard) int {
	score := 0
	aces := 0

	for _, card := range cards {
		score += card.value
		if card.name == "A" {
			aces++
		}
	}

	for score > 21 && aces > 0 {
		score -= 10
		aces--
	}

	return score
}

func blackjackHand(cards []blackjackCard) string {
	parts := make([]string, 0, len(cards))

	for _, card := range cards {
		parts = append(parts, card.name)
	}

	return strings.Join(parts, " · ")
}

func BetBlackjack(evt *events.Message, client *whatsmeow.Client, args []string) {
	amount, ok := startBet(evt, client, args, 1, "/blackjack <cantidad>")
	if !ok {
		return
	}

	player := []blackjackCard{drawBlackjackCard(), drawBlackjackCard()}
	dealer := []blackjackCard{drawBlackjackCard(), drawBlackjackCard()}

	playerNatural := blackjackScore(player) == 21
	dealerNatural := blackjackScore(dealer) == 21

	for blackjackScore(player) < 17 {
		player = append(player, drawBlackjackCard())
	}

	for blackjackScore(dealer) < 17 {
		dealer = append(dealer, drawBlackjackCard())
	}

	playerScore := blackjackScore(player)
	dealerScore := blackjackScore(dealer)
	payout := 0
	result := "La banca gana."

	switch {
	case dealerNatural && !playerNatural:
		result = "Blackjack de la banca."
	case playerNatural && !dealerNatural:
		payout = roundedPayout(amount, 2.5)
		result = "¡Blackjack natural!"
	case playerScore > 21:
		result = "Te pasaste de 21."
	case dealerScore > 21 || playerScore > dealerScore:
		payout = amount * 2
		result = "¡Ganaste la ronda!"
	case playerScore == dealerScore:
		payout = amount
		result = "Empate. Se devuelve tu apuesta."
	}

	details := fmt.Sprintf(
		"🃏 Tú: *%s* = %d\n🏦 Banca: *%s* = %d\n%s",
		blackjackHand(player),
		playerScore,
		blackjackHand(dealer),
		dealerScore,
		result,
	)

	settleBet(evt, client, "BLACKJACK", amount, payout, details)
}

// ============================================================
// CARRERA
// /race <cantidad> <caballo 1-5>
// ============================================================

func BetRace(evt *events.Message, client *whatsmeow.Client, args []string) {
	amount, ok := startBet(evt, client, args, 2, "/race <cantidad> <caballo 1-5>")
	if !ok {
		return
	}

	chosen, err := strconv.Atoi(args[1])
	if err != nil || chosen < 1 || chosen > 5 {
		sendBetText(client, evt, "Elige un caballo del 1 al 5.")
		return
	}

	winner := rand.Intn(5) + 1
	payout := 0

	if chosen == winner {
		payout = amount * 4
	}

	settleBet(evt, client, "CARRERA", amount, payout,
		fmt.Sprintf("🏇 Elegiste el caballo *%d*.\n🏁 Ganó el caballo *%d*.", chosen, winner))
}

// ============================================================
// PLINKO
// /plinko <cantidad> <bajo|medio|alto>
// ============================================================

func BetPlinko(evt *events.Message, client *whatsmeow.Client, args []string) {
	amount, ok := startBet(evt, client, args, 2, "/plinko <cantidad> <bajo|medio|alto>")
	if !ok {
		return
	}

	choice := strings.ToLower(args[1])
	var multipliers []int

	switch choice {
	case "bajo", "low":
		multipliers = []int{0, 1, 1, 1, 1, 2, 1, 1, 1, 0}
		choice = "bajo"
	case "medio", "med", "medium":
		multipliers = []int{0, 0, 1, 2, 2, 0, 2, 1, 1, 0}
		choice = "medio"
	case "alto", "high":
		multipliers = []int{0, 0, 0, 0, 5, 0, 0, 3, 0, 0}
		choice = "alto"
	default:
		sendBetText(client, evt, "Elige riesgo *bajo*, *medio* o *alto*.")
		return
	}

	index := rand.Intn(len(multipliers))
	multiplier := multipliers[index]
	payout := amount * multiplier

	settleBet(evt, client, "PLINKO", amount, payout,
		fmt.Sprintf("📍 Riesgo: *%s*\n🟡 Casilla: *%d/10*\n✖️ Multiplicador: *x%d*", choice, index+1, multiplier))
}

// ============================================================
// LIMBO
// /limbo <cantidad> <multiplicador 1.10-10>
// ============================================================

func BetLimbo(evt *events.Message, client *whatsmeow.Client, args []string) {
	amount, ok := startBet(evt, client, args, 2, "/limbo <cantidad> <multiplicador 1.10-10>")
	if !ok {
		return
	}

	multiplier, valid := parseMultiplier(args[1], 1.10, 10)
	if !valid {
		sendBetText(client, evt, "El multiplicador debe estar entre *1.10x* y *10.00x*.")
		return
	}

	chance := 95.0 / multiplier
	roll := rand.Float64() * 100
	payout := 0

	if roll < chance {
		payout = roundedPayout(amount, multiplier)
	}

	settleBet(evt, client, "LIMBO", amount, payout,
		fmt.Sprintf("🎯 Objetivo: *%.2fx*\n📈 Probabilidad aproximada: *%.2f%%*\n🎲 Resultado: *%.2f*", multiplier, chance, roll))
}

// ============================================================
// NÚMERO
// /numberbet <cantidad> <0-9>
// ============================================================

func BetNumber(evt *events.Message, client *whatsmeow.Client, args []string) {
	amount, ok := startBet(evt, client, args, 2, "/numberbet <cantidad> <0-9>")
	if !ok {
		return
	}

	guess, err := strconv.Atoi(args[1])
	if err != nil || guess < 0 || guess > 9 {
		sendBetText(client, evt, "Elige un dígito del 0 al 9.")
		return
	}

	result := rand.Intn(10)
	payout := 0

	if guess == result {
		payout = amount * 9
	}

	settleBet(evt, client, "NÚMERO", amount, payout,
		fmt.Sprintf("🔢 Tu número: *%d*\n🎯 Resultado: *%d*", guess, result))
}

// ============================================================
// CRASH
// /crash <cantidad> <retirada 1.10-10>
// ============================================================

func BetCrash(evt *events.Message, client *whatsmeow.Client, args []string) {
	amount, ok := startBet(evt, client, args, 2, "/crash <cantidad> <retirada 1.10-10>")
	if !ok {
		return
	}

	target, valid := parseMultiplier(args[1], 1.10, 10)
	if !valid {
		sendBetText(client, evt, "El objetivo debe estar entre *1.10x* y *10.00x*.")
		return
	}

	crashPoint := 0.99 / (1 - rand.Float64())
	if crashPoint > 100 {
		crashPoint = 100
	}

	payout := 0
	if crashPoint >= target {
		payout = roundedPayout(amount, target)
	}

	settleBet(evt, client, "CRASH", amount, payout,
		fmt.Sprintf("🚀 Querías retirar en *%.2fx*.\n💥 El juego explotó en *%.2fx*.", target, crashPoint))
}

// ============================================================
// PAR O IMPAR
// /evenodd <cantidad> <par|impar>
// ============================================================

func BetEvenOdd(evt *events.Message, client *whatsmeow.Client, args []string) {
	amount, ok := startBet(evt, client, args, 2, "/evenodd <cantidad> <par|impar>")
	if !ok {
		return
	}

	choice := strings.ToLower(args[1])
	if choice != "par" && choice != "impar" && choice != "p" && choice != "i" {
		sendBetText(client, evt, "Debes elegir *par* o *impar*.")
		return
	}

	result := rand.Intn(100) + 1
	label := "impar"

	if result%2 == 0 {
		label = "par"
	}

	chosen := choice
	if chosen == "p" {
		chosen = "par"
	}
	if chosen == "i" {
		chosen = "impar"
	}

	payout := 0
	if chosen == label {
		payout = roundedPayout(amount, 1.9)
	}

	settleBet(evt, client, "PAR O IMPAR", amount, payout,
		fmt.Sprintf("🎯 Elegiste *%s*.\n🔢 Salió *%d (%s)*.", chosen, result, label))
}

// ============================================================
// MINAS
// /minesbet <cantidad> <minas 2-20> <casilla 1-25>
// ============================================================

func BetMines(evt *events.Message, client *whatsmeow.Client, args []string) {
	amount, ok := startBet(evt, client, args, 3, "/minesbet <cantidad> <minas 2-20> <casilla 1-25>")
	if !ok {
		return
	}

	mineCount, err := strconv.Atoi(args[1])
	if err != nil || mineCount < 2 || mineCount > 20 {
		sendBetText(client, evt, "El número de minas debe estar entre 2 y 20.")
		return
	}

	cell, err := strconv.Atoi(args[2])
	if err != nil || cell < 1 || cell > 25 {
		sendBetText(client, evt, "Elige una casilla del 1 al 25.")
		return
	}

	mines := rand.Perm(25)[:mineCount]
	hitMine := false

	for _, mine := range mines {
		if mine+1 == cell {
			hitMine = true
			break
		}
	}

	payout := 0
	multiplier := 0.95 * 25 / float64(25-mineCount)

	if !hitMine {
		payout = roundedPayout(amount, multiplier)
	}

	result := "💎 Casilla segura."
	if hitMine {
		result = "💣 ¡Encontraste una mina!"
	}

	settleBet(evt, client, "MINAS", amount, payout,
		fmt.Sprintf("🧨 Minas: *%d/25*\n📦 Casilla elegida: *%d*\n%s\n✖️ Multiplicador: *%.2fx*", mineCount, cell, result, multiplier))
}

// ============================================================
// JACKPOT
// /jackpot <cantidad>
// ============================================================

func BetJackpot(evt *events.Message, client *whatsmeow.Client, args []string) {
	amount, ok := startBet(evt, client, args, 1, "/jackpot <cantidad>")
	if !ok {
		return
	}

	roll := rand.Intn(100) + 1
	payout := 0

	if roll == 100 {
		payout = amount * 80
	}

	settleBet(evt, client, "JACKPOT", amount, payout,
		fmt.Sprintf("🎰 Número ganador: *100*\n🎲 Tu resultado: *%d/100*\n💎 Premio: *x80*", roll))
}

// ============================================================
// RULETA DE PREMIOS
// /wheel <cantidad>
// ============================================================

func BetWheel(evt *events.Message, client *whatsmeow.Client, args []string) {
	amount, ok := startBet(evt, client, args, 1, "/wheel <cantidad>")
	if !ok {
		return
	}

	multipliers := []int{0, 0, 1, 1, 1, 2, 0, 3, 1, 0}
	index := rand.Intn(len(multipliers))
	multiplier := multipliers[index]
	payout := amount * multiplier

	settleBet(evt, client, "RULETA DE PREMIOS", amount, payout,
		fmt.Sprintf("🎡 Sector: *%d/10*\n✖️ Premio: *x%d*", index+1, multiplier))
}

// ============================================================
// BACCARAT
// /baccarat <cantidad> <jugador|banca|empate>
// ============================================================

func BetBaccarat(evt *events.Message, client *whatsmeow.Client, args []string) {
	amount, ok := startBet(evt, client, args, 2, "/baccarat <cantidad> <jugador|banca|empate>")
	if !ok {
		return
	}

	choice := strings.ToLower(args[1])
	if choice != "jugador" && choice != "banca" && choice != "empate" {
		sendBetText(client, evt, "Elige *jugador*, *banca* o *empate*.")
		return
	}

	player := (rand.Intn(10) + rand.Intn(10)) % 10
	banker := (rand.Intn(10) + rand.Intn(10)) % 10

	winner := "empate"
	if player > banker {
		winner = "jugador"
	} else if banker > player {
		winner = "banca"
	}

	payout := 0
	if choice == winner {
		switch winner {
		case "jugador":
			payout = amount * 2
		case "banca":
			payout = roundedPayout(amount, 1.95)
		case "empate":
			payout = amount * 8
		}
	}

	settleBet(evt, client, "BACCARAT", amount, payout,
		fmt.Sprintf("🧍 Jugador: *%d*\n🏦 Banca: *%d*\n🎯 Elegiste: *%s*\n🏆 Resultado: *%s*", player, banker, choice, winner))
}

// ============================================================
// GUERRA
// /war <cantidad>
// ============================================================

func BetWar(evt *events.Message, client *whatsmeow.Client, args []string) {
	amount, ok := startBet(evt, client, args, 1, "/war <cantidad>")
	if !ok {
		return
	}

	player := rand.Intn(13) + 1
	banker := rand.Intn(13) + 1
	payout := 0
	result := "La banca gana."

	if player > banker {
		payout = roundedPayout(amount, 1.9)
		result = "Tu carta es mayor."
	} else if player == banker {
		payout = amount
		result = "Empate. Se devuelve la apuesta."
	}

	settleBet(evt, client, "GUERRA", amount, payout,
		fmt.Sprintf("🃏 Tu carta: *%d*\n🏦 Carta de la banca: *%d*\n%s", player, banker, result))
}

// ============================================================
// AYUDA DE APUESTAS
// /bethelp
// ============================================================

func BetHelp(evt *events.Message, client *whatsmeow.Client, _ []string) {
	message := `╭〔 🎲 JUEGOS DE APUESTAS 〕
> Apuesta mínima: 10 cometas
> Apuesta máxima: 10.000.000 cometas

🎲 /dice <cantidad> <1-6>
📈 /highlow <cantidad> <alto|bajo>
✂️ /rpsbet <cantidad> <piedra|papel|tijera>
🃏 /blackjack <cantidad>
🏇 /race <cantidad> <caballo 1-5>
🟡 /plinko <cantidad> <bajo|medio|alto>
🎯 /limbo <cantidad> <1.10-10>
🔢 /numberbet <cantidad> <0-9>
💥 /crash <cantidad> <1.10-10>
⚖️ /evenodd <cantidad> <par|impar>
💣 /minesbet <cantidad> <minas 2-20> <casilla 1-25>
💎 /jackpot <cantidad>
🎡 /wheel <cantidad>
🏦 /baccarat <cantidad> <jugador|banca|empate>
🃏 /war <cantidad>

También existen: /flip, /roulette y /slots.
Consulta /einfo para ver la economía completa.
╰━━━━━━━━━━━━━━━━━━━━`

	sendBetText(client, evt, message)
}