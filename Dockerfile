FROM golang:1.19 as builder

# uncomment the line below to use goproxy.io
# RUN go env -w GOPROXY=https://goproxy.io,direct

WORKDIR /app
COPY . .
RUN make all

FROM scratch as exporter
COPY --from=builder /app/build/* .
