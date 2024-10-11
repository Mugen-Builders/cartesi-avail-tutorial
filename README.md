# Developing with Cartesi and Avail: A Step-by-Step Guide

## Introduction

This guide will walk you through setting up a Cartesi dApp using Avail on your local machine. You will learn how to send transactions either directly or via Paio using EIP-712, and how to inspect the dApp state and outputs through NoNodo.

## Prerequisites

- Nonodo: A node responsible for communication with Avail and Cartesi.
- Cartesi Machine: Backend logic engine for the dApp.

## Install Tools

### 1. **NONODO:**

- Install globally via npm using the command:

```bash
npm i -g nonodo
```

- Clone the Nonodo repo by using the command:

```bash
git clone https://github.com/Calindra/nonodo
```

### 2. **CARTESI MACHINE:**

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

### 3. **CARTESI CLI**

You can Install the Cartesi Cli using the command;

```bash
npm i -g @cartesi/cli
```

## Development on Local

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

## Interacting via the CLI

- Send a dummy anvil execution that should be picked up by your dapp running in the Cartesi Machine.

```bash
curl --location 'http://localhost:8080/transactions' \
--header 'Content-Type: application/json' \
--data '{
    "signature": "0x373ca4d18d48e1fe3698971968e068e964effb1d36bab4b0204fa8aa1c7449f3517bb4b9b964554e6c81467399355580e0840b426a7855bd0c8e178368c584d61c",
    "typedData": "eyJ0eXBlcyI6eyJDYXJ0ZXNpTWVzc2FnZSI6W3sibmFtZSI6ImFwcCIsInR5cGUiOiJhZGRyZXNzIn0seyJuYW1lIjoibm9uY2UiLCJ0eXBlIjoidWludDY0In0seyJuYW1lIjoibWF4X2dhc19wcmljZSIsInR5cGUiOiJ1aW50MTI4In0seyJuYW1lIjoiZGF0YSIsInR5cGUiOiJzdHJpbmcifV0sIkVJUDcxMkRvbWFpbiI6W3sibmFtZSI6Im5hbWUiLCJ0eXBlIjoic3RyaW5nIn0seyJuYW1lIjoidmVyc2lvbiIsInR5cGUiOiJzdHJpbmcifSx7Im5hbWUiOiJjaGFpbklkIiwidHlwZSI6InVpbnQyNTYifSx7Im5hbWUiOiJ2ZXJpZnlpbmdDb250cmFjdCIsInR5cGUiOiJhZGRyZXNzIn1dfSwicHJpbWFyeVR5cGUiOiJDYXJ0ZXNpTWVzc2FnZSIsImRvbWFpbiI6eyJuYW1lIjoiQXZhaWxNIiwidmVyc2lvbiI6IjEiLCJjaGFpbklkIjoiMHg3YTY5IiwidmVyaWZ5aW5nQ29udHJhY3QiOiIweENjQ0NjY2NjQ0NDQ2NDQ0NDQ0NjQ2NDY2NDY0NDQ2NDY2NjY2NjY0MiLCJzYWx0IjoiIn0sIm1lc3NhZ2UiOnsiYXBwIjoiMHhhYjc1MjhiYjg2MmZiNTdlOGEyYmNkNTY3YTJlOTI5YTBiZTU2YTVlIiwiZGF0YSI6IkdNIiwibWF4X2dhc19wcmljZSI6IjEwIiwibm9uY2UiOiIxIn19"
}'
```

## Interacting via the Frontend

Interacting with your Cartesi dApp via avail utilises EIP712 to sign typed data which is relayed on the users behalf to the avail testnet. We have a demo example available which you can clone, and integrate into the dapp running on your local machine very easily. You can choose to modify this dApp to fit and match your ideal implementation and design.

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
