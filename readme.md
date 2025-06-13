# Banking Ledger v1.0

A simple Bank REST API that implements the Ledger Pattern. The interesting part here is the integration with
Keycloak for Authorization and ELK Stack for observability. For generating the IBAN for the accounts, the application implements
the MOD 97-10 algorithm that allows to compute the Domestic Control digits and the Control Digits of the IBAN itself. Also,
an IBAN verification method is implemented.

For testing purpose, the testcontainers library is being used for everything running in docker, as it provides a simple interface that makes it simple.

Regarding the API documentation, gin-swagger is being used to generate the docs and yaml files.

## Keycloak Configuration

The Keycloak service needs some little configuration in order to work along with the Ledger app. 

1. Go to Realm settings in the left lateral menu.
2. Press User Profile
3. Create the `client_id` attribute (BE SURE TO MAKE IT NOT REQUIRED).

![](doc/img/realm-settings.png)

4. The client_id user attribute needs to be added to the JWT payload. For this purpose, click on Client scopes and select profile.
5. Once in the `profile` scope click on the bottom arrow from the `Add mapper` button and select `by configuration`.
6. Create a mapper for the option `User Attribute`
7. Give it a name, then search it and click on it to the configuration.
8. The mapper should look like the following picture.

![](doc/img/keycloak-mapper.png)


## Authorization

The app uses Keycloak as an Authorization Service. Some environment variables [(env)](src/.env) involving Keycloak need to be set before using the Golang app.

```.env
# Keycloak
HOST=  YOUR HOST
ADMIN_USER=  YOUR ADMIN USER
ADMIN_PASSWORD=  YOUR ADMIN PASSWORD
CLIENT_ID=   YOUR CLIENT ID
CLIENT_SECRET= YOUR CLIENT SECRET
# USUALLY THE FOLLOWING ENDPOINTS ARE THE ONES BELOW
JWK_URL=/realms/ledger/protocol/openid-connect/certs
GET_TOKEN_URL=/realms/ledger/protocol/openid-connect/token
CREATE_USER_URL=/admin/realms/ledger/users
AUTH_SERVICES_DISCOVER_URL=/realms/ledger/.well-known/openid-configuration
```


 Keycloak makes use of the asymmetric RSA Encryption for the `access_token` payload. The Golang app implements RFC 7517 to use the Keycloak JWK (JSON Web Key) to decode the JWT. For further information visit [the Keycloack models.go file](src/api/keycloak/models.go).


## Keycloak Documentation
### Users
* [Users management documentation](https://www.keycloak.org/docs-api/latest/rest-api/index.html#_users)

* [Credential Representation](https://www.keycloak.org/docs-api/latest/rest-api/index.html#CredentialRepresentation)

* [User Representation](https://www.keycloak.org/docs-api/latest/rest-api/index.html#UserRepresentation)

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
