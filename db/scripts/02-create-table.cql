CREATE TABLE if not exists test.counters (
	key text,
	value counter,
	PRIMARY KEY (key)
)
WITH gc_grace_seconds = 172800;