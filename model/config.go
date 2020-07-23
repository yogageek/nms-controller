package model

type CustomConfig struct {
	QueryName       string `json:"query_name"`
	Target          string `json:"target"`
	Interval        int    `json:"interval"`
	Response        string `json:"response"`
	HaveParameters  bool   `json:"have_parameters"`
	QueryParameters []struct {
		Name  string   `json:"name"`
		Value []string `json:"value"`
	} `json:"query_parameters,omitempty"`
	Metrics []Metric `json:"metrics"`
}

type Metric struct {
	Header string `json:"header"`
	Help   string `json:"help"`
	Path   string `json:"path"`
	Type   string `json:"type"`
	Labels []struct {
		Name  string `json:"name"`
		Path  string `json:"path,omitempty"`
		Value string `json:"value,omitempty"`
	} `json:"labels"`
}

type Label struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
