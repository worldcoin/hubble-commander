package metrics

import "github.com/Worldcoin/hubble-commander/utils"

func EventNameToMetricsEventFilterCallLabel(eventName string) string {
	return utils.CamelCaseToSnakeCase(eventName) + "_filter_log_call"
}
