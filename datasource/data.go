package datasource

type DataAccessor interface {
	Insert(data Data) (interface{}, error)
	Update(id interface{}, data Data) (interface{}, error)
	Delete(id interface{}, data Data) (interface{}, error)
	Find(data Data, query map[string]interface{}, order []string, results interface{}) error
	FindById(data Data, id interface{}, result interface{}) error
	FindPaging(data Data, query map[string]interface{}, order []string, page, limit int, results interface{}) error
}

type Data interface {
	PersistenceName() string
}