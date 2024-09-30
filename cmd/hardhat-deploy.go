package cmd

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/heimdahl-xyz/heimdahl-cli/internal"
	"github.com/spf13/cobra"
	"io"
	"log"
	"math/big"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

type ABI struct {
	Inputs []struct {
		InternalType string `json:"internalType"`
		Name         string `json:"name"`
		Type         string `json:"type"`
	} `json:"inputs"`
	StateMutability string `json:"stateMutability,omitempty"`
	Type            string `json:"type"`
	Anonymous       bool   `json:"anonymous,omitempty"`
	Name            string `json:"name,omitempty"`
	Outputs         []struct {
		InternalType string `json:"internalType"`
		Name         string `json:"name"`
		Type         string `json:"type"`
	} `json:"outputs,omitempty"`
}

type ContractArtifact struct {
	Format                 string   `json:"_format"`
	ContractName           string   `json:"contractName"`
	SourceName             string   `json:"sourceName"`
	Abi                    []ABI    `json:"abi"`
	Bytecode               string   `json:"bytecode"`
	DeployedBytecode       string   `json:"deployedBytecode"`
	LinkReferences         struct{} `json:"linkReferences"`
	DeployedLinkReferences struct{} `json:"deployedLinkReferences"`
}

type ContractDeployRequest struct {
	ContractName   string `json:"contract_name"`
	TransactionHex string `json:"transaction_hex"`
	Chain          string `json:"chain"`
	Network        string `json:"network"`
}

type ContractMeta struct {
	ContractName string
	ABI          abi.ABI
	Bytecode     string
}

func readContractArtifact(path string) (*ContractMeta, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("Error reading contract artifact file: %w", err)
	}

	var contr ContractArtifact
	err = json.Unmarshal(b, &contr)
	if err != nil {
		return nil, fmt.Errorf("Error parsing contract artifact file: %w", err)
	}

	b, err = json.Marshal(contr.Abi)
	if err != nil {
		return nil, fmt.Errorf("Error parsing contract artifact file: %w", err)
	}

	ab, err := abi.JSON(bytes.NewReader(b))
	if err != nil {
		return nil, fmt.Errorf("Error parsing ABI: %w", err)
	}

	return &ContractMeta{
		ContractName: contr.ContractName,
		ABI:          ab,
		Bytecode:     contr.Bytecode,
	}, nil
}

func readConfig(configFile string) (*Config, error) {
	b, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("Error reading config file: %v. "+
			"You can init it with heim-cli hardhat-init...\n", err)
	}

	var config Config
	err = json.NewDecoder(bytes.NewReader(b)).Decode(&config)
	if err != nil {
		fmt.Printf("Error parsing config file: %v", err)
		return nil, err
	}
	return &config, nil
}

type NonceResponse struct {
	Nonces map[string]uint64 `json:"nonces"`
}

func getNonces(address string) (*NonceResponse, error) {
	url := fmt.Sprintf("%s/v1/address/nonce?addresses=%s", host, address)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("Error creating get request: %s", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error getting nonce: %s", err)
	}
	defer resp.Body.Close()

	if !(resp.StatusCode == http.StatusCreated || resp.StatusCode == http.StatusOK) {
		return nil, fmt.Errorf("Failed to get nonce: %s", resp.Status)
	}

	var nonceResponse NonceResponse
	err = json.NewDecoder(resp.Body).Decode(&nonceResponse)
	if err != nil {
		return nil, fmt.Errorf("Error decoding nonce response: %s", err)
	}

	return &nonceResponse, nil
}

type DeployResult struct {
	ContractAddress string `json:"contract_address"`
	TransactionHash string `json:"transaction_hash"`
	BlockHash       string `json:"block_hash"`
	BlockNumber     uint64 `json:"block_number"`
	GasSpent        uint64 `json:"gas_spent"`
	GasPrice        uint64 `json:"gas_price"`
	Chain           string `json:"chain"`
	Network         string `json:"network"`
}

type DeploymentResponse struct {
	DeploymentResult *DeployResult `json:"deployment_result"`
	Error            string        `json:"error"`
}

func getAddress(pkk string) string {
	var pk string
	if pkk[:2] == "0x" {
		pk = pkk[2:]
	} else {
		pk = pkk
	}

	privateKey, err := crypto.HexToECDSA(pk)
	if err != nil {
		log.Fatalf("Failed to load private key: %v", err)
	}

	// Create an authorized transactor
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatalf("Failed to get public key from private key")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	return fromAddress.String()
}

func getContractFiles(rootFolder string) ([]string, error) {
	contractFiles := []string{}
	err := filepath.WalkDir(rootFolder+"/artifacts/contracts", func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// Check if the file has a .json extension and it's not a directory
		if !d.IsDir() && filepath.Ext(d.Name()) == ".json" && filepath.Ext(d.Name()) != ".dbg.json" {
			if strings.Contains(d.Name(), ".dbg.json") {
				return nil
			}
			contractFiles = append(contractFiles, path)
		}
		return nil
	})

	if err != nil {
		log.Println("Error walking the directory:", err)
		return nil, err
	}

	return contractFiles, nil
}

