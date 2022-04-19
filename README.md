# Nocut API

URL shortener API.

## Development deployment

Development deployment consists of docker-compose, that runs 

- mongo
- mongo-express
- nocut-api-air

### Quick start

Clone the repository and cd to it
```bash
git clone https://github.com/furrygem/nocut-api.git
cd nocut-api
```

Create and populate files containing credentials for mongodb and mongo-express

```bash
touch mongo_root_password.txt
touch mongo_root_username.txt

echo <Secure Password> > mongo_root_password.txt
echo <Username> > mongo_root_username.txt
```

Start docker compose

```bash
docker-compose -f docker-compose-dev.yml up
```

If startup is successful, 3 services will be ran.

| Service       | Description                                                        | Exposed ports |
|---------------|--------------------------------------------------------------------|---------------|
| mongo         | MongoDB                                                            | 27017:27017   |
| mongo-express | Express-based web interface for mongodb                            | 8081:8081     |
| nocut-api-air | Air golang live application reloader, running the golang nocut API | 8080:8080     |


### Configuration

Development deployment is configured using ``.air.toml``, ``dev-config.yaml`` (or any other application configuration file) and ``docker-compose-dev.yml`` (or any other docker-compose file).

#### Application configuration file

Fields marekd with ``*`` do not have default value and require configuration.

```yaml
bind_addr: <IP address to host application on> # Default: 0.0.0.0
bind_port: <Port number to host application on> # Default: 8080
log_level: <Logrus-supported logging level> # Default: info
mongodb:
 host: <MongoDB host> # Default: 127.0.0.1
 port: <MongoDB port> # Default: 27017
 database: <Database name> # Default: nocut
 auth_db: <MongoDB authentication source> # Default: admin
 collection: <MongoDB collection name> # Default: links
 username_file: <Path to a file containing db username> # *
 password_file: <Path to a file containing db password> # *
 link_ttl: <Duration for link to persist in the system> # Default 3m
```

[Available Logging Levels](https://github.com/sirupsen/logrus#level-logging)
