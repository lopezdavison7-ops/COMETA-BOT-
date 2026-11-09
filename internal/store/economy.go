package store

import (
"context"
"database/sql"
"fmt"
"time"

_ "modernc.org/sqlite"
)

var GlobalEconomy *EconomyStore

type EconomyStore struct {
db *sql.DB
}

func InitEconomy(dbPath string) error {
db, err := sql.Open("sqlite", fmt.Sprintf("file:%s?_pragma=busy_timeout(10000)&_pragma=journal_mode(WAL)&_pragma=synchronous(NORMAL)&_pragma=foreign_keys(1)", dbPath))
if err != nil {
return err
}

_, err = db.Exec(`
CREATE TABLE IF NOT EXISTS economy (
jid TEXT PRIMARY KEY,
balance INTEGER DEFAULT 0,
last_work INTEGER DEFAULT 0
);
CREATE INDEX IF NOT EXISTS idx_balance ON economy(balance DESC);
`)
if err != nil {
return err
}

GlobalEconomy = &EconomyStore{db: db}
return nil
}

func (s *EconomyStore) Work(jid string, amount int, cooldownSecs int) (string, error) {
ctx := context.Background()
now := time.Now().Unix()

var lastWork int64
var currentBal int
err := s.db.QueryRowContext(ctx, "SELECT balance, last_work FROM economy WHERE jid = ?", jid).Scan(&currentBal, &lastWork)

if err == sql.ErrNoRows {
_, err = s.db.ExecContext(ctx, "INSERT INTO economy (jid, balance, last_work) VALUES (?, ?, ?)", jid, amount, now)
if err != nil {
return "", err
}
return fmt.Sprintf("¡Bienvenido al sistema! Trabajaste y ganaste *%d cometas* ☄️", amount), nil
}
if err != nil {
return "", err
}

if now-lastWork < int64(cooldownSecs) {
remaining := int64(cooldownSecs) - (now - lastWork)
mins := remaining / 60
secs := remaining % 60
return fmt.Sprintf("⏳ ¡Tranquilo! Debes esperar *%dm %ds* para volver a trabajar.", mins, secs), nil
}

newBal := currentBal + amount
_, err = s.db.ExecContext(ctx, "UPDATE economy SET balance = ?, last_work = ? WHERE jid = ?", newBal, now, jid)
if err != nil {
return "", err
}

return fmt.Sprintf("💼 Trabajaste duro y ganaste *%d cometas* ☄️\n💰 Balance actual: *%d*", amount, newBal), nil
}

func (s *EconomyStore) GetBalance(jid string) (int, error) {
ctx := context.Background()
var balance int
err := s.db.QueryRowContext(ctx, "SELECT balance FROM economy WHERE jid = ?", jid).Scan(&balance)
if err == sql.ErrNoRows {
_, _ = s.db.ExecContext(ctx, "INSERT INTO economy (jid, balance) VALUES (?, 0)", jid)
return 0, nil
}
return balance, err
}

type TopUser struct {
JID     string
Balance int
Rank    int
}

func (s *EconomyStore) GetTopUsers(limit int) ([]TopUser, error) {
ctx := context.Background()
rows, err := s.db.QueryContext(ctx, "SELECT jid, balance FROM economy ORDER BY balance DESC LIMIT ?", limit)
if err != nil {
return nil, err
}
defer rows.Close()

var tops []TopUser
rank := 1
for rows.Next() {
var jid string
var balance int
if err := rows.Scan(&jid, &balance); err != nil {
continue
}
tops = append(tops, TopUser{JID: jid, Balance: balance, Rank: rank})
rank++
}
return tops, nil
}
