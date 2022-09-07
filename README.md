# blockchain-indexer

docker exec -it kafka kafka-topics.sh --bootstrap-server localhost:29092 --topic blocks-topic-bsc --create --replication-factor 1 --partitions 1

docker exec -it kafka kafka-topics.sh --delete -topic blocks-topic-bsc --bootstrap-server localhost:9092

docker exec -it kafka kafka-topics.sh --list --bootstrap-server localhost:9092

docker exec -it kafka kafka-console-consumer.sh --topic blocks-topic-bsc --from-beginning --bootstrap-server localhost:9092
