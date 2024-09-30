package cmd

import (
	"encoding/hex"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/heimdahl-xyz/heimdahl-cli/internal"
	"github.com/spf13/cobra"
	"log"
)

var hardhatDeployCmd = &cobra.Command{
	Use:   "hardhat-deploy",
	Short: "Init hardhat project configuration",
	Run: func(cmd *cobra.Command, args []string) {

		configFile := rootFolder + "/heimdahl.json"

		config, err := internal.ReadConfig(configFile)
		if err != nil {
			log.Println("Error reading config file:", err)
			return
		}

		abis := make(map[string]abi.ABI)
		metas, err := internal.GetContractMetas(rootFolder)
		if err != nil {
			log.Println("Error getting contract metas:", err)
			return
		}

		for _, meta := range metas {
			abis[meta.ContractName] = meta.ABI
		}

		var deployRequests = make(map[string][]internal.ContractDeployRequest)
		for _, network := range config.Networks {
			contractArgs := internal.GetContractsArgs(network.ContractParams)

			addr := internal.GetAddress(network.PrivateKey)

			noncesResp, err := internal.GetNonces(addr, host, apiKey)
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

				deployRequests[addr] = append(deployRequests[addr], internal.ContractDeployRequest{
					ContractName:   meta.ContractName,
					Chain:          network.Chain,
					Network:        network.Network,
					TransactionHex: str,
				})
			}
		}

		for addr, reqs := range deployRequests {
			log.Println("deploying contracts for address:", addr)
			res, err := internal.DeployContracts(reqs, host, apiKey)
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
