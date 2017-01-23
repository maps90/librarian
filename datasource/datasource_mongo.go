package datasource

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type MongodbAccess struct {
	Server   string
	Database string
	Session  *mgo.Session
}

func newMongoRepository(server, database string) (DataAccessor, error) {
	repo := &MongodbAccess{Server: server, Database: database}
	session, err := mgo.Dial(server)
	if err != nil {
		return nil, err
	}
	repo.Session = session
	return repo, nil
}

func (r *MongodbAccess) getCollection(session *mgo.Session, n string) *mgo.Collection {
	return session.DB(r.Database).C(n)
}

func (r MongodbAccess) Insert(data Data) (interface{}, error) {
	session := r.Session.Copy()

	defer session.Close()

	collectionName := data.PersistenceName()

	c := r.getCollection(session, collectionName)
	if err := c.Insert(data); err != nil {
		return nil, err
	}
	msg := fmt.Sprintf("Data inserted to %s collection", collectionName)
	return struct{ msg string }{msg}, nil
}

func (r MongodbAccess) Update(id interface{}, data Data) (interface{}, error) {
	session := r.Session.Copy()

	defer session.Close()
	cName := data.PersistenceName()
	c := r.getCollection(session, cName)
	if err := c.Update(bson.M{"_id": id}, data); err != nil {
		return nil, err
	}
	msg := fmt.Sprintf("Collection %s with objectId(%s) is updated", cName, id)
	return struct{ msg string }{msg}, nil
}

func (r MongodbAccess) Delete(id interface{}, data Data) (interface{}, error) {
	session := r.Session.Copy()

	defer session.Close()
	cName := data.PersistenceName()
	c := r.getCollection(session, cName)

	if err := c.RemoveId(id); err != nil {
		return nil, err
	}
	msg := fmt.Sprintf("Collection %s with objectId(%s) is deleted", cName, id)
	return struct{ msg string }{msg}, nil
}

func (r MongodbAccess) Find(data Data, query map[string]interface{}, order []string, results interface{}) error {
	session := r.Session.Copy()

	defer session.Close()
	c := r.getCollection(session, data.PersistenceName())

	return c.Find(query).Sort(order...).All(results)
}

func (r MongodbAccess) FindPaging(data Data, query map[string]interface{}, order []string, page, limit int, results interface{}) error {
	session := r.Session.Copy()
	defer session.Close()

	c := r.getCollection(session, data.PersistenceName())

	return c.Find(query).Sort(order...).Limit(limit).Skip((page - 1) * limit).All(results)
}

func (r MongodbAccess) FindById(data Data, objectId interface{}, result interface{}) error {
	session := r.Session.Copy()

	defer session.Close()
	c := r.getCollection(session, data.PersistenceName())
	return c.FindId(objectId).One(result)
}
