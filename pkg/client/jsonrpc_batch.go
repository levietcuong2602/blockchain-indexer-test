package client

type RPCRequestMapper func(interface{}) RPCRequest

func MakeBatchRequests(elements []interface{}, batchSize int, mapper RPCRequestMapper) (requests []RPCRequests) {
	batches := MakeBatches(elements, batchSize)
	for _, batch := range batches {
		var reqs RPCRequests
		for _, ele := range batch {
			mapped := mapper(ele)
			reqs = append(reqs, &mapped)
		}
		requests = append(requests, reqs)
	}

	return
}

func MakeBatches(elements []interface{}, batchSize int) (batches [][]interface{}) {
	batch := make([]interface{}, 0)
	size := 0

	for _, ele := range elements {
		if size >= batchSize {
			batches = append(batches, batch)
			size = 0
			batch = make([]interface{}, 0)
		}
		size++
		batch = append(batch, ele)
	}

	batches = append(batches, batch)

	return
}
