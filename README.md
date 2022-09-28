# blockchain-indexer

## TODO

- Add block producer (kafka) ✅
- Add block consumer (kafka) ✅
- Add metrics ✅
- Add Postgres ✅
- Add nodes monitoring ✅
- Add nodes backup mechanism ✅
- Add queues with Rabbit MQ for workers: ✅
    -- Worker that saves txs in Postgres;
- Add more chains
- Add API
- Check if consumers are scalable
- Draw a diagram of the architecture
- Update README.md

## Problems

- Block trackers are not working correctly (it doesn't set Current block num, only for ETHEREUM and the it resets it)
- Cosmos txs are not parsed
