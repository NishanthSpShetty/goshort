## URL shortner in go

    Implementation of url shortner micro service using hexagonal architecture in Go.

    Using Mongodb or Redis as a datastore of choice, the service will create a short url.


### Run 

Need to update the config file related to the datastore, URL_DB indicates the type of data store which is already setup in its relevent config

#### With mongodb as datastore
    
    Update the mongo.ini file with mongodb specific config and run the following
    ```
     source mongo.ini && go run main.go
    ```
#### With redis as datastore
    
    Update the redis.ini file with redis url config and run the following
    ```
     source redis.ini && go run main.go
    ```

### Usage

Send a post request as follows to the service
```
curl -X POST localhost:8080 -H "Content-Type: application/json" \
         -d '{"url":"https://github.com/NishanthSpShetty"}'
```

This will return the following

```
{"code":"akrNeyQGR","url":"https://github.com/NishanthSpShetty","create_at":1616693447}%
```


Go over the brower and copy the Code returned by above call and replace the code in the followig URL
```
http://localhost:8080/uberCswMg
```
