#include <cuda_runtime.h>
#include <curand_kernel.h>
#include <stdint.h>
#include <stdio.h>
#include <time.h>

#define NUM_ROUNDS 32
#define NUM_SEGMENTS 32 // For a 2048-bit number (32 * 64 bits)

extern "C" {

// CUDA kernel to initialize random states
__global__ void initCurand(curandState *state, unsigned long seed) {
    int idx = threadIdx.x + blockIdx.x * blockDim.x;
    curand_init(seed, idx, 0, &state[idx]);
}

// Device function to add two large numbers modulo n
__device__ void addMod(uint64_t *a, uint64_t *b, uint64_t *mod, uint64_t *result) {
    uint64_t carry = 0;
    for (int i = 0; i < NUM_SEGMENTS; i++) {
        uint64_t sum = a[i] + b[i] + carry;
        result[i] = sum % mod[i];
        carry = sum / mod[i];
    }
}

// Device function to multiply two large numbers modulo n
__device__ void mulMod(uint64_t *a, uint64_t *b, uint64_t *mod, uint64_t *result) {
    uint64_t temp[NUM_SEGMENTS] = {0};

    for (int i = 0; i < NUM_SEGMENTS; i++) {
        __uint128_t carry = 0;
        for (int j = 0; j < NUM_SEGMENTS - i; j++) {
            // Perform multiplication and add carry, then apply modular reduction
            __uint128_t prod = (__uint128_t)a[i] * b[j] + temp[i + j] + carry;

            // Apply modular reduction on the product within each segment
            temp[i + j] = (uint64_t)(prod % mod[i + j]);
            carry = prod / mod[i + j];
        }
    }

    // Final reduction of temp array to fit within mod
    for (int k = 0; k < NUM_SEGMENTS; k++) {
        result[k] = temp[k] % mod[k];
    }
}

// Device function to perform modular exponentiation on large numbers
__device__ void modExp(uint64_t *base, uint64_t *exp, uint64_t *mod, uint64_t *result) {
    uint64_t temp[NUM_SEGMENTS] = {1}; // Initialize result as 1
    uint64_t baseTemp[NUM_SEGMENTS];
    for (int i = 0; i < NUM_SEGMENTS; i++) baseTemp[i] = base[i];

    for (int i = NUM_SEGMENTS * 64 - 1; i >= 0; i--) {
        mulMod(temp, temp, mod, temp); // result = (result * result) % mod

        // Debug: Print temp after squaring
        // printf("modExp: temp after squaring: %llu\n", temp[0]);

        if ((exp[i / 64] & (1ULL << (i % 64))) != 0) {
            mulMod(temp, baseTemp, mod, temp); // result = (result * base) % mod

            // Debug: Print temp after multiplication
            // printf("modExp: temp after multiplication: %llu\n", temp[0]);
        }
    }

    for (int k = 0; k < NUM_SEGMENTS; k++) {
        result[k] = temp[k] % mod[k];
    }
}

// CUDA kernel for Miller-Rabin primality test on segmented large numbers
__global__ void millerRabinTestKernel(uint64_t *numbers, int *results, curandState *state, int numCandidates) {
    int idx = threadIdx.x + blockIdx.x * blockDim.x;
    if (idx >= numCandidates) return;

    curandState localState = state[idx];
    results[idx] = 1; // Assume prime initially

    uint64_t *n = &numbers[idx * NUM_SEGMENTS];
    uint64_t d[NUM_SEGMENTS], x[NUM_SEGMENTS], a[NUM_SEGMENTS];
    int r = 0;

    // Initialize d as n - 1
    for (int i = 0; i < NUM_SEGMENTS; i++) d[i] = n[i];
    d[0] -= 1;  // Subtract 1 (n - 1)

    // Factor d as d * 2^r
    while ((d[0] & 1) == 0) {
        for (int i = 0; i < NUM_SEGMENTS; i++) d[i] >>= 1;
        r++;
    }

    // Miller-Rabin rounds
    for (int round = 0; round < NUM_ROUNDS; round++) {
        // Generate a random base a in the range [2, n-2]
        for (int i = 0; i < NUM_SEGMENTS; i++) a[i] = curand(&localState) % n[i];
        a[0] = 2;

        modExp(a, d, n, x); // Compute x = a^d % n

        // Debug: Print x after modExp
        // printf("Round %d, x after modExp: %llu\n", round, x[0]);

        if (x[0] == 1 || x[0] == n[0] - 1) continue; // Possibly prime

        int continueLoop = 0;
        for (int i = 1; i < r; i++) {
            modExp(x, x, n, x); // x = (x * x) % n
            if (x[0] == n[0] - 1) {
                continueLoop = 1;
                break;
            }
            // Debug: Print x during squaring
            printf("Inner loop %d, x during squaring: %llu\n", i, x[0]);
        }
        if (!continueLoop) {
            results[idx] = 0; // Composite
            return;
        }
    }
}

void initCurandWrapper(curandState *d_state, int numCandidates) {
    initCurand<<<(numCandidates + 255) / 256, 256>>>(d_state, time(0));
}

// Run Miller-Rabin primality test on a range of large numbers (in segments)
void millerRabinTestRange(uint64_t *numbers, int *results, int numCandidates) {
    curandState *d_state;
    cudaMalloc(&d_state, numCandidates * sizeof(curandState));
    initCurandWrapper(d_state, numCandidates);

    uint64_t *d_numbers;
    int *d_results;
    cudaMalloc(&d_numbers, numCandidates * NUM_SEGMENTS * sizeof(uint64_t));
    cudaMalloc(&d_results, numCandidates * sizeof(int));

    cudaMemcpy(d_numbers, numbers, numCandidates * NUM_SEGMENTS * sizeof(uint64_t), cudaMemcpyHostToDevice);

    millerRabinTestKernel<<<(numCandidates + 255) / 256, 256>>>(d_numbers, d_results, d_state, numCandidates);

    cudaMemcpy(results, d_results, numCandidates * sizeof(int), cudaMemcpyDeviceToHost);

    cudaFree(d_state);
    cudaFree(d_numbers);
    cudaFree(d_results);
}
}
