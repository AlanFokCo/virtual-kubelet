package apis

const (
	VOL_TYPE_NFS              = "NFSVolume"
	VOL_TYPE_EMPTYDIR         = "EmptyDirVolume"
	VOL_TYPE_CONFIGFILEVOLUME = "ConfigFileVolume"
)

type Volume struct {
	Type                 string             `name:"Type"`
	Name                 string             `name:"Name"`
	NfsVolumePath        string             `name:"NFSVolume.Path"`
	NfsVolumeServer      string             `name:"NFSVolume.Server"`
	NfsVolumeReadOnly    string   			`name:"NFSVolume.ReadOnly"`
	EmptyDirVolumeEnable string   			`name:"EmptyDirVolume.Enable"`
	ConfigFileToPaths    []ConfigFileToPath `name:"ConfigFileVolume.ConfigFileToPath" type:"Repeated"`
}