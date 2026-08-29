package repository

// imports you'll need:
import (
	"errors"
	"math"
	"math/big"
	"nailzbydardo/internal/db/sqlc"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// ---- UUID ----
func pgUUIDToString(u pgtype.UUID) (string, error) {
	if !u.Valid {
		return "", errors.New("invalid uuid")
	}
	return uuid.UUID(u.Bytes).String(), nil
}

func stringToPgUUID(s string) (pgtype.UUID, error) {
	parsedUUID, err := uuid.Parse(s)

	if err != nil {
		return pgtype.UUID{}, errors.New("error parsing string to uuid")
	}
	return pgtype.UUID{Bytes: parsedUUID, Valid: true}, nil
}

// ---- Timestamptz ----
func pgTimestamptzToTime(t pgtype.Timestamptz) time.Time {
	return t.Time
}

func pgTimestamptzToTimePtr(t pgtype.Timestamptz) *time.Time {
	if !t.Valid {
		return nil
	}
	return &t.Time
}

func timeToPgTimestamptz(t time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{
		Time:  t,
		Valid: true,
	}
}

func timePtrToPgTimestamptz(t *time.Time) pgtype.Timestamptz {
	if t == nil {
		return pgtype.Timestamptz{
			Valid: false,
		}
	}
	return timeToPgTimestamptz(*t)
}

// ---- Numeric (money, stored as cents) ----
func pgNumericToCents(n pgtype.Numeric) int64 {
    if !n.Valid || n.Int == nil {
        return 0
    }

    // value = n.Int * 10^n.Exp
    // cents  = value * 100 = n.Int * 10^(n.Exp + 2)
    rat := new(big.Rat).SetInt(n.Int)
    exp := n.Exp + 2

    pow := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(abs32(exp))), nil)
    powRat := new(big.Rat).SetInt(pow)

    if exp >= 0 {
        rat.Mul(rat, powRat)
    } else {
        rat.Quo(rat, powRat)
    }

    f, _ := rat.Float64()
    return int64(math.Round(f))
}

func pgNumericToCentsPtr(n pgtype.Numeric) *int64 {
    if !n.Valid {
        return nil
    }
    cents := pgNumericToCents(n)
    return &cents
}

func centsToPgNumeric(cents int64) pgtype.Numeric {
    return pgtype.Numeric{
        Int:   big.NewInt(cents),
        Exp:   -2,
        Valid: true,
    }
}

func centsPtrToPgNumeric(cents *int64) pgtype.Numeric {
    if cents == nil {
        return pgtype.Numeric{Valid: false}
    }
    return centsToPgNumeric(*cents)
}

func abs32(n int32) int32 {
    if n < 0 {
        return -n
    }
    return n
}

// ---- Text ----
func pgTextToStringPtr(t pgtype.Text) *string {
	if !t.Valid {
		return nil
	}
	return &t.String
}

func stringPtrToPgText(s *string) pgtype.Text {
	if s == nil {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: *s, Valid: true}
}

// ---- Date ----
func pgDateToTime(d pgtype.Date) time.Time {
	return d.Time
}

func pgDateToTimePtr(d pgtype.Date) *time.Time {
	if !d.Valid {
		return nil
	}
	return &d.Time
}

func timeToPgDate(t time.Time) pgtype.Date {
	return pgtype.Date{
		Time:  t,
		Valid: true,
	}
}

func timePtrToPgDate(t *time.Time) pgtype.Date {
	if t == nil {
		return pgtype.Date{
			Valid: false,
		}
	}
	return timeToPgDate(*t)
}

func nullPaymentMethodToPtr(p sqlc.NullPaymentMethod) *sqlc.PaymentMethod {
    if !p.Valid {
        return nil
    }
    return &p.PaymentMethod
}

func ptrToNullPaymentMethod(p *sqlc.PaymentMethod) sqlc.NullPaymentMethod {
    if p == nil {
        return sqlc.NullPaymentMethod{Valid: false}
    }
    return sqlc.NullPaymentMethod{PaymentMethod: *p, Valid: true}
}

// Converts a pgtype.Numeric into a plain int64 (no ×100 scaling) —
// use for numeric columns representing a plain integer value, like
// a percentage, rather than money.
func pgNumericToInt(n pgtype.Numeric) int64 {
    if !n.Valid || n.Int == nil {
        return 0
    }
    rat := new(big.Rat).SetInt(n.Int)
    if n.Exp >= 0 {
        pow := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(n.Exp)), nil)
        rat.Mul(rat, new(big.Rat).SetInt(pow))
    } else {
        pow := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(-n.Exp)), nil)
        rat.Quo(rat, new(big.Rat).SetInt(pow))
    }
    f, _ := rat.Float64()
    return int64(math.Round(f))
}

// Converts a plain int64 into pgtype.Numeric (no ×100 scaling) —
// the inverse of pgNumericToInt.
func intToPgNumeric(i int64) pgtype.Numeric {
    return pgtype.Numeric{Int: big.NewInt(i), Exp: 0, Valid: true}
}