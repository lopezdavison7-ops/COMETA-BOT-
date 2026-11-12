package store

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"
)

type ActivityResult struct {
	Action       string
	Title        string
	Emoji        string
	Detail       string
	Reward       int
	XP           int
	MissionBonus int
	Skipped      bool
	Remaining    int64
	Level        int
	LeveledUp    bool
}

type MissionProgress struct {
	ID          string
	Name        string
	Description string
	Progress    int
	Target      int
	Reward      int
	Completed   bool
}

var activityCooldowns = map[string]int64{
	"work":      600,
	"adventure": 1200,
	"mine":      900,
	"fish":      900,
	"beg":       300,
	"slut":      1200,
	"hunt":      600,
	"dig":       480,
	"talar":     1500,
	"crime":     1800,
}

var activityColumns = map[string]string{
	"work":      "last_work",
	"adventure": "last_adventure",
	"mine":      "last_mine",
	"fish":      "last_fish",
	"beg":       "last_beg",
	"slut":      "last_slut",
	"hunt":      "last_hunt",
	"dig":       "last_dig",
	"talar":     "last_talar",
	"crime":     "last_crime",
}

var dailyMissionDefinitions = []MissionProgress{
	{
		ID:          "activities",
		Name:        "Nuevas actividades",
		Description: "Completa 4 actividades nuevas: caza, excavación, tala o crimen",
		Target:      4,
		Reward:      500,
	},
	{
		ID:          "gather",
		Name:        "Recolector",
		Description: "Completa 3 actividades de caza, excavación o tala",
		Target:      3,
		Reward:      350,
	},
	{
		ID:          "earn",
		Name:        "Buen negocio",
		Description: "Consigue 600 cometas con las nuevas actividades",
		Target:      600,
		Reward:      700,
	},
}

func (s *EconomyStore) InitExtraTables() error {
	statements := []string{
		`CREATE TABLE IF NOT EXISTS economy_effects (
			jid TEXT PRIMARY KEY,
			xp_boost INTEGER NOT NULL DEFAULT 0,
			shield INTEGER NOT NULL DEFAULT 0
		)`,
		`CREATE TABLE IF NOT EXISTS economy_daily_streaks (
			jid TEXT PRIMARY KEY,
			streak INTEGER NOT NULL DEFAULT 0,
			last_claim_date TEXT NOT NULL DEFAULT ''
		)`,
		`CREATE TABLE IF NOT EXISTS economy_mission_progress (
			jid TEXT NOT NULL,
			mission_date TEXT NOT NULL,
			mission_id TEXT NOT NULL,
			progress INTEGER NOT NULL DEFAULT 0,
			completed INTEGER NOT NULL DEFAULT 0,
			PRIMARY KEY (jid, mission_date, mission_id)
		)`,
	}

	for _, statement := range statements {
		if _, err := s.db.Exec(statement); err != nil {
			return err
		}
	}

	return nil
}

func (s *EconomyStore) ActivateEffect(jid, effect string) error {
	if effect != "xp_boost" && effect != "shield" {
		return fmt.Errorf("efecto no reconocido")
	}

	s.ensureUser(jid)

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, err := s.db.Exec(
		`INSERT OR IGNORE INTO economy_effects (jid) VALUES (?)`,
		jid,
	); err != nil {
		return err
	}

	column := effect
	_, err := s.db.Exec(
		`UPDATE economy_effects SET `+column+` = 1 WHERE jid = ?`,
		jid,
	)

	return err
}

func (s *EconomyStore) ConsumeEffect(jid, effect string) bool {
	if effect != "xp_boost" && effect != "shield" {
		return false
	}

	s.ensureUser(jid)

	s.mu.Lock()
	defer s.mu.Unlock()

	result, err := s.db.Exec(
		`UPDATE economy_effects
		SET `+effect+` = 0
		WHERE jid = ? AND COALESCE(`+effect+`, 0) > 0`,
		jid,
	)
	if err != nil {
		return false
	}

	changed, err := result.RowsAffected()
	return err == nil && changed > 0
}

