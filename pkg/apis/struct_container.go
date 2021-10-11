package apis

type Container struct {
	Name            string           `json:"Name" xml:"Name" `
	Image           string           `json:"Image" xml:"Image"`
	Memory          float64          `json:"Memory" xml:"Memory"`
	Cpu             float64          `json:"Cpu" xml:"Cpu"`
	RestartCount    int              `json:"RestartCount" xml:"RestartCount"`
	WorkingDir      string           `json:"WorkingDir" xml:"WorkingDir"`
	ImagePullPolicy string           `json:"ImagePullPolicy" xml:"ImagePullPolicy"`
	Commands        []string         `json:"Commands" xml:"Commands"`
	Args            []string         `json:"Args" xml:"Args"`
	PreviousState   ContainerState   `json:"PreviousState" xml:"PreviousState"`
	CurrentState    ContainerState   `json:"CurrentState" xml:"CurrentState"`
	VolumeMounts    []VolumeMount    `json:"VolumeMounts" xml:"VolumeMounts"`
	Ports           []ContainerPort  `json:"Ports" xml:"Ports"`
	EnvironmentVars []EnvironmentVar `json:"EnvironmentVars" xml:"EnvironmentVars"`
}

