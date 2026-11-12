package store

import (
	"fmt"
	"time"
)

var allowedCooldownActions = map[string]bool{
	"work": true, "daily": true, "rob": true, "beg": true,
	"mine": true, "adventure": true, "fish": true, "slut": true,
	"hunt": true, "dig": true, "talar": true, "crime": true, "allw": true,
}

func (s *EconomyStore) GetCooldown(jid, action string) int64 {
	if !allowedCooldownActions[action] {
		return 0
	}

	var lastUse int64
	query := `SELECT COALESCE(last_` + action + `, 0) FROM users WHERE jid = ?`
	_ = s.db.QueryRow(query, jid).Scan(&lastUse)

	return lastUse
}

func (s *EconomyStore) SetCooldown(jid, action string) {
	if !allowedCooldownActions[action] {
		return
	}

	s.ensureUser(jid)

	now := time.Now().Unix()
	query := `UPDATE users SET last_` + action + ` = ? WHERE jid = ?`
	_, _ = s.db.Exec(query, now, jid)
}

func (s *EconomyStore) CheckCooldown(jid, action string, cooldownSec int) (int64, bool) {
	lastUse := s.GetCooldown(jid, action)
	now := time.Now().Unix()
	elapsed := now - lastUse

	if elapsed < int64(cooldownSec) {
		return int64(cooldownSec) - elapsed, false
	}

	return 0, true
}

func (s *EconomyStore) AddColumnIfNotExists(columnName string) {
	if !allowedCooldownActions[columnNameToAction(columnName)] {
		return
	}

	rows, err := s.db.Query(`PRAGMA table_info(users)`)
	if err != nil {
		return
	}
	defer rows.Close()

	cols := make(map[string]bool)

	for rows.Next() {
		var cid, notnull, pk int
		var name, ctype string
		var dflt interface{}

		if err := rows.Scan(&cid, &name, &ctype, &notnull, &dflt, &pk); err != nil {
			return
		}

		cols[name] = true
	}

	if err := rows.Err(); err != nil {
		return
	}

	if err := rows.Close(); err != nil {
		return
	}

	if !cols[columnName] {
		_, _ = s.db.Exec(`ALTER TABLE users ADD COLUMN ` + columnName + ` INTEGER DEFAULT 0`)
	}
}

func columnNameToAction(column string) string {
	if len(column) > 5 && column[:5] == "last_" {
		return column[5:]
	}

	return ""
}

func (s *EconomyStore) InitCooldownColumns() {
	columns := []string{
		"last_mine",
		"last_adventure",
		"last_fish",
		"last_slut",
		"last_hunt",
		"last_dig",
		"last_talar",
		"last_crime",
		"last_allw",
	}

	for _, column := range columns {
		s.AddColumnIfNotExists(column)
	}

	if err := s.InitExtraTables(); err != nil {
		fmt.Println("❌ Error inicializando tablas extra de economía:", err)
	}
}