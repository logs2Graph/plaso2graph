package Extractor

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"log"
	. "plaso2graph/master/src/Entity"
	"sync"
	"time"
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

		case []ScriptBlock:
			InsertScriptBlocksNeo4j(con, d.([]ScriptBlock))
			break

		case []User:
			InsertUsersNeo4j(con, d.([]User))
			break

		case []Group:
			InsertGroupsNeo4j(con, d.([]Group))
			break

		case []File:
			InsertFilesNeo4j(con, d.([]File))
			break

		case []ScheduledTask:
			InsertTasksNeo4j(con, d.([]ScheduledTask))
			break

		case []Service:
			InsertServicesNeo4j(con, d.([]Service))
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
		case []Registry:
			InsertRegistriesNeo4j(con, d.([]Registry))
			break
		}

	}
	wg.Wait()
	/*if args["verbose"].(bool) {
		log.Println("Neo4j Extractor finished")
	}*/
}

func Neo4jPostProcessing(args map[string]interface{}) {

	wg := sync.WaitGroup{}

	con := args["connector"].(Neo4JConnector)
	fmt.Println("Linking processes...")

	go func() {
		wg.Add(1)
		linkProcess(con)
		wg.Done()
	}()

	fmt.Println("Linking user to process...")

	go func() {
		wg.Add(1)
		wg.Done()
		linkUsers(con)
	}()

	fmt.Println("Linking ScriptBlocks...")
	go func() {
		wg.Add(1)
		linkScriptBlock(con)
		wg.Done()
	}()

	fmt.Println("Linking computers...")
	go func() {
		wg.Add(1)
		linkComputers(con)
		wg.Done()

	}()

	fmt.Println("Linking Connections...")

	go func() {
		wg.Add(1)
		handleConnections(con)
		wg.Done()
	}()

	fmt.Println("Processing Events...")
	go func() {
		wg.Add(1)
		handleEvents(con)
		wg.Done()
	}()

	time.Sleep(2 * time.Second)
	wg.Wait()

	updateIds(con)
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

func InsertScriptBlocksNeo4j(con Neo4JConnector, sb []ScriptBlock) {
	for _, s := range sb {
		InsertScriptBlockNeo4j(con, s)
	}
}

func InsertScriptBlockNeo4j(con Neo4JConnector, s ScriptBlock) {
	sess := con.Driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	_, err := sess.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		return persistScriptBlock(tx, s)
	})
	handleErr(err)
}

func persistScriptBlock(tx neo4j.Transaction, s ScriptBlock) (interface{}, error) {
	query := `CREATE (:ScriptBlock {date: $date, timestamp: $timestamp, scriptblockid: $scriptblockid, scriptblocktext: $scriptblocktext, context: $context, 
		process_id: $processid, message_number: $message_number, message_total: $message_total, path: $path, computer: $computer, evidence: $evidence})`
	parameters := map[string]interface{}{
		"date":            s.Date,
		"timestamp":       s.Timestamp,
		"scriptblockid":   s.ScriptBlockID,
		"scriptblocktext": s.Text,
		"processid":       s.ProcessID,
		"message_number":  s.MessageNumber,
		"message_total":   s.MessageTotal,
		"path":            s.Path,
		"computer":        s.Computer,
		"context":         s.Context,
		"evidence":        s.Evidence,
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
	query := "CREATE (:User {fullname: $fullname, username: $username, comments: $comments, sid: $sid, domain: $domain})"
	parameters := map[string]interface{}{
		"fullname": u.FullName,
		"username": u.Username,
		"comments": u.Comments,
		"sid":      u.SID,
		"domain":   u.Domain,
	}
	_, err := tx.Run(query, parameters)
	return nil, err
}

func InsertGroupsNeo4j(con Neo4JConnector, groups []Group) {
	for _, g := range groups {
		InsertGroupNeo4j(con, g)
	}
}

func InsertGroupNeo4j(con Neo4JConnector, g Group) {
	sess := con.Driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	_, err := sess.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		return persistGroup(tx, g)
	})
	handleErr(err)
}

