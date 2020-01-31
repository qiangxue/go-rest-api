version: '2.1'
services:
  server:
    image: server
    build:
      context: .
      dockerfile: cmd/server/Dockerfile
    volumes:
      - /tmp/app:/var/log/app
    ports:
      - "8080:8080"
    environment:
      - APP_ENV=local
      - APP_DSN=postgres://db/go_restful?sslmode=disable&user=postgres&password=postgres
    depends_on:
      db:
        condition: service_healthy
  db:
    image: "postgres:alpine"
    restart: always
    environment:
      POSTGRES_USER: "postgres"
      POSTGRES_PASSWORD: "postgres"
      POSTGRES_DB: "go_restful"
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 10s
      timeout: 5s
      retries: 5
