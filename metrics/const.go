package metrics

const (
	namespace = "hubble"

	// Subsystems
	apiSubsystem    = "api"
	rollupSubsystem = "rollup"
)

// API metrics
const (
	// Transaction labels
	TransferLabel = "transfer"
	C2TLabel      = "create2transfer"

	// Transaction statuses
	AcceptedTxStatus = "accepted"
	RejectedTxStatus = "rejected"

	// Batch labels
	TransferBatchLabel = "transfer"
	C2TBatchLabel      = "create2transfer"
	DepositBatchLabel  = "deposit"
)
