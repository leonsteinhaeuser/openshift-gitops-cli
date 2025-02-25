package menu

type CarrierCreateCluster struct {
	Environment string
	Stage       string
	ClusterName string
	Properties  map[string]any
}

type CarrierCreateStage struct {
	Environment string
	StageName   string
}
