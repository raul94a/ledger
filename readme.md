# Banking Ledger




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
