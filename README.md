## How to run
`
docker build -t infura-exporter .
`
```
docker run --rm -it -p 9877:9877 -e INFURA__URL=https://mainnet.infura.io/v3/ -e INFURA__TOKEN=abcd1234 axieinfinity/infura-exporter
or
docker run --rm -it -p 9877:9877 -e INFURA__TOKEN=abcd1234 axieinfinity/infura-exporter
```