func persistGroup(tx neo4j.Transaction, g Group) (interface{}, error) {
	query := "CREATE (:Group {name: $name, domain: $domain, computer:$computer, evidence: $evidence})"
	parameters := map[string]interface{}{
		"name":     g.Name,
		"domain":   g.Domain,
		"computer": g.Computer,
		"evidence": g.Evidence,
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
	query := "CREATE (:ScheduledTask {application: $application, user: $user, comment: $comment, trigger: $trigger, computer: $computer, evidence: $evidence})"
	parameters := map[string]interface{}{
		"application": t.Application,
		"user":        t.User,
		"comment":     t.Comment,
		"trigger":     t.Trigger,
		"computer":    t.Computer,
		"evidence":    t.Evidence,
	}
	_, err := tx.Run(query, parameters)
	return nil, err
}

func InsertServicesNeo4j(con Neo4JConnector, services []Service) {
	for _, service := range services {
		InsertServiceNeo4j(con, service)
	}
}

func InsertServiceNeo4j(con Neo4JConnector, service Service) {
	sess := con.Driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	_, err := sess.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		return persistService(tx, service)
	})
	handleErr(err)
}

func persistService(tx neo4j.Transaction, s Service) (interface{}, error) {
	query := "CREATE (:Service {name: $name, filename: $filename, service_type: $service_type, start_type: $start_type, error_control: $error_control, user: $user, computer: $computer, dll: $dll, evidence: $evidence})"
	parameters := map[string]interface{}{
		"name":          s.Name,
		"filename":      s.Filename,
		"service_type":  s.ServiceType,
		"start_type":    s.StartType,
		"error_control": s.ErrorControl,
		"dll":           s.Dll,
		"user":          s.User,
		"computer":      s.Computer,
		"evidence":      s.Evidence,
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
	query := "CREATE (:WebHistory {url: $url, title: $title, visit_count: $visit_count, last_visit_time: $last_visit_time, timestamp: $timestamp, path: $path, evidence: $evidence, user: $user, computer: $computer})"
	parameters := map[string]interface{}{
		"url":             h.Url,
		"title":           h.Title,
		"visit_count":     h.VisitCount,
		"last_visit_time": h.LastTimeVisited,
		"path":            h.Path,
		"evidence":        h.Evidence,
		"user":            h.User,
		"computer":        h.Computer,
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
	query := "CREATE (:File {fullpath: $fullpath, filename: $filename, extension: $extension, is_allocated: $is_allocated, date: $date, timestamp: $timestamp, timestamp_desc: $timestamp_desc, evidence: $evidence, date: $date, computer: $computer})"
	parameters := map[string]interface{}{
		"fullpath":       f.FullPath,
		"filename":       f.Filename,
		"extension":      f.Extension,
		"is_allocated":   f.IsAllocated,
		"date":           f.Date,
		"timestamp":      f.Timestamp,
		"timestamp_desc": f.TimestampDesc,
		"computer":       f.Computer,
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
		domain_destination: $domain_destination, group: $group, group_domain: $group_domain, process_source: $process_source, process_source_id: $process_source_id,
		process_target: $process_target, process_target_id: $process_target_id, fullpath: $fullpath, filename: $filename,
		extension: $extension, computer: $computer,evidence: $evidence})`
	parameters := map[string]interface{}{
		"timestamp":          e.Timestamp,
		"date":               e.Date,
		"title":              e.Title,
		"event_type":         e.Type,
		"user_source":        e.UserSource,
		"user_destination":   e.UserDestination,
		"domain_source":      e.UserSourceDomain,
		"domain_destination": e.UserDestinationDomain,
		"group":              e.GroupName,
		"group_domain":       e.GroupDomain,
		"process_source":     e.ProcessSource,
		"process_source_id":  e.ProcessSourceId,
		"process_target":     e.ProcessTarget,
		"process_target_id":  e.ProcessTargetId,
		"fullpath":           e.FullPath,
		"filename":           e.Filename,
		"extension":          e.Extension,
		"evidence":           e.Evidence,
		"computer":           e.Computer,
	}
	_, err := tx.Run(query, parameters)
	return nil, err
}

func InsertRegistriesNeo4j(con Neo4JConnector, registry []Registry) {
	for _, r := range registry {
		InsertRegistryNeo4j(con, r)
	}
}

func InsertRegistryNeo4j(con Neo4JConnector, r Registry) {
	sess := con.Driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})

	_, err := sess.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		return persistRegistry(tx, r)
	})
	handleErr(err)
}

func persistRegistry(tx neo4j.Transaction, r Registry) (interface{}, error) {
	query := "CREATE (:Registry {timestamp: $timestamp, date: $date, key: $key, value: $value, computer: $computer, evidence: $evidence})"
	parameters := map[string]interface{}{
		"timestamp": r.LastModifictationTimestamp,
		"date":      r.LastModificationTime,
		"key":       r.Path,
		"value":     r.Entries,
		"computer":  r.Computer,
		"evidence":  r.Evidence,
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

func linkGroup(con Neo4JConnector) {
	sess := con.Driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})

	// Link group to event to prepare for future query
	_, err := sess.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		query := `match (e:Event) where e.group <> ""
		match (g:Group) where e.group = g.name and e.group_domain = g.domain
		merge (e)-[:ABOUT]->(g)`
		var param map[string]interface{}
		return tx.Run(query, param)
	})

	handleErr(err)

	_, err = sess.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		query := `match (u:User)-->(e:Event)-->(g:Group) where e.event_type =~ "(?i).*enable.*"
		merge (u)-[:ENABLE{timestamp:e.timestamp, date: e.date}]->(g)`
		var param map[string]interface{}
		return tx.Run(query, param)
	})
	handleErr(err)

	_, err = sess.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		query := `match (u:User)-->(e:Event)-->(g:Group) where e.event_type =~ "(?i).*disable.*"
		merge (u)-[:DISABLE{timestamp:e.timestamp, date: e.date}]->(g)`
		var param map[string]interface{}
		return tx.Run(query, param)
	})
	handleErr(err)

	_, err = sess.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		query := `match (u:User)<--(e:Event)-->(g:Group) match (e)<--(source:User)
		where e.event_type =~ "(?i).*added.*"
		merge (u)-[:AddedTo{timestamp: e.timestamp, date:e.date}]->(g)`
		var param map[string]interface{}
		return tx.Run(query, param)
	})
	handleErr(err)

	_, err = sess.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		query := `match (u:User)<--(e:Event)-->(g:Group) match (e)<--(source:User)
		where e.event_type =~ "(?i).*removed.*"
		merge (u)-[:RemovedFrom{timestamp: e.timestamp, date:e.date}]->(g)`
		var param map[string]interface{}
		return tx.Run(query, param)
	})
	handleErr(err)

}

func linkUsers(con Neo4JConnector) {
	sess := con.Driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	_, err := sess.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		var param map[string]interface{}
		return tx.Run(`match (u:User) match (p) where u.fullname = p.user or u.username = p.user and p.user <> "" and u.fullname <> "-" merge (u)-[:BY]->(p)`, param)
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

func linkScriptBlock(con Neo4JConnector) {
	//create link based on pid, ppid and name. Quick Filter to avoid some duplicates
	sess := con.Driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	_, err := sess.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		query := `match (s:ScriptBlock) match(p:Process) where p.pid = s.process_id and p.timestamp > s.timestamp
		merge (p)-[:EXECUTE]->(s)`
		var param map[string]interface{}
		return tx.Run(query, param)
	})
	handleErr(err)

	//Remove duplicates
	_, err = sess.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		var param map[string]interface{}
		return tx.Run("match (m)-[r:EXECUTE]->(n:ScriptBlock)<-[s:EXECUTE]-(o) where n.timestamp - m.timestamp < n.timestamp - o.timestamp delete s", param)
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
		match (p:Process) where p.fullpath = c.process and p.pid = c.process_id and p.timestamp < c.timestamp and c.computer = p.computer
		match (h:Host) where h.ip = c.ip_destination 
		merge (p)-[r:CONNECT{port_source:c.port_source,port_destination:c.port_destination, ip_source:c.ip_source, timestamp:c.timestamp, date:c.date}]->(h)`

		parameters := map[string]interface{}{}
		_, err := tx.Run(query, parameters)
		return nil, err
	})
	handleErr(err)

}

