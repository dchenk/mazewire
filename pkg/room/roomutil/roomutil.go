// Package roomutil provides a few utility functions and types for changing a room.Tree structure.
package roomutil

import (
	"fmt"

	"github.com/dchenk/mazewire/pkg/room"
)

// ColumnsPerRow says how many columns each type of row must have.
var ColumnsPerRow = map[string]int{
	"full":                         1,
	"halfhalf":                     2,
	"thirdthirdthird":              3,
	"thirdtwothirds":               2,
	"twothirdsthird":               2,
	"quarterquarterquarterquarter": 4,
	"quarterquarterhalf":           3,
	"quarterhalfquarter":           3,
	"halfquarterquarter":           3,
}

// ChangedTree is an ordered collection of changed section elements.
type ChangedTree []struct {
	Type    string            `msgp:"type"`
	Dyn     int64             `msgp:"dyn"`
	Options map[string]string `msgp:"options"`
	Rows    []struct {
		Type    string            `msgp:"type"`
		Dyn     int64             `msgp:"dyn"`
		Options map[string]string `msgp:"options"`
		Modules [][]struct {
			Type string `msgp:"type"`
			Dyn  int64  `msgp:"dyn"`
			Data []byte `msgp:"data"`
		} `msgp:"modules"`
	} `msgp:"rows"`
}

// Validate loops over the tree to ensure that all the new contents in the tree exist in the newContents slice and are
// of the right role.
// This function checks if each row has the right number of columns, based on the row's type.
func (ct ChangedTree) Validate() error {
	for i := range ct {
		for ii := range ct[i].Rows {
			r := &ct[i].Rows[ii]
			numCols, ok := ColumnsPerRow[r.Type]
			if !ok {
				return ErrInvalidElement(fmt.Sprintf("type given is %q", r.Type))
			}
			if numCols != len(r.Modules) {
				return ErrInvalidElement(fmt.Sprintf("got wrong number of columns; expected %d for row type %q but got %d",
					numCols, r.Type, len(r.Modules)))
			}
			for col := range ct[i].Rows[ii].Modules {
				_ = col
				// TODO: Is there anything to check?
			}
		}
	}
	return nil
}

// ErrInvalidElement reports that an element in a ChangedTree is not valid.
type ErrInvalidElement string

// Error implements the error interface for ErrMissingContent.
func (e ErrInvalidElement) Error() string {
	return fmt.Sprintf("roomutil: an element is invalid (some context: %v)", e)
}

// PublishDynamicElement publishes a copy of the given element's current state (passed in as the data parameter) as
// a dynamic element and creates another copy of the state to be the first version of the element's dynamic settings.
// elemType must be one of the following: "section", "row", "module", "module-all".
// The id value returned is the ID with which the element was published.
func PublishDynamicElement(ds room.DataStore, elemType string, data []byte) (id int64, err error) { // TODO: not implemented

	//if len(data) == 0 {
	//	return 0, errors.New("provided element body is empty")
	//}
	//
	//switch elemType {
	//case "section", "row", "module":
	//	elemType += "_dyn"
	//case "module-all":
	//	elemType = "module_dyn_all"
	//default:
	//	return 0, errors.New("room: unexpected elemType")
	//}
	//
	//id, err = ds.Store(&room.Datum{Data: data})
	//if err != nil {
	//	return
	//}
	//
	//// Save the first version of the dynamic settings.
	//_, err = ds.Store(&room.Datum{Role: elemType + "_sett", K: id, Data: data})
	return

}
