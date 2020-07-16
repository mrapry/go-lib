# go-lib

Go libraries mrapry

#### Install
```shell
$ go get github.com/mrapry/go-lib
```

#### Basic environment example
```
# Basic env configuration
ENVIRONMENT=[string]

# Service Handlers
## Server
USE_REST=[bool]
USE_GRPC=[bool]
USE_GRAPHQL=[bool]
## Worker
USE_KAFKA_CONSUMER=[bool]
USE_CRON_SCHEDULER=[bool]
USE_REDIS_SUBSCRIBER=[bool]

REST_HTTP_PORT=[int]
GRAPHQL_HTTP_PORT=[int]
GRPC_PORT=[int]

BASIC_AUTH_USERNAME=[string]
BASIC_AUTH_PASS=[string]

# optional if using mongo database
MONGODB_HOST_WRITE=[string]
MONGODB_HOST_READ=[string]
MONGODB_DATABASE_NAME=[string]

# optional if using sql database
SQL_DRIVER_NAME=[string]
SQL_DB_READ_HOST=[string]
SQL_DB_READ_USER=[string]
SQL_DB_READ_PASSWORD=[string]
SQL_DB_WRITE_HOST=[string]
SQL_DB_WRITE_USER=[string]
SQL_DB_WRITE_PASSWORD=[string]
SQL_DATABASE_NAME=[string]

# optional if using redis
REDIS_READ_HOST=[string]
REDIS_READ_PORT=[string]
REDIS_READ_AUTH=[string]
REDIS_WRITE_HOST=[string]
REDIS_WRITE_PORT=[string]
REDIS_WRITE_AUTH=[string]

KAFKA_BROKERS=[string],[string]
KAFKA_CLIENT_ID=[string]
KAFKA_CONSUMER_GROUP=[string]

JAEGER_TRACING_HOST=[string]
GRAPHQL_SCHEMA_DIR=[string]
JSON_SCHEMA_DIR=[string]
```
