package Extractor

import (
	"context"
	//"fmt"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"log"
	. "plaso2graph/master/src/Entity"
	"sync"
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

func InitializeNeo4jExtractor(args map[string]interface{}) map[string]interface{} {
	//fmt.Println("Initializing Neo4j Extractor")

	if args["username"] == nil {
		log.Fatal("Username is required")
	}

	if args["password"] == nil {
		log.Fatal("Password is required")
	}

	if args["url"] == nil {
		log.Fatal("Url is required")
	}

	args["connector"] = Neo4jConnect(args["username"].(string), args["password"].(string), args["url"].(string))

	return args
}

func Neo4jExtract(data []interface{}, args map[string]interface{}) {

	if args["verbose"] == nil {
		args["verbose"] = false
	}

	con := args["connector"].(Neo4JConnector)

	var wg sync.WaitGroup
	for _, d := range data {
		switch d.(type) {
		case []Process:
			InsertProcessesNeo4j(con, d.([]Process))
			break

		case []User:
			InsertUsersNeo4j(con, d.([]User))
			break

		case []File:
			InsertFilesNeo4j(con, d.([]File))
			break

		case []ScheduledTask:
			InsertTasksNeo4j(con, d.([]ScheduledTask))
			break

		case []Computer:
			InsertComputersNeo4j(con, d.([]Computer))
			break

		case []Domain:
			InsertDomainsNeo4j(con, d.([]Domain))
			break

		case []WebHistory:
			InsertWebHistoriesNeo4j(con, d.([]WebHistory))
			break

		case []Connection:
			InsertConnectionsNeo4j(con, d.([]Connection))
			break

		case []Event:
			InsertEventsNeo4j(con, d.([]Event))
			break
		}

	}
	wg.Wait()
	if args["verbose"].(bool) {
		log.Println("Neo4j Extractor finished")
	}
}

func ParrallelNeo4jExtract(data []interface{}, args map[string]interface{}) {

	if args["verbose"] == nil {
		args["verbose"] = false
	}

	con := args["connector"].(Neo4JConnector)

	var wg sync.WaitGroup
	for _, d := range data {
		switch d.(type) {
		case []Process:
			go func() {
				wg.Add(1)
				defer wg.Done()
				InsertProcessesNeo4j(con, d.([]Process))
			}()
			break

		case []User:
			go func() {
				wg.Add(1)
				defer wg.Done()
				InsertUsersNeo4j(con, d.([]User))
			}()
			break

		case []File:
			go func() {
				wg.Add(1)
				defer wg.Done()
				InsertFilesNeo4j(con, d.([]File))
			}()
			break

		case []ScheduledTask:
			go func() {
				wg.Add(1)
				defer wg.Done()
				InsertTasksNeo4j(con, d.([]ScheduledTask))
			}()
			break

		case []Computer:
			go func() {
				wg.Add(1)
				defer wg.Done()
				InsertComputersNeo4j(con, d.([]Computer))
			}()
			break

		case []Domain:
			go func() {
				wg.Add(1)
				defer wg.Done()
				InsertDomainsNeo4j(con, d.([]Domain))
			}()
			break

		case []WebHistory:
			go func() {
				wg.Add(1)
				defer wg.Done()
				InsertWebHistoriesNeo4j(con, d.([]WebHistory))
			}()
			break
		}

	}
	wg.Wait()
	if args["verbose"].(bool) {
		log.Println("Neo4j Extractor finished")
	}
}

func Neo4jPostProcessing(args map[string]interface{}) {
	con := args["connector"].(Neo4JConnector)
	linkProcess(con)
	linkUsers(con)
	linkComputers(con)
	handleConnections(con)
	handleEvents(con)
}

func InsertProcessesNeo4j(con Neo4JConnector, ps []Process) {
	for _, p := range ps {
		InsertProcessNeo4j(con, p)
	}
}

func InsertProcessNeo4j(con Neo4JConnector, p Process) {
	sess := con.Driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	_, err := sess.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		return persistProcess(tx, p)
	})
	handleErr(err)
}

func persistProcess(tx neo4j.Transaction, p Process) (interface{}, error) {
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

func InsertUsersNeo4j(con Neo4JConnector, users []User) {
	for _, u := range users {
		InsertUserNeo4j(con, u)
	}
}

func InsertUserNeo4j(con Neo4JConnector, u User) {
	sess := con.Driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	_, err := sess.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		return persistUser(tx, u)
	})
	handleErr(err)
}

func persistUser(tx neo4j.Transaction, u User) (interface{}, error) {
	query := "CREATE (:User {name: $name, sid: $sid, domain: $domain})"
	parameters := map[string]interface{}{
		"name":   u.Name,
		"sid":    u.SID,
		"domain": u.Domain,
	}
	_, err := tx.Run(query, parameters)
	return nil, err
}

