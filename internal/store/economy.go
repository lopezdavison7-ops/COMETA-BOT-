package store

import (
"database/sql"
"fmt"
"math/rand"
"sort"
"strconv"
"sync"
"time"

_ "modernc.org/sqlite"
)

type EconomyStore struct {
db *sql.DB
mu sync.Mutex
}

var GlobalEconomy *EconomyStore

func InitEconomy(dbPath string) error {
db, err := sql.Open("sqlite", dbPath)
if err != nil {
return err
}

db.Exec(`PRAGMA journal_mode=WAL;`)
db.Exec(`PRAGMA synchronous=NORMAL;`)

GlobalEconomy = &EconomyStore{db: db}
mustMigrateEconomy(db)

schema := `
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
);`

if _, err := db.Exec(schema); err != nil {
return err
}

db.Exec(`ALTER TABLE users ADD COLUMN bank INTEGER DEFAULT 0`)

db.Exec(`UPDATE users SET 
balance = COALESCE(balance, 0), 
bank = COALESCE(bank, 0), 
xp = COALESCE(xp, 0), 
level = COALESCE(level, 1),
last_work = COALESCE(last_work, 0),
last_daily = COALESCE(last_daily, 0),
last_rob = COALESCE(last_rob, 0),
last_beg = COALESCE(last_beg, 0)`)

GlobalEconomy.InitCooldownColumns()
GlobalEconomy.InitGroupsTable()

return nil
}

func (s *EconomyStore) ensureUser(jid string) {
s.db.Exec(`INSERT OR IGNORE INTO users (jid, balance, bank, xp, level) VALUES (?, 0, 0, 0, 1)`, jid)
s.db.Exec(`UPDATE users SET balance = COALESCE(balance, 0), bank = COALESCE(bank, 0), xp = COALESCE(xp, 0), level = COALESCE(level, 1) WHERE jid = ?`, jid)
}

func (s *EconomyStore) GetBalance(jid string) int {
s.ensureUser(jid)
var balance int
s.db.QueryRow(`SELECT balance FROM users WHERE jid = ?`, jid).Scan(&balance)
return balance
}

func (s *EconomyStore) GetBank(jid string) int {
s.ensureUser(jid)
var bank int
s.db.QueryRow(`SELECT COALESCE(bank, 0) FROM users WHERE jid = ?`, jid).Scan(&bank)
return bank
}

func (s *EconomyStore) GetProfile(jid string) (int, int) {
s.ensureUser(jid)
var level, xp int
s.db.QueryRow(`SELECT COALESCE(level, 1), COALESCE(xp, 0) FROM users WHERE jid = ?`, jid).Scan(&level, &xp)
return level, xp
}

func (s *EconomyStore) AddBalance(jid string, amount int) {
s.ensureUser(jid)
s.mu.Lock()
defer s.mu.Unlock()
s.db.Exec(`UPDATE users SET balance = balance + ? WHERE jid = ?`, amount, jid)
}

func (s *EconomyStore) AddBank(jid string, amount int) {
s.ensureUser(jid)
s.mu.Lock()
defer s.mu.Unlock()
s.db.Exec(`UPDATE users SET bank = COALESCE(bank, 0) + ? WHERE jid = ?`, amount, jid)
}

