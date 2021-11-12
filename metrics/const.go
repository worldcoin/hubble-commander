package metrics

const (
	namespace = "hubble"

	// Subsystems
	apiSubsystem        = "api"
	batchSubsystem      = "batch"
	commitmentSubsystem = "commitment"
)

// API metrics
const (
	// Transaction labels
	TransferLabel = "transfer"
	C2TLabel      = "create2transfer"

	// Transaction statuses
	AcceptedTxStatus = "accepted"
	RejectedTxStatus = "rejected"
)
