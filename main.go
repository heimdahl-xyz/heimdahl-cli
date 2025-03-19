package main

import "github.com/heimdahl-xyz/heimdahl-cli/cmd"

func main() {
	//var j = ` {"timestamp":1739105584,"from_address":"0xDFd5293D8e347dFe59E90eFd55b2956a1343963d","to_addresss":"0x82D5bfC6E075101e813105374a6EE697437eBdDf","amount":13849870074,"token_address":"0xdAC17F958D2ee523a2206206994597C13D831ec7","symbol":"USDT","chain":"ethereum","network":"mainnet","tx_hash":"0xebc94f69638f2eb152652ff077f34492784d22bd2ce3539744b1ba4972c2a66f","decimals":6,"position":21809080}`
	//
	//var t lib.FungibleTokenTransfer
	//
	//err := json.Unmarshal([]byte(j), &t)
	//if err != nil {
	//	log.Fatal("invalid binary %s", err)
	//}
	//log.Println("t ", t.Amount)
	cmd.Execute() // Run the root command
}