func (s *EconomyStore) ApplyXPBoost(jid string, baseXP int) (bool, int, bool, error) {
	if baseXP <= 0 {
		return false, 0, false, nil
	}

	s.ensureUser(jid)

	s.mu.Lock()
	defer s.mu.Unlock()

	tx, err := s.db.Begin()
	if err != nil {
		return false, 0, false, err
	}
	defer tx.Rollback()

	var active int

	if err := tx.QueryRow(
		`SELECT COALESCE(xp_boost, 0)
		FROM economy_effects
		WHERE jid = ?`,
		jid,
	).Scan(&active); err != nil {
		if err == sql.ErrNoRows {
			return false, 0, false, nil
		}

		return false, 0, false, err
	}

	if active <= 0 {
		return false, 0, false, nil
	}

	if _, err := tx.Exec(
		`UPDATE economy_effects SET xp_boost = 0 WHERE jid = ?`,
		jid,
	); err != nil {
		return false, 0, false, err
	}

	level, leveled, err := addXPToTx(tx, jid, baseXP)
	if err != nil {
		return false, 0, false, err
	}

	if err := tx.Commit(); err != nil {
		return false, 0, false, err
	}

	return true, level, leveled, nil
}

func (s *EconomyStore) HasEffect(jid, effect string) bool {
	if effect != "xp_boost" && effect != "shield" {
		return false
	}

	var active int

	err := s.db.QueryRow(
		`SELECT COALESCE(`+effect+`, 0)
		FROM economy_effects
		WHERE jid = ?`,
		jid,
	).Scan(&active)

	return err == nil && active > 0
}

func (s *EconomyStore) GetDailyStreak(jid string) (int, string) {
	s.ensureUser(jid)

	var streak int
	var lastDate string

	err := s.db.QueryRow(
		`SELECT COALESCE(streak, 0), COALESCE(last_claim_date, '')
		FROM economy_daily_streaks
		WHERE jid = ?`,
		jid,
	).Scan(&streak, &lastDate)

	if err != nil {
		return 0, ""
	}

	return streak, lastDate
}

