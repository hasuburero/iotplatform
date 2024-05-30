# sample worker api
```
worker_id = POST /worker

loop {

    job_id = POST /worker/$worker_id/contract

    if job_id not exist
        continue

    job_context = GET /job/$job_id

    input_data = GET /data/$job_context.input_data/blob

    exec ./bin < cat input_data > output_data

    data_id = POST /data output_data

    POST /job/$job_id set status=finished,output_data=data_id

}
```

## worker_id = POST /worker
```
request body
content-type = application/json
{
    "runtime": [
        "default",
        "dlabLua54"
    ]
}

response body
content-type = application/json
{
    "code": 0,
    "status": "ok",
    "message": "?",
    "id": "78204124132",
    "runtime": [
        "???",
        "???"
    ]
}
```

## job_id = POST /worker/$worker_id/contract
```
request body
content-type = application/json
{
    "id": "78204124132",
    "tags": [],
    "timeout": 20
}

response body
content-type = application/json
{
    "code": 0,
    "status": "ok",
    "message": "no job"
    "id": "78204124132"
}
```

## job_context = GET /job/$job_id
```
request body
content-type = application/json
{
    "job_id": "78204124132"
}

response body
content-type = application/json
{
    "code": 0,
    "status": "ok",
    "message": "???",
    "id": "78204124132",
    "job_status": "Job found",
    "job_input_id": "78204124132"
    "job_output_id": "12345678900"
    "functio": "78204124132"
    "runtime": "runtime type"

    "tags": [
        "job tag"
    ]
    "lambda": {
        "id": "78204124132",
        "codex": "78204124132",
        "runtime": "default"
    },
    "input": {
        "id": "78204124132"
    },
    "output": {
        "id": "78204124132"
    },
    "state": "Running"
}
```

## input_data = GET /data/$job_context.input_data/blob
## input_data = GET /data/{dID}/blob
```
response body

if correct request
content-type = application/octet-stream    
{
    "file content goes there"
    "octet-stream means any file (binary or not)"
}
elif not found target file -> Status code == 404
content-type = application/json
{
    "code": 601,
    "status": "error",
    "message": "Object Not Found"
}
```

## data_id = POST /data
```
request body
content-type = multipart/form-data
{
    "file content goes there"    
    "multipart/form-data means any file (binary or not)"
}

response body
content-type = application/json
{
    "code": 0,
    "status": "ok",
    "message": "http: no such file",
    "id": "78204124132",
    "checksum": "78204124132"
}
```

## POST /job/$job_id
```
request body
content-type = application/json
{
    "state": "finished",
    "output": "78204124132"

    # extra scheme for next job like pipelines
    "input": "78204124132"
}

response body
content-type = application/json
{
    "code": 0,
    "status": "ok",
    "message": "Job found",
}
```
