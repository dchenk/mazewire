package specs

import (
	"encoding/json"
	"fmt"

	"github.com/dchenk/mazewire/pkg/types/version"
)

// A Spec is a specification of a plugin's requirement for its use.
type Spec struct {
	// Type is the type of specification defined in this PluginSpec.
	Type SpecType

	// Value is the particular value for the spec type the plugin needs.
	// The Value is interpreted in light of what the Type is and may even be ignored for some spec types.
	Value fmt.Stringer
}

// A SpecType defines a specification by a plugin, which customizes the way in which a plugin is invoked.
// Whether a plugin may be invoked for certain kinds of hooks or even register custom hooks depends entirely
// on the kinds of SpecType specifications it has made and been authorized for.
type SpecType uint32

const (
	// DatabaseCapability indicates a kind of database capability a plugin needs.
	DatabaseCapability SpecType = 1

	// Dependency indicates a particular dependency on the version of the main app system or another plugin at
	// a particular version.
	// The value must be a JSON-encoded PluginMetadata object.
	Dependency SpecType = 2

	// RegisterHooks indicates that a plugin needs to register custom hook types.
	RegisterHooks SpecType = 3

	// Async indicates that the plugin needs to run asynchronously upon invocation by a hook. The main application
	// must not wait for the plugin to return control.
	Async SpecType = 4

	// Email indicates a kind of email (sending or email log checking) capability a plugin needs.
	Email SpecType = 5

	// PageRender indicates that a plugin needs to modify the rendering of pages.
	PageRender SpecType = 6

	// Tasks indicates that a plugin needs to register tasks.
	Tasks SpecType = 7

	// Plugins indicates that a plugin needs to manage the system's plugins.
	Plugins SpecType = 8
)

// A DependencySpec specifies a plugin's dependency.
// To specify a minimum version for the main app system, set ID to "-system-"; otherwise, set it to the ID of
// the plugin.
type DependencySpec struct {
	ID  string
	Ver version.Version
}

// String implements fmt.Stringer for the DependencySpec type.
func (ds *DependencySpec) String() string {
	jsonBytes, err := json.Marshal(ds)
	if err != nil {
		panic(fmt.Sprintf("specs: could not JSON marshal DependencySpec; %v", err))
	}
	return string(jsonBytes)
}

// AsyncSpec specifies that a hook must run asynchronously from the main application.
// An AsyncSpec should be set for a plugin only if the plugin is in fact intended to be run asynchronously;
// that is, there is no "false" for AsyncSpec.
type AsyncSpec struct{}

// String implements fmt.Stringer for the AsyncSpec type.
func (as *AsyncSpec) String() string {
	return "true"
}

type DatabaseCapabilityType uint8

const (
	Read DatabaseCapabilityType = 1 << iota
	Write
)

// WithDatabaseCapability is a helper function for constructing a database access capability spec.
func WithDatabaseCapability(tableName string, access DatabaseCapabilityType) Spec {
	// TODO: loop over table names and set up the Value accordingly
	return Spec{Type: DatabaseCapability}
}

func WithEmailSending() Spec {
	return Spec{} // TODO
}

type emailSendingSpec struct {
}

//
//type ContentPermissions struct {
//	data.ContentGetter
//	data.ContentInserter
//	data.ContentDeleter
//	data.ContentManager
//}
