package driver

import (
	mgo "gopkg.in/mgo.v2"
	bson "gopkg.in/mgo.v2/bson"
	rspb "k8s.io/helm/pkg/proto/hapi/release"
)

// CosmosDBDriverName is the string name of this driver.
const CosmosDBDriverName = "CosmosDB"

// CosmosDB is the in-memory storage driver implementation.
type CosmosDB struct {
	db       *mgo.Database
	collName string
}

// NewCosmosDB initializes a new memory driver.
func NewCosmosDB(db *mgo.Database, collName string) *CosmosDB {
	return &CosmosDB{db: db, collName: collName}
}

// Name returns the name of the driver.
func (c *CosmosDB) Name() string {
	return CosmosDBDriverName
}

// Get returns the release named by key or returns ErrReleaseNotFound.
func (c *CosmosDB) Get(key string) (*rspb.Release, error) {
	rel := new(rsbp.Release)
	if err := c.coll().Find(bson.M{"Name": key}).One(rel); err != nil {
		return nil, err
	}
	return rel, nil
}

// List returns the list of all releases such that filter(release) == true
func (c *CosmosDB) List(filter func(*rspb.Release) bool) ([]*rspb.Release, error) {
	query := c.coll().Find(nil)
	lst := []*rsbp.Release{}
	// TODO: is there a way that we can translate the filter into a mongo selector?
	// Doing so would decrease network traffic (possibly greatly)
	if err := query.All(lst); err != nil {
		return nil, err
	}
	ret := []*rsbp.Release{}
	for _, elt := range lst {
		if filter(elt) {
			ret = append(ret, elt)
		}
	}
	return ret, nil
}

// Query returns the set of releases that match the provided set of labels
func (c *CosmosDB) Query(keyvals map[string]string) ([]*rspb.Release, error) {
	// TODO
	return nil, nil
}

// Create creates a new release or returns ErrReleaseExists.
func (c *CosmosDB) Create(key string, rls *rspb.Release) error {
	return nil
}

// Update updates a release or returns ErrReleaseNotFound.
func (c *CosmosDB) Update(key string, rls *rspb.Release) error {
	return nil
}

// Delete deletes a release or returns ErrReleaseNotFound.
func (c *CosmosDB) Delete(key string) (*rspb.Release, error) {
	return nil, nil
}

func (c *CosmosDB) coll() *mgo.Collection {
	return c.db.C(c.collName)
}
