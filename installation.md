# Installation
## Prerequisites

- Nonodo: A node responsible for communication with Avail and Cartesi.
- Cartesi Machine: Backend logic engine for the dApp.

## Install Tools

You can install nonodo in multiple ways. Either through an NPM Package or building from source

### 1. Nonodo

- Install brunodo, the experimental version of nonodo through npm using the command:

```bash
npm install -g brunodo
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

:::warning

For Mac users, running this command might trigger a prompt from Mac built in security feauture. Check this part of the [troubleshooting section](./troubleshooting#1-cartesi-machine-blocked-by-mac-security-feautures) for more guidelines on how to resolve this.

:::

### 3. Cartesi CLI

You can Install the Cartesi Cli using the command;

```bash
npm i -g @cartesi/cli
```