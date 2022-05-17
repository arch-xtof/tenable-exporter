package models

type ServerVulnerabilityCount map[string]map[string]*struct {
	Region string
	Group  string

	Count int
}

type WorkstationVulnerabilityCount map[string]map[string]*struct {
	OS  string
	App string

	Count int
}
