package main

/*
#cgo LDFLAGS: -lgmp
#include <gmp.h>
#include <stdlib.h>

// Function to generate a random prime using GMP
void generate_random_prime(mpz_t result, int bits) {
    gmp_randstate_t state;
    gmp_randinit_default(state);

    // Seed the random state with the current time
    mpz_t seed;
    mpz_init(seed);
    mpz_set_ui(seed, time(NULL)); // Use the current time as the seed
    gmp_randseed(state, seed);
    mpz_clear(seed);

    mpz_urandomb(result, state, bits); // Generate a random number with specified bits length
    mpz_nextprime(result, result);     // Find the next prime greater than the random number

	gmp_randclear(state);
}
*/
import "C"
import (
	"fmt"
	"time"
	"unsafe"
)

// generatePrimeGMP generates a prime number with the given bit length using GMP.
func generatePrimeGMP(bitLength int) string {
	var prime C.mpz_t            // Declare prime as a C.mpz_t array (not a pointer)
	C.mpz_init(&prime[0])        // Initialize the mpz_t variable
	defer C.mpz_clear(&prime[0]) // Clear mpz_t memory when done

	C.generate_random_prime(&prime[0], C.int(bitLength))

	// Convert the prime number to a string in hexadecimal format
	primeStr := C.mpz_get_str(nil, 16, &prime[0])
	defer C.free(unsafe.Pointer(primeStr))
	return C.GoString(primeStr)
}

func main() {
	bitLength := 4096

	startTime := time.Now()
	// Generate a 4096-bit prime number using GMP
	var primeHex string
	var primeHexArr [10]string
	for i := 0; i < 10; i++ {
		primeHex = generatePrimeGMP(bitLength)
		present := false
		for j := 0; j < i; j++ {
			if primeHex == primeHexArr[j] {
				fmt.Println("The generated prime number is already in the array.")
				present = true
				break
			}
		}
		if !present {
			primeHexArr[i] = primeHex
			fmt.Println("Generated Prime (Hex):", primeHex)
		}
	}
	// primeHex := generatePrimeGMP(bitLength)
	elapsedTime := time.Since(startTime)
	fmt.Println("Elapsed time:", elapsedTime)

	// fmt.Println("Generated Prime (Hex):", primeHex)
}
