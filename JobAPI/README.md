# Job API lambda


stucture 
1. main.go
2. router -> router.go (return routing engine with all middleware added if needed) also handles the context parsing and passing to controllers, so that in future we plan to update the framework from gin to xyz we just have to update the router
there we will implement adapter pattern to have a liberty to switch to any Rest API framework in future

3. controller -> this package contains all the handlers mapped to the router. we will receive a common request structure that contains all request data like query param, body, path param or token data 
this also does any validation or 
4. Service -> this layer decides where we have to get data from or basically have all bisiness logic

5. Persistance - containes interface to interact with all database, cache, file storage