# color-lizard

![](https://i.pinimg.com/originals/db/1c/10/db1c10086897ec6bed7ef20d8480fca8.jpg)

Target version of more open mockey

This allows developers to mock their GET, POST, PUT endpoints for desired response status code, response headers and response

This is configuration driven. Below is sample configuration:
Additionally you can use localhost:8881/add with the configuration as the body to POST your configuration at runtime.


```
{
    "/oauth/token":{
    "method":"GET",
    "status":200,
    "response":"{\"result\": \"test1\"}",
    "headers":{
      "header1":"value1",
      "header2":"value2"
    }
  },
   "/test2/test2":{
    "method":"POST",
    "status":201,
    "response":"{\"result\": \"test2\"}",
    "headers":{
      "header1":"value1",
      "header2":"value2"
    }
  }
}
```

```
curl -X GET 'http://localhost:8881/colorlizard/oauth/token?test=1'

{"result": "test1"}
```

# How to add new endpoints

```
curl -X POST \
  https://colorlizard-tgt.dev.target.com/add \
  -H 'postman-token: 67319743-3171-d815-371c-e1de2aba087e' \
  -d '  {"/promo": {
    "method": "PUT",
    "status": 200,
    "response": "{\"Validation error\": \"ok\"}",
    "headers": {
      "header1": "value1",
      "header2": "value2"
    }
  }}

```
Above will create new endpoint 

```curl -X PUT \
    http://colorlizard-tgt.dev.target.com/colorlizard/promo
```

which will give result as 

```
{
    "Validation error": "ok"
}
``` 
