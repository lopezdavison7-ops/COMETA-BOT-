package store

import (
	"errors"
	"fmt"
	"strings"
)

const MaxBetAmount = 10_000_000

var (
	ErrInvalidBetAmount = errors.New("la apuesta debe estar entre 10 y 10.000.000 cometas")
	ErrInsufficientBetBalance = errors.New("no tienes suficientes cometas en la billetera")
)

// SettleBet descuenta la apuesta y acredita el retorno en una sola operación SQL.
// payout es el retorno total, incluida la apuesta original cuando corresponde.
func (s *EconomyStore) SettleBet(jid string, wager, payout int) error {
	if s == nil || s.db == nil {
		return errors.New("la economía no está inicializada")
	}

	jid = strings.TrimSpace(jid)
	if jid == "" || wager < 10 || wager > MaxBetAmount || payout < 0 {
		return ErrInvalidBetAmount
	}

	s.ensureUser(jid)
	s.mu.Lock()
	defer s.mu.Unlock()

	result, err := s.db.Exec(
		`UPDATE users
		 SET balance = COALESCE(balance, 0) - ? + ?
		 WHERE jid = ? AND COALESCE(balance, 0) >= ?`,
		wager, payout, jid, wager,
	)
	if err != nil {
		return fmt.Errorf("no se pudo liquidar la apuesta: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("no se pudo comprobar el saldo: %w", err)
	}
	if rows == 0 {
		return ErrInsufficientBetBalance
	}

	return nil
}