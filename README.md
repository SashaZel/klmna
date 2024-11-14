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

sudo docker run -d --rm -p 80:80 cr.yandex/<container_registry_id>/klmna:0.2.1
```

### Local docker run

```
docker build -t klmna:0.2.1 .

docker run -it --rm -p 80:80 klmna:0.1.13

docker ps

docker stop <container_hash>
```

### Local DB run

```
docker-compose up

psql -d "host=localhost port=5432 dbname=klmna-db user=klmna-user"
```

remove DB volumes

```
docker-compose down --volumes
```