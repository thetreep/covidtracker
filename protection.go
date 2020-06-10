package covidtracker

type Protection struct {
	ID       ProtectionID
	Type     ProtectionType
	Name     string
	Quantity int
}

type ProtectionID string

type ProtectionType string

const (
	Mask         ProtectionType = "mask"
	MaskSewn     ProtectionType = "mask-sewn"
	MaskSurgical ProtectionType = "mask-surgical"
	MaskFFPX     ProtectionType = "mask-ffpx"
	Gel          ProtectionType = "gel"
)
