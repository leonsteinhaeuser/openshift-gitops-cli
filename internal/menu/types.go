package menu

type CarrierCreateCluster struct {
	Environment string
	Stage       string
	ClusterName string
	Addons      map[string]map[string]any
	Properties  map[string]string
}

type CarrierCreateStage struct {
	Environment string
	StageName   string
	Properties  map[string]string
}

type CarrierCreateEnvironment struct {
	EnvironmentName string
	Properties      map[string]string
}
