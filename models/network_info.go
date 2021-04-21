package models

type NetworkInfo struct {
	ChainState
	BlockNumber          uint32  `json:"blockNumber"`
	LatestBatch          *string `json:"latestBatch"`
	LatestFinalisedBatch *string `json:"latestFinalisedBatch"`
}
