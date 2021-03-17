# HTTP in-memory key-value storage  

## Run server

While starting, you can use the following keys:

```
Usage 
  -addr string
        HTTP network address (default ":5555")
  -allowEmptyDBOnStart
        allow empty db if backup file didn't find (default true)
  -buInterval int
        backup interval in milliseconds (default 1000)
  -file string
        path to backup file (default "./db.json")
```

For example:

```bash
$ go run ./cmd -file=./db.json -addr=:5550 -buInterval=1000 -allowEmptyDBOnStart=false
```


## API examples

Service responds with JSON data.

### List

List current storage state.

```bash
$ curl -X GET http://localhost:5550/list
{
 "data": {
  "key1": "val1",
  "key11": "1aeoua",
  "key2": "bababa",
  "key3": "val3",
  "key4": "val4",
  "key5": "",
 },
 "err": ""
}
```

### GET

Request the value by the key "key7".

```bash
$ curl -X GET http://localhost:5550/get/key7
{
 "data": "val7",
 "err": ""
}
```

Getting a value for a nonexistent key returns an error.

```bash
$ curl -X GET http://localhost:5550/get/nonexistentKey
{
 "data": null,
 "err": "storage: there is no such key in the storage"
}
```


### UPSERT

Accepts both **PUT** and **POST** HTTP verbs. Returns upserted records.

```bash
$ curl -X POST 'http://localhost:5550/upsert?key6=val6&key7=val7&key8=val8'
{
 "data": {
  "key6": "val6",
  "key7": "val7",
  "key8": "val8"
 },
 "err": ""
}
```

### DELETE

Accepts both **DELETE** and **POST** HTTP verbs.

```bash
$ curl -X DELETE http://localhost:5550/delete/key5
{
 "data": "",
 "err": ""
}
$ curl -X POST http://localhost:5550/delete/key6
{
 "data": "val6",
 "err": ""
}
```

Trying to delete a value for a nonexistent key returns an error.

```bash
$ curl -X DELETE http://localhost:5550/delete/nonexistentKey
{
 "data": null,
 "err": "storage: there is no such key in the storage"
}
```
