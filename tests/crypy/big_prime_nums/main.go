package main

/*
#cgo LDFLAGS: -L. -lprime_gen -lcudart
#include <stdint.h>

// Declaration of the CUDA functions
void millerRabinTestRange(uint64_t *numbers, int *results, int numCandidates);
*/
import "C"
import (
	"crypto/rand"
	"fmt"
	"math/big"
	"sync"
	"sync/atomic"
	"unsafe"
)

// Generate a random big integer with the specified number of bits
func generateRandomBigInt(bitSize int) *big.Int {
	n, _ := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), uint(bitSize)))
	return n.SetBit(n, bitSize-1, 1) // Ensure the number has the required bit length
}

// Convert a big.Int to an array of uint64 segments
func bigIntToUint64Array(n *big.Int, numSegments int) []uint64 {
	bytes := n.Bytes()
	segments := make([]uint64, numSegments)
	for i := 0; i < len(bytes); i++ {
		segmentIndex := (len(bytes) - i - 1) / 8
		segments[segmentIndex] = (segments[segmentIndex] << 8) | uint64(bytes[i])
	}
	return segments
}

func millerRabinCPU(n *big.Int, rounds int) bool {
	if n.Cmp(big.NewInt(2)) < 0 {
		return false
	}
	if n.Cmp(big.NewInt(2)) == 0 || n.Cmp(big.NewInt(3)) == 0 {
		return true
	}
	if new(big.Int).Mod(n, big.NewInt(2)).Cmp(big.NewInt(0)) == 0 {
		return false
	}

	// Factor n-1 as d * 2^r
	d := new(big.Int).Sub(n, big.NewInt(1))
	r := 0
	for new(big.Int).And(d, big.NewInt(1)).Cmp(big.NewInt(0)) == 0 {
		d.Rsh(d, 1)
		r++
	}

	for i := 0; i < rounds; i++ {
		a, _ := rand.Int(rand.Reader, new(big.Int).Sub(n, big.NewInt(4)))
		a.Add(a, big.NewInt(2))

		x := new(big.Int).Exp(a, d, n)
		if x.Cmp(big.NewInt(1)) == 0 || x.Cmp(new(big.Int).Sub(n, big.NewInt(1))) == 0 {
			continue
		}

		composite := true
		for j := 1; j < r; j++ {
			x.Exp(x, big.NewInt(2), n)
			if x.Cmp(new(big.Int).Sub(n, big.NewInt(1))) == 0 {
				composite = false
				break
			}
		}
		if composite {
			return false
		}
	}
	return true
}

