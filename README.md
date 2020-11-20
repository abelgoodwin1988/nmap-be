# nmap-be
An API for scanning open ports given addresses. The only endpoint is `/portscan` which must have a body like `'{"addresses":"google.com","127.0.0.1"}'`. The addresses provided must be an ip, or a fqdn.

An example curl to the service could be
```bash
curl -H "Content-Type: application/json" \
 --request GET \
 --data '{"addresses":"google.com","127.0.0.1"}' \
 localhost:8080/portscan
 ```

# Quickstart

### Requirements

- Must have one of the last two minor versions of go installed on your machine.
- Must have docker installed.

## Instructions

- Navigate to the base directory of the project, and start a local mysql container by running
```bash
docker-compose -f deployments/docker-compose.local.yml up
```
- This will start a simple mysql container that has applied the schema found in `/build/database/schema.sql`. There is no seeded data. However, subsequently stopping and starting the container will have persistent data so long as the volume speicified in `/deployements/docker-compose.local.yml` is not removed.
- Start the api by running `go run main.go`
- Validate that the service is functioning by opening another terminal and running the following:
```bash
curlie -H "Content-Type: application/json" \
 --request GET \
 --data '{"addresses":"google.com","127.0.0.1"}' \
 localhost:8080
```
 