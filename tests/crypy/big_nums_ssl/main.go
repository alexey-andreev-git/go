package main

/*
#cgo LDFLAGS: -lcrypto
#include <openssl/bn.h>
#include <stdlib.h>

char* generate_random_prime(int bits) {
    BIGNUM *prime = BN_new();
    BN_generate_prime_ex(prime, bits, 0, NULL, NULL, NULL);

    // Convert BIGNUM to hexadecimal string
    char *prime_str = BN_bn2hex(prime);
    BN_free(prime);
    return prime_str;
}
*/
import "C"
import (
	"fmt"
	"sync"
	"time"
	"unsafe"
)

// generatePrimeOpenSSL generates a prime number with the given bit length using OpenSSL.
func generatePrimeOpenSSL(bitLength int) string {
	primeStr := C.generate_random_prime(C.int(bitLength))
	defer C.free(unsafe.Pointer(primeStr))
	return C.GoString(primeStr)
}

func main() {
	bitLength := 4096

	startTime := time.Now()
	// Generate a 4096-bit prime number using OpenSSL
	// var primeHex string
	var primeHexes [10]string
	wg := sync.WaitGroup{}
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func(primeHex *string) {
			defer wg.Done()
			*primeHex = generatePrimeOpenSSL(bitLength)
			fmt.Println("Generated Prime (Hex):", *primeHex)
		}(&primeHexes[i])
		// primeHexes[i] = generatePrimeOpenSSL(bitLength)
		// fmt.Println("Generated Prime (Hex):", primeHexes[i])
	}
	wg.Wait()
	elapsedTime := time.Since(startTime)
	fmt.Println("Time taken:", elapsedTime)
	// primeHex := generatePrimeOpenSSL(bitLength)

	// fmt.Println("Generated Prime (Hex):", primeHex)
}
