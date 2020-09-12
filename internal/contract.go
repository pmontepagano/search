package contract

// GlobalContract represents a requirements contract that characterizes a channel
type GlobalContract interface {
	ProjectLocalContract(participant string) PartialContract
}

// PartialContract represents a projection of a GlobalContract to a participant
type PartialContract interface {
	IsSatisfiedBy(candidate LocalContract) bool
}

// LocalContract represents a provides contract
type LocalContract interface{}
