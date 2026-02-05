# giftCode usages report

![giftCode report](./gift.report.png)

## description
1. user place GiftCode in request url and calls `GiftCodeReport` endpoint from the giftCode component.
2. giftCode component checks if GiftCode exist and returns request's response.(gift code usage list including : wallet ID, date)

# Api contract

## GiftCode report
```
Name:   GiftCodeReport
Method: Get
Url:    https://localhost:7878/gift/report/:giftCode
Headers:
Body:
Errors:
   - code: 404
     Name: not found
     Body:
         {
            "error" : "GiftCode does not exist",
         }
   - code: 500
     Name: internal server error
     Body:
   - code: 400
     Name: bad request
     Body: 
         {
            "error" : "invalid request url, please enter GiftCode correctly",
         }
Responses:
    - code: 200
      Name: ok 
      Body:
          {
            "walletID" : string,
            "date" : (Time.time)
            "--------------------------------"
          }
```
