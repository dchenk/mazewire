package instance

// A Manager can manage the system's server instances.
// Implementations are particular to the underlying hosting or cloud service on which the
// instances are running.
type Manager interface {
	GetInstances()
}