func InsertComputersNeo4j(con Neo4JConnector, computers []Computer) {
	for _, c := range computers {
		InsertComputerNeo4j(con, c)
	}
}

func InsertComputerNeo4j(con Neo4JConnector, c Computer) {
	sess := con.Driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	_, err := sess.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		return persistComputer(tx, c)
	})
	handleErr(err)
}

func persistComputer(tx neo4j.Transaction, c Computer) (interface{}, error) {
	query := "CREATE (:Computer {name: $name, domain: $domain})"
	parameters := map[string]interface{}{
		"name":   c.Name,
		"domain": c.Domain,
	}
	_, err := tx.Run(query, parameters)
	return nil, err
}

func InsertTasksNeo4j(con Neo4JConnector, tasks []ScheduledTask) {
	for _, t := range tasks {
		InsertTaskNeo4j(con, t)
	}
}

func InsertTaskNeo4j(con Neo4JConnector, t ScheduledTask) {
	sess := con.Driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	_, err := sess.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		return persistTask(tx, t)
	})
	handleErr(err)
}

func persistTask(tx neo4j.Transaction, t ScheduledTask) (interface{}, error) {
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

func InsertDomainsNeo4j(con Neo4JConnector, domains []Domain) {
	for _, d := range domains {
		InsertDomainNeo4j(con, d)
	}
}

func InsertDomainNeo4j(con Neo4JConnector, d Domain) {
	sess := con.Driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	_, err := sess.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		return persistDomain(tx, d)
	})
	handleErr(err)
}

func persistDomain(tx neo4j.Transaction, d Domain) (interface{}, error) {
	query := "CREATE (:Domain {name: $name})"
	parameters := map[string]interface{}{
		"name": d.Name,
	}
	_, err := tx.Run(query, parameters)
	return nil, err
}

func InsertWebHistoriesNeo4j(con Neo4JConnector, history []WebHistory) {
	for _, h := range history {
		InsertWebHistoryNeo4j(con, h)
	}
}

func InsertWebHistoryNeo4j(con Neo4JConnector, h WebHistory) {
	sess := con.Driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	_, err := sess.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		return persistWebHistory(tx, h)
	})
	handleErr(err)
}

func persistWebHistory(tx neo4j.Transaction, h WebHistory) (interface{}, error) {
	query := "CREATE (:WebHistory {url: $url, title: $title, visit_count: $visit_count, last_visit_time: $last_visit_time, timestamp: $timestamp, path: $path, evidence: $evidence, user: $user})"
	parameters := map[string]interface{}{
		"url":             h.Url,
		"title":           h.Title,
		"visit_count":     h.VisitCount,
		"last_visit_time": h.LastTimeVisited,
		"path":            h.Path,
		"evidence":        h.Evidence,
		"user":            h.User,
		"timestamp":       h.Timestamp,
	}
	_, err := tx.Run(query, parameters)
	return nil, err
}

func InsertFilesNeo4j(con Neo4JConnector, files []File) {
	for _, f := range files {
		InsertFileNeo4j(con, f)
	}
}

func InsertFileNeo4j(con Neo4JConnector, f File) {
	sess := con.Driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	_, err := sess.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		return persistFile(tx, f)
	})
	handleErr(err)
}

func persistFile(tx neo4j.Transaction, f File) (interface{}, error) {
	query := "CREATE (:File {fullpath: $fullpath, filename: $filename, extension: $extension, is_allocated: $is_allocated, date: $date, timestamp: $timestamp, timestamp_desc: $timestamp_desc, evidence: $evidence, date: $date})"
	parameters := map[string]interface{}{
		"fullpath":       f.FullPath,
		"filename":       f.Filename,
		"extension":      f.Extension,
		"is_allocated":   f.IsAllocated,
		"date":           f.Date,
		"timestamp":      f.Timestamp,
		"timestamp_desc": f.TimestampDesc,
		"evidence":       f.Evidence,
	}
	_, err := tx.Run(query, parameters)
	return nil, err
}

func InsertConnectionsNeo4j(con Neo4JConnector, connections []Connection) {
	for _, c := range connections {
		InsertConnectionNeo4j(con, c)
	}
}

func InsertConnectionNeo4j(con Neo4JConnector, c Connection) {
	sess := con.Driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})

	_, err := sess.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		return persistConnection(tx, c)
	})
	handleErr(err)
}

