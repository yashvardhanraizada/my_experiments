# db_experiments
Prototype code pieces from my personal database experiments.

### Load Leaking:

1. This requires WSL2, Ubuntu on your Windows machine. Set it up.
2. Download, install and setup Redis on Ubuntu WSL.
3. Fire the following command on WSL terminal.

```bash
sudo service redis-server start
sudo service redis-server restart
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