func (s *EconomyStore) Work(jid string, amount int, cooldown int) (string, error) {
s.ensureUser(jid)
s.mu.Lock()
defer s.mu.Unlock()

var lw int64
s.db.QueryRow(`SELECT COALESCE(last_work, 0) FROM users WHERE jid = ?`, jid).Scan(&lw)

now := time.Now().Unix()
if now-lw < int64(cooldown) {
r := int64(cooldown) - (now - lw)
return fmt.Sprintf("Tranquilo chacho, descansá. Podés volver a trabajar en %d minutos y %d segundos.", r/60, r%60), fmt.Errorf("cooldown")
}

s.db.Exec(`UPDATE users SET balance = balance + ?, last_work = ? WHERE jid = ?`, amount, now, jid)

var cx, cl int
s.db.QueryRow(`SELECT COALESCE(xp, 0), COALESCE(level, 1) FROM users WHERE jid = ?`, jid).Scan(&cx, &cl)

nxp := cx + 15
need := cl * 100
up := false
if nxp >= need {
cl++
nxp -= need
up = true
}
s.db.Exec(`UPDATE users SET xp = ?, level = ? WHERE jid = ?`, nxp, cl, jid)

msg := fmt.Sprintf("¡Trabajaste duro y ganaste *%d* cometas! 💰\n+15 XP", amount)
if up {
msg += "\n🎉 ¡Subiste al nivel " + strconv.Itoa(cl) + "!"
}
return msg, nil
}

func (s *EconomyStore) Daily(jid string) (string, error) {
s.ensureUser(jid)
s.mu.Lock()
defer s.mu.Unlock()

var ld int64
s.db.QueryRow(`SELECT COALESCE(last_daily, 0) FROM users WHERE jid = ?`, jid).Scan(&ld)

now := time.Now().Unix()
if now-ld < 86400 {
r := 86400 - (now - ld)
return fmt.Sprintf("Ya reclamaste tu premio diario. Volvé en %d horas y %d minutos.", r/3600, (r%3600)/60), fmt.Errorf("cooldown")
}

amount := rand.Intn(501) + 500
s.db.Exec(`UPDATE users SET balance = balance + ?, last_daily = ? WHERE jid = ?`, amount, now, jid)

var cx, cl int
s.db.QueryRow(`SELECT COALESCE(xp, 0), COALESCE(level, 1) FROM users WHERE jid = ?`, jid).Scan(&cx, &cl)
nxp := cx + 50
need := cl * 100
up := false
if nxp >= need {
cl++
nxp -= need
up = true
}
s.db.Exec(`UPDATE users SET xp = ?, level = ? WHERE jid = ?`, nxp, cl, jid)

msg := fmt.Sprintf("¡Reclamaste tu recompensa diaria de *%d* cometas! 🎁\n+50 XP", amount)
if up {
msg += "\n🎉 ¡Subiste al nivel " + strconv.Itoa(cl) + "!"
}
return msg, nil
}

func (s *EconomyStore) Rob(thiefJID, victimJID string) (string, error) {
s.ensureUser(thiefJID)
s.ensureUser(victimJID)
s.mu.Lock()
defer s.mu.Unlock()

var lr int64
s.db.QueryRow(`SELECT COALESCE(last_rob, 0) FROM users WHERE jid = ?`, thiefJID).Scan(&lr)

now := time.Now().Unix()
if now-lr < 3600 {
r := 3600 - (now - lr)
return fmt.Sprintf("La policía te está vigilando. Intentá robar de nuevo en %d minutos.", r/60), fmt.Errorf("cooldown")
}

var vb int
s.db.QueryRow(`SELECT COALESCE(balance, 0) FROM users WHERE jid = ?`, victimJID).Scan(&vb)
if vb < 100 {
return "Ese usuario es más pobre que tú, chacho. No vale la pena robarle. 😂", fmt.Errorf("poor")
}

if rand.Intn(100) < 60 {
stolen := rand.Intn(vb/10) + 50
if stolen > vb {
stolen = vb
}
s.db.Exec(`UPDATE users SET balance = MAX(0, balance - ?), last_rob = ? WHERE jid = ?`, stolen, now, victimJID)
s.db.Exec(`UPDATE users SET balance = balance + ? WHERE jid = ?`, stolen, thiefJID)

var cx, cl int
s.db.QueryRow(`SELECT COALESCE(xp, 0), COALESCE(level, 1) FROM users WHERE jid = ?`, thiefJID).Scan(&cx, &cl)
nxp := cx + 30
need := cl * 100
up := false
if nxp >= need {
cl++
nxp -= need
up = true
}
s.db.Exec(`UPDATE users SET xp = ?, level = ? WHERE jid = ?`, nxp, cl, thiefJID)

msg := fmt.Sprintf("¡Robaste *%d* cometas exitosamente! 🥷💰\n+30 XP", stolen)
if up {
msg += "\n🎉 ¡Subiste al nivel " + strconv.Itoa(cl) + "!"
}
return msg, nil
}

s.db.Exec(`UPDATE users SET balance = MAX(0, balance - 100), last_rob = ? WHERE jid = ?`, now, thiefJID)
return "¡Te atraparon! La policía te confiscó *100* cometas por intento de robo. 👮‍♂️🚫", fmt.Errorf("failed")
}

