# color-lizard


Target version of more open mockey

This allows developers to mock their GET, POST, PUT endpoints for desired response status code, response headers and response

This is configuration driven. Below is sample configuration:

  


`[
  {
    "path": "oauth/token",
    "status": 200,
    "response": "{\"result\": \"test1\"}", 
    "method": "GET", 
    "headers": {
      "header1": "value1",
      "header2": "value2"
    }
  }
]
`


 