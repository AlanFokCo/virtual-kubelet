package apis

type ContainerState struct {
	State        string `json:"State" xml:"State"`
	DetailStatus string `json:"DetailStatus" xml:"DetailStatus"`
	ExitCode     int    `json:"ExitCode" xml:"ExitCode"`
	StartTime    string `json:"StartTime" xml:"StartTime"`
	FinishTime   string `json:"FinishTime" xml:"FinishTime"`
}