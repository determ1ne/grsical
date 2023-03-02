# grsical

Parse ZJU grs class timetable and generate iCalender file for you.

Undergraduate(ugrs) version: [cxz66666/ugrsical](https://github.com/cxz66666/ugrsical)

本科生版： [cxz66666/ugrsical](https://github.com/cxz66666/ugrsical)


## How to compile

### Using Docker

```
make docker-all
```

### Manually

Built and tested using go 1.19

```
make all
```

## As a client

- Copy `configs/upfile.json.example` to `configs/upfile.json` and fill it
- `./grsical -i ./configs/upfile.json -c ./configs/config.json -t ./configs/tweaks.json`

## As a server

**Environment variables:**

- GRSICALSRV_CFG: config file url
- GRSICALSRV_TWEAKS: tweaks config url
- GRSICALSRV_ENCKEY: AES encryption key
- GRSICALSRV_HOST: server hostname
- GRSICALSRV_IP_HEADER: from which IP should be read
- GRSICALSRV_REDIS_ADDR: redis server address
- GRSICALSRV_REDIS_PASS: redis server password

Run `./grsicalsrv`