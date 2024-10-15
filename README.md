**Developing with Cartesi and Avail: A Step-by-Step Guide**
1. [Introduction](#introduction)
2. [Prerequisites](#prerequisites)
3. [Install Tools](#install-tools)
   - 3.1. [Nonodo](#1-nonodo)
   - 3.2. [Cartesi Machine](#2-cartesi-machine)
   - 3.3. [Cartesi CLI](#3-cartesi-cli)
4. [Building and Developing Locally](#building-and-developing-locally)
   - 4.1. [Using Brunodo](#using-brunodo)
   - 4.2. [Python](#python)
   - 4.3. [Rust](#rust)
   - 4.4. [Golang](#golang)
5. [Testing Your Node Locally with a Testnet](#testing-your-node-locally-with-a-testnet)
   - 5.1. [Build Nonodo from Source](#build-nonodo-from-source)
   - 5.3. [Python](#python-1)
   - 5.4. [Rust](#rust-1)
   - 5.5. [Golang](#golang-1)
   
**Interacting with Your dApp**
1. [Interacting via the CLI](#interacting-via-the-cli)
2. [Interacting via the Frontend Template](#interacting-via-the-frontend-template)
3. [Interacting via the NPM Package](#interacting-via-the-npm-package)
4. [Inspecting and Reading Outputs](#inspecting-and-reading-outputs)
   - 4.1. [Inspecting State](#inspecting-state)
   - 4.2. [Querying Outputs](#querying-outputs)

# Developing with Cartesi and Avail: A Step-by-Step Guide

## Introduction

This guide will walk you through setting up a Cartesi dApp using Avail on your local machine. You will learn how to send transactions either directly (through Cartesi Rollups Smart Contracts deployed on an L1) or through Avail DA using EIP-712 signed messages. Also how to inspect the dApp state and outputs through the APIs provided by the Cartesi Rollups Framework.

## Prerequisites

- Nonodo: A node responsible for communication with Avail and Cartesi.
- Cartesi Machine: Backend logic engine for the dApp.

## Install Tools
You can install nonodo in multiple ways. Either through an NPM Package or building from source
### 1. Nonodo

- Install globally via npm using the command:

```bash
npm i -g nonodo
```

- Install brunodo, the experimental version of nonodo through npm using the command:

```bash
npm i -g brunodo
```

- Build from source: Clone the Nonodo repo by using the command:

```bash
git clone https://github.com/Calindra/nonodo
```

### 2. Cartesi Machine

- Download the Cartesi machine for your OS from [this link](https://github.com/edubart/cartesi-machine-everywhere/releases).

**For Linux and macOS:**

- Extract the tar.xz file:

```bash
   tar -xf <filename>.tar.xz
```

Replace `\<filename\>` with the actual name of the file you downloaded.

- Navigate to the extracted directory, rename the extracted folder to `cartesi-machine`:

```bash
cd <cartesi-machine>
```

- Set up environment variables for the Cartesi Machine. You'll need to add the `bin` directory to your system’s PATH so that you can run the Cartesi Machine binaries from anywhere. For Linux or macOS, you can do this by adding the following line to your `\~/.bashrc`, `\~/.bash_profile`, or `\~/.zshrc` file, depending on your shell:

```bash
  export PATH=$PATH:/path/to/cartesi-machine/bin
```

Replace `/path/to/cartesi-machine/` with the actual path to the `bin` folder inside the extracted directory, you can get this by running the command in your terminal while inside the cartesi machine folder: `pwd`. This should print out the path to location of the cartesi-machine folder.

- After adding the line, refresh your terminal configuration by running:

```bash
   source ~/.bashrc
```

Or, if you're using zsh:

```bash
   source ~/.zshrc
```

- Verify the installation by checking if the Cartesi Machine binary is available. You can do this by running:

```bash
   cartesi-machine --help
```

This should display the available options for the Cartesi Machine, indicating that it’s correctly set up.

### 3. Cartesi CLI

You can Install the Cartesi Cli using the command;

```bash
npm i -g @cartesi/cli
```

## Building and developing locally
To test and develop your app locally it is higly recommended that you use `Brunodo` since it contains all the experiemental features from `Nonodo` such as simulating Avail inputs.

With it you can skip a lot of the setup and send EIP-712 messages directly to the node using the nonce and submit endpoints that will be running on `localhost:8080/nonce` and `localhost:8080/submit`

- Create and build a new Cartesi dApp using the following command:

1. **Python**

```bash
cartesi create my-dapp --template python
cd my-dapp
cartesi build
```

- Start **brunodo** using the command;

```bash
brunodo
```

- Run the Cartesi Machine Locally on bare metal using the command;

```bash
cartesi-machine --network --flash-drive=label:root,filename:.cartesi/image.ext2 \
--volume=.:/mnt --env=ROLLUP_HTTP_SERVER_URL=http://10.0.2.2:5004 --workdir=/mnt -- python dapp.py
```

2. **Rust**

```bash
cartesi create my-dapp --template rust
cd my-dapp
cartesi build
```

- Start **brunodo** using the command;

```bash
brunodo
```

- Run the Cartesi Machine Locally on bare metal using the command;

```bash
cartesi-machine --network --flash-drive=label:root,filename:.cartesi/image.ext2 --env=ROLLUP_HTTP_SERVER_URL=http://10.0.2.2:5004/ -- /opt/cartesi/dapp/dapp
```

3. **Golang**

```bash
cartesi create my-dapp --template go
cd my-dapp
cartesi build
```

- Start **brunodo** using the command;

```bash
brunodo
```

- Run the Cartesi Machine Locally on bare metal using the command;

```bash
cartesi-machine --network --flash-drive=label:root,filename:.cartesi/image.ext2 --env=ROLLUP_HTTP_SERVER_URL=http://10.0.2.2:5004 -- /opt/cartesi/dapp/dapp
```

## Testing your node locally with a testnet

For this you will need to build nonodo from source

- cd into the **_nonodo_** cloned repo, then update the content of the .env file with the following.

```bash
AVAIL_RPC_URL<Avail RPC Url>
AVAIL_MNEMONIC=<Mnemonic for wallet>
```

The RPC can be any avail testnet supporting rpc, the default **_nonodo_** uses is:  
` wss://turing-rpc.avail.so/ws`. Then the Mnemonics should be for a wallet that has sufficient Avail test tokens which **_nonodo_** will be using to relay your transactions to the avail network.

- Start **_Nonodo_** using the command;

```bash
./nonodo --avail-enabled -d --sqlite-file db.sqlite3
```

- Create and build a new Cartesi dApp using the following command:

1. **Python**

```bash
cartesi create my-dapp --template python
cd my-dapp
cartesi build
```

- Run the Cartesi Machine Locally on bare metal using the command;

```bash
cartesi-machine --network --flash-drive=label:root,filename:.cartesi/image.ext2 \
--volume=.:/mnt --env=ROLLUP_HTTP_SERVER_URL=http://10.0.2.2:5004 --workdir=/mnt -- python dapp.py
```

2. **Rust**

```bash
cartesi create my-dapp --template rust
cd my-dapp
cartesi build
```

- Run the Cartesi Machine Locally on bare metal using the command;

```bash
cartesi-machine --network --flash-drive=label:root,filename:.cartesi/image.ext2 --env=ROLLUP_HTTP_SERVER_URL=http://10.0.2.2:5004/ -- /opt/cartesi/dapp/dapp
```

3. **Golang**

```bash
cartesi create my-dapp --template go
cd my-dapp
cartesi build
```

- Run the Cartesi Machine Locally on bare metal using the command;

```bash
cartesi-machine --network --flash-drive=label:root,filename:.cartesi/image.ext2 --env=ROLLUP_HTTP_SERVER_URL=http://10.0.2.2:5004 -- /opt/cartesi/dapp/dapp
```

# Interacting with your dApp

## Interacting via the CLI

- Sending transactions such as deposits or generic messages through the layer 1 is done in the same ways as Cartesi Rollups standalone. You can use `cast`, the `cartesi cli` or other approaches. You can follow them here in the [docs](https://docs.cartesi.io/cartesi-rollups/1.5/development/send-requests/)

- To send a dummy anvil execution that should be picked up by your dapp running in the Cartesi Machine you can run the following command:

```bash
curl --location 'http://localhost:8080/submit' \
--header 'Content-Type: application/json' \
--data '{
    "signature": "0x373ca4d18d48e1fe3698971968e068e964effb1d36bab4b0204fa8aa1c7449f3517bb4b9b964554e6c81467399355580e0840b426a7855bd0c8e178368c584d61c",
    "typedData": "eyJ0eXBlcyI6eyJDYXJ0ZXNpTWVzc2FnZSI6W3sibmFtZSI6ImFwcCIsInR5cGUiOiJhZGRyZXNzIn0seyJuYW1lIjoibm9uY2UiLCJ0eXBlIjoidWludDY0In0seyJuYW1lIjoibWF4X2dhc19wcmljZSIsInR5cGUiOiJ1aW50MTI4In0seyJuYW1lIjoiZGF0YSIsInR5cGUiOiJzdHJpbmcifV0sIkVJUDcxMkRvbWFpbiI6W3sibmFtZSI6Im5hbWUiLCJ0eXBlIjoic3RyaW5nIn0seyJuYW1lIjoidmVyc2lvbiIsInR5cGUiOiJzdHJpbmcifSx7Im5hbWUiOiJjaGFpbklkIiwidHlwZSI6InVpbnQyNTYifSx7Im5hbWUiOiJ2ZXJpZnlpbmdDb250cmFjdCIsInR5cGUiOiJhZGRyZXNzIn1dfSwicHJpbWFyeVR5cGUiOiJDYXJ0ZXNpTWVzc2FnZSIsImRvbWFpbiI6eyJuYW1lIjoiQXZhaWxNIiwidmVyc2lvbiI6IjEiLCJjaGFpbklkIjoiMHg3YTY5IiwidmVyaWZ5aW5nQ29udHJhY3QiOiIweENjQ0NjY2NjQ0NDQ2NDQ0NDQ0NjQ2NDY2NDY0NDQ2NDY2NjY2NjY0MiLCJzYWx0IjoiIn0sIm1lc3NhZ2UiOnsiYXBwIjoiMHhhYjc1MjhiYjg2MmZiNTdlOGEyYmNkNTY3YTJlOTI5YTBiZTU2YTVlIiwiZGF0YSI6IkdNIiwibWF4X2dhc19wcmljZSI6IjEwIiwibm9uY2UiOiIxIn19"
}'
```

## Interacting via the Frontend Template

We have a demo example available which you can clone, and integrate into the dapp running on your local machine very easily. You can choose to modify this dApp to fit and match your ideal implementation and design. 
It contains a ways to send many different types of input. Including interacting with your Cartesi dApp via avail, which utilises EIP-712 to sign typed data which is relayed on the users behalf to the avail testnet. 

### Installation

- Clone the frontend repo integrated with EIP 712 by using this command:

```bash
git clone https://github.com/Calindra/frontend-web-cartesi
```

- Install all the necessary dependencies by running this commands:

```bash
cd frontend-web-cartesi
git checkout feature/refactor-simplification
yarn install
```

- Generate the necessary rollup data’s by running this command;

```bash
yarn codegen
```

- Start the frontend application by running:

```bash
yarn dev --port 3000
```

- Finally open your browser and navigate to the URL where your frontend dapp is running, you can now interact with your dapp running on local by signing and sending data to your dapp via Avail.
- To send data via avail use the “Send L2 EIP-712 Input” form in the Input section.

## Interacting via the the NPM Package

Interacting with your Cartesi dApp using the `@mugen-builders/client` npm package allows you to send data via EIP-712 or directly using signed inputs. This package simplifies the process of relaying data to Cartesi dApps, providing flexibility to work with both EIP-712 formatted data and standard inputs. 
You can check the description of the function in the package's [page](https://www.npmjs.com/package/@mugen-builders/client)

### Installation

- Install the npm package by running the following command:

```bash
npm install @mugen-builders/client@0.1.2-rc1.0
```

### Usage

To integrate the package into the front-end of your dApp, use the `advanceEIP712` and `advanceInput` methods to handle both EIP-712 typed data and simple input data(which goes through the L1). Below is an example implementation:

```javascript
import { advanceInput, advanceEIP712 } from "@mugen-builders/client";

const addInput = async (_input) => {
    const provider = new ethers.providers.Web3Provider(wallet.provider);
    const signer = provider.getSigner();

    // For EIP-712 input
    let availInput = await advanceEIP712(
        signer, 
        provider, 
        dappAddress, 
        _input, 
        { cartesiNodeUrl: "http://localhost:8080" }
    );

    // For simple input
    let l1Input = await advanceInput(
        signer, 
        dappAddress, 
        _input, 
        { inputBoxAddress: "0x58Df21fE097d4bE5dCf61e01d9ea3f6B81c2E1dB" }
    );
};
```
The return of `advanceEIP712` will be the same as `advanceInput`. Both methods will return lists with reports, notices and vouchers generated from that input, allowing you to interact with your dApp using the provided data. 

This simplifies interaction with your dApp, providing an easy way to handle both types of inputs.

## Inspecting and reading outputs

### Inspecting state
Inspecting the state of your dApp though `handle_inspect` function is done in the same way as using Cartesi Rollups standalone. You can refer to the [docs](https://docs.cartesi.io/cartesi-rollups/1.5/development/send-requests/#make-inspect-calls)

### Querying outputs
Querying outputs directly is the exact same as using Cartesi Rollups standalone. You can refer to the [docs](https://docs.cartesi.io/cartesi-rollups/1.5/rollups-apis/graphql/overview/)

To query outputs from a specific the process is very similar to using Cartesi Rollups Standalone. You can refer to the [docs](https://docs.cartesi.io/cartesi-rollups/1.5/rollups-apis/graphql/overview/) to read more.
The big difference is the output format. Instead of querying inputs through the `index` field, you query them through an `id` field.

This id field can come in two ways:
- It is a hex value returned from `/submit` endpoint when the input comes from and EIP-712 signed message 
- It is string containing a scalar integer value that can be found inside the events emitted by the `inputBox` contract when sending the transaction through the layer 1.

#### Example Queries
##### Listing inputs
```graphql
query {
  inputs(first: 30) {
    edges {
      node {
        id
        index
        status
        blockTimestamp
        msgSender
        payload     
      }
    }
  }
}
```

##### Getting a specific input through its `id`
```graphql
query {
  input(id: "0x67b8ca2e94df99f1c8b31f073c2e2aea1c6e86b979b9fe1178270a963b466cb6") {
    id
    index
    status
    blockTimestamp
    msgSender
    payload
    notices{
      edges{
        node {
          payload
        }
      }
    }
    reports{
      edges{
        node {
          payload
        }
      }
    }
    vouchers{
      edges{
        node {
          payload
        }
      }
    }
  }
}
```
