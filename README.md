mysql-slave-healthcheck-agent
=============================

Description
------------

`mysql-slave-healthcheck-agent` is a HTTP daemon to show an information of "SHOW SLAVE STATUS" in JSON format.

This is useful for HAProxy's health check access instead of `option mysql-check`.

```
# haproxy.cfg
listen mysql
       bind     127.0.0.1:3306
       mode     tcp
       balance  roundrobin
       option   httpchk
       server   slave1 192.168.1.11:3306 check port 5000
       server   slave2 192.168.1.12:3306 check port 5000
```

Usage
------

```
$ mysql-slave-healthcheck-agent
2014/05/21 00:04:28 Listing port 5000
2014/05/21 00:04:28 dsn root:@tcp(127.0.0.1:3306)/?charset=utf8
```

```
$ curl localhost:5000
OK

* The query "SHOW SLAVE STATUS" was succeeded, return HTTP status 200.
* If could not connect to the MySQL or the MySQL is not a slave, return HTTP status 500.
* If slave is not running or the replication lag exceeds the limit, return HTTP status 500

Options
-------

* -port : http listen port number. default 5000.
* -dsn : Data Source Name for MySQL. default "root:@tcp(127.0.0.1:3306)/?charset=utf8"
  See also https://github.com/go-sql-driver/mysql#dsn-data-source-name
  mysql's user must have a privilege of "REPLICATION CLIENT".
* --lag-limit : max replication lag for healthy slave

How to build
------------

    $ go get github.com/fujiwara/mysql-slave-healthcheck-agent

How to build (cross compile)
------------

Install gox.

    $ go get github.com/mitchellh/gox
    $ gox -build-toolchain

And make.

    $ make

The binary files will be build in `pkg` directory.

```
pkg
├── darwin_amd64
│   └── mysql-slave-healthcheck-agent
└── linux_amd64
    └── mysql-slave-healthcheck-agent
```

License
-------
Apache License 2.0

Author
-------
fujiwara.shunichiro@gmail.com
