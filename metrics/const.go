package metrics

const (
	namespace = "hubble"

	// Subsystems
	apiSubsystem = "api"
)

// API metrics
const (
	// Transaction labels
	TransferTxLabel = "transfer"
	C2TTxLabel      = "create2transfer"

	// Transaction statuses
	AcceptedTxStatus = "accepted"
	RejectedTxStatus = "rejected"
)
