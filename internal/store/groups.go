package store

import (
"sync"
)

var (
groupsMu sync.Mutex
)

func (s *EconomyStore) InitGroupsTable() {
s.db.Exec(`
CREATE TABLE IF NOT EXISTS groups (
jid TEXT PRIMARY KEY,
welcome INTEGER DEFAULT 0,
goodbye INTEGER DEFAULT 0,
welcome_msg TEXT DEFAULT '',
goodbye_msg TEXT DEFAULT '',
antibot INTEGER DEFAULT 0,
antilink INTEGER DEFAULT 0,
antifake INTEGER DEFAULT 0
)`)
}

func (s *EconomyStore) ensureGroup(jid string) {
s.db.Exec(`INSERT OR IGNORE INTO groups (jid) VALUES (?)`, jid)
}

func (s *EconomyStore) SetGroupSetting(jid, setting string, value int) {
groupsMu.Lock()
defer groupsMu.Unlock()
s.ensureGroup(jid)
s.db.Exec(`UPDATE groups SET `+setting+` = ? WHERE jid = ?`, value, jid)
}

func (s *EconomyStore) GetGroupSetting(jid, setting string) int {
s.ensureGroup(jid)
var val int
s.db.QueryRow(`SELECT COALESCE(`+setting+`, 0) FROM groups WHERE jid = ?`, jid).Scan(&val)
return val
}

func (s *EconomyStore) SetGroupMessage(jid, field, msg string) {
groupsMu.Lock()
defer groupsMu.Unlock()
s.ensureGroup(jid)
s.db.Exec(`UPDATE groups SET `+field+` = ? WHERE jid = ?`, msg, jid)
}

func (s *EconomyStore) GetGroupMessage(jid, field string) string {
s.ensureGroup(jid)
var val string
s.db.QueryRow(`SELECT COALESCE(`+field+`, '') FROM groups WHERE jid = ?`, jid).Scan(&val)
return val
}
