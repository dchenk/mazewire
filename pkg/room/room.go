// Package room provides a tool to compile and build HTML pages out of data structures defined here.
package room

import "bytes"

// A DataStore represents a database: it is able to save Datum entities and retrieve them by ID.
// An instance of a DataStore should cache internally the Datum entities already retrieved.
type DataStore interface {
	// Store saves the Data field of a Datum in a data store and returns the ID it was given.
	Store([]byte) (int64, error)

	// Get retrieves a single Datum. This function may be called repeatedly with a single ID throughout the
	// compiling of a Tree or building of a CompiledTree, so all data already queried should be cached for
	// efficient access.
	Get(int64) (*Datum, error)

	// GetMulti retrieves the Datum entities identified by their IDs.
	GetMulti([]int64) ([]Datum, error)
}

// A ModuleCustomizer is able to customize what ModuleCompiler is used to compile a Module and what ModuleBuilder
// is used to build its HTML.
type ModuleCustomizer interface {
	// GetModuleCompiler returns a ModuleCompiler to customize the way a Module is compiled.
	// This function should not unmarshall the Data field in the Module because the same Module will be passed in
	// to the ModuleCompiler's Compile method.
	GetModuleCompiler(*Module) ModuleCompiler

	// GetModuleBuilder returns a ModuleBuilder to customize the HTML output of a Module.
	// This function should not unmarshall the Data field in the Module because the same Module will be passed in
	// to the ModuleBuilder's BuildHTML method.
	GetModuleBuilder(*Module) ModuleBuilder
}

// A ModuleCompiler compiles a Module as much as possible to minimize building HTML upon requests to display the module.
type ModuleCompiler interface {
	// Compile compiles a Module. The Type field of the returned *Module must be StaticElementType if the compiled module
	// can be saved as is, using the Data field as the pre-built HTML not needing any building.
	//
	// A ModuleCustomizer is passed in to make recursive Tree compilation possible.
	//
	// The *Module passed in points to the same Module that the ModuleCustomizer got at first. Compile can decode the
	// Data field as needed.
	//
	// As much of the CSS as possible should be written out to PageCSS at compile time.
	Compile(DataStore, ModuleCustomizer, *Module, *PageCSS) (*Module, error)
}

// A ModuleBuilder is able to build out the HTML of a Module. Each Module will be built dynamically upon each request
// to display if it is not built out into its final static form at compile time.
type ModuleBuilder interface {
	// BuildHTML writes to the Buffer the HTML output of the Module. The CSS for the module that should be used
	// with the page is written to the PageCSS buffer.
	BuildHTML(DataStore, ModuleCustomizer, *Module, *bytes.Buffer, *PageCSS) error
}

// A Datum represents a generic container in which the dynamic data for a Section, Row, or Module may be saved.
// The Id of each Datum is a unique identifier.
type Datum struct {
	Id   int64
	Data []byte
}
