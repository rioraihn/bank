package http

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

type HealthResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

type DatabaseHealthResponse struct {
	Status              string         `json:"status"`
	Message             string         `json:"message,omitempty"`
	Timestamp           string         `json:"timestamp"`
	WalletDatabase      DatabaseStatus `json:"wallet_database"`
	TransactionDatabase DatabaseStatus `json:"transaction_database"`
}

type DatabaseStatus struct {
	Connected bool   `json:"connected"`
	Message   string `json:"message"`
	Latency   string `json:"latency,omitempty"`
}
