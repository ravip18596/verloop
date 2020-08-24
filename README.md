Install Cassandra
-----------------
This installs Apache Cassandra:

```Shell
brew install cassandra
```

Starting/Stopping Cassandra
---------------------------
To have launchd start cassandra now and restart at login:

```Shell
brew services start cassandra
```

Install cqlsh
```Shell
brew install cql
```

Setting Up Tables
----------

```Shell
echo "CREATE KEYSPACE stories WITH replication = {'class': 'SimpleStrategy', 'replication_factor': '1'}  AND durable_writes = true;" | cqlsh;

echo "create table stories.story(story_id bigint, created_at timestamp, updated_at timestamp, title List<text>, primary key(story_id));" | cqlsh;

echo "create table stories.paragraph(story_id bigint, paragraph_id bigint,sentence_id int sentences List<text>, primary key(story_id,paragraph_id,sentence_id));" | cqlsh;
```
Starting Redis Server
------------
```Shell
$ wget http://download.redis.io/releases/redis-6.0.6.tar.gz
$ tar xzf redis-6.0.6.tar.gz
$ cd redis-6.0.6
$ make
```
 The binaries that are now compiled are available in the src directory. Run Redis with:
```Shell
$ src/redis-server
```
It is running at port 6379 - according to defualt config
You can interact with Redis using the built-in client:
```Shell
$ src/redis-cli
redis> set foo bar
OK
redis> get foo
"bar"
```

Running Web Server
------------
running at port 8050
```Shell
go get -v ./...
export GOPATH=<verloop directory location>
export VERLOOP_DEBUG=true
export VERLOOP_DSN=localhost
go run main.go 
or 
go build -o verloop && ./verloop
```

Changes Pending
---------

```text
- K6 Load Test with CPU and memory profiler running
```