func handleFileCreate(con Neo4JConnector) {
	sess := con.Driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	// Create File Based On Events "CreateFile" and "DeleteFile"
	_, err := sess.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		query := `match (e:Event) where e.event_type = "File Created" with collect(e) as events
		UNWIND events as event
		create (:File {fullpath: event.fullpath, filename: event.filename, extension:event.extension, 
		timestamp:event.timestamp, date:event.date, timestamp_desc:"Creation Time", computer:event.computer})`
		parameters := map[string]interface{}{}
		_, err := tx.Run(query, parameters)
		return nil, err
	})
	handleErr(err)

	fmt.Println("Link Processes to Files based on 'File Create' and 'File Delete' Events")
	_, err = sess.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		query := `match (e:Event) where e.event_type = "File Created"
		match (p:Process) where e.process_source = p.fullpath and e.process_source_id = p.pid and e.computer = p.computer
		match (f:File) where f.fullpath = e.fullpath and f.computer = e.computer
		merge (p)-[:CREATE]->(f)`
		parameters := map[string]interface{}{}
		_, err := tx.Run(query, parameters)
		return nil, err
	})
	handleErr(err)

}

func handleFileDelete(con Neo4JConnector) {
	sess := con.Driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})

	_, err := sess.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		query := `match (e:Event) where e.event_type = "File Deleted" with collect(e) as events
		UNWIND events as event
		create (:File {fullpath: event.fullpath, filename: event.filename, extension:event.extension, 
		timestamp:event.timestamp, date:event.date, timestamp_desc:"Deletion Time", computer:event.computer})`
		parameters := map[string]interface{}{}
		_, err := tx.Run(query, parameters)
		return nil, err
	})
	handleErr(err)

	_, err = sess.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		query := `match (e:Event) where e.event_type = "File Deleted"
		match (p:Process) where e.process_source = p.fullpath and e.process_source_id = p.pid and e.computer = p.computer
		match (f:File) where f.fullpath = e.fullpath and f.computer = e.computer
		merge (p)-[:DELETE]->(f)`
		parameters := map[string]interface{}{}
		_, err := tx.Run(query, parameters)
		return nil, err
	})
	handleErr(err)

}

