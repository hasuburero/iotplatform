## api package

### worker

```
GET worker/contract
POST /worker
GET /worker
DELETE /worker
```

GET /worker/contract

```
application/json
{
  worker_id string
}

application/json
{
  worker_id string
  job_id string
  data1_id string
  data2_id string
  function_id string
  runtime string
}
```

POST /worker

```
application/plane ?
{
}

application/json
{
  worker_id string
}
```

GET /worker

```
application/json
{
  worker_id string
}

application/json
{
  worker_id string
  Runtime []string
}
```

### job

### data

```
POST /data/reg
GET /data
POST /data
PUT /data
DELETE /data
```
