# cassandra-migrator
Cassandra migrator is a tool that automates the execution of CQL scripts on a cassandra server.

# Requirements

This tool needs cqlsh (cassandra shell) to do execution of the scritps. To install it you'll need to install pip.

## pip
```console
sudo apt install python-pip
```

## cqlsh
```console
pip install cqlsh
```

# Installation

```console
go get samueltorres/cassandra-migrator
```

# Usage

```
Usage of cassandra-migrator:

  cassandra-migrator [options]

  -cqlversion string
        cql version used to run the scripts (default "3.4.4")
  
  -dirs string
        list of directories where cassandra scripts are located (default "/db/scripts")
  
  -hosts string
        list of cassandra hosts (default "localhost:9042")
  
  -keyspace string
        keyspace in which the scripts will be executed
  
  -password string
        password used to be authenticated on cassandra
  
  -retries int
        number of retries connecting to cassandra hosts (default 20)
  
  -retryInterval duration
        interval in milliseconds for each retry (default 1µs)
  
  -username string
        username used to be authenticated on cassandra
```

Instead of passing flags you can use enviroment variables prefixed with `CM_`.

### Example

```console
cassandra-migrator --hosts=localhost:9042 --dirs=/db/admin_scripts:/db/scripts --username=cassandra --password=cassandra
```

# Usage with Docker

```console
docker run -v ./db/scripts:/db/scripts -e CM_HOSTS=localhost:9042 samueltorres/cassandra-migrator
```

# License

This project is licensed under the MIT open source license.