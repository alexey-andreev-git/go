package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"sync"
	"time"
)

var (
	prime_nums = []string{
		"de5b679e01f588063c3da59ff270dbb2baa72f4eb93ceefe06bbde9fca9d888d739c6621b4b8758d55a13dab4730cca6b0e0dda87ea68836cd6049e8a352cd9632c14740c24aa0e24e3a2d662c807edbd0c6392650112075a1fe68fd7367a8ff31b917f713ea772daa3f9bb7e89c17c299d2b1b4421e855b32dac504ac9baa6e037f16038d1cabb17c08be4fbaeb242a9f76568cdfc2666a6a54d430c7b3ff1da4b90da607f806752151681c65d5231edd82b2032063866c388cf5f2493a38a8c292149ed572d1fc534eb57ac22c6d8a3e688102d6d3f251f274d0620a90d8956b52dbc8b21c3225fe060ae051865a076c24de0f144f069db2493036fb8e26a3c1255343d79f74655c7890135a770d3dad53a9625eb6e3f7cba1955dad9f2bb1bf11c9a047b795a967f678ab862fed73e32d3c2068d63bad2c51b032b9aa1ac16dbc6bc3883e703bdb8f6eaeb4b95078c5010d94541be98ea2a3f675ff7c5d457a5714d6b5a40927d8527303f4e476603257e9bf7dbb73adec32841182380a54d6072961e905ecc356a14292a2bdca2082422aa433fadf680bc22512244be8ca800abd9b16595068f52e3d67effd97f00514bd38280411f5584641c35d5cbb1461724e61163024ccb21262807060d4f7661ffdf4ef19c91196540bbd4d67e3be146e682388e59e33c004ef96221fae01beab009bc8cb511112bc7345125c959d",
		"df6562fb84ccf1d3d008beb7086590d9fbc940966eecae7b7f00ff99cb993e3903713a04e6491031680b96e157fbc8b8e893e8e597ab1340d13aa5ad3f52b6613afc8a5d8b9e270b8fb084fdf5983d6f5bd4256b079fb549dd98b74a647cdcd75310ce5c392bbadb31960ec9674f60aba1b61890ad50acc9b21305041195cf460b772748a2a4bc000a13bccbd789603ef6a2ce979783e0ac8edbe73e3668cf5f61027dedf9273f2072c47a8b0f429bcee7ce59feeb3136486d4452ad5522fe75998c76e7ee6b86212ee51efc250798bdeff93a190ca0b94233fed1ff0cc061cd7e2737c31b459588016034e172c5a825a3d45d665833ea45b85ed39ec1e34e948a227caf58a1c28da85ef0a373e9f5ce3c8e34ff4bab376a28f9c244ef7ec0be9a5e4629e7045f7ca015deeba2b8b171419cff463918ad3a48bdc67bf98335dec56309dbffbe0fa362a2ba352d45c8174a008ec71c8bfc8fe01bd40ac2d392ca558781719706afd119e2f8cb752e25de6ef68c047b4cf5c1c15a482799fb10383e5eec6d42b189d4f074c16cc70a7dec2c11efe6d6591199f1d32baf76e4f0c8332abd3a34b14893d1003399aa327cde73880393f7cc2d675ba0d20fa574de9c457c62da5d43785e094325299408554eb46c452991097f12dddadc421ede7cc36dc79741a24822580cc1dc97558213686498849e8647aad74cb694cdc8f027ef",
		"cd2627d888968faa5c09947cb7696a38a1b90cae4bd1ae5bd6c318d2e186a93e5b1f01d6b872adf4036eb22d78a46111e618fe5e019727bb3e73ea378242a20c7e4334a4e6e328ed030f327fb1be9788b6cecec0a15b0e527d63df22fb165647310370f3656376edb7709656c172945fdb754995335e140d95cfff98720dbb40db4d82f31946127e1eaba3c182c0c7a158bef4937a6d61e6a9c0eb6d0a65acd0e5b9af9840089475f725623e27b168357de159a61797b9b982ef617400db60609324c79ad110a54f6e8cb25e74bd65655814257b00fc4647ce9ca7e15a36db3d2ff2d62a4fc556841a110bcd965299fadba82bddc1e86ce182b90727cb5b80f8cf406288bb8c2d40e3342393cb9d67d1f14bcd53ca17a56b61a4e8a0c5dda43bf3c598dbb08ebe06790e1748970d0e2d9b44d08be11a4672463008319a84487662c0a67bcdb569c38152629e2ff99196d0e881579c1d366ba994c8fd4cbaeea402682875e1daff7c4cc160f6693d031644d1c855c1dac3e45c1c0350329b8e422576862f55e968b46065769e0042d62d1126a47969ea32a5d68c83657974046ad24e97d02b5d0f54776b0f38adb2db3706248cded2223934939c8c91734ad37721040c25b8bf7580d14fc4261e7e1ce0306a96b4ab7b0284b39591286cee98e13d9c8fd171c680d8064b7e676c481b79d3dcd2ff2964527a2040f2863239383b",
		"c12db86840a328db4e0a9ca82d2a34b0b1bb22e6c2249bb707bd30074cacd14e72c2758471cb8864573dbe5d85374f31722fb1d9314de17e4fd19ecea78cf6e5a30f5d09d874c2995b64be77dd62a69412d4bb75292c86ec6d8429ab6d60c886e891eae8ead7f6222cbad8b8efaff69d6540494af50de13c7f9f3cda80f0fdfbff9d6952e5273fb9501ccfc250ed1e211eb250361f19833798dcccaa3eeb144c2f1372fb4c2931ecd0c8bb89da9babb5aecc171aa4a8f1489024dca843a3805b018bde55d21628d6356a37834e3e20032c4c125870c7294914069777abc3e00ff256f247530ebea0d8cb61c768231a08953c9f611e6efbac3c0400da03bb3f00b527233455406db69b1bb1895339bce2cb0e43c967e02f7343392e6dac87314e7cfa437267d5f8368a78af5d4e01ee1fbd2c676b49f47e17814b9b4314d71111bc0b336322343ed8660ca38269729eaaae3013618a7c4d6c93ab81110116259e7b0860cae3e577da3b0894280667e39859b9ed75121a31a02e220fc1379673938189eef6b0f2d38c3cb22224686044fd1138eda06bd75969fe1a117630e0d591f32e6ae45dbd26b74a919c7cb7a3b61428b38a9f16272d4d49f16f11ebab0fe232598115a99e7059220603a96af8c23cf408a4e6708cc7f4fbc61e69d2bae0c974eb228659947e51b761f8ce3b4639dce0b611d4c010497a18e03516b56f0b53",
		"cb7e0f7c22fd0bd4e6dce7d1deecaf88c754a2be79e2370710c6332ff4dc2b59adb100c362ce0805fc0bc752aa9d57d0c22f9d329044addb6651aa4521df068d69f574ac0a73d9b2a72509ba7203d4d192bf830bd38c38de87f13ba6c024557c3f196e65109c6bcd5515a9a02fa9a250a08657c8428ad39d2d14c31ef3d47ddce5534a0b83ec00b42dd766399bc8035691db2012ddefa9716d6612b74aa674253af4fb2b26c96839b17d00df12f62a6275ed7607f3fc8d340ed63619cf8aa01bc67ea678d853890acfc2658925ed11ca0908571985f21407d1a9bc22e35eb274e535c55c4718c882c959069335a352b81ca249f227319d0924d42b37e7f1598657ea43f52d029cd33ba6858236492f704237aaea7a5b3757c927cb057990c5f401f7a509f698f0bd10c7629401be32503484355522965cab87610938c3109f87e6addd2146907451b81158d55d7267a92226293f47987c82df8af0dce3a24d8b67150572b5e719a9fcd7ce43b50afa2863103a85d2a1896fcef1b7adabc596e37fa76fb8fcc3e0b17a341a4fcd3c4369a58482e0aa60f2217bff975e9e25a4d964abccdff37c2cc3510085bbfb3380527e25d4d297ac185f94b99b3c3d9a904ba1f85e55e445b41d6d321907d42f63c143ed594a478140912e0581ee39717f5a99cd4f4ec1aed19babea91736b80825eeeebde8e5c747db627b5ef552daefbb3",
	}
	primeFile = "prime_numbers.txt"
	mutex     sync.Mutex
)

