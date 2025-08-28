package fraudlabspro

// The Configuration struct stores the FraudLabs Pro API key.
type Configuration struct {
	apiKey        string
	source        string
	sourceVersion string
}

// OpenConfiguration initializes with the FraudLabs Pro API key
func OpenConfiguration(apikey string) (*Configuration, error) {
	var config = &Configuration{}
	config.apiKey = apikey
	config.source = "sdk-go"
	config.sourceVersion = "1.0.0"
	return config, nil
}
