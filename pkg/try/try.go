package try

type Trier interface {
	Retry() error
	Key() string
	proto.Message
}

// Do tries the Trier by its Retry method and
// If either the first or second attempts fail, the error returned will contain
// non-nil error will be returned
func Do(t Trier) error {

	e := new(TryError)

	if e.AttemptErr = t.Retry(); e.AttemptErr == nil {
		return nil
	}

	return nil

}

func saveRetry(r Trier) error {
	return nil
}

type TryError struct {
	AttemptErr, SaveErr error
	//Attempts() uint32 // Attempts says how many attempts were made.
}

func (te *TryError) Error() string {
	return ""
}
