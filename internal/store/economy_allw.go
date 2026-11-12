package store

import (
	"fmt"
	"time"
)

func (s *EconomyStore) AllWork(jid string) ([]ActivityResult, error) {
	s.ensureUser(jid)

	s.mu.Lock()
	defer s.mu.Unlock()

	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	now := time.Now().Unix()

	var lastAllw int64
	if err := tx.QueryRow(
		`SELECT COALESCE(last_allw, 0) FROM users WHERE jid = ?`,
		jid,
	).Scan(&lastAllw); err != nil {
		return nil, err
	}

	if now-lastAllw < 3600 {
		return nil, fmt.Errorf("cooldown:%d", 3600-(now-lastAllw))
	}

	boost := xpBoostAvailableTx(tx, jid)

	actions := []string{
		"work",
		"adventure",
		"mine",
		"fish",
		"beg",
		"slut",
		"hunt",
		"dig",
		"talar",
		"crime",
	}

	results := make([]ActivityResult, 0, len(actions))
	executed := 0

	for _, action := range actions {
		result, didExecute, err := runActivityTx(tx, jid, action, now, &boost)
		if err != nil {
			return nil, err
		}

		if didExecute {
			executed++
		}

		results = append(results, result)
	}

	if executed == 0 {
		return results, nil
	}

	if _, err := tx.Exec(
		`UPDATE users SET last_allw = ? WHERE jid = ?`,
		now,
		jid,
	); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return results, nil
}