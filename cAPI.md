# sample requester api
```
#make lambda
func_blob_id = POST /data

func_id = POST /lambda

#make job
input_id = POST /data

job_id = POST /job input=input_id

loop{
    status = GET /job/$job_id

    if status = finished
        break
}
```

## func_blob_id = POST /data
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
    "id": "78204124132"
    "hash": "78204124132"
}
```

## func_id = POST lambda
```
request body
content-type = application/json
{
    "codex": "78204124132",
    "runtime": "78204124132"
}

response body
content-type = application/json
{
    "code": 0,
    "status": "ok",
    "id": "78204124132"
}
```

## POST job
```
request body
content-type = application/json; charset=utf-8
{
    "input": "78204124132", //input data id
    "lambda": "78204124132", //lambda id
    "tags": [
        "+a"
    ]
}

response body
content-type = application/json; charset=utf-8
{
    "code": "78204124132"
    "status": "ok"
    "message": "job found"
    "id": "78204124132"
}
```

