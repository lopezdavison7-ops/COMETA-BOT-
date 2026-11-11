
package store

import (
	"database/sql"
	"fmt"
)

func mustMigrateEconomy(db *sql.DB) {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			jid TEXT PRIMARY KEY,
			balance INTEGER DEFAULT 0,
			bank INTEGER DEFAULT 0,
			last_work INTEGER DEFAULT 0,
			last_daily INTEGER DEFAULT 0,
			last_weekly INTEGER DEFAULT 0,
			last_monthly INTEGER DEFAULT 0,
			last_rob INTEGER DEFAULT 0,
			last_beg INTEGER DEFAULT 0,
			xp INTEGER DEFAULT 0,
			level INTEGER DEFAULT 1
		)
	`)
	if err != nil {
		fmt.Println("❌ Error creando tabla de economía:", err)
		return
	}

	rows, err := db.Query(`PRAGMA table_info(users)`)
	if err != nil {
		fmt.Println("❌ Error leyendo columnas de economía:", err)
		return
	}

	cols := make(map[string]bool)

	for rows.Next() {
		var cid int
		var name, ctype string
		var notnull int
		var dflt interface{}
		var pk int

		if err := rows.Scan(&cid, &name, &ctype, &notnull, &dflt, &pk); err != nil {
			fmt.Println("❌ Error leyendo columna de economía:", err)
			rows.Close()
			return
		}

		cols[name] = true
	}

	if err := rows.Err(); err != nil {
		fmt.Println("❌ Error recorriendo columnas de economía:", err)
		rows.Close()
		return
	}

	rows.Close()

	needed := map[string]string{
		"balance":      "INTEGER DEFAULT 0",
		"bank":         "INTEGER DEFAULT 0",
		"last_work":    "INTEGER DEFAULT 0",
		"last_daily":   "INTEGER DEFAULT 0",
		"last_weekly":  "INTEGER DEFAULT 0",
		"last_monthly": "INTEGER DEFAULT 0",
		"last_rob":     "INTEGER DEFAULT 0",
		"last_beg":     "INTEGER DEFAULT 0",
		"xp":           "INTEGER DEFAULT 0",
		"level":        "INTEGER DEFAULT 1",
	}

	for col, def := range needed {
		if !cols[col] {
			if _, err := db.Exec(`ALTER TABLE users ADD COLUMN ` + col + ` ` + def); err != nil {
				fmt.Println("❌ Error agregando columna", col+":", err)
			} else {
				fmt.Println("✅ Columna agregada:", col)
			}
		}
	}

	_, err = db.Exec(`
		UPDATE users SET
			balance = COALESCE(balance, 0),
			bank = COALESCE(bank, 0),
			last_work = COALESCE(last_work, 0),
			last_daily = COALESCE(last_daily, 0),
			last_weekly = COALESCE(last_weekly, 0),
			last_monthly = COALESCE(last_monthly, 0),
			last_rob = COALESCE(last_rob, 0),
			last_beg = COALESCE(last_beg, 0),
			xp = COALESCE(xp, 0),
			level = COALESCE(level, 1)
	`)

	if err != nil {
		fmt.Println("❌ Error actualizando datos de economía:", err)
		return
	}

	fmt.Println("🧩 Migración de economía lista")
}
