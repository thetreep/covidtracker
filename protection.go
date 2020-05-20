package covidtracker

type Protection struct {
	ID   ProtectionID
	Type ProtectionType
	Name string
}

type ProtectionID string

type ProtectionType string

const (
	Mask ProtectionType = "mask"
	Gel  ProtectionType = "gel"
)
