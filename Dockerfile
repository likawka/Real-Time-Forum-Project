FROM golang:latest
WORKDIR /app
COPY . .
WORKDIR /app/backend
RUN go mod download
EXPOSE 8080
CMD ["go", "run", "main.go"]