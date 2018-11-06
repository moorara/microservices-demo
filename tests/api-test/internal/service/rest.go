package service

type (
	// RESTTest is for testing a service with RESTful API
	RESTTest struct {
		Name           string `yaml:"name"`
		ServiceAddress string `yaml:"service_address"`
		APIBase        string `yaml:"api_base"`
	}
)