func (s *EconomyStore) Pay(senderJID, receiverJID string, amount int) (string, error) {
if amount <= 0 {
return "El monto debe ser mayor a 0.", fmt.Errorf("invalid")
}

s.ensureUser(senderJID)
s.ensureUser(receiverJID)
s.mu.Lock()
defer s.mu.Unlock()

var b int
s.db.QueryRow(`SELECT COALESCE(balance, 0) FROM users WHERE jid = ?`, senderJID).Scan(&b)
if b < amount {
return "No tenés suficientes cometas para hacer esa transferencia. 😅", fmt.Errorf("insufficient")
}

s.db.Exec(`UPDATE users SET balance = balance - ? WHERE jid = ?`, amount, senderJID)
s.db.Exec(`UPDATE users SET balance = balance + ? WHERE jid = ?`, amount, receiverJID)

return fmt.Sprintf("✅ Transferiste *%d* cometas exitosamente.", amount), nil
}

func (s *EconomyStore) Beg(jid string) (string, error) {
s.ensureUser(jid)
s.mu.Lock()
defer s.mu.Unlock()

var lb int64
s.db.QueryRow(`SELECT COALESCE(last_beg, 0) FROM users WHERE jid = ?`, jid).Scan(&lb)

now := time.Now().Unix()
if now-lb < 300 {
r := 300 - (now - lb)
return fmt.Sprintf("La gente ya te dio hoy. Volvé a mendigar en %d minutos y %d segundos.", r/60, r%60), fmt.Errorf("cooldown")
}

amount := rand.Intn(41) + 10
s.db.Exec(`UPDATE users SET balance = balance + ?, last_beg = ? WHERE jid = ?`, amount, now, jid)

msgs := []string{
"Una abuelita te dio *%d* cometas por ser buen chico. 👵",
"Encontraste *%d* cometas tirados en la vereda. 🪙",
"Alguien te lanzó *%d* cometas por lástima. 🥺",
"Vendiste una botella y ganaste *%d* cometas. ♻️",
}
return fmt.Sprintf(msgs[rand.Intn(len(msgs))], amount), nil
}

type TopUser struct {
JID     string
Balance int
Level   int
}

func (s *EconomyStore) GetTop(limit int) []TopUser {
var t []TopUser
rows, _ := s.db.Query(`SELECT jid, COALESCE(balance, 0), COALESCE(level, 1) FROM users WHERE COALESCE(balance, 0) > 0 ORDER BY balance DESC LIMIT ?`, limit)
defer rows.Close()
for rows.Next() {
var u TopUser
rows.Scan(&u.JID, &u.Balance, &u.Level)
t = append(t, u)
}
sort.Slice(t, func(i, j int) bool { return t[i].Balance > t[j].Balance })
return t
}

func (s *EconomyStore) GetLevelTop(limit int) []TopUser {
var t []TopUser
rows, _ := s.db.Query(`SELECT jid, COALESCE(level, 1) FROM users ORDER BY level DESC, xp DESC LIMIT ?`, limit)
defer rows.Close()
for rows.Next() {
var u TopUser
rows.Scan(&u.JID, &u.Level)
t = append(t, u)
}
return t
}
