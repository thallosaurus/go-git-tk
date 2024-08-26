FROM golang:1.22.3 as builder

WORKDIR /app
COPY . .
RUN go build -o main cmd/main/main.go

FROM alpine
COPY --from=builder /app/main /bin/git-frontend
CMD ["/bin/git-frontend"]