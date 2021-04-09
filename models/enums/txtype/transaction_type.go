package txtype

type TransactionType uint8

const (
	Transfer        TransactionType = 1
	Create2Transfer TransactionType = 3
	MassMigration   TransactionType = 5
)

var TransactionTypes = map[TransactionType]string{
	Transfer:        "TRANSFER",
	Create2Transfer: "CREATE2TRANSFER",
	MassMigration:   "MASS_MIGRATION",
}

func (s TransactionType) Ref() *TransactionType {
	return &s
}

func (s TransactionType) String() string {
	msg, exists := TransactionTypes[s]
	if !exists {
		return "UNKNOWN"
	}
	return msg
}
