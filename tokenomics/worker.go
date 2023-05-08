// SPDX-License-Identifier: ice License 1.0

package tokenomics

func (r *repository) distributedWorkerIndices(concurrency, partitionCount, batchSize int) map[int][]uint16 {
	workerIndices := make(map[int][]uint16, concurrency)
	for ix := 0; ix < partitionCount; ix++ {
		if _, found := workerIndices[ix%(concurrency)]; !found {
			workerIndices[ix%(concurrency)] = make([]uint16, 0, partitionCount/concurrency)
		}
		workerIndices[ix%(concurrency)] = append(workerIndices[ix%(concurrency)], uint16(ix))
	}

	return workerIndices
}
