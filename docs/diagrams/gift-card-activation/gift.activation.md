# gift activation

![gift activation](./gift.activation.png)

## Description

1. user call `giftCardActivation` endpoint from giftCode component
2. giftCode component call `addCredit` endpoint from wallet component
3. wallet component return `addCredit` request response
4. giftCode component return `giftCardActivation` request response

# Api contract

## gift code

```
Name:   giftCardActivation
Method: POST
Url:    http://localhost:7878/gift/use
Headers: no content
Body:
    {
        "phone": (string),
        "code":  (string)
    }
Errors:
    - code: 500
      Name: Internal Server Error
      Body: no content
Response:
    - code: 200
      Name: StatusOK
      Body:
          {
            "message" : "success",
          }
```

## wallet

```

Name:   addCredit
Method: Post
Url:    http://localhost:9898/wallet/gift
Headers: no content
Body:
    {
       "mobileNumber" : (int),
    }
Errors:
    - code: 500
      Name: Internal Server Error
      Body:
          {
            "error" : "failed to add credit",
          }
    - code: 400
      Name: Bad Request
      Body:
          {
            "error" : "invalid gift code",
          }
Response:
    - code: 200
      Name: statusOK
      Body:
          {
            "message" : "success",
          }
```