func handleEventUsers(con Neo4JConnector) {
	sess := con.Driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	_, err := sess.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		query := `match (u:User) match (e:Event) where e.user_destination = u.fullname or e.user_destination = u.username and u.username <> ""
		merge (e)<-[:ACTS]-(u)`
		parameters := map[string]interface{}{}
		_, err := tx.Run(query, parameters)
		return nil, err
	})
	handleErr(err)

	_, err = sess.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		query := `match (u:User) match (e:Event) where e.user_source = u.fullname or e.user_source = u.username and u.username <> ""
		merge (e)-[:ON]->(u)`
		parameters := map[string]interface{}{}
		_, err := tx.Run(query, parameters)
		return nil, err
	})
	handleErr(err)
}

func handleRawAccessRead(con Neo4JConnector) {
	sess := con.Driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})

	_, err := sess.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		query := `match (e:Event) where e.event_type = "Raw Access Read"
		with collect(e) as events
		unwind events as event
		create (:File{fullpath: event.fullpath, filename: event.filename, computer: event.computer})`
		parameters := map[string]interface{}{}
		_, err := tx.Run(query, parameters)
		return nil, err
	})
	handleErr(err)

	// Link Process and File based on "Raw Access Read" Events
	_, err = sess.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		query := `match (e:Event) where e.event_type = "Raw Access Read"
		with collect(e) as events
		unwind events as event
		match (p:Process) where p.fullpath = event.process_source and p.pid = event.process_source_id and p.timestamp <= event.timestamp and p.computer = event.computer
		match (f:File) where f.fullpath = event.fullpath and f.computer = event.computer
		merge (p)-[:LOAD{timestamp:event.timestamp, date: event.date}]->(f)`
		parameters := map[string]interface{}{}
		_, err := tx.Run(query, parameters)
		return nil, err
	})

}

func handleMemoryAccess(con Neo4JConnector) {
	sess := con.Driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})

	_, err := sess.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		query := `match (event:Event) where event.event_type = "Process's Memory Access"
		with collect(event) as events
		unwind events as e
		match (source:Process) where source.fullpath = e.process_source and source.pid = e.process_source_id and source.timestamp <= e.timestamp and e.computer = source.computer
		match (target: Process) where target.fullpath = e.process_target and target.pid = e.process_target_id and target.timestamp <= e.timestamp and e.computer = target.computer
		merge (source)-[:MEMORY_ACCESS{timestamp: e.timestamp, date:e.date}]->(target)`
		parameters := map[string]interface{}{}
		_, err := tx.Run(query, parameters)
		return nil, err
	})
	handleErr(err)

	// Delete Duplicate (due to Pid collision on reboots)

	_, err = sess.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		query := `match (m)-[r:MEMORY_ACCESS]->(n)<-[s:MEMORY_ACCESS]-(o) 
		where n.timestamp - m.timestamp < n.timestamp - o.timestamp and m.pid = o.pid and m.fullpath=o.fullpath delete s`
		parameters := map[string]interface{}{}
		_, err := tx.Run(query, parameters)
		return nil, err
	})
	handleErr(err)

}

