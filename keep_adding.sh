#!/bin/bash

while true
do
    timestamp=$(echo $(( $(date +%s) * 1000 )))
    curl -X POST -H "Content-Type: application/json" -d "{\"created\": ${timestamp}}" http://localhost:5984/mydb -u "admin:YOURPASSWORD"
    sleep 1
done