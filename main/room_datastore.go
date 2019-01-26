package main

import (
	"github.com/dchenk/mazewire/pkg/data"
	"github.com/dchenk/mazewire/pkg/room"
)

// dataStore implements room.DataStore for this project. The room.Datum type is represented as a Blob that
// is stored in a particular site's blobs table.
type dataStore struct {
	s *data.Site // Applicable site for each request
}

// Store stores the Data field of a room.Datum with the role of "room_elem".
func (ds *dataStore) Store(d []byte) (insertID int64, err error) {
	insertID, err = ds.s.BlobInsert("room_elem", 0, d)
	return
}

// GetMulti retrieves a room.Datum element, which is represented as the Blob type in our datastore.
func (ds *dataStore) Get(id int64) (*room.Datum, error) {
	blob, err := ds.s.BlobByID(id)
	if err != nil {
		return nil, err
	}
	return &room.Datum{Id: blob.Id, Data: blob.V}, nil
}

// GetMulti retrieves multiple room.Datum elements, which are represented as the Blob type in our datastore.
func (ds *dataStore) GetMulti(ids []int64) ([]room.Datum, error) {
	blobs, err := ds.s.BlobsIdIn(ids)
	roomData := make([]room.Datum, len(blobs))
	for i := range blobs {
		roomData[i] = room.Datum{Id: blobs[i].Id, Data: blobs[i].V}
	}
	return roomData, err
}
