package session

type Attribute int32

const (
	DanceMove Attribute = iota
	Accuracy
	EpochMs
	All
)

func (a Attribute) String() string {
	switch a {
	case DanceMove:
		return "/dance+move"
	case Accuracy:
		return "/accuracy"
	case EpochMs:
		return "/epoch+ms"
	case All:
		return "/"
	}
	return "unknown"
}