func main() {
	numCandidates := 1 // Number of random numbers to generate and test
	bitSize := 4096    // Bit length of each number (e.g., 1024 or 2048)
	numSegments := bitSize / 64

	// Generate random large numbers and convert them to uint64 arrays
	numbers := make([]uint64, numCandidates*numSegments)
	results := make([]int32, numCandidates)

	n := new(big.Int)
	// testNum := "e3c3b046e90acc3b20f74e6adf9eae83359922a2a661d8f32a4c8039dac363526c77cb5ca3302c8145cf50e917f327575da0ef4a4b59fec7fb385b687d996db218eb453bde1725a01885830a5911966113665f6dd5df9f51ab6b651ab60aef61524de48f40ce568770ac9c5dd34c8b2be8938e4f972bf7fcadff86bb44ea6a2bf6629981b3b8276dd6c36fc0baccf4a85e99ec25917c3f2df6b93284f1ce5323783e9ab7649ea256b4d4ca363c01b43e611df1eea268be10539cd680877e2c1e9d933f57fbac12019ed8a7953d7005555bab200a3dad3e452608f903869717cd1d09cafc99c0f6067f6333a37372fa038087bf66577abc325e4fcf945f1441980b025e7b39fa3ded22ee4dd48651fdd69cd111afcc6bbe2bf5f25d1a854e15088f4b995b9af1d8969f5bde404da8d32eb72da614f6f4c35b7c9de8e6ecb3ad8ccd9296c8f4d05894b7c68f5ec2eff57b704c0f67dc78322c4878abc65210df93a174ffe1be800f8a68fc4667ad1daaba8e8e2dde11d72b495c4122c6288067b30adb246c91b2310fc9f81b648cc0c84fde1f59a2775cc6ffb48a5615ef3e83c4f931c293fe628936ce2c6567ea7013bfa4e4fbd2207492549b2a1e9e4f0b766f28aec1dcc7334bd210e8000fcf3c183d17cd15fbc6e9f8092fe03d50b1845ce81e8996aaf64bd3dc48fc6360106f462119594a9257c4df83d057ef883a246e7f"
	testNum := "D74D76EEA45A6E5A077C7B84FA74B702A5290107352F5FE9B54358760FB33552CCE7FAD1C44E14BAE2343171319D990A63406AB52A9CD9A97AD211C9E0155093D33445E54BB71A712C240E06221FFA573CEF12ED166A03C7CEEACF9E898CBADCCC198DED4F4C4A2844268079861738EFB354047BBD0334552BB6C2D873250B9F32ECDC05E2BBDEC037C59DB93215F950F6275B43F523F6ACC8B320FE168F527435B379C4BF3FA0F7E7B84F490107F8B2FC04E7DC9B2878267CE419D5454AF5CF0A24C9F55549E634F800949D2CC47E4C8E619B9B2D2BFC61F227417C46FA85B19F6910AA039E67DD6CCFBC3AC48CBC83A63B062D5405C660D94AADE0DE489C79007FD3D5972B5FF9466CFF9D7927D6D7D347E6D533AC860A17655BC8960830B0438411038DE528E67551B40357011DF905331E58002F64A81918C87775BD60E304C02EC29CFE9D6D63F0829D7E22061A242DE7E3465C1C28B67A31941F4BB43600F793A0285502A59C4CC3FD9BFAF10962DDA3C41798F900A9A57D922AF406420552FF0351E14D828EE1FDCDD1D8A5131ACE53188AEEAA9F05E45E1023CBE87939D4FEAB33441D43AE5CE9B5AE4A80728AB2766FE42E6F3D7437D23BE6461D433D1C5DCAE26960EE067988F4B3933B1E10079419100C174491650B882BE974DBE8DE0C4078D79AFD9C4F50FA0248010C75FD5BFC6DDB38053A29BD0DA58A6673"
	testNumBigInt := new(big.Int)
	testNumBigInt.SetString(testNum, 16)
	halfTestNum := new(big.Int).Div(testNumBigInt, big.NewInt(2))
	var numRange [3]*big.Int
	numRange[0] = new(big.Int).Div(halfTestNum, big.NewInt(4))
	numRange[1] = new(big.Int).Div(halfTestNum, big.NewInt(2))
	numRange[2] = new(big.Int).Sub(halfTestNum, numRange[0])
	var isComposite atomic.Bool
	isComposite.Store(false)
	startBound := big.NewInt(2)
	endBound := big.NewInt(2)
	wg := sync.WaitGroup{}
	for k := 0; k < 3; k++ {
		startBound = new(big.Int).Set(endBound)
		endBound = new(big.Int).Set(numRange[k])
		wg.Add(1)
		go func(startBound, endBound *big.Int) {
			for i := new(big.Int).Set(startBound); i.Cmp(endBound) == -1; i = new(big.Int).Add(i, big.NewInt(1)) {
				if new(big.Int).Mod(testNumBigInt, i).Cmp(big.NewInt(0)) == 0 {
					fmt.Printf("Number %d is composite.\n", k)
					isComposite.Store(true)
					break
				}
			}
			if new(big.Int).Mod(testNumBigInt, endBound).Cmp(big.NewInt(0)) != 0 {
				fmt.Printf("Number %d is probably prime.\n", k)
			}
		}(startBound, endBound)
	}
	wg.Wait()
	if !isComposite.Load() {
		fmt.Printf("Number %d is probably prime.\n", 1)
	}
	// i := new(big.Int)
	// isComposite := false
	// for i = big.NewInt(2); i.Cmp(halfTestNum) == -1; i = new(big.Int).Add(i, big.NewInt(1)) {
	// 	if new(big.Int).Mod(testNumBigInt, i).Cmp(big.NewInt(0)) == 0 {
	// 		isComposite = true
	// 		fmt.Printf("Number %d is composite.\n", 1)
	// 		break
	// 	}
	// }
	// if !isComposite {
	// 	fmt.Printf("Number %d is probably prime.\n", 1)
	// }
	// testNum := "17"
	for i := 0; i < numCandidates; i++ {
		if i == 0 {
			n.SetString(testNum, 16)
		} else {
			n = generateRandomBigInt(bitSize)
		}
		segments := bigIntToUint64Array(n, numSegments)
		copy(numbers[i*numSegments:(i+1)*numSegments], segments)
		fmt.Printf("Generated number %d: %s\n", i+1, n.Text(16))
		// if millerRabinCPU(n, 10) {
		// 	fmt.Printf("Number %d is probably prime.\n", i+1)
		// } else {
		// 	fmt.Printf("Number %d is composite.\n", i+1)
		// }
	}

	// Convert Go slices to C pointers
	numbersPtr := (*C.uint64_t)(unsafe.Pointer(&numbers[0]))
	resultsPtr := (*C.int)(unsafe.Pointer(&results[0]))

	// Call the CUDA function
	C.millerRabinTestRange(numbersPtr, resultsPtr, C.int(numCandidates))
	byteArray := make([]byte, numSegments*8)
	for i := 0; i < numSegments; i++ {
		for j := 0; j < 8; j++ {
			byteArray[i*8+j] = byte(numbers[i] >> (8 * (7 - j)))
		}
	}
	n = new(big.Int).SetBytes(byteArray)
	if millerRabinCPU(n, 32) {
		fmt.Printf("Number %d is probably prime.\n", 1)
	}

	// Output the results
	for i := 0; i < numCandidates; i++ {
		// if results[i] == 1 {
		// 	fmt.Printf("Number %d is probably prime.\n", i+1)
		// } else {
		// 	fmt.Printf("Number %d is composite.\n", i+1)
		// }
		if results[i] == 1 {
			fmt.Printf("Number %d is probably prime.\n", i+1)
		}
	}
}
