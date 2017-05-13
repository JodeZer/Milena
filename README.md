# Milena
kafka probe tool

While using kafka as a prodcuer, it is a problem to know what if prodcuer pushs right message to it (At least kafka-manager doesn't show topic contents); 
Or when producer wants to know whether consumers can pull right msg from kafka , producers may has to write another test-consumer.

I wrote this toy tool just for proving  to testers that my producer program works very well !

Once Milena started , she can watch your  every new messgae of kafka topic confed in conf file and write it to a log file. And she also supports duration for offsets 
which means U can always restart Milena and don't need to worry about reconsume history messages.

## download

``` shell
go get github.com/JodeZer/Milena
```

## compile

```shell
make build
```
## config
config is yml file like this:

```yml
datadir : data

servers:
- name: cluster1
  brokers:
  - "192.168.1.x:9092"
  - "192.168.1.x:9092"
  - "192.168.1.x:9092"
  topics:
  - name: "xxx"
    partitions:
    - partition: 0
      start: 7
  - name: "xxxxxx"
```

## run
TODO: -f command to point config file soon;and now, it find `Milena.yml` in it's dir
config your own Milena.yml for configuration only

```shell
cd bin && ./Milena
```

## stop
TODO: this will gen a command soon
```shell
kill -2 $(pid of Milena)
```

## thanks to

- [sarama](https://github.com/Shopify/sarama)
- [leveldb](https://github.com/syndtr/goleveldb)
- [configor](https://github.com/jinzhu/configor)
- [logrus](https://github.com/sirupsen/logrus) (though my own fork actually)

## TODO
- optimize log append engine(cur can't be called engine, it's just a working shit)
- testify and optimize stop mechanism 
- gen more cmd to help
  - stop signal -s
  - reload signal -s
  - config option -f
  - clean command
- add travis.ci to help
- more powerful makefile
- golang deps manage
- fix shit code(always)
