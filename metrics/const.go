package metrics

const (
	namespace = "hubble"

	// Subsystems
	apiSubsystem     = "api"
	rollupSubsystem  = "rollup"
	syncingSubsystem = "syncing"
)

// API metrics
const (
	// Transaction labels
	TransferTxLabel = "transfer"
	C2TTxLabel      = "create2transfer"

	// Transaction statuses
	AcceptedTxStatus = "accepted"
	RejectedTxStatus = "rejected"

	// Batch labels
	TransferBatchLabel = "transfer"
	C2TBatchLabel      = "create2transfer"
	DepositBatchLabel  = "deposit"
)

// Syncing metrics
const (
	SyncAccountsMethod = "sync_accounts"
)
