# State channels Smart Contracts

## Download dependencies

1. Run `npm install` in this directory

## Deploy NitroAdjudicator

1. Open console 1. Run `npm run contracts:node` in this directory. This will start Hardhat Network.
2. Open console 2. Run `npm run contracts:deploy-localhost` in this directory. It will deploy NitroAdjudicator on localhost network and write its address to addresses.json.
3. Don't close console 1. While it is running, you can communicate with the contract deployed.

> NOTE: deployed contract addresses available in `addresses.json` file in such format:

```json
{
  "chainId_value": [
    {
      "chainId": "string",
      "name": "string",
      "contracts": {
        "contractName": {
          "address": "hex"
        }
      }
    }
  ]
}
```
