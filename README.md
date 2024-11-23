# klmna

### Local development

```
go run main.go
```

### Deploy to EC2

[YC Console](https://console.yandex.cloud)

```
docker build --platform linux/amd64 -t klmna:0.2.1 .
```

[YC docker login](https://yandex.cloud/ru/docs/container-registry/operations/authentication)

```
docker tag klmna:0.2.1 cr.yandex/<container_registry_id>/klmna:0.2.1

docker push cr.yandex/<container_registry_id>/klmna:0.2.1
```

create VM, connent via SSH, install Docker, [login to YC container registry](https://yandex.cloud/ru/docs/container-registry/tutorials/run-docker-on-vm/console#run)

```
ssh -i <path_to_ssh_key>/<ssh_key_name> <vm_user_name>@<vm_public_ip>

sudo docker pull cr.yandex/<container_registry_id>/klmna:0.2.1

sudo docker run -it --rm -e POSTGRES_DB=$POSTGRES_DB \
    -e POSTGRES_USER=$POSTGRES_USER \
    -e POSTGRES_PASSWORD=$POSTGRES_PASSWORD \
    -e POSTGRES_HOST=$POSTGRES_HOST \
    -e POSTGRES_PORT=$POSTGRES_PORT \
    -p 80:80 cr.yandex/<container_registry_id>/klmna:0.2.1
```

### Local docker run


```
docker build -t klmna:0.2.1 .

export POSTGRES_DB=<DB_nam>
export POSTGRES_USER=<DB_user>
POSTGRES_PASSWORD=<DB_password>
export POSTGRES_HOST=<DB_IP>
export POSTGRES_PORT=<DB_port>

docker run -it --rm -e POSTGRES_DB=$POSTGRES_DB \
    -e POSTGRES_USER=$POSTGRES_USER \
    -e POSTGRES_PASSWORD=$POSTGRES_PASSWORD \
    -e POSTGRES_HOST=$POSTGRES_HOST \
    -e POSTGRES_PORT=$POSTGRES_PORT \
    -p 80:80 klmna:0.2.5

docker ps

docker stop <container_hash>
```

docker run -it --rm -e POSTGRES_DB=$POSTGRES_DB \
    -e POSTGRES_USER=$POSTGRES_USER \
    -e POSTGRES_PASSWORD=$POSTGRES_PASSWORD \
    -e POSTGRES_HOST=$POSTGRES_HOST \
    -e POSTGRES_PORT=$POSTGRES_PORT \
    -p 80:80 klmna:0.2.1

### Migrations

Migrations in folder `internal/db/migrations` and run at app start

### Local DB run

```
docker-compose up

psql -d "host=localhost port=5432 dbname=klmna-db user=klmna-user"
```

remove DB volumes

```
docker-compose down --volumes
```

### Local development

- Rename `example.env` to `.env` or provide own environmental variable.
- Run DB
- `go run main.go`

### Keep app

- formatted `gofmt -s -w .`
- tidy `go mod tidy`
