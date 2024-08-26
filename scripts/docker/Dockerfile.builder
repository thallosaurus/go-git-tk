FROM fedora:40 as builder
RUN dnf install -y golang @development-tools
WORKDIR /app
COPY . .
RUN make build
CMD ["/app/build/main"]