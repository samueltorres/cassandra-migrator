CREATE KEYSPACE IF NOT EXISTS test
  WITH REPLICATION = {
   'class' : 'NetworkTopologyStrategy',
   'datacenter1' : 1
  };