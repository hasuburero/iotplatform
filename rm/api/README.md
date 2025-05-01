# api package

## worker

```
GET /worker/contract
POST /worker
GET /worker
DELETE /worker
```

### GET /worker/contract

request

```
X-Worker_Id: xxxxxxxx // 8 bytes?
application/json

{
  worker_id string
}
```

response

```
application/json

{
  worker_id string
  job_id string
  data1_id string // データ受信用
  data2_id string // データ送信用，requester, worker共に設定可能．すでに設定してあればそこに送信
  function_id string
  runtime string
}
```

### POST /worker

request

```
application/json

{
  Runtime []string
}
```

response

```
application/json

{
  worker_id string
}
```

### GET /worker

request

```
application/json

{
  worker_id string
}

```

response

```
application/json

{
  worker_id string
  Runtime []string
}
```

## Job

```
GET /job
POST /job
DELETE /job
```

### GET /job

request

```
X-Job-Id: xxxxxxxx // 8 bytes?

no contents
```

response

```
application/json

{
  job_id string
  data1_id string
  data2_id string
  function_id string
  runtime string
}
```

### POST /job

request

```
application/json

{
  data_id string
  function_id string
  runtime string
}
```

response

```
application/json

{
  job_id string
}
```

## DELETE /job

request

```
X-Job-Id: xxxxxxxx // 8 bytes?

no content
```

response

```
no content
```

## Data

```
POST /data/reg
GET /data
POST /data
PUT /data
DELETE /data
```
