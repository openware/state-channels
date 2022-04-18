# State channels Smart Contracts

## Download dependencies

1. Run `npm install` in this directory

## Deploy NitroAdjudicator to localhost

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

## .env file

For tests to run an `.env` file must be present with the following variables:

```bash
DEFAULT_GAS=6721975             # as of 03.2022
DEFAULT_GAS_PRICE=20000000000   # as of 03.2022
GANACHE_HOST=0.0.0.0            # localhost
GANACHE_PORT=8561               # ganache port by default
CHAIN_NETWORK_ID=9001           # EVMOS chain id
DEV_HTTP_SERVER_PORT=3000

# These contract addresses get defined in the global jest setup
NITRO_ADJUDICATOR_ADDRESS= 0x0000000000000000000000000000000000000000
COUNTING_APP_ADDRESS= 0x0000000000000000000000000000000000000000
SINGLE_ASSET_PAYMENT_ADDRESS= 0x0000000000000000000000000000000000000000
TRIVIAL_APP_ADDRESS= 0x0000000000000000000000000000000000000000
TEST_FORCE_MOVE_ADDRESS= 0x0000000000000000000000000000000000000000
TEST_NITRO_ADJUDICATOR_ADDRESS= 0x0000000000000000000000000000000000000000
TEST_TOKEN_ADDRESS= 0x0000000000000000000000000000000000000000
```

## Deploy NitroAdjudicator to rinkeby network

Before deploying NitroAdjudicator to rinkeby network, some preparation is needed:

1. set `RINKEBY_URL` env variable to valid link (use infura account)
2. set `RINKEBY_PRIVATE_KEY` env variable as private key from address that will be charged for contract deployment (`Note:` you need to fulfill you address before deploying smart-contract)
3. In rinkeby-accounts.json file you have list of rinkeby addresses, be free to modify this list as you want
(to interact with golang examples, please do not forget to fulfill addresses with funds)

After:

- Open console. Run `npm run contracts:deploy-rinkeby` in this directory. It will deploy NitroAdjudicator on rinkeby network and write its address to addresses.json.

> NOTE: deployed contract addresses will be available in `addresses.json`.