func handleImageLoaded(con Neo4JConnector) {
	sess := con.Driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})

	_, err := sess.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		query := `match (e:Event) where e.event_type = "Image Loaded"
		with collect(e) as events
		unwind events as event
		create (:File{fullpath: event.fullpath, filename: event.filename, computer: event.computer})`
		parameters := map[string]interface{}{}
		_, err := tx.Run(query, parameters)
		return nil, err
	})
	handleErr(err)

	// Link File and Process from Events "Image Loaded"

	_, err = sess.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		query := `match (e:Event) where e.event_type = "Image Loaded"
		with collect(e) as events
		unwind events as event
		match (p:Process) where p.fullpath = event.process_source and p.pid = event.process_source_id and p.computer = event.computer and p.timestamp <= event.timestamp
		match (f:File) where f.fullpath = event.fullpath and f.computer = event.computer
		merge (p)-[:LOAD{timestamp:event.timestamp, date: event.date}]->(f)`
		parameters := map[string]interface{}{}
		_, err := tx.Run(query, parameters)
		return nil, err
	})
	handleErr(err)

}

func handleDisableUserEvents(con Neo4JConnector) {
	sess := con.Driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})

	_, err := sess.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		query := `match (e:Event) where e.event_type = "User Account Disabled"
		with collect(e) as events
		unwind events as event
		match (actor:User)-[ACTS]->(event)-[:ON]->(target:User)
		merge (actor)-[:DISABLE{timestamp:event.timestamp, date: event.date}]->(target)`
		parameters := map[string]interface{}{}
		_, err := tx.Run(query, parameters)
		return nil, err
	})
	handleErr(err)

}

func handleEnableUserEvents(con Neo4JConnector) {
	sess := con.Driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})

	_, err := sess.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		query := `match (e:Event) where e.event_type = "User Account Enabled"
		with collect(e) as events
		unwind events as event
		match (actor:User)-[ACTS]->(event)-[:ON]->(target:User)
		merge (actor)-[:ENABLE{timestamp:event.timestamp, date: event.date}]->(target)`
		parameters := map[string]interface{}{}
		_, err := tx.Run(query, parameters)
		return nil, err
	})
	handleErr(err)
}

func handleDeleteUserEvents(con Neo4JConnector) {
	sess := con.Driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})

	_, err := sess.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		query := `match (e:Event) where e.event_type = "User Account Deleted"
		with collect(e) as events
		unwind events as event
		match (actor:User)-[ACTS]->(event)-[:ON]->(target:User)
		merge (actor)-[:DELETE{timestamp:event.timestamp, date: event.date}]->(target)`
		parameters := map[string]interface{}{}
		_, err := tx.Run(query, parameters)
		return nil, err
	})
	handleErr(err)
}

func handleCreateUserEvents(con Neo4JConnector) {
	sess := con.Driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})

	_, err := sess.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		query := `match (e:Event) where e.event_type = "User Account Created"
		with collect(e) as events
		unwind events as event
		match (actor:User)-[ACTS]->(event)-[:ON]->(target:User)
		merge (actor)-[:CREATE{timestamp:event.timestamp, date: event.date}]->(target)`
		parameters := map[string]interface{}{}
		_, err := tx.Run(query, parameters)
		return nil, err
	})
	handleErr(err)
}

func handleChangeUserEvents(con Neo4JConnector) {
	sess := con.Driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})

	_, err := sess.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		query := `match (e:Event) where e.event_type = "User Account Changed"
		with collect(e) as events
		unwind events as event
		match (actor:User)-[ACTS]->(event)-[:ON]->(target:User)
		merge (actor)-[:CHANGE{timestamp:event.timestamp, date: event.date}]->(target)`
		parameters := map[string]interface{}{}
		_, err := tx.Run(query, parameters)
		return nil, err
	})
	handleErr(err)
}

