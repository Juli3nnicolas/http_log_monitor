# Nginx container

Contains a slightly modified nginx conf to log traffic to /tmp/access.log

## Build and run the container
```
docker build -t test_nginx .
docker run -d -p 80:80 --name test_nginx test_nginx
```

## Check access logs
```
docker exec test_nginx tail -f /tmp/access.log
```

## Query the server (to trigger a log write)
```
curl --request GET --url http://127.0.0.1/
```