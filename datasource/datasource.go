package datasource

import (
	"errors"
	"strings"
)

const (
	MYSQL = "mysqlaccess"
	MONGO = "mongoaccess"
)

func NewDataAccessFactory(dbapps, database, server string) (*DataAccessor, error) {

	switch dbapps {
	case MYSQL:
		return NewMysqlRepository(server)
	case MONGO:
		return NewMongoRepository(database, server)
	}
	dbstrings := []string{
		MYSQL, MONGO,
	}

	return nil, errors.New("Invalid db apps " + dbapps + ", please use (" + strings.Join(dbstrings, "|")+ ")" )

}