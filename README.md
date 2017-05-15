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

the default config file is ${curdir}/Milena.yml
```shell
cd bin && ./Milena -f ${conf}.yml
```

## stop

in the start cmd dir
```shell
./Milena -s stop
```


## results

topic contents will append to a file which is name by

```shell
${datadir}/${clustername}/${topicName}.log
```
conetnt will like this:

```shell
ts=>[${timestamp}] p:${parition} o:${offset} =>${value}
```

there is a metadata dir in `${datadir}/${clustername}`, never delete it unless lt is broken.It stored offsets already consumed.

## thanks to

- [sarama](https://github.com/Shopify/sarama)
- [leveldb](https://github.com/syndtr/goleveldb)
- [configor](https://github.com/jinzhu/configor)
- [logrus](https://github.com/sirupsen/logrus) (though my own fork actually)

## TODO
- optimize log append engine(cur can't be called engine, it's just a working shit)
- testify and optimize stop mechanism
- gen more cmd to help
  - reload signal -s
  - repair command
- add travis.ci to help
- more powerful makefile
- fix shit code(always)
