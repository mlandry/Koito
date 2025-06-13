---
title: Installation
description: Guide on how to install Koito to start tracking your listening history.
---

## Docker
By far the easiest way to get up and running with Koito is using docker. Here is an example Docker Compose file to get you up and running in minutes:
```yaml title="compose.yaml"
services:
  koito:
    image: gabehf/koito:latest
    container_name: koito
    depends_on:
      - db
    user: 1000:1000
    environment:
      - KOITO_DATABASE_URL=postgres://postgres:secret_password@db:5432/koitodb
      - KOITO_ALLOWED_HOSTS=koito.example.com,192.168.0.100
    ports:
      - "4110:4110"
    volumes:
      - ./koito-data:/etc/koito
    restart: unless-stopped

  db:
    image: postgres:16
    container_name: psql
    restart: unless-stopped
    environment:
      POSTGRES_DB: koitodb
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: secret_password
    volumes:
      - ./db-data:/var/lib/postgresql/data

```
Or, you can use `docker run` commands. First, you need to run Postgres:
```sh
docker run \
    --name psql \
    --restart unless-stopped \
    -e POSTGRES_DB=koitodb \
    -e POSTGRES_USER=postgres \
    -e POSTGRES_PASSWORD=secret_password \
    -v ./db-data:/var/lib/postgresql/data \
    -p 5432:5432 \
    -d \
    postgres:16
```
Then, run Koito:
```sh
docker run \
    --name koito \
    -u 1000:1000 \
    -e KOITO_DATABASE_URL=postgres://postgres:secret_password@postgres_ip:5432/koitodb \
    -e KOITO_ALLOWED_HOSTS=koito.example.com,192.168.0.100 \
    -p 4110:4110 \
    -v ./koito-data:/etc/koito \
    --restart unless-stopped \
    -d \
    gabehf/koito:latest
```
Be sure to replace `secret_password` with a random password of your choice, and set `KOITO_ALLOWED_HOSTS` to include the domain name or IP address you will be accessing Koito 
from when using either of the Docker methods described above.

Those are the two required environment variables. You can find a full list of configuration options in the [configuration reference](/reference/configuration).

When using `docker run`, you will also need to fill in the IP address of your postgres instance.

:::caution
Setting `KOITO_ALLOWED_HOSTS=*` will allow requests from any host, but this is not recommended as it introduces security vulnerabilities.
:::

## Build from source

If you don't want to use docker, you can also build the application from source.

First, you need to install dependencies. Koito relies on `libvips-dev` to process images.

```sh
sudo apt install libvips-dev
```

If you aren't installing on an Ubuntu or Debian based system, you can find other ways to install `libvips-dev` on the [libvips wiki](https://github.com/libvips/libvips/wiki/)

Then, clone the repository and execute the build command using the included Makefile:

```sh
git clone https://github.com/gabehf/koito && cd koito
make build
```

When the build is finished, you can run the executable at the root of the directory. You'll also need to defined the required environment variables.

```sh
KOITO_DATABASE_URL=postgres://postgres:secret_password@postgres_ip:5432/koitodb \
KOITO_ALLOWED_HOSTS=koito.example.com,192.168.0.100 \
./koito
```

Then, navigate your browser to `localhost:4110` to enter your Koito instance.

:::note
You will need to provide your own Postgres instance. You can find downloads to Postgres [here](https://www.postgresql.org/download/).
:::