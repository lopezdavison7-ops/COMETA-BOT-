
package store

import (
	"fmt"
	"math/rand"
	"time"
)

func (s *EconomyStore) Weekly(jid string) (string, error) {
	return s.claimTimedReward(
		jid,
		"last_weekly",
		7*24*60*60,
		3500,
		7000,
		350,
		"semanal",
	)
}

func (s *EconomyStore) Monthly(jid string) (string, error) {
	return s.claimTimedReward(
		jid,
		"last_monthly",
		30*24*60*60,
		15000,
		30000,
		1500,
		"mensual",
	)
}

func (s *EconomyStore) claimTimedReward(
	jid string,
	lastColumn string,
	cooldown int64,
	minComets int,
	maxComets int,
	xpReward int,
	period string,
) (string, error) {
	if lastColumn != "last_weekly" && lastColumn != "last_monthly" {
		return "Tipo de recompensa no válido.", fmt.Errorf("invalid_reward")
	}

	s.ensureUser(jid)

	s.mu.Lock()
	defer s.mu.Unlock()

	tx, err := s.db.Begin()
	if err != nil {
		return "No se pudo iniciar la recompensa.", err
	}
	defer tx.Rollback()

	var lastClaim int64

	query := fmt.Sprintf(
		`SELECT COALESCE(%s, 0) FROM users WHERE jid = ?`,
		lastColumn,
	)

	if err := tx.QueryRow(query, jid).Scan(&lastClaim); err != nil {
		return "No se pudo consultar tu recompensa.", err
	}

	now := time.Now().Unix()
	remaining := cooldown - (now - lastClaim)

	if lastClaim > 0 && remaining > 0 {
		return fmt.Sprintf(
			"Ya reclamaste tu premio %s. Volvé en %s.",
			period,
			formatRewardDuration(remaining),
		), fmt.Errorf("cooldown")
	}

	amount := rand.Intn(maxComets-minComets+1) + minComets

	updateBalance := fmt.Sprintf(
		`UPDATE users SET balance = COALESCE(balance, 0) + ?, %s = ? WHERE jid = ?`,
		lastColumn,
	)

	if _, err := tx.Exec(updateBalance, amount, now, jid); err != nil {
		return "No se pudo entregar tu recompensa.", err
	}

	var currentXP, currentLevel int

	if err := tx.QueryRow(
		`SELECT COALESCE(xp, 0), COALESCE(level, 1)
		 FROM users WHERE jid = ?`,
		jid,
	).Scan(&currentXP, &currentLevel); err != nil {
		return "No se pudo consultar tu experiencia.", err
	}

	if currentLevel < 1 {
		currentLevel = 1
	}

	newXP := currentXP + xpReward
	levelsGained := 0

	for {
		requiredXP := currentLevel * 100

		if newXP < requiredXP {
			break
		}

		newXP -= requiredXP
		currentLevel++
		levelsGained++
	}

	if _, err := tx.Exec(
		`UPDATE users SET xp = ?, level = ? WHERE jid = ?`,
		newXP,
		currentLevel,
		jid,
	); err != nil {
		return "No se pudo guardar tu experiencia.", err
	}

	if err := tx.Commit(); err != nil {
		return "No se pudo finalizar tu recompensa.", err
	}

	msg := fmt.Sprintf(
		"¡Reclamaste tu recompensa %s de *%d* cometas! 🎁\n+%d XP",
		period,
		amount,
		xpReward,
	)

	if levelsGained == 1 {
		msg += fmt.Sprintf(
			"\n🎉 ¡Subiste al nivel %d!",
			currentLevel,
		)
	} else if levelsGained > 1 {
		msg += fmt.Sprintf(
			"\n🎉 ¡Subiste %d niveles y ahora sos nivel %d!",
			levelsGained,
			currentLevel,
		)
	}

	return msg, nil
}

func formatRewardDuration(seconds int64) string {
	if seconds < 0 {
		seconds = 0
	}

	days := seconds / 86400
	hours := (seconds % 86400) / 3600
	minutes := (seconds % 3600) / 60
	remainingSeconds := seconds % 60

	switch {
	case days > 0:
		return fmt.Sprintf(
			"%d días, %d horas y %d minutos",
			days,
			hours,
			minutes,
		)

	case hours > 0:
		return fmt.Sprintf(
			"%d horas y %d minutos",
			hours,
			minutes,
		)

	case minutes > 0:
		return fmt.Sprintf(
			"%d minutos y %d segundos",
			minutes,
			remainingSeconds,
		)

	default:
		return fmt.Sprintf("%d segundos", seconds)
	}
}
