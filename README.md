# GOGEO:

CLI and http server, for reversed geocoding, address lookup, and image rendering. for varies providers, in a common format. 

![Alt text](gogeo.png?raw=true "Gogeo")

## gogeo http

command flags

* port - server address for the http server to listen on.
* bing-key - if specifed request using bing doesn't need to provider a key, otherwise it will default to env: GOGEO_BING
* google-key - if specifed request using google doesn't need to provider a key, otherwise it will default to env: GOGEO_GOOGLE
* mapquest-key - if specifed request using maprequest doesn't need to provider a key, otherwise it will default to env: GOGEO_MAPQUEST

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

 
  
