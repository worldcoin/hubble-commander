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
)

// Rollup metrics
const (
	TransferBatchLabel = "transfer"
	C2TBatchLabel      = "create2transfer"
	DepositBatchLabel  = "deposit"
)

// Syncing metrics
const (
	SyncAccountsMethod                  = "sync_accounts"
	SyncBatchesMethod                   = "sync_batches"
	SyncRangeMethod                     = "sync_range"
	SyncDepositsWithNoNewSubTreesMethod = "sync_deposits_no_new_sub_trees"
	SyncDepositsWithNewSubTreesMethod   = "sync_deposits_with_new_sub_trees"
	SyncTokensMethod                    = "sync_tokens"
)
