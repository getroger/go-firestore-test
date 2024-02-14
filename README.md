# go-firestore-test

Experimenting with Google Firestore in Go.

## Compile
```bash
CGO_ENABLED=0 go build -o firestore-test .
```

## Usage
```
Usage of ./firestore-test:
  -collection string
    	collection name (default "tasks")
  -db string
    	database name
  -get-by-ref string
    	get document by ref
  -project string
    	project id
  -refs-only
    	print document refs instead
  -timeout int
    	timeout in seconds (default 30)
```

## Run
All examples use the optional `--collection` flag to specify the collection to query, which defaults to `tasks`.


### Get all ref IDs in a collection
```bash
./firestore-test --refs-only --db '(default)' --project live-332912
```

Example output:
```
....
zy0l4GhLcGT32XAn6ZEf
zzFlCcHYTb0EAmH1P1GZ
zzMV283rqYyQ7eRaIXzW

-------------------------
total documents processed: 6387
last document ref: zzMV283rqYyQ7eRaIXzW
elapsed time: 5.328804541s
-------------------------
```

### Get all documents in a collection
```bash
./firestore-test --db '(default)' --project live-332912
```

Example output:

```bash
u2AjaSPpQ3JeFMJRykdu
u3HJQ5VTHJJHSo6Q2mCU
u3K3tO1YCo3ZIAqV74rg
u4i1z1RG3VbJ69EiJx72
error getting document: rpc error: code = DeadlineExceeded desc = context deadline exceeded


-------------------------
total documents processed: 5738
last document ref: u4i1z1RG3VbJ69EiJx72
elapsed time: 29.99822875s
```

### Get a single document by ID
```bash
./firestore-test --get-by-ref u4yRw3p8tMxOqobJeS6K --db '(default)' --project live-332912 --timeout 5
```

Example output:
```
error getting document by ref: rpc error: code = DeadlineExceeded desc = context deadline exceeded
```