package models

type ExportStatus struct {
	Uuid            string `json:"uuid"`
	Status          string `json:"status"`
	ChunksAvailable []int  `json:"chunks_available"`
	TotalChunks     int    `json:"total_chunks"`
	FinishedChunks  int    `json:"finished_chunks"`
}

type Chunk []byte

type Asset struct {
	Id   string `json:"Id"`
	Tags []Tag  `json:"tags"`
}

type Tag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Vuln struct {
	Severity                 string `json:"severity"`
	SeverityModificationType string `json:"severity_modification_type"`
	Output                   string `json:"output"`
	Asset                    struct {
		Hostname string `json:"hostname"`
		Uuid     string `json:"uuid"`
	} `json:"asset"`
	Plugin struct {
		Name string   `json:"name"`
		Cpe  []string `json:"cpe"`
	} `json:"plugin"`
}
