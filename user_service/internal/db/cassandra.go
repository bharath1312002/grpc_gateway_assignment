package db

import (
	"github.com/gocql/gocql"
	"log"
	"user_service/config"
)

type CassandraDetailsSvc struct {
	cfg *config.Config
}

func NewCassandraDetailsSvc(cfg *config.Config) *CassandraDetailsSvc {
	return &CassandraDetailsSvc{
		cfg: cfg,
	}
}

func (a *CassandraDetailsSvc) ConnectCassandra() *gocql.Session {

	cluster := gocql.NewCluster(a.cfg.CassandraDetails.Address)
	cluster.Keyspace = a.cfg.CassandraDetails.KeySpace
	cluster.Consistency = gocql.Quorum
	cluster.Port = a.cfg.CassandraDetails.Port
	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatalf("Failed to connect to Cassandra: %v", err)
	}
	return session
}
