## Trying to figure out problem with continous deletion of data through _design in Couchdb

Start a database locally:
```
docker compose up -d
```

Create design doc `_design/nonce` with index name: `older_than`, reduce: `NONE`

Map function: 

```javascript
function(doc) 
{ 
    var now = Date.now() - (5 * 1000);
    if (doc.created < now) {
        emit(doc.created, doc._rev); 
    }
}
```


Query documents using view: 
```
curl -X GET -u "admin:YOURPASSWORD" "http://localhost:5984/mydb/_design/nonce/_view/older_than?limit=100"
```