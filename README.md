# db_experiments
Prototype code pieces from my personal database experiments.

### Load Leaking
An experiment to explore and demonstrate as to how leaking a random proportion (say 20%) of your read query load deliberately to database (even if the key is present in your cache) in a typical setup of Server - Cache - Database can help you reduce overall average response times of the system.

#### Overall Setup:

1. This requires WSL2, Ubuntu on your Windows machine. Set it up.
2. Download, install and setup Redis on Ubuntu WSL.
3. Fire the following command on WSL terminal.

```bash
sudo service redis-server start
```

4. Download, install and setup MySQL on your machine.
5. Start the service MYSQL83 (or similar) on your Local device and then use the following commands.

```bash
# Login as a root user
mysql -u root -p

# Create a new database
mysql> create database dll_experiment;

# Create a new user and set its password
mysql> create user 'dllexpuser'@'%' identified by 'dllexppwd';

# Gives all privileges to the new user on the newly created database
mysql> grant all on dll_experiment.* to 'dllexpuser'@'%';
```

6. Set your connection string in the main.go file as per your database and user setup.
7. Create a simple key, value store table, say "keyvalue" and push some data to it. Say 10000 rows.
8. Set maxmemory and maxmemory-policy params to "1100000" and "allkeys-lru" in the redis.conf file on your WSL to set a limit on the number of keys your cache can store and to setup an eviction policy for keys, which in our case is LRU for all keys.
9. Now, in your own terminal, use the following commands to run the experiment.

```bash
cd load_leaking
go run main.go
```

#### Configuring Redis:

To access redis.conf file, use the following commands on WSL:

```bash
# To run your shell as a privileged shell
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
# Open Redis CLI
redis-cli

# To test if server up and running
PING

# To get a view of all the keys in cache
KEYS *

# To get or set the value of a param from redis.conf (temporarily)
config get param
config set param value
```
