# kuberdbs
Get ondemand databases on top of Kubernetes using `curl`.

## Databases
This project currently supports the following databases:
* Redis
* MySQL

## Usage
Edit kuberdbs kubernetes deployment and replace `--redis-addr` and `--mysql-addr` with your redis and mysql database hostnames in the kubernetes cluster

```bash
vim kubernetes/kuberdbs-deployment.yaml
....
        args:
          - "--redis-addr"
          - "redis.redis:6379" # change this to your redis server address
          - "--mysql-addr"
          - "mysql.mysql:3306" # change this to your mysql server address
....
```

Edit kuberdbs secret file

```bash
vim kubernetes/kuberdbs-secret.yaml
apiVersion: v1
data:
  rootpw: "MYSQL_ROOT_PASSWORD" # change this to your mysql password
kind: Secret
metadata:
  name: mysql-secret
```


Create deployment

```bash
kubectl apply -f kubernetes/kuberdbs-deployment.yaml
```

Create `kuberdbs` service

```bash
kubectl apply -f kubernetes/kuberdbs-service.yaml
```

Get load balancer address

```bash
kubectl describe service kuberdbs
Name:			kuberdbs
Namespace:		default
Labels:			app=kuberdbs
Annotations:		kubectl.kubernetes.io/last-applied-configuration=
Selector:		app=kuberdbs
Type:			LoadBalancer
IP:			10.0.0.0
LoadBalancer Ingress:	kuberdbs.region.elb.amazonaws.com
Port:			kuberdbs	80/TCP
NodePort:		kuberdbs	32122/TCP
Endpoints:		10.0.0.1:8080
Session Affinity:	None
Events:			<none>
```

Get a new Redis database
```bash
curl http://kuberdbs.region.elb.amazonaws.com/redis
```

Which will return something like:

```
REDIS_URL=redis://redis.redis:9999/4524
```

Get a new MySQL database:
```bash
curl http://kuberdbs.region.elb.amazonaws.com/mysql
```

Which will return something like:

```
DATABASE_URL=mysql://Username:Password@mysql.mysql:3306/DatabaseName
```