package api

import (
	"crypto/sha256"
	"math/big"
)

func GetPartition(id string, totalPartitions int) int {
	hash := sha256.Sum256([]byte(id))
	hashInt := new(big.Int).SetBytes(hash[:])
	return int(hashInt.Mod(hashInt, big.NewInt(int64(totalPartitions))).Int64())
}
