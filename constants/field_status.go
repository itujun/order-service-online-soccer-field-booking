package constants

type FieldStatusString string

const (
	AvailableStatus FieldStatusString = "pending"
	BookedStatus    FieldStatusString = "settlement"
)

func (p FieldStatusString) String() string {
	return string(p)
}
