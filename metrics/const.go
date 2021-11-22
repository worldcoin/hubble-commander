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
	SyncAccountsMethod = "sync_accounts"
	SyncBatchesMethod  = "sync_batches"
	SyncRangeMethod    = "sync_range"
	SyncDepositsMethod = "sync_deposits"
	SyncTokensMethod   = "sync_tokens"
)

// Blockchain metrics
const (
	SinglePubKeyRegisteredLogRetrievalCall = "single_pubkey_registered_log_retrieval_call"
	BatchPubKeyRegisteredLogRetrievalCall  = "batch_pubkey_registered_log_retrieval_call"
	NewBatchLogRetrievalCall               = "new_batch_log_retrieval_call"
	DepositsFinalisedLogRetrievalCall      = "deposits_finalized_log_retrieval_call"
	DepositQueuedLogRetrievalCall          = "deposit_queued_log_retrieval_call"
	DepositSubTreeReadyLogRetrievalCall    = "deposit_sub_tree_ready_log_retrieval_call"
	RegisteredTokenLogRetrievalCall        = "registered_token_log_retrieval_call"
)
