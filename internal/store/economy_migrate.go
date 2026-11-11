package store

import (
"database/sql"
"fmt"
)

func mustMigrateEconomy(db *sql.DB) {
db.Exec(`
CREATE TABLE IF NOT EXISTS users (
jid TEXT PRIMARY KEY,
balance INTEGER DEFAULT 0,
bank INTEGER DEFAULT 0,
last_work INTEGER DEFAULT 0,
last_daily INTEGER DEFAULT 0,
last_rob INTEGER DEFAULT 0,
last_beg INTEGER DEFAULT 0,
xp INTEGER DEFAULT 0,
level INTEGER DEFAULT 1
)`)

rows, err := db.Query(`PRAGMA table_info(users)`)
if err != nil {
fmt.Println("❌ Error leyendo columnas de economía:", err)
return
}
defer rows.Close()

cols := make(map[string]bool)
for rows.Next() {
var cid int
var name string
var ctype string
var notnull int
var dflt interface{}
var pk int
rows.Scan(&cid, &name, &ctype, &notnull, &dflt, &pk)
cols[name] = true
}

needed := map[string]string{
"balance":    "INTEGER DEFAULT 0",
"bank":       "INTEGER DEFAULT 0",
"last_work":  "INTEGER DEFAULT 0",
"last_daily": "INTEGER DEFAULT 0",
"last_rob":   "INTEGER DEFAULT 0",
"last_beg":   "INTEGER DEFAULT 0",
"xp":         "INTEGER DEFAULT 0",
"level":      "INTEGER DEFAULT 1",
}

for col, def := range needed {
if !cols[col] {
_, err := db.Exec(`ALTER TABLE users ADD COLUMN ` + col + ` ` + def)
if err != nil {
fmt.Println("❌ Error agregando columna", col+":", err)
} else {
fmt.Println("✅ Columna agregada:", col)
}
}
}

db.Exec(`
UPDATE users SET
balance = COALESCE(balance, 0),
bank = COALESCE(bank, 0),
last_work = COALESCE(last_work, 0),
last_daily = COALESCE(last_daily, 0),
last_rob = COALESCE(last_rob, 0),
last_beg = COALESCE(last_beg, 0),
xp = COALESCE(xp, 0),
level = COALESCE(level, 1)
`)

fmt.Println("🧩 Migración de economía lista")
}
