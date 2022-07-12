package storage

type Storage interface {
	StoreLog(Log) error
	GetLogs(SearchOpt) ([]Log, error)

	Close() error
}

// Log defines the basic log info pushed by the server
type Log struct {
	JobID    int    `json:"jobID"`
	LogData  string `json:"logData"`
	Date     string `json:"date"`
	LogLevel string `json:"logLevel"`
}

// SearchOpt defines the different options to filter with
type SearchOpt struct {
	Query     string `json:"query"`
	LogLevel  string `json:"logLevel"`
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
}
