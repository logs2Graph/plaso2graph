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

func Links(con Neo4JConnector) {
	linkProcess(con)
	linkUsers(con)
	linkComputers(con)
}

func InsertProcesses(con Neo4JConnector, ps []Process) {
	for _, p := range ps {
		InsertProcess(con, p)
	}
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

func InsertUsers(con Neo4JConnector, users []User) {
	for _, u := range users {
		InsertUser(con, u)
	}
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

func InsertComputers(con Neo4JConnector, computers []Computer) {
	for _, c := range computers {
		InsertComputer(con, c)
	}
}

func InsertComputer(con Neo4JConnector, c Computer) {
	sess := con.Driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	_, err := sess.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		return PersistComputer(tx, c)
	})
	handleErr(err)
}

func PersistComputer(tx neo4j.Transaction, c Computer) (interface{}, error) {
	query := "CREATE (:Computer {name: $name, domain: $domain})"
	parameters := map[string]interface{}{
		"name":   c.Name,
		"domain": c.Domain,
	}
	_, err := tx.Run(query, parameters)
	return nil, err
}

func InsertScheduledTasks(con Neo4JConnector, tasks []ScheduledTask) {
	for _, t := range tasks {
		InsertScheduledTask(con, t)
	}
}

func InsertScheduledTask(con Neo4JConnector, t ScheduledTask) {
	sess := con.Driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	_, err := sess.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		return PersistScheduledTask(tx, t)
	})
	handleErr(err)
}

func PersistScheduledTask(tx neo4j.Transaction, t ScheduledTask) (interface{}, error) {
	query := "CREATE (:ScheduledTask {application: $application, user: $user, comment: $comment, trigger: $trigger})"
	parameters := map[string]interface{}{
		"application": t.Application,
		"user":        t.User,
		"comment":     t.Comment,
		"trigger":     t.Trigger,
	}
	_, err := tx.Run(query, parameters)
	return nil, err
}

func InsertDomains(con Neo4JConnector, domains []Domain) {
	for _, d := range domains {
		InsertDomain(con, d)
	}
}

func InsertDomain(con Neo4JConnector, d Domain) {
	sess := con.Driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	_, err := sess.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		return PersistDomain(tx, d)
	})
	handleErr(err)
}

func PersistDomain(tx neo4j.Transaction, d Domain) (interface{}, error) {
	query := "CREATE (:Domain {name: $name})"
	parameters := map[string]interface{}{
		"name": d.Name,
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

func linkUsers(con Neo4JConnector) {
	sess := con.Driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	_, err := sess.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		var param map[string]interface{}
		return tx.Run("match (u:User) match (p) where u.name = p.user merge (u)-[:BY]->(p)", param)
	})

	handleErr(err)
}

func linkComputers(con Neo4JConnector) {
	sess := con.Driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	_, err := sess.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		var param map[string]interface{}
		return tx.Run("match (c:Computer) match (p) where c.name = p.computer merge (c)-[:ON]->(p)", param)
	})
	handleErr(err)
}
