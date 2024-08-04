package telemetry

type Service struct {
	Name        string `yaml:"name"`
	Version     string `yaml:"version"`
	Environment string `yaml:"environment"`
}

type Config struct {
	Service Service `yaml:"service"`
}
