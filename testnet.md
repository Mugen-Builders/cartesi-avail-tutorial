# Running on a Testnet

## Testing your node locally with a testnet

:::warning

This version of the integration with Avail supports only Sepolia (chainId: 11155111) as the testnet.

:::

- Register your dApp Address

To run on the testnet environment you will need a dApp address on the network. For this we prepared a web page where you can resgister an address for your dApp
https://address.mugen.builders
In the above link you can connect with your wallet and using you public key generate a **dApp address** that will be used in the command that follows.

- Start **_brunodo_** using the command with the flag with the flag that enables integration with Avail;

```bash
brunodo --avail-enabled -d --contracts-input-box-block 6850934 --rpc-url https://sepolia.drpc.org --epoch-blocks 0 --contracts-application-address <dapp-address>
```

## Running the Machine

- In another terminal, create and build a new Cartesi dApp using the following command:

### 1. **Python**

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

### 2. **Rust**

```bash
cartesi create my-dapp --template rust
cd my-dapp
cartesi build
```

- Run the Cartesi Machine Locally on bare metal using the command;

```bash
cartesi-machine --network --flash-drive=label:root,filename:.cartesi/image.ext2 --env=ROLLUP_HTTP_SERVER_URL=http://10.0.2.2:5004 -- /opt/cartesi/dapp/dapp
```

### 3. **Golang**

```bash
cartesi create my-dapp --template go
cd my-dapp
cartesi build
```

- Run the Cartesi Machine Locally on bare metal using the command;

```bash
cartesi-machine --network --flash-drive=label:root,filename:.cartesi/image.ext2 --env=ROLLUP_HTTP_SERVER_URL=http://10.0.2.2:5004 -- /opt/cartesi/dapp/dapp
```

### 4. **Javascript**

```bash
cartesi create my-dapp --template javascript
cd my-dapp
cartesi build
```

- Run the Cartesi Machine Locally on bare metal using the command;

```bash
cartesi-machine --network --flash-drive=label:root,filename:.cartesi/image.ext2 \
--volume=.:/mnt --env=ROLLUP_HTTP_SERVER_URL=http://10.0.2.2:5004 --workdir=/opt/cartesi/dapp -- node index
```
