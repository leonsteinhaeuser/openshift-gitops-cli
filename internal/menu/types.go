package menu

type CarrierCreateCluster struct {
	Environment string
	Stage       string
	ClusterName string
}

type CarrierCreateStage struct {
	Environment string
	StageName   string
}
