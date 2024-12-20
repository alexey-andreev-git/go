// big_math.cu
#include <cuda_runtime.h>
#include <stdint.h>
#include <stdio.h>

typedef struct modExpParams {
    uint64_t base;
    uint64_t exp;
    uint64_t mod;
    uint64_t result;
} modExpParams;

// CUDA kernel for modular exponentiation (simple example)
__global__ void modExp(modExpParams *params, uint64_t numParams) {
    int idx = blockIdx.x * blockDim.x + threadIdx.x;
    if (idx < numParams) {
        uint64_t res = 1;
        uint64_t b = params[idx].base;
        uint64_t exp = params[idx].exp;
        uint64_t mod = params[idx].mod;
        while (exp > 0) {
            if (exp % 2 == 1) {
                res = (res * b) % mod;
            }
            b = (b * b) % mod;
            exp /= 2;
        }
        params[idx].result = res;
    }
    // uint64_t res = 1;
    // uint64_t b = *base;
    // while (exp > 0) {
    //     if (exp % 2 == 1) {
    //         res = (res * b) % mod;
    //     }
    //     b = (b * b) % mod;
    //     exp /= 2;
    // }
    // *result = res;
}

// Wrapper function to call the kernel
extern "C" void cudaModExp(modExpParams *params, uint64_t numParams) {
    // uint64_t *d_base, *d_result;
    // cudaMalloc(&d_base, sizeof(uint64_t));
    // cudaMalloc(&d_result, sizeof(uint64_t));
    modExpParams *d_params;
    cudaMalloc(&d_params, numParams * sizeof(modExpParams));
    
    // cudaMemcpy(d_base, base, sizeof(uint64_t), cudaMemcpyHostToDevice);
    cudaMemcpy(d_params, params, numParams * sizeof(modExpParams), cudaMemcpyHostToDevice);

    // modExp<<<1, 1>>>(d_base, exp, mod, d_result);  // Launch kernel on 1 block, 1 thread
    modExp<<<(numParams + 255) / 256, 256>>>(d_params, numParams);

    // cudaMemcpy(result, d_result, sizeof(uint64_t), cudaMemcpyDeviceToHost);
    cudaMemcpy(params, d_params, numParams * sizeof(modExpParams), cudaMemcpyDeviceToHost);

    // cudaFree(d_base);
    // cudaFree(d_result);
    cudaFree(d_params);
}
