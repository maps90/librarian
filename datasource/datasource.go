package datasource

import (
	"errors"
	"strings"
)

const (
	MYSQL = "mysqlaccess"
	MONGO = "mongoaccess"
)

func NewDatasourceFactory(dbapps, database, server string) (interface{}, error) {
	switch dbapps {
	case MYSQL:
		return newMysqlRepository(server)
	case MONGO:
		return newMongoRepository(database, server)
	}
	dbstrings := []string{
		MYSQL, MONGO,
	}

	return nil, errors.New("Invalid db apps " + dbapps + ", please use (" + strings.Join(dbstrings, "|") + ")")

}
