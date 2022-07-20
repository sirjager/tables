# Build stage
FROM golang:1.18.4-alpine3.16 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

# Run stage
FROM alpine:3.16.1
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder ./app/scripts/start.sh .
COPY --from=builder ./app/scripts/wait-for.sh .
COPY --from=builder ./app/migrations ./migrations
RUN chmod +x /app/main /app/migrate /app/start.sh /app/wait-for.sh 


EXPOSE 8080
CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]
