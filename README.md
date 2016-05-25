# GOGEO:

a simple wrapper for google maps api.

## gogeo http

command flags

* port - server address for the http server to listen on.
* bing-key - if specifed request using bing doesn't need to provider a key
* google-key - if specifed request using google doesn't need to provider a key
* maprequest-key - if specifed request using maprequest doesn't need to provider a key

---

  GET /{goggle,bing or mapquest}/{json,yml or xml}

  parameters:
  * addr - The street address that you want to geocode.
  * loc - format: {latitude,longitude} location to lookup
  * key - (optional) api key can be set though the command line

---
  note: ecurrently google is the only implemented provider with images. working on the rest.

  GET /google/png

  parameters:
  * addr - The street address that you want to geocode.
  * size (optional) image size: can be specified as {width}x{height} or {size}
  * zoom (optional)
  * scale (optional)
  * key - (optional) api key

 
  
