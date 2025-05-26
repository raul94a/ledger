# Banking Ledger v1.0

A simple Bank REST API that implements the Ledger Pattern. The interesting part here is the integration with
Keycloak for Authorization and ELK Stack for observability. For generating the IBAN for the accounts, the application implements
the MOD 97-10 algorithm that allows to compute the Domestic Control digits and the Control Digits of the IBAN itself. Also,
an IBAN verification method is implemented.

For testing purpose, the testcontainers library is being used for everything running in docker, as it provides a simple interface that makes it simple.

Regarding the API documentation, gin-swagger is being used to generate the docs and yaml files.

## Keycloak Documentation
### Users
[Users management documentation](https://www.keycloak.org/docs-api/latest/rest-api/index.html#_users)
[Credential Representation](https://www.keycloak.org/docs-api/latest/rest-api/index.html#CredentialRepresentation)
[User Representation](https://www.keycloak.org/docs-api/latest/rest-api/index.html#UserRepresentation)

## Tech Stack

- Docker and Docker Compose
- Postgres 17
- Golang 1.21.6
- gin v1.10.1
- TestContainers
- Swagger/OpenAPI
- ELK Stack
- Keycloak




## Tips
### Using migrate:

- USER: Your db user
- PASSWORD: Your db password
- HOST_IP: The IP where the database is running in DEV, usually localhost
- PORT: The port where Postgres is running, usually 5432
- DATABASE: The name of your database

```sh
migrate -database "postgres://USER:PASSWORD@HOST_IP:PORT/DATABASE?sslmode=disable" -path ./migrations up
```

### Gin tutorial

[Mastering backend Gin tutorial](https://masteringbackend.com/posts/gin-framework#the-framework)
