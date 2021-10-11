package apis

type Event struct {
	Count          int    `json:"Count" xml:"Count"`
	Type           string `json:"Type" xml:"Type"`
	Name           string `json:"Name" xml:"Name"`
	Message        string `json:"Message" xml:"Message"`
	FirstTimestamp string `json:"FirstTimestamp" xml:"FirstTimestamp"`
	LastTimestamp  string `json:"LastTimestamp" xml:"LastTimestamp"`
}