
package store

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

type allwOutcome struct {
	name   string
	min    int
	max    int
	weight int
}

func (s *EconomyStore) AllWork(jid string) (string, error) {
	s.ensureUser(jid)

	s.mu.Lock()
	defer s.mu.Unlock()

	tx, err := s.db.Begin()
	if err != nil {
		return "No se pudo iniciar ALL WORK.", err
	}
	defer tx.Rollback()

	now := time.Now().Unix()

	var lastAllw int64
	if err := tx.QueryRow(
		`SELECT COALESCE(last_allw, 0) FROM users WHERE jid = ?`,
		jid,
	).Scan(&lastAllw); err != nil {
		return "No se pudo consultar el cooldown de ALL WORK.", err
	}

	if lastAllw > 0 && now-lastAllw < 3600 {
		remaining := 3600 - (now - lastAllw)
		return fmt.Sprintf(
			"ALL WORK ya fue utilizado. Volvé en %s.",
			allwDuration(remaining),
		), fmt.Errorf("cooldown")
	}

	var balance, currentXP, currentLevel int

	if err := tx.QueryRow(
		`SELECT COALESCE(balance, 0), COALESCE(xp, 0), COALESCE(level, 1)
		 FROM users WHERE jid = ?`,
		jid,
	).Scan(&balance, &currentXP, &currentLevel); err != nil {
		return "No se pudieron consultar tus datos económicos.", err
	}

	if currentLevel < 1 {
		currentLevel = 1
	}

	lines := make([]string, 0, 6)
	totalCoins := 0
	totalXP := 0

	check := func(column string, cooldown int64) (int64, bool, error) {
		var last int64

		query := fmt.Sprintf(
			`SELECT COALESCE(%s, 0) FROM users WHERE jid = ?`,
			column,
		)

		if err := tx.QueryRow(query, jid).Scan(&last); err != nil {
			return 0, false, err
		}

		remaining := cooldown - (now - last)

		if last > 0 && remaining > 0 {
			return remaining, false, nil
		}

		return 0, true, nil
	}

	mark := func(column string) error {
		query := fmt.Sprintf(
			`UPDATE users SET %s = ? WHERE jid = ?`,
			column,
		)

		_, err := tx.Exec(query, now, jid)
		return err
	}

	// TRABAJO
	if remaining, ok, err := check("last_work", 600); err != nil {
		return "No se pudo consultar el cooldown de trabajo.", err
	} else if !ok {
		lines = append(lines, "❌ TRABAJO: le falta "+allwDuration(remaining))
	} else {
		reward := rand.Intn(150) + 50

		totalCoins += reward
		totalXP += 15

		lines = append(lines, fmt.Sprintf(
			"💼 TRABAJO: +%d cometas › +15 XP",
			reward,
		))

		if err := mark("last_work"); err != nil {
			return "No se pudo guardar el cooldown de trabajo.", err
		}
	}

	// AVENTURA
	if remaining, ok, err := check("last_adventure", 1200); err != nil {
		return "No se pudo consultar el cooldown de aventura.", err
	} else if !ok {
		lines = append(lines, "❌ AVENTURA: le falta "+allwDuration(remaining))
	} else {
		outcomes := []allwOutcome{
			{"Encontraste un cofre olvidado", 100, 250, 25},
			{"Derrotaste a un bandido", 150, 300, 20},
			{"Rescataste a un mercader", 200, 400, 15},
			{"Descubriste una cueva secreta", 250, 500, 12},
			{"Completaste una misión épica", 400, 800, 8},
			{"Encontraste un tesoro legendario", 600, 1200, 5},
			{"Te perdiste y no ganaste nada", 0, 0, 15},
		}

		chosen := pickAllwOutcome(outcomes)
		reward := randomAllw(chosen.min, chosen.max)

		totalCoins += reward

		if reward > 0 {
			lines = append(lines, fmt.Sprintf(
				"🗺️ AVENTURA: %s › +%d cometas",
				chosen.name,
				reward,
			))
		} else {
			lines = append(lines, "🗺️ AVENTURA: "+chosen.name+" › +0 cometas")
		}

		if err := mark("last_adventure"); err != nil {
			return "No se pudo guardar el cooldown de aventura.", err
		}
	}

	// MINAR
	if remaining, ok, err := check("last_mine", 900); err != nil {
		return "No se pudo consultar el cooldown de minería.", err
	} else if !ok {
		lines = append(lines, "❌ MINAR: le falta "+allwDuration(remaining))
	} else {
		resources := []allwOutcome{
			{"🪨 Piedra", 10, 30, 30},
			{"⚫ Carbón", 20, 50, 25},
			{"⚙️ Hierro", 40, 80, 20},
			{"🥇 Oro", 80, 150, 15},
			{"💎 Diamante", 150, 300, 8},
			{"💚 Esmeralda", 200, 400, 2},
		}

		chosen := pickAllwOutcome(resources)
		reward := randomAllw(chosen.min, chosen.max)

		totalCoins += reward

		lines = append(lines, fmt.Sprintf(
			"⛏️ MINAR: %s › +%d cometas",
			chosen.name,
			reward,
		))

		if err := mark("last_mine"); err != nil {
			return "No se pudo guardar el cooldown de minería.", err
		}
	}

	// PESCAR
	if remaining, ok, err := check("last_fish", 900); err != nil {
		return "No se pudo consultar el cooldown de pesca.", err
	} else if !ok {
		lines = append(lines, "❌ PESCAR: le falta "+allwDuration(remaining))
	} else {
		if rand.Intn(100) < 10 {
			lines = append(lines, "🎣 PESCAR: se rompió la caña; no atrapaste nada")
		} else {
			fishes := []allwOutcome{
				{"🐟 Sardina", 15, 35, 30},
				{"🐠 Trucha", 30, 60, 25},
				{"🍣 Salmón", 50, 100, 20},
				{"🐡 Atún", 80, 150, 15},
				{"⚔️ Pez espada", 120, 250, 8},
				{"✨ Pez dorado", 200, 400, 2},
			}

			chosen := pickAllwOutcome(fishes)
			reward := randomAllw(chosen.min, chosen.max)

			totalCoins += reward

			lines = append(lines, fmt.Sprintf(
				"🎣 PESCAR: %s › +%d cometas",
				chosen.name,
				reward,
			))
		}

		if err := mark("last_fish"); err != nil {
			return "No se pudo guardar el cooldown de pesca.", err
		}
	}

	// MENDIGAR
	if remaining, ok, err := check("last_beg", 300); err != nil {
		return "No se pudo consultar el cooldown de mendigar.", err
	} else if !ok {
		lines = append(lines, "❌ MENDIGAR: le falta "+allwDuration(remaining))
	} else {
		reward := rand.Intn(41) + 10

		totalCoins += reward

		lines = append(lines, fmt.Sprintf(
			"🙏 MENDIGAR: +%d cometas",
			reward,
		))

		if err := mark("last_beg"); err != nil {
			return "No se pudo guardar el cooldown de mendigar.", err
		}
	}

	// SLUT
	if remaining, ok, err := check("last_slut", 1200); err != nil {
		return "No se pudo consultar el cooldown de SLUT.", err
	} else if !ok {
		lines = append(lines, "❌ SLUT: le falta "+allwDuration(remaining))
	} else {
		outcomes := []allwOutcome{
			{"Un cliente generoso te dio propina", 50, 150, 30},
			{"Trabajaste toda la noche", 100, 250, 25},
			{"Un VIP te pagó extra", 200, 400, 20},
			{"Ganaste un concurso de baile", 300, 600, 15},
			{"Te contrataron para un evento exclusivo", 500, 1000, 8},
			{"La policía te multó", -100, -50, 2},
		}

		chosen := pickAllwOutcome(outcomes)
		reward := randomAllw(chosen.min, chosen.max)

		if reward < 0 && -reward > balance {
			reward = -balance
		}

		totalCoins += reward

		if reward >= 0 {
			lines = append(lines, fmt.Sprintf(
				"💋 SLUT: %s › +%d cometas",
				chosen.name,
				reward,
			))
		} else {
			lines = append(lines, fmt.Sprintf(
				"💋 SLUT: %s › -%d cometas",
				chosen.name,
				-reward,
			))
		}

		if err := mark("last_slut"); err != nil {
			return "No se pudo guardar el cooldown de SLUT.", err
		}
	}

	newBalance := balance + totalCoins
	if newBalance < 0 {
		newBalance = 0
	}

	newXP := currentXP + totalXP
	levelsGained := 0

	for newXP >= currentLevel*100 {
		newXP -= currentLevel * 100
		currentLevel++
		levelsGained++
	}

	if _, err := tx.Exec(
		`UPDATE users
		 SET balance = ?, xp = ?, level = ?, last_allw = ?
		 WHERE jid = ?`,
		newBalance,
		newXP,
		currentLevel,
		now,
		jid,
	); err != nil {
		return "No se pudieron guardar los resultados de ALL WORK.", err
	}

	if err := tx.Commit(); err != nil {
		return "No se pudieron confirmar los resultados de ALL WORK.", err
	}

	result := strings.Join(lines, "\n")
	result += fmt.Sprintf("\n\n💰 GANANCIA NETA › %+d cometas", totalCoins)
	result += fmt.Sprintf("\n✨ EXPERIENCIA TOTAL › +%d XP", totalXP)

	if levelsGained == 1 {
		result += fmt.Sprintf("\n⭐ ¡SUBISTE AL NIVEL %d!", currentLevel)
	} else if levelsGained > 1 {
		result += fmt.Sprintf(
			"\n⭐ ¡SUBISTE %d NIVELES! Ahora sos nivel %d.",
			levelsGained,
			currentLevel,
		)
	}

	result += "\n\n> Podés volver a usar ALLW en 1 hora."

	return result, nil
}

func pickAllwOutcome(outcomes []allwOutcome) allwOutcome {
	totalWeight := 0

	for _, outcome := range outcomes {
		totalWeight += outcome.weight
	}

	roll := rand.Intn(totalWeight)

	for _, outcome := range outcomes {
		if roll < outcome.weight {
			return outcome
		}

		roll -= outcome.weight
	}

	return outcomes[len(outcomes)-1]
}

func randomAllw(min, max int) int {
	if max <= min {
		return min
	}

	return rand.Intn(max-min+1) + min
}

func allwDuration(seconds int64) string {
	if seconds < 0 {
		seconds = 0
	}

	days := seconds / 86400
	hours := (seconds % 86400) / 3600
	minutes := (seconds % 3600) / 60
	remainingSeconds := seconds % 60

	if days > 0 {
		return fmt.Sprintf("%dd %dh %dm", days, hours, minutes)
	}

	if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	}

	if minutes > 0 {
		return fmt.Sprintf("%dm %ds", minutes, remainingSeconds)
	}

	return fmt.Sprintf("%ds", remainingSeconds)
}
