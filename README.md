## Maker checker Service

### Prerequisites

To launch the entire service run the command from **dev** folder

```
docker-compose up
```

After launched you can make all requests

### Run Unit Tests

```
 go test ./... -v
```

### Example Requests / Responses

#### Create message

```
    curl -X POST http://localhost:9090/messages \
     -H "Content-Type: application/json" \
     -d '{
            "senderId" : "1",
            "recipientId" : "2",
            "content" : "Hi"
          }'
```

Get a response:

```
{
    "content": "Hi",
    "id": "ce222716-edf5-4d2e-803a-70dea8758ae4",
    "recipientId": "2",
    "senderId": "1",
    "status": "PENDING"
}
```

#### For the secured endpoints we use `Authorization` header and for simplicity i hardcoded it in ENV files. For prod obviously we use JWT generators
#### Get Messages

```
curl -X GET http://localhost:9090/messages?status=all \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer alpheya"
```

Get a response:

```
{
    "messages": [
        {
            "Id": "ce222716-edf5-4d2e-803a-70dea8758ae4",
            "SenderId": "1",
            "RecipientId": "2",
            "Content": "Hi",
            "Status": "PENDING",
            "Ts": "2025-01-10T10:08:58.193608Z"
        }
    ]
}

```

#### Update Messages

```
curl -X PUT http://localhost:9090/messages/ce222716-edf5-4d2e-803a-70dea8758ae4/reject \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer alpheya"
```

```
curl -X PUT http://localhost:9090/messages/ce222716-edf5-4d2e-803a-70dea8758ae4/approve \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer alpheya"
```
