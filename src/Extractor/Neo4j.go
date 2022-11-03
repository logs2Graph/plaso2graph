package Extractor

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"log"
	. "plaso2graph/master/src/Entity"
)

func handleErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type Neo4JConnector struct {
	Username string
	Password string
	Url      string
	Driver   neo4j.Driver
	Context  context.Context
}

func Neo4jConnect(username string, password string, url string) Neo4JConnector {
	var con Neo4JConnector
	driver, err := neo4j.NewDriver(url, neo4j.BasicAuth(username, password, ""))
	handleErr(err)
	con.Driver = driver
	ctx := context.Background()
	con.Context = ctx
	con.Url = url
	con.Username = username
	con.Password = password

	return con
}

func InsertProcesses(con Neo4JConnector, ps []Process) {
	for _, p := range ps {
		InsertProcess(con, p)
	}
	linkProcess(con)
}

func InsertProcess(con Neo4JConnector, p Process) {
	sess := con.Driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	_, err := sess.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		return PersistProcess(tx, p)
	})
	handleErr(err)
}

func PersistProcess(tx neo4j.Transaction, p Process) (interface{}, error) {
	query := "CREATE (:Process {created_time: $created_time, timestamp: $timestamp, name: $name, pid: $pid,commandline: $commandline, "
	query += "ppid: $ppid, pprocess_name: $pprocess_name, pprocess_commandline: $pprocess_commandline, "
	query += "user: $user, user_domain: $user_domain, computer: $computer, logonid: $logonid, evidence: $evidence})"
	parameters := map[string]interface{}{
		"created_time":         p.CreatedTime,
		"timestamp":            p.Timestamp,
		"name":                 p.Name,
		"pid":                  p.PID,
		"commandline":          p.Commandline,
		"ppid":                 p.PPID,
		"pprocess_name":        p.Pprocess_name,
		"pprocess_commandline": p.Pprocess_commandline,
		"user":                 p.User,
		"user_domain":          p.User_Domain,
		"logonid":              p.LogonID,
		"computer":             p.Computer,
		"evidence":             p.Evidence,
	}
	_, err := tx.Run(query, parameters)
	return nil, err
}

func linkProcess(con Neo4JConnector) {

	//create link based on pid, ppid and name. Quick Filter to avoid some duplicates
	sess := con.Driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	_, err := sess.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		var param map[string]interface{}
		return tx.Run("match (n) match (m) where n.ppid = m.pid and n.pprocess_name = m.name and n.timestamp > m.timestamp merge (m)-[:EXECUTE]->(n)", param)
	})
	handleErr(err)

	//Remove duplicates
	_, err = sess.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		var param map[string]interface{}
		return tx.Run("match (m)-[r:EXECUTE]->(n)<-[s:EXECUTE]-(o) where n.timestamp - m.timestamp < n.timestamp - o.timestamp delete s", param)
	})
	handleErr(err)

}
