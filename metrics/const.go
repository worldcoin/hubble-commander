package metrics

const (
	namespace = "hubble"

	// Subsystems
	apiSubsystem        = "api"
	rollupSubsystem     = "rollup"
	syncingSubsystem    = "syncing"
	blockchainSubsystem = "blockchain"
)

// API metrics
const (
	// Transaction labels
	TransferTxLabel = "transfer"
	C2TTxLabel      = "create2transfer"
	MMTxLabel       = "mass_migration"

	// Transaction statuses
	AcceptedTxStatus = "accepted"
	RejectedTxStatus = "rejected"
)

// Rollup metrics
const (
	TransferBatchLabel = "transfer"
	C2TBatchLabel      = "create2transfer"
	MMBatchLabel       = "mass_migration"
	DepositBatchLabel  = "deposit"
)

// Syncing metrics
const (
	SyncAccountsMethod = "sync_accounts"
	SyncBatchesMethod  = "sync_batches"
	SyncRangeMethod    = "sync_range"
	SyncDepositsMethod = "sync_deposits"
	SyncTokensMethod   = "sync_tokens"
	SyncSpokesMethod   = "sync_spokes"
)
