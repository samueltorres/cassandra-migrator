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
		hosts         = fs.String("hosts", "localhost:9042", "list of cassandra hosts")
		keyspace      = fs.String("keyspace", "", "keyspace in which the scripts will be executed")
		dirs          = fs.String("dirs", "/db/scripts", "list of directories where cassandra scripts are located")
		username      = fs.String("username", "", "username used to be authenticated on cassandra")
		password      = fs.String("password", "", "password used to be authenticated on cassandra")
		cqlVersion    = fs.String("cqlversion", "3.4.4", "cql version used to run the scripts")
		retries       = fs.Int("retries", 20, "number of retries connecting to cassandra hosts")
		retryInterval = fs.Duration("retryInterval", time.Millisecond*1000, "interval in milliseconds for each retry")
	)
	ff.Parse(fs, os.Args[1:], ff.WithEnvVarPrefix("CM"))

	executor := CqlExecutor{
		Hosts:      *hosts,
		Keyspace:   *keyspace,
		Username:   *username,
		Password:   *password,
		CQLVersion: *cqlVersion,
	}

	ss, err := createSession(*hosts, *retries, *retryInterval)
	if err != nil {
		log.Fatal("Error creating session", err)
	}
	defer ss.Close()

	dd := strings.Split(*dirs, ":")
	for _, d := range dd {

		files, err := ioutil.ReadDir(d)
		if err != nil {
			log.Fatal(err)
		}

		var filesToExecute []string
		for _, file := range files {
			if strings.HasSuffix(file.Name(), ".cql") {

				absDir, err := filepath.Abs(d)
				if err != nil {
					log.Fatal("Invalid directory", d)
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

func createSession(addr string, retryCount int, retryInterval time.Duration) (*gocql.Session, error) {
	cluster := gocql.NewCluster(addr)
	var session *gocql.Session
	n := 1

	for n < retryCount {
		session, err := cluster.CreateSession()
		if err == nil {
			return session, err
		}

		fmt.Println("Host not available, retrying.")
		time.Sleep(retryInterval)
		n++
	}

	return session, errors.New("Could not could not create session")
}
