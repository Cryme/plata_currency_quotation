package env

type Environment string

const (
	Local   Environment = "local"
	Dev     Environment = "dev"
	Preprod Environment = "preprod"
	Prod    Environment = "prod"
)
