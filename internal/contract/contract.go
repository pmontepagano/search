package contract

// GlobalContract represents a requirements contract that characterizes a channel
type GlobalContract struct {
	contract string
	participants []string
}

// PartialContract represents a projection of a GlobalContract to a participant
type PartialContract struct {
	contract string
	participants []string
}

// LocalContract represents a provides contract
type LocalContract struct {
	contract string
	participants []string
}

// Projects GC into PartialContract. For now, this is a dummy implementation that only
// considers contracts of two participants, so we simply copy contract and participant fields
func (gc *GlobalContract) ProjectPartialContract(participant string) PartialContract {
	return PartialContract{
		contract: gc.contract,
		participants: gc.participants,
	}
}

// Determines if candidate LocalContract satisfies the PartialContract. For now, this is
// a dummy implementation that always returns TRUE
// TODO: do the check, this is a dummy implementation
func (pc *PartialContract) IsSatisfiedBy(candidate LocalContract) bool {
	return true
}
