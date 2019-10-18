package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/gocql/gocql"
	"github.com/peterbourgon/ff"
)

func main() {

	fs := flag.NewFlagSet("cassandra-migrator", flag.ExitOnError)
	var (
		hosts      = fs.String("hosts", "localhost:9042", "list of cassandra hosts")
		keyspace   = fs.String("keyspace", "", "the keyspace the scripts will be executed against")
		dirs       = fs.String("dirs", "/db/scripts", "list of directories where cassandra scripts are located")
		username   = fs.String("username", "", "username used to be authenticated on cassandra")
		password   = fs.String("password", "", "password used to be authenticated on cassandra")
		cqlVersion = fs.String("cqlversion", "3.4.4", "cql version used to run the scripts")
	)
	ff.Parse(fs, os.Args[1:], ff.WithEnvVarPrefix("CQLM"))

	executor := CqlExecutor{
		Hosts:      *hosts,
		Keyspace:   *keyspace,
		Username:   *username,
		Password:   *password,
		CQLVersion: *cqlVersion,
	}

	ss, err := createSession(*hosts)
	if err != nil {
		fmt.Println("error creating session", err)
		return
	}
	defer ss.Close()

	dd := strings.Split(*dirs, ":")
	for _, dir := range dd {

		files, err := ioutil.ReadDir(dir)
		if err != nil {
			log.Fatal(err)
		}

		var filesToExecute []string
		for _, file := range files {
			if strings.HasSuffix(file.Name(), ".cql") {

				absDir, err := filepath.Abs(dir)
				if err != nil {
					log.Fatal("Invalid Directory", dir)
				}

				filesToExecute = append(filesToExecute, path.Join(absDir, file.Name()))
			}
		}
		sort.Strings(filesToExecute)

		for _, file := range filesToExecute {
			err = executor.Execute(file)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func createSession(addr string) (*gocql.Session, error) {
	cluster := gocql.NewCluster(addr)
	var session *gocql.Session
	attempts := 20
	n := 1

	for n < attempts {
		session, err := cluster.CreateSession()
		if err == nil {
			return session, err
		}

		fmt.Println("retrying one more time")
		time.Sleep(1 * time.Second)
		n++
	}

	return session, errors.New("could not create session")
}
