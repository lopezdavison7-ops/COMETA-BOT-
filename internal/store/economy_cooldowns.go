
package store

import "time"

func (s *EconomyStore) GetCooldown(jid, action string) int64 {
	var lastUse int64

	query := `SELECT COALESCE(last_` + action + `, 0) FROM users WHERE jid = ?`
	s.db.QueryRow(query, jid).Scan(&lastUse)

	return lastUse
}

func (s *EconomyStore) SetCooldown(jid, action string) {
	now := time.Now().Unix()

	query := `UPDATE users SET last_` + action + ` = ? WHERE jid = ?`
	s.db.Exec(query, now, jid)
}

func (s *EconomyStore) CheckCooldown(jid, action string, cooldownSec int) (int64, bool) {
	lastUse := s.GetCooldown(jid, action)
	now := time.Now().Unix()

	if now-lastUse < int64(cooldownSec) {
		return int64(cooldownSec) - (now - lastUse), false
	}

	return 0, true
}

func (s *EconomyStore) AddColumnIfNotExists(columnName string) {
	rows, err := s.db.Query(`PRAGMA table_info(users)`)
	if err != nil {
		return
	}
	defer rows.Close()

	cols := make(map[string]bool)

	for rows.Next() {
		var cid int
		var name, ctype string
		var notnull int
		var dflt interface{}
		var pk int

		if rows.Scan(&cid, &name, &ctype, &notnull, &dflt, &pk) == nil {
			cols[name] = true
		}
	}

	if !cols[columnName] {
		s.db.Exec(`ALTER TABLE users ADD COLUMN ` + columnName + ` INTEGER DEFAULT 0`)
	}
}

func (s *EconomyStore) InitCooldownColumns() {
	s.AddColumnIfNotExists("last_mine")
	s.AddColumnIfNotExists("last_adventure")
	s.AddColumnIfNotExists("last_fish")
	s.AddColumnIfNotExists("last_slut")
	s.AddColumnIfNotExists("last_allw")
}
