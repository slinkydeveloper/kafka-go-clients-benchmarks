# Kafka Go Clients benchmarks

## Producer

Run kafka performance consumer 

```
./rdkafka_performance -b localhost:9092 -C -G consumer-test-1 -t test-a -V > consumer_log.csv
```

Run sarama producer

```
go run sarama/main.go -mode producer -broker localhost:9092 -topic test-a 
```

## Consumer

Run kafka performance producer

```
./rdkafka_performance -b localhost:9092 -P -t test-a -V -r 200000 -s 1024
```
