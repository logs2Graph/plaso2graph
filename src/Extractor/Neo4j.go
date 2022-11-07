package Extractor

import (
	"context"
	//"fmt"
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
	query := "CREATE (:Process {created_time: $created_time, timestamp: $timestamp, filename: $filename, fullpath: $fullpath,pid: $pid,commandline: $commandline, "
	query += "ppid: $ppid, pprocess_name: $pprocess_name, pprocess_commandline: $pprocess_commandline, "
	query += "user: $user, user_domain: $user_domain, computer: $computer, logonid: $logonid, evidence: $evidence})"
	//fmt.Println("Created time:" + fmt.Sprint(p.CreatedTime))
	//fmt.Println(fmt.Sprint(p.Evidence))

	parameters := map[string]interface{}{
		"created_time":         p.CreatedTime,
		"timestamp":            p.Timestamp,
		"fullpath":             p.FullPath,
		"filename":             p.Filename,
		"pid":                  p.PID,
		"commandline":          p.Commandline,
		"ppid":                 p.PPID,
		"pprocess_name":        p.ParentProcessName,
		"pprocess_commandline": p.ParentProcessCommandline,
		"user":                 p.User,
		"user_domain":          p.UserDomain,
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
		return tx.Run("match (n) match (m) where n.pid <> 0 and m.pid <> 0 and n.ppid = m.pid and n.pprocess_name = m.fullpath and n.timestamp > m.timestamp merge (m)-[:EXECUTE]->(n)", param)
	})
	handleErr(err)

	//Remove duplicates
	_, err = sess.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		var param map[string]interface{}
		return tx.Run("match (m)-[r:EXECUTE]->(n)<-[s:EXECUTE]-(o) where n.timestamp - m.timestamp < n.timestamp - o.timestamp delete s", param)
	})
	handleErr(err)

}

func InsertUsers(con Neo4JConnector, users []User) {
	for _, u := range users {
		InsertUser(con, u)
	}
	LinkUsers(con)
}

func InsertUser(con Neo4JConnector, u User) {
	sess := con.Driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	_, err := sess.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		return PersistUser(tx, u)
	})
	handleErr(err)
}

func PersistUser(tx neo4j.Transaction, u User) (interface{}, error) {
	query := "CREATE (:User {name: $name, sid: $sid, domain: $domain})"
	parameters := map[string]interface{}{
		"name":   u.Name,
		"sid":    u.SID,
		"domain": u.Domain,
	}
	_, err := tx.Run(query, parameters)
	return nil, err
}

func LinkUsers(con Neo4JConnector) {
	//TODO: Link User With Process
}
