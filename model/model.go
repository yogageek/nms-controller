package model

import "fmt"

type ConfigIdDetail struct {
	ID           int64          `json:"id,omitempty"`
	CustomConfig []CustomConfig `json:"configs"`
}

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
	UidPath string   `json:"uid_path"`
	Metrics []Metric `json:"metrics"`
}

type CustomConfigData struct {
	QueryNameKeys string
	Label         map[string]string
	Value         interface{}
}

//一個cfg可能有多個queryKeys
func (c *CustomConfig) GetQueryKeys() []string {
	var queryNameKeys []string
	queryName := c.QueryName
	for _, metric := range c.Metrics {
		queryKey := queryName + metric.Header
		queryNameKeys = append(queryNameKeys, queryKey)
	}
	return queryNameKeys
}

//一個cfg可能有多個queryKeys
func (c *CustomConfig) PrintQueryNameAndUrl() {
	queryNameAndUrl := c.QueryName + c.Target
	fmt.Println(queryNameAndUrl)
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

type RedisMetric struct {
	UID    string  `json:"uid"`
	Value  float64 `json:"value"`
	Labels []Label `json:"labels"`
}

type Err struct {
	Error_ string `json:"error,omitempty"`
}

//exporter
type MetricE struct {
	UID    string  `json:"uid"`
	Value  float64 `json:"value"`
	Labels []struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"labels"`
}

type MetricsE []Metric
