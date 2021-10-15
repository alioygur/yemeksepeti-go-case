# yemeksepeti-go-case

yemeksepeti golang case study

## Build & Run

```
# run tests
make test

# run dev
make run-dev

# build
make build

# build docker image
make build-docker

# run
make run
```

## Usage

**Store a key**
```
curl --location --request POST 'localhost:8080/set' \
--header 'Content-Type: application/json' \
--data-raw '{
    "key": "adi",
    "value": "veli"
}'
```

**Get a key**
```
curl --location --request GET 'localhost:8080/get?key=adi'
```