// Function to check if a number exists in the prime_nums array
func isNumberInArray(number string) bool {
	for _, prime := range prime_nums {
		if prime == number {
			return true
		}
	}
	return false
}

// Function to create a big integer from a 512-byte array
func bigIntFromBytes(bytes []byte) *big.Int {
	n := new(big.Int)
	n.SetBytes(bytes)
	return n
}

func savePrimeToFile(primeHex string) {
	mutex.Lock()
	defer mutex.Unlock()

	file, err := os.OpenFile(primeFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("\"%s\",\n", primeHex))
	if err != nil {
		fmt.Println("Error writing to file:", err)
	}
}

func generateAndCheckPrimes(wg *sync.WaitGroup, bitLength int) {
	defer wg.Done()
	// bitLength := 4096

	for i := 0; i < 50; i++ {
		prime, err := rand.Prime(rand.Reader, bitLength)
		if err != nil {
			fmt.Println("Error generating prime:", err)
			continue
		}

		primeHex := prime.Text(16)

		mutex.Lock()
		if isNumberInArray(primeHex) {
			mutex.Unlock()
			continue
		}
		prime_nums = append(prime_nums, primeHex)
		mutex.Unlock()

		savePrimeToFile(primeHex)
	}
}

func test_div() {
	// Example: 512-byte numbers (random values)
	// Here we use fixed byte arrays for demonstration purposes.
	// In a real case, the numbers could be larger and read from files, input, etc.
	num1Hex := "FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF"
	num2Hex := "02" // A simple number for division

	// Convert 512-byte hex strings to big.Int
	num1Bytes, err := hex.DecodeString(num1Hex)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	num2Bytes, err := hex.DecodeString(num2Hex)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	num1 := bigIntFromBytes(num1Bytes)
	num2 := bigIntFromBytes(num2Bytes)

	// Check for zero to avoid division by zero
	if num2.Cmp(big.NewInt(0)) == 0 {
		fmt.Println("Error: divisor is zero")
		return
	}

	// Perform division
	quotient := new(big.Int).Div(num1, num2)
	remainder := new(big.Int).Mod(num1, num2)

	// Print results
	fmt.Println("Quotient:", quotient.Text(16))   // Output quotient in hexadecimal format
	fmt.Println("Remainder:", remainder.Text(16)) // Output remainder in hexadecimal format
}

func test_rnd(bitLength int) {
	// Generate a 4096-bit prime number
	prime, err := rand.Prime(rand.Reader, bitLength)
	if err != nil {
		fmt.Println("Error generating prime:", err)
		return
	}

	// Convert the prime number to hexadecimal string
	numberHex := prime.Text(16)

	// Check if the generated prime number is already in the array
	if isNumberInArray(numberHex) {
		fmt.Println("The generated prime number is already in the array.")
	} else {
		fmt.Println("The generated prime number is not in the array.")
	}

	// Output the prime number in hexadecimal format
	fmt.Println("Generated 4096-bit prime number:", numberHex)
}

func main() {
	// Test division
	test_div()

	startTime := time.Now()
	// Define the bit length of the prime number
	bitLength := 4096
	test_rnd(bitLength)
	elapsedTime := time.Since(startTime)
	fmt.Println("Elapsed time:", elapsedTime)

	// var wg sync.WaitGroup
	// Initialize 100 goroutines, each generating 50 primes
	// for i := 0; i < 100; i++ {
	// 	wg.Add(1)
	// 	go generateAndCheckPrimes(&wg, bitLength)
	// }

	// wg.Wait()
	fmt.Println("Prime generation completed.")

}
