package drainer

// Orchestrator orchestrates the queue drainer and the payer
type Orchestrator struct {
	drainer *Drainer
}

// NewOrchestrator creates an orechestrator
func NewOrchestrator(apiURL string, pollTime int) (*Orchestrator, error) {
	creditChan := make(chan int)

	drainer, err := NewDrainer(apiURL, pollTime, creditChan)
	if err != nil {
		return nil, err
	}

	orchestrator := Orchestrator{
		drainer: drainer,
	}
	return &orchestrator, nil
}

// Start starts all the orchestration channels
func (o *Orchestrator) Start() {
	go o.drainer.Drain()
}

// Stop stops all orechestration channels
func (o *Orchestrator) Stop() {
	o.drainer.Stop()
}
