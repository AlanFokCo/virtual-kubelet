package apis

type ContainerGroup struct {
	ContainerGroupId   string      `json:"ContainerGroupId" xml:"ContainerGroupId"`
	ContainerGroupName string      `json:"ContainerGroupName" xml:"ContainerGroupName"`
	Memory             float64     `json:"Memory" xml:"Memory"`
	Cpu                float64     `json:"Cpu" xml:"Cpu"`
	RestartPolicy      string      `json:"RestartPolicy" xml:"RestartPolicy"`
	IntranetIp         string      `json:"IntranetIp" xml:"IntranetIp"`
	Status             string      `json:"Status" xml:"Status"`
	InternetIp         string      `json:"InternetIp" xml:"InternetIp"`
	CreationTime       string      `json:"CreationTime" xml:"CreationTime"`
	SucceededTime      string      `json:"SucceededTime" xml:"SucceededTime"`
	Tags               []Tag       `json:"Tags" xml:"Tags"`
	Events             []Event     `json:"Events" xml:"Events"`
	Containers         []Container `json:"Containers" xml:"Containers"`
	Volumes            []Volume    `json:"Volumes" xml:"Volumes"`
}
