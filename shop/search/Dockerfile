FROM golang as builder
RUN mkdir /build 
ADD . /build/
WORKDIR /build 
ENV CGO_ENABLED=0 
ENV GOOS=linux
RUN go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o main main.go
# Runtime Image
FROM scratch
ENV GIN_MODE=release
COPY --from=builder /build/main /app/
COPY config.yaml /app/
COPY products.json /app/
WORKDIR /app
CMD ["./main"]