func handleLogonEvents(con Neo4JConnector) {
	sess := con.Driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})

	_, err := sess.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		query := `match (e:Event) where e.event_type = "Logon"
		with collect(e) as events
		unwind events as event
		match (actor:User)-[ACTS]->(event)-[:ON]->(target:User)
		where actor.fullname <> "-"
		merge (actor)-[:LOGON{timestamp:event.timestamp, date: event.date}]->(target)
		`
		parameters := map[string]interface{}{}
		_, err := tx.Run(query, parameters)
		return nil, err
	})

	handleErr(err)

	_, err = sess.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		query := `match (e:Event) where e.event_type = "Logon"
		with collect(e) as events
		unwind events as event
		match (event)-[:ON]->(target:User)
		match (computer:Computer) where computer.name = event.computer
		merge (target)-[:LOGON{timestamp:event.timestamp, date: event.date}]->(computer)
		`
		parameters := map[string]interface{}{}
		_, err := tx.Run(query, parameters)
		return nil, err
	})

	handleErr(err)
}

func handleLogoffEvents(con Neo4JConnector) {
	sess := con.Driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})

	_, err := sess.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		query := `match (e:Event) where e.event_type = "Logoff"
		with collect(e) as events
		unwind events as event
		match (event)-[:ON]->(target:User)
		match (computer:Computer) where computer.name = event.computer
		merge (target)-[:LOGOFF{timestamp:event.timestamp, date: event.date}]->(computer)
		`
		parameters := map[string]interface{}{}
		_, err := tx.Run(query, parameters)
		return nil, err
	})

	handleErr(err)
}

func handleEvents(con Neo4JConnector) {

	wg := sync.WaitGroup{}

	// Create and Link File with Process based on Events "File Create" and "File Delete"
	go func() {
		wg.Add(1)
		handleFileCreate(con)
		wg.Done()
	}()

	go func() {
		wg.Add(1)
		handleFileDelete(con)
		wg.Done()
	}()

	// Link Events to Users
	go func() {
		wg.Add(1)
		handleEventUsers(con)
		wg.Done()
	}()

	time.Sleep(3 * time.Second)
	wg.Wait()

	// Create File based on "RawAccessRead" Events
	go func() {
		wg.Add(1)
		handleRawAccessRead(con)
		wg.Done()

	}()
	// Link Process with Process based on "MemoryAccess" Events

	go func() {
		wg.Add(1)
		handleMemoryAccess(con)
		wg.Done()
	}()
	// Create File from Events "Image Loaded"
	go func() {
		wg.Add(1)
		handleImageLoaded(con)
		wg.Done()
	}()

	// Handle User -> User Events
	go func() {
		wg.Add(1)
		handleCreateUserEvents(con)
		wg.Done()
	}()

	time.Sleep(3 * time.Second)
	wg.Wait()

	go func() {
		wg.Add(1)
		handleDeleteUserEvents(con)
		wg.Done()
	}()

	go func() {
		wg.Add(1)
		handleEnableUserEvents(con)
		wg.Done()
	}()

	go func() {
		wg.Add(1)
		handleDisableUserEvents(con)
		wg.Done()
	}()

	time.Sleep(3 * time.Second)
	wg.Wait()

	// handle Logon Events

	go func() {
		wg.Add(1)
		handleLogonEvents(con)
		wg.Done()
	}()

	// handle Logoff Events

	go func() {
		wg.Add(1)
		handleLogoffEvents(con)
		wg.Done()
	}()

	// handle Change User Events
	go func() {
		wg.Add(1)
		handleChangeUserEvents(con)
		wg.Done()
	}()

	// Link Events to Groups
	go func() {
		wg.Add(1)
		linkGroup(con)
		wg.Done()
	}()

	time.Sleep(2 * time.Second)
	wg.Wait()
}

func updateIds(con Neo4JConnector) {
	sess := con.Driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})

	_, err := sess.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		query := `match (n) set n.objectid = id(n)`
		parameters := map[string]interface{}{}
		_, err := tx.Run(query, parameters)
		return nil, err
	})
	handleErr(err)
}
