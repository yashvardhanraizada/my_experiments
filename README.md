# db_experiments
Prototype code pieces from my personal database experiments.

### Load Leaking:

1. This requires WSL2, Ubuntu on your Windows machine. Set it up.
2. Download, install and setup Redis on Ubuntu WSL.
3. Fire the following command on WSL terminal.

```bash
sudo service redis-server start
```

4. In your own terminal, use the following commands to run the experiment.

```bash
cd load_leaking
go run main.go
```

To access redis.conf file, use the following commands on WSL:

```bash
sudo -s
cd /etc/redis
nano redis.conf
```
After making changes to redis.conf, don't forget to restart redis-server:

```bash
sudo service redis-server restart
```

You can access redis keys and conf params via redis-cli:

```bash
redis-cli
PING // to test if server up and running
KEYS * // to get a view of all the keys in cache
config get param // to get the value of a param from redis.conf
config set param value // to set a value to a param of redis.conf temporarily
```