func deployContracts(reqs []ContractDeployRequest) (string, error) {
	url := fmt.Sprintf("%s/v1/contracts/deploy", host) // Use the global host variable

	reqd, err := json.Marshal(reqs)
	if err != nil {
		log.Println("Error marshalling deploy requests:", err)
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqd))
	if err != nil {
		log.Printf("Error creating post request %s:", err)
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Error sending deploy requests:", err)
		return "", err
	}
	defer resp.Body.Close()

	if !(resp.StatusCode == http.StatusCreated || resp.StatusCode == http.StatusOK) {
		log.Println("Failed to deploy contracts", resp.Status)
		return "", err
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		return "", err
	}

	return string(b), nil
}

func getContractMetas(rootFolder string) ([]*ContractMeta, error) {
	contractFiles, err := getContractFiles(rootFolder)
	if err != nil {
		log.Println("Error walking the directory:", err)
		return nil, err
	}

	metas := []*ContractMeta{}
	for _, contractFile := range contractFiles {
		contractMeta, err := readContractArtifact(contractFile)
		if err != nil {
			fmt.Println("Error reading contract artifact:", err)
			return nil, err
		}

		metas = append(metas, contractMeta)
	}
	return metas, nil
}

func ConvertToSolidityType(value interface{}) interface{} {
	switch v := value.(type) {
	case int, int8, int16, int32, int64:
		return big.NewInt(reflect.ValueOf(v).Int()) // Convert to *big.Int
	case uint, uint8, uint16, uint32, uint64:
		return new(big.Int).SetUint64(reflect.ValueOf(v).Uint()) // Convert to *big.Int
	case string:
		return v // Strings are directly compatible
	case bool:
		return v // Booleans are directly compatible
	case float32, float64:
		// Solidity does not have floating-point types. This is a potential issue.
		// Converting to big.Int might cause loss of precision, so handle it carefully.
		log.Printf("Warning: Solidity does not support floating-point types. Converting %v to int.", v)
		return big.NewInt(int64(reflect.ValueOf(v).Float())) // Convert to *big.Int
	case *big.Int:
		return v // Already the correct type
	default:
		log.Fatalf("Unsupported type: %T. Only integers, strings, booleans, and *big.Int are supported.", value)
		return nil
	}
}

func getContractsArgs(params []ContractParams) map[string][]interface{} {
	out := make(map[string][]interface{})
	for _, param := range params {
		for _, arg := range param.Arguments.Inputs {
			out[param.ContractName] = append(out[param.ContractName], ConvertToSolidityType(arg.Value))
		}
	}
	return out
}

var hardhatDeployCmd = &cobra.Command{
	Use:   "hardhat-deploy",
	Short: "Init hardhat project configuration",
	Run: func(cmd *cobra.Command, args []string) {

		configFile := rootFolder + "/heimdahl.json"

		config, err := readConfig(configFile)
		if err != nil {
			log.Println("Error reading config file:", err)
			return
		}

		abis := make(map[string]abi.ABI)
		metas, err := getContractMetas(rootFolder)
		if err != nil {
			log.Println("Error getting contract metas:", err)
			return
		}

		for _, meta := range metas {
			abis[meta.ContractName] = meta.ABI
		}

		var deployRequests = make(map[string][]ContractDeployRequest)
		for _, network := range config.Networks {
			contractArgs := getContractsArgs(network.ContractParams)

			addr := getAddress(network.PrivateKey)

			noncesResp, err := getNonces(addr)
			if err != nil {
				log.Println("Error getting nonces:", err)
				return
			}

			nonce := noncesResp.Nonces[addr+"."+network.Chain+"."+network.Network]
			for _, meta := range metas {
				contractByteCode := common.FromHex(meta.Bytecode)

				if len(contractArgs[meta.ContractName]) > 0 {

					argBytes, err := abis[meta.ContractName].Pack("", contractArgs[meta.ContractName]...)
					if err != nil {
						log.Println("Error packing contract arguments:", err)
						return
					}
					contractByteCode = append(contractByteCode, argBytes...)
				}

				str, err := internal.CreateSignedTransactionHex(internal.TransactionInput{
					BytecodeHex:   hex.EncodeToString(contractByteCode),
					NetworkID:     int64(network.ChainID),
					PrivateKeyHex: network.PrivateKey,
					GasLimit:      1000000,
					GasPrice:      1000000000,
					Nonce:         nonce,
				})
				nonce += 1
				if err != nil {
					log.Println("Error creating signed transaction hex:", err)
					return
				}

				deployRequests[addr] = append(deployRequests[addr], ContractDeployRequest{
					ContractName:   meta.ContractName,
					Chain:          network.Chain,
					Network:        network.Network,
					TransactionHex: str,
				})
			}
		}

		for addr, reqs := range deployRequests {
			log.Println("deploying contracts for address:", addr)
			res, err := deployContracts(reqs)
			if err != nil {
				log.Println("Error deploying contracts:", err)
				return
			}
			log.Println("Deployment response:", res)
		}
	},
}

func init() {
	hardhatDeployCmd.Flags().StringVarP(&rootFolder, "rootFolder", "r", ".", "Hardhat project root folder (optional)")
	_ = hardhatDeployCmd.MarkFlagRequired("networks")
}
