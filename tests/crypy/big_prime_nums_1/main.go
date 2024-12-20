package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	bitLength := 64
	// testNum := "D74D76EEA45A6E5A077C7B84FA74B702A5290107352F5FE9B54358760FB33552CCE7FAD1C44E14BAE2343171319D990A63406AB52A9CD9A97AD211C9E0155093D33445E54BB71A712C240E06221FFA573CEF12ED166A03C7CEEACF9E898CBADCCC198DED4F4C4A2844268079861738EFB354047BBD0334552BB6C2D873250B9F32ECDC05E2BBDEC037C59DB93215F950F6275B43F523F6ACC8B320FE168F527435B379C4BF3FA0F7E7B84F490107F8B2FC04E7DC9B2878267CE419D5454AF5CF0A24C9F55549E634F800949D2CC47E4C8E619B9B2D2BFC61F227417C46FA85B19F6910AA039E67DD6CCFBC3AC48CBC83A63B062D5405C660D94AADE0DE489C79007FD3D5972B5FF9466CFF9D7927D6D7D347E6D533AC860A17655BC8960830B0438411038DE528E67551B40357011DF905331E58002F64A81918C87775BD60E304C02EC29CFE9D6D63F0829D7E22061A242DE7E3465C1C28B67A31941F4BB43600F793A0285502A59C4CC3FD9BFAF10962DDA3C41798F900A9A57D922AF406420552FF0351E14D828EE1FDCDD1D8A5131ACE53188AEEAA9F05E45E1023CBE87939D4FEAB33441D43AE5CE9B5AE4A80728AB2766FE42E6F3D7437D23BE6461D433D1C5DCAE26960EE067988F4B3933B1E10079419100C174491650B882BE974DBE8DE0C4078D79AFD9C4F50FA0248010C75FD5BFC6DDB38053A29BD0DA58A6673"
	// testNum := "FFFF"
	// testNumBigInt := new(big.Int)
	// testNumBigInt.SetString(testNum, 16)
	testNumBigInt, err := rand.Prime(rand.Reader, bitLength)
	if err != nil {
		fmt.Println("Error generating prime:", err)
		return
	}
	halfTestNum := new(big.Int).Div(testNumBigInt, big.NewInt(2))

	// Используем срез вместо массива для numRange
	numRange := []*big.Int{
		new(big.Int).Div(halfTestNum, big.NewInt(4)),
		new(big.Int).Div(halfTestNum, big.NewInt(2)),
		new(big.Int).Sub(halfTestNum, new(big.Int).Div(halfTestNum, big.NewInt(4))),
		new(big.Int).Set(halfTestNum),
	}

	var isComposite atomic.Bool
	isComposite.Store(false)
	startBound := big.NewInt(2)
	endBound := big.NewInt(2)
	wg := sync.WaitGroup{}
	startTime := time.Now()

	for k := 0; k < len(numRange); k++ {
		startBound = new(big.Int).Set(endBound)
		// endBound = new(big.Int).Add(startBound, numRange[k])
		endBound = new(big.Int).Set(numRange[k])
		wg.Add(1)
		go func(k int, startBound, endBound *big.Int) {
			defer wg.Done()
			// totalIterations := new(big.Int).Sub(endBound, startBound)
			var iterationCounter uint32 = 0
			iterationsDone := new(big.Int).SetInt64(0)
			fmt.Printf("Goroutine %d started. Range: %s - %s\n", k, startBound.String(), endBound.String())

			for i := new(big.Int).Set(startBound); i.Cmp(endBound) == -1 && !isComposite.Load(); i = new(big.Int).Add(i, big.NewInt(1)) {
				if new(big.Int).Mod(testNumBigInt, i).Cmp(big.NewInt(0)) == 0 {
					fmt.Printf("Number %d is composite. Divisible by %s.\n", k, i.String())
					isComposite.Store(true)
					break
				}

				iterationCounter++
				iterationCounter &= 0xFFFFFF
				if iterationCounter == 0 {
					// Достигли 65536 * 256 итераций
					iterationsDone.Add(iterationsDone, big.NewInt(65536))
					// Вычисляем и выводим процент выполненной работы
					percentage := new(big.Float).Quo(new(big.Float).SetInt(i), new(big.Float).SetInt(endBound))
					percentage.Mul(percentage, big.NewFloat(100))
					percentStr := fmt.Sprintf("%3.9f", percentage)
					duration := time.Since(startTime)
					fmt.Printf("Goroutine %d progress: %s%% completed. Duration: %v\n", k, percentStr, duration)
				}
			}

			// Добавляем оставшиеся итерации
			iterationsDone.Add(iterationsDone, big.NewInt(int64(iterationCounter)))

			if !isComposite.Load() && new(big.Int).Mod(testNumBigInt, endBound).Cmp(big.NewInt(0)) != 0 {
				fmt.Printf("Number %d is probably prime in this range.\n", k)
			}
		}(k, startBound, endBound)
	}

	wg.Wait()
	if !isComposite.Load() {
		fmt.Println("Number is probably prime.")
	}
}
