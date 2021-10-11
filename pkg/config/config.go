package config

type providerConfig struct {
	CPU             string
	Memory          string
	Pods            string
	ClusterName     string
}

func GetDefaultConfig() providerConfig {
	return providerConfig{
		CPU: 			"20",
		Memory: 		"100Gi",
		Pods: 			"20",
		ClusterName: 	"default",
	}
}