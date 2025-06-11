FROM node AS frontend

WORKDIR /client
COPY ./client/package.json ./client/yarn.lock ./
RUN yarn install
COPY ./client .
ENV BUILD_TARGET=docker
RUN yarn run build


FROM golang:1.23 AS backend

WORKDIR /app

RUN apt-get update && \
	apt-get install -y libvips-dev pkg-config && \
	rm -rf /var/lib/apt/lists/*

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -o app ./cmd/api


FROM debian:bookworm-slim AS final

WORKDIR /app

RUN apt-get update && \
	apt-get install -y libvips42 && \
	rm -rf /var/lib/apt/lists/*

COPY --from=backend /app/app ./app
COPY --from=frontend /client/build ./client/build
COPY ./client/public ./client/public
COPY ./assets ./assets
COPY ./db ./db

EXPOSE 4110

ENTRYPOINT ["./app"]
