package grafana

type AlertModel struct {
	State       string
	EvalMatches []EvalMatchModel
}

type EvalMatchModel struct {
	Value float64
	Tags  EvalMatchModelTag
}
type EvalMatchModelTag struct {
	Namespace string
	Pod       string
	Container string
	Node      string
}

func (m AlertModel) IsOk() bool {
	return m.State == "ok"
}

func (m AlertModel) IsAlerting() bool {
	return m.State == "alerting"
}