func persistConnection(tx neo4j.Transaction, c Connection) (interface{}, error) {
	query := "CREATE (:Connection {timestamp: $timestamp, date:$date, protocol: $protocol, ip_source: $ip_source, ip_destination: $ip_destination, port_source: $port_source, port_destination: $port_destination, user: $user, user_domain: $user_domain, computer: $computer, process: $process, process_id: $process_id})"
	parameters := map[string]interface{}{
		"timestamp":        c.Timestamp,
		"date":             c.Date,
		"protocol":         c.Protocol,
		"ip_source":        c.SourceIP,
		"ip_destination":   c.DestinationIP,
		"port_source":      c.SourcePort,
		"port_destination": c.DestinationPort,
		"user":             c.User,
		"user_domain":      c.UserDomain,
		"computer":         c.Computer,
		"process":          c.ProcessName,
		"process_id":       c.ProcessId,
	}
	_, err := tx.Run(query, parameters)
	return nil, err
}

func InsertEventsNeo4j(con Neo4JConnector, events []Event) {
	for _, e := range events {
		InsertEventNeo4j(con, e)
	}
}

func InsertEventNeo4j(con Neo4JConnector, e Event) {
	sess := con.Driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})

	_, err := sess.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		return persistEvent(tx, e)
	})
	handleErr(err)
}

func persistEvent(tx neo4j.Transaction, e Event) (interface{}, error) {
	query := `CREATE (:Event {timestamp: $timestamp, date: $date, event_type: $event_type, title: $title,
		user_source: $user_source, user_destination: $user_destination, domain_source: $domain_source,
		domain_destination: $domain_destination, process: $process, process_id: $process_id, fullpath: $fullpath, filename: $filename,
		extension: $extension})`
	parameters := map[string]interface{}{
		"timestamp":          e.Timestamp,
		"date":               e.Date,
		"title":              e.Title,
		"event_type":         e.Type,
		"user_source":        e.UserSource,
		"user_destination":   e.UserDestination,
		"domain_source":      e.UserDomainSource,
		"domain_destination": e.UserDestinationDomain,
		"process":            e.Process,
		"process_id":         e.ProcessId,
		"fullpath":           e.FullPath,
		"filename":           e.Filename,
		"extension":          e.Extension,
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

func handleConnections(con Neo4JConnector) {
	sess := con.Driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})

	//Create Hosts Nodes based on Connection's IP destination
	_, err := sess.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		query := "match (n:Connection) with collect(distinct n.ip_destination) as ip_dests FOREACH (ip IN ip_dests | Create (:Host {domain: \"\", ip:ip}))"
		parameters := map[string]interface{}{}
		_, err := tx.Run(query, parameters)
		return nil, err
	})
	handleErr(err)

	// Convert Connection nodes to relationships between Hosts and Processes
	_, err = sess.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		query := `match (n:Connection) with collect(n) as connections
		UNWIND connections as c
		match (p:Process) where p.fullpath = c.process and p.pid = c.process_id and p.timestamp < c.timestamp
		match (h:Host) where h.ip = c.ip_destination 
		merge (p)-[r:CONNECT{port_source:c.port_source,port_destination:c.port_destination, ip_source:c.ip_source, timestamp:c.timestamp, date:c.date}]->(h)`

		parameters := map[string]interface{}{}
		_, err := tx.Run(query, parameters)
		return nil, err
	})
	handleErr(err)

}

func handleEvents(con Neo4JConnector) {
	sess := con.Driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})

	// Create File Based On Events "CreateFile" and "DeleteFile"
	_, err := sess.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		query := `match (e:Event) where e.event_type = "File Created" with collect(e) as events
		UNWIND events as event
		create (:File {fullpath: event.fullpath, filename: event.filename, extension:event.extension, 
		timestamp:event.timestamp, date:event.date, timestamp_desc:"Creation Time"})
		with event
		match (f:File) match (p:Process) where f.fullpath = event.fullpath and p.fullpath = event.process and p.pid = event.process_id
		merge (p)-[:CREATE]->(f)`
		parameters := map[string]interface{}{}
		_, err := tx.Run(query, parameters)
		return nil, err
	})
	handleErr(err)

	_, err = sess.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		query := `match (e:Event) where e.event_type = "File Deleted" with collect(e) as events
		UNWIND events as event
		create (:File {fullpath: event.fullpath, filename: event.filename, extension:event.extension, 
		timestamp:event.timestamp, date:event.date, timestamp_desc:"Deletion Time"})
		with event
		match (f:File) match (p:Process) where f.fullpath = event.fullpath and p.fullpath = event.process and p.pid = event.process_id
		merge (p)-[:Delete]->(f)`
		parameters := map[string]interface{}{}
		_, err := tx.Run(query, parameters)
		return nil, err
	})
	handleErr(err)

}