func (s *EconomyStore) RecordDailyClaim(jid string) error {
	s.ensureUser(jid)

	s.mu.Lock()
	defer s.mu.Unlock()

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	now := time.Now()
	today := now.Format("2006-01-02")
	yesterday := now.AddDate(0, 0, -1).Format("2006-01-02")

	var streak int
	var lastDate string

	err = tx.QueryRow(
		`SELECT COALESCE(streak, 0), COALESCE(last_claim_date, '')
		FROM economy_daily_streaks
		WHERE jid = ?`,
		jid,
	).Scan(&streak, &lastDate)

	if err == sql.ErrNoRows {
		streak = 0
		lastDate = ""
	} else if err != nil {
		return err
	}

	if lastDate != today {
		if lastDate == yesterday {
			streak++
		} else {
			streak = 1
		}
	}

	_, err = tx.Exec(
		`INSERT INTO economy_daily_streaks (jid, streak, last_claim_date)
		VALUES (?, ?, ?)
		ON CONFLICT(jid) DO UPDATE SET
			streak = excluded.streak,
			last_claim_date = excluded.last_claim_date`,
		jid,
		streak,
		today,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *EconomyStore) GetDailyMissions(jid string) ([]MissionProgress, error) {
	s.ensureUser(jid)

	s.mu.Lock()
	defer s.mu.Unlock()

	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	today := time.Now().Format("2006-01-02")

	for _, mission := range dailyMissionDefinitions {
		if _, err := tx.Exec(
			`INSERT OR IGNORE INTO economy_mission_progress
			(jid, mission_date, mission_id, progress, completed)
			VALUES (?, ?, ?, 0, 0)`,
			jid,
			today,
			mission.ID,
		); err != nil {
			return nil, err
		}
	}

	rows, err := tx.Query(
		`SELECT mission_id, progress, completed
		FROM economy_mission_progress
		WHERE jid = ? AND mission_date = ?`,
		jid,
		today,
	)
	if err != nil {
		return nil, err
	}

	progressMap := make(map[string]struct {
		Progress  int
		Completed bool
	})

	for rows.Next() {
		var id string
		var progress, completed int

		if err := rows.Scan(&id, &progress, &completed); err != nil {
			rows.Close()
			return nil, err
		}

		progressMap[id] = struct {
			Progress  int
			Completed bool
		}{
			Progress:  progress,
			Completed: completed != 0,
		}
	}

	if err := rows.Err(); err != nil {
		rows.Close()
		return nil, err
	}

	rows.Close()

	result := make([]MissionProgress, 0, len(dailyMissionDefinitions))

	for _, definition := range dailyMissionDefinitions {
		state := progressMap[definition.ID]

		definition.Progress = state.Progress
		definition.Completed = state.Completed

		result = append(result, definition)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return result, nil
}

func missionDate() string {
	return time.Now().Format("2006-01-02")
}

func recordMissionProgressTx(tx *sql.Tx, jid, action string, earned int) (int, error) {
	today := missionDate()

	for _, mission := range dailyMissionDefinitions {
		if _, err := tx.Exec(
			`INSERT OR IGNORE INTO economy_mission_progress
			(jid, mission_date, mission_id, progress, completed)
			VALUES (?, ?, ?, 0, 0)`,
			jid,
			today,
			mission.ID,
		); err != nil {
			return 0, err
		}
	}

	increments := make(map[string]int)

	isNewActivity := action == "hunt" ||
		action == "dig" ||
		action == "talar" ||
		action == "crime"

	if isNewActivity {
		increments["activities"] = 1

		if action == "hunt" || action == "dig" || action == "talar" {
			increments["gather"] = 1
		}

		if earned > 0 {
			increments["earn"] = earned
		}
	}

	bonus := 0

	for id, increment := range increments {
		var progress, completed int
		var definition *MissionProgress

		for i := range dailyMissionDefinitions {
			if dailyMissionDefinitions[i].ID == id {
				definition = &dailyMissionDefinitions[i]
				break
			}
		}

		if definition == nil {
			continue
		}

		if err := tx.QueryRow(
			`SELECT progress, completed
			FROM economy_mission_progress
			WHERE jid = ? AND mission_date = ? AND mission_id = ?`,
			jid,
			today,
			id,
		).Scan(&progress, &completed); err != nil {
			return bonus, err
		}

		if completed != 0 {
			continue
		}

		next := progress + increment
		if next > definition.Target {
			next = definition.Target
		}

		isCompleted := next >= definition.Target
		completedValue := 0

		if isCompleted {
			completedValue = 1
			bonus += definition.Reward

			if _, err := tx.Exec(
				`UPDATE users
				SET balance = COALESCE(balance, 0) + ?
				WHERE jid = ?`,
				definition.Reward,
				jid,
			); err != nil {
				return bonus, err
			}
		}

		if _, err := tx.Exec(
			`UPDATE economy_mission_progress
			SET progress = ?, completed = ?
			WHERE jid = ? AND mission_date = ? AND mission_id = ?`,
			next,
			completedValue,
			jid,
			today,
			id,
		); err != nil {
			return bonus, err
		}
	}

	return bonus, nil
}

func addXPToTx(tx *sql.Tx, jid string, amount int) (int, bool, error) {
	if amount <= 0 {
		var level int

		err := tx.QueryRow(
			`SELECT COALESCE(level, 1) FROM users WHERE jid = ?`,
			jid,
		).Scan(&level)

		return level, false, err
	}

	var xp, level int

	if err := tx.QueryRow(
		`SELECT COALESCE(xp, 0), COALESCE(level, 1)
		FROM users
		WHERE jid = ?`,
		jid,
	).Scan(&xp, &level); err != nil {
		return level, false, err
	}

	oldLevel := level
	xp += amount

	for level > 0 && xp >= level*100 {
		xp -= level * 100
		level++
	}

	_, err := tx.Exec(
		`UPDATE users SET xp = ?, level = ? WHERE jid = ?`,
		xp,
		level,
		jid,
	)

	return level, level > oldLevel, err
}

func xpBoostAvailableTx(tx *sql.Tx, jid string) bool {
	var active int

	if err := tx.QueryRow(
		`SELECT COALESCE(xp_boost, 0)
		FROM economy_effects
		WHERE jid = ?`,
		jid,
	).Scan(&active); err != nil {
		return false
	}

	return active > 0
}

func consumeXPBoostTx(tx *sql.Tx, jid string) error {
	_, err := tx.Exec(
		`UPDATE economy_effects SET xp_boost = 0 WHERE jid = ?`,
		jid,
	)

	return err
}

func rollActivity(action string) ActivityResult {
	result := ActivityResult{Action: action}

	switch action {
	case "work":
		result.Title, result.Emoji = "TRABAJO", "💼"
		result.Reward = rand.Intn(150) + 50
		result.XP = 15
		result.Detail = "Cumpliste tu jornada y cobraste el sueldo."

	case "adventure":
		result.Title, result.Emoji = "AVENTURA", "🗺️"

		outcomes := []weightedOutcome{
			{"Encontraste un cofre abandonado", 100, 250, 25},
			{"Ahuyentaste a unos bandidos y tomaste su botín", 150, 300, 20},
			{"Ayudaste a un mercader", 200, 400, 15},
			{"Exploraste una cueva", 250, 500, 12},
			{"Completaste una misión épica", 400, 800, 8},
			{"Hallaste un tesoro legendario", 600, 1200, 5},
			{"Te perdiste y regresaste con las manos vacías", 0, 0, 15},
		}

		chosen := pickWeightedOutcome(outcomes)
		result.Detail = chosen.text

		if chosen.max > 0 {
			result.Reward = rand.Intn(chosen.max-chosen.min+1) + chosen.min
		}

	case "mine":
		result.Title, result.Emoji = "MINAR", "⛏️"

		outcomes := []weightedOutcome{
			{"Encontraste piedra", 10, 30, 30},
			{"Encontraste carbón", 20, 50, 25},
			{"Encontraste hierro", 40, 80, 20},
			{"Encontraste oro", 80, 150, 15},
			{"Encontraste diamantes", 150, 300, 8},
			{"Encontraste una esmeralda", 200, 400, 2},
		}

		chosen := pickWeightedOutcome(outcomes)
		result.Detail = chosen.text
		result.Reward = rand.Intn(chosen.max-chosen.min+1) + chosen.min

	case "fish":
		result.Title, result.Emoji = "PESCA", "🎣"

		if rand.Intn(100) < 10 {
			result.Detail = "Se rompió tu caña. No atrapaste nada."
			break
		}

		outcomes := []weightedOutcome{
			{"Pescaste una sardina", 15, 35, 30},
			{"Pescaste una trucha", 30, 60, 25},
			{"Pescaste un salmón", 50, 100, 20},
			{"Pescaste un atún", 80, 150, 15},
			{"Pescaste un pez espada", 120, 250, 8},
			{"Pescaste un pez dorado", 200, 400, 2},
		}

		chosen := pickWeightedOutcome(outcomes)
		result.Detail = chosen.text
		result.Reward = rand.Intn(chosen.max-chosen.min+1) + chosen.min

	case "beg":
		result.Title, result.Emoji = "MENDIGAR", "🥺"
		result.Reward = rand.Intn(41) + 10

		phrases := []string{
			"Una abuelita te dio unas monedas.",
			"Encontraste cometas tirados en la vereda.",
			"Alguien te ayudó por lástima.",
			"Reciclaste botellas y las vendiste.",
		}

		result.Detail = phrases[rand.Intn(len(phrases))]

	case "slut":
		result.Title, result.Emoji = "SLUT", "💋"

		outcomes := []weightedOutcome{
			{"Un cliente generoso te dio propina", 50, 150, 30},
			{"Trabajaste toda la noche", 100, 250, 25},
			{"Un VIP te pagó extra", 200, 400, 20},
			{"Ganaste un concurso de baile", 300, 600, 15},
			{"Te contrataron para un evento exclusivo", 500, 1000, 8},
			{"La policía te multó", -100, -50, 2},
		}

		chosen := pickWeightedOutcome(outcomes)
		result.Detail = chosen.text
		result.Reward = rand.Intn(chosen.max-chosen.min+1) + chosen.min

	case "hunt":
		result.Title, result.Emoji = "CAZA", "🏹"

		outcomes := []weightedOutcome{
			{"Cazaste un conejo", 45, 90, 30},
			{"Cazaste un zorro", 70, 130, 25},
			{"Cazaste un jabalí", 100, 190, 20},
			{"Cazaste un ciervo", 140, 260, 15},
			{"Encontraste una presa rara", 220, 380, 10},
		}

		chosen := pickWeightedOutcome(outcomes)
		result.Detail = chosen.text
		result.Reward = rand.Intn(chosen.max-chosen.min+1) + chosen.min
		result.XP = 20

	case "dig":
		result.Title, result.Emoji = "EXCAVACIÓN", "🪏"

		outcomes := []weightedOutcome{
			{"Encontraste piedras antiguas", 20, 60, 30},
			{"Encontraste cobre", 40, 90, 25},
			{"Encontraste hierro", 80, 150, 20},
			{"Encontraste oro enterrado", 150, 280, 15},
			{"Encontraste un diamante", 300, 500, 10},
		}

		chosen := pickWeightedOutcome(outcomes)
		result.Detail = chosen.text
		result.Reward = rand.Intn(chosen.max-chosen.min+1) + chosen.min
		result.XP = 15

	case "talar":
		result.Title, result.Emoji = "TALA", "🪓"

		outcomes := []weightedOutcome{
			{"Cortaste madera común", 60, 120, 45},
			{"Conseguiste madera resistente", 100, 190, 35},
			{"Encontraste un árbol de madera rara", 180, 320, 20},
		}

		chosen := pickWeightedOutcome(outcomes)
		result.Detail = chosen.text
		result.Reward = rand.Intn(chosen.max-chosen.min+1) + chosen.min
		result.XP = 25

	case "crime":
		result.Title, result.Emoji = "CRIMEN", "🕵️"

		if rand.Intn(100) < 50 {
			result.Detail = "El golpe salió bien y escapaste con el botín."
			result.Reward = rand.Intn(1151) + 350
			result.XP = 50
		} else {
			result.Detail = "La policía te atrapó y te impuso una multa."
			result.Reward = -(rand.Intn(191) + 60)
		}
	}

	return result
}

type weightedOutcome struct {
	text             string
	min, max, weight int
}

func pickWeightedOutcome(outcomes []weightedOutcome) weightedOutcome {
	total := 0

	for _, item := range outcomes {
		total += item.weight
	}

	roll := rand.Intn(total)

	for _, item := range outcomes {
		roll -= item.weight
		if roll < 0 {
			return item
		}
	}

	return outcomes[len(outcomes)-1]
}

func runActivityTx(
	tx *sql.Tx,
	jid string,
	action string,
	now int64,
	xpBoost *bool,
) (ActivityResult, bool, error) {
	result := ActivityResult{Action: action}

	column, ok := activityColumns[action]
	cooldown, hasCooldown := activityCooldowns[action]

	if !ok || !hasCooldown {
		return result, false, fmt.Errorf("actividad no válida: %s", action)
	}

	var lastUse int64

	if err := tx.QueryRow(
		`SELECT COALESCE(`+column+`, 0) FROM users WHERE jid = ?`,
		jid,
	).Scan(&lastUse); err != nil {
		return result, false, err
	}

	if now-lastUse < cooldown {
		result.Skipped = true
		result.Remaining = cooldown - (now - lastUse)
		return result, false, nil
	}

	result = rollActivity(action)

	if result.Reward < 0 {
		var balance int

		if err := tx.QueryRow(
			`SELECT COALESCE(balance, 0) FROM users WHERE jid = ?`,
			jid,
		).Scan(&balance); err != nil {
			return result, false, err
		}

		if -result.Reward > balance {
			result.Reward = -balance
		}
	}

	if _, err := tx.Exec(
		`UPDATE users
		SET `+column+` = ?,
			balance = MAX(0, COALESCE(balance, 0) + ?)
		WHERE jid = ?`,
		now,
		result.Reward,
		jid,
	); err != nil {
		return result, false, err
	}

	if result.XP > 0 {
		if xpBoost != nil && *xpBoost {
			result.XP *= 2

			if err := consumeXPBoostTx(tx, jid); err != nil {
				return result, false, err
			}

			*xpBoost = false
		}

		level, leveled, err := addXPToTx(tx, jid, result.XP)
		if err != nil {
			return result, false, err
		}

		result.Level = level
		result.LeveledUp = leveled
	}

	missionBonus, err := recordMissionProgressTx(tx, jid, action, result.Reward)
	if err != nil {
		return result, false, err
	}

	result.MissionBonus = missionBonus

	return result, true, nil
}

func (s *EconomyStore) RunActivity(jid, action string) (ActivityResult, error) {
	if _, ok := activityColumns[action]; !ok {
		return ActivityResult{}, fmt.Errorf("actividad no reconocida")
	}

	s.ensureUser(jid)

	s.mu.Lock()
	defer s.mu.Unlock()

	tx, err := s.db.Begin()
	if err != nil {
		return ActivityResult{}, err
	}
	defer tx.Rollback()

	boost := xpBoostAvailableTx(tx, jid)

	result, executed, err := runActivityTx(
		tx,
		jid,
		action,
		time.Now().Unix(),
		&boost,
	)
	if err != nil {
		return ActivityResult{}, err
	}

	if !executed {
		return result, nil
	}

	if err := tx.Commit(); err != nil {
		return ActivityResult{}, err
	}

	return result, nil
}