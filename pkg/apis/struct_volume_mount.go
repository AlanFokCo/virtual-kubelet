package apis

type VolumeMount struct {
	MountPath string   	`name:"MountPath"`
	ReadOnly  string 	`name:"ReadOnly"`
	Name      string    `name:"Name"`
}
