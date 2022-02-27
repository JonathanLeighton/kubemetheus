
FROM golang:1.16

COPY . /opt/kubmetheus

WORKDIR /opt/kubmetheus
RUN go build -o out/kubmetheus .


FROM debian:bullseye

WORKDIR /opt
COPY --from=0 /opt/kubmetheus/out/kubmetheus ./
CMD ["./kubmetheus"]
