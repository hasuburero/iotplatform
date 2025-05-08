# api package

## worker

```
POST /worker
GET /worker
DELETE /worker
GET /worker/contract
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
X-Worker-Id: xxxxxxxx // 8 bytes?
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

### DELETE /worker

request

```
X-Worker-Id: xxxxxxxx // 8bytes?
application/json

{
  worker_id string
}
```

response

```
no content
```

### GET /worker/contract

request

```
X-Worker-Id: xxxxxxxx // 8 bytes?
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

## Job

```
GET /job
POST /job
PUT /job
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
  status string
}
```

### POST /job

request

```
application/json

{
  data1_id string
  data2_id string
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

### PUT /job

request

```
X-Job-Id: xxxxxxxx // job_id, 8 bytes?
application/json

{
  job_id string
  data1_id string
  data2_id string
  function_id string
  runtime string
  status string
}
```

response

```
no content
```

### DELETE /job

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
GET /data
POST /data
PUT /data
DELETE /data
POST /data/reg
```

### GET /data

request

```
X-Data-Id: xxxxxxxx // data_id, 8 bytes?

no content
```

response

```
multipart/form-data; // fieldnameは特に規定していない．デフォルトは"file"とする

content
```

### POST /data

request

```
multipart/form-data; // fieldnameは特に規定していない．デフォルトは"file"とする

data // ファイル名は規定なし．デフォルトは"file", サーバ上ではdata_idをファイル名とする
```

response

```
text/plain

xxxxxxxx // data_id, 8 bytes?
```

### PUT /data

request

```
X-Data-Id: xxxxxxxx // data_id, 8 bytes?
multipart/form-data;

data // ファイル名は規定なし，デフォルトは"file", サーバ上ではdata_idを　ファイル名とする
```

response

```
text/plain

xxxxxxxx // data_id, 8 bytes?
```

### DELETE /data

request

```
X-Data-Id: xxxxxxxx // data_id, 8 bytes?

no content
```

response

```
no content

no content
```

### POST /data/reg

request

```
no content

no content
```

response

```
text/plain

xxxxxxxx // data_id, 8 bytes?
```
