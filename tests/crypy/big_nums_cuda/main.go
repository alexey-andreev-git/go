package main

/*
#cgo LDFLAGS: -L/home/alexan/projects/temp/tests/go/big_nums_cuda
#cgo LDFLAGS: -lbig_math -lcudart
#include <stdint.h>

// Declaration of the C function from the CUDA library
typedef struct modExpParams {
    uint64_t base;
    uint64_t exp;
    uint64_t mod;
    uint64_t result;
} modExpParams;

void cudaModExp(modExpParams *params, uint64_t numParams);
*/
import "C"
import (
	"fmt"
	"time"
	"unsafe"
)

func main() {
	// Define parameters for modular exponentiation
	base := uint64(1234567890123456)
	exp := uint64(9876543210654321)
	mod := uint64(1000000000000007)
	nums := 10000 // Number of calculations

	// Create a slice of modExpParams to hold parameters for each calculation
	params := make([]C.modExpParams, nums)
	for i := 0; i < nums; i++ {
		params[i] = C.modExpParams{
			base:   C.uint64_t(base),
			exp:    C.uint64_t(exp),
			mod:    C.uint64_t(mod),
			result: 0,
		}
	}

	// Loop until a key is pressed
	fmt.Println("Press 'Enter' to stop calculations...")
	stop := make(chan bool)
	go func() {
		fmt.Scanln()
	}()

	for {
		select {
		case <-stop:
			break
		default:
			startTime := time.Now()
			// Call the CUDA function for modular exponentiation on the GPU
			for i := 0; i < nums; i++ {
				C.cudaModExp((*C.modExpParams)(unsafe.Pointer(&params[0])), C.uint64_t(nums))
				if i == 0 || i == nums-1 {
					elapsedTime := time.Since(startTime)
					fmt.Printf("Elapsed time: %s\n", elapsedTime)
				}
			}
			elapsedTime := time.Since(startTime)
			fmt.Printf("Elapsed time: %s\n", elapsedTime)

			// Print results
			// for i := 0; i < nums; i++ {
			// 	result := uint64(params[i].result)
			// 	fmt.Printf("Result of %d^%d mod %d = %d\n", base, exp, mod, result)
			// }
			fmt.Printf("Nums: %d x %d = %d\n", nums, nums, nums*nums)
			result := uint64(params[0].result)
			fmt.Printf("Result of %d^%d mod %d = %d\n", base, exp, mod, result)
		}
	}
}
