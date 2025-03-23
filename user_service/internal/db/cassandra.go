package db

import (
	"github.com/gocql/gocql"
	"log"
)

func ConnectCassandra() *gocql.Session {
	cluster := gocql.NewCluster("127.0.0.1") // Ensure this matches Cassandra's rpc_address
	cluster.Keyspace = "user_service"        // Ensure this keyspace exists
	cluster.Consistency = gocql.Quorum
	cluster.Port = 9042 // Ensure this matches Cassandra's native_transport_port
	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatalf("Failed to connect to Cassandra: %v", err)
	}
	return session
}
