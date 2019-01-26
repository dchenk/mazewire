package room

import "fmt"

// ErrMissingModuleCompiler is an error type that reports that a ModuleCompiler could not be found to
// compile a particular Module.
type ErrMissingModuleCompiler struct {
	*Module
}

// Error implements the error interface for ErrMissingModuleCompiler.
func (em ErrMissingModuleCompiler) Error() string {
	return fmt.Sprintf("room: missing ModuleCompiler implementation for Module with Type %q", em.Type)
}

// ErrMissingModuleBuilder is an error type that reports that a ModuleCompiler could not be found to
// compile a particular Module.
type ErrMissingModuleBuilder struct {
	*Module
}

// Error implements the error interface for ErrMissingModuleBuilder.
func (em ErrMissingModuleBuilder) Error() string {
	return fmt.Sprintf("room: missing ModuleBuilder implementation for Module with Type %q", em.Type)
}
