# Getting Started

Welcome to the documentation about Zarf, the air-gap tool! Let's get you started on your Zarfing journey!

## Installing Zarf

There are multiple ways to get the Zarf CLI onto your machine including installing from the Defense Unicorns Homebrew Tap, downloading a prebuilt binary from our GitHub releases, or even building the CLI from scratch yourself.

### Installing from the Defense Unicorns Homebrew Tap

[Homebrew](https://brew.sh/) is an open-source software package manager that simplifies the installation of software on macOS and Linux. With Homebrew, installing Zarf is super simple!

```bash
brew tap defenseunicorns/tap
brew install zarf
```

The above command detects your OS and system architecture and installs the correct Zarf CLI binary for your machine. Thanks to the magic of Homebrew, the CLI should be installed onto your $PATH and ready for immediate use.

### Downloading a prebuilt binary from our GitHub releases

All [Zarf releases](https://github.com/defenseunicorns/zarf/releases) on GitHub include prebuilt binaries that you can download and use. We offer a small range of combinations of OS and architecture for you to choose from. Once downloaded, you can install the binary onto your $PATH by moving the binary to the `/usr/local/bin` directory.

```bash
mv ./path/to/downloaded/{ZARF_FILE} /usr/local/bin/zarf
```

### Building the CLI from scratch

If you want to build the CLI from scratch, you can do that too! Our local builds depend on [Go 1.19.x](https://golang.org/doc/install) and are built using [make](https://www.gnu.org/software/make/).
::: note The `make` build-cli` command builds a binary for each combinations of OS and architecture. If you want to shorten the build time, you can use an alternative command to only build the binary you need:

- `make build-cli-mac-intel`
- `make build-cli-mac-apple`
- `make build-cli-linux-amd`
- `make build-cli-linux-arm`

You can learn more about building [here](./4-user-guide/1-the-zarf-cli/1-building-your-own-cli.md).
:::

---

## Verifying Zarf Install

Now that you have installed Zarf onto your path, let's verify that it is working by checking two things! First, we'll check the version of Zarf that has been installed with the command:

```bash
zarf version

# Expected output should look similar to the following
vX.X.X  # X.X.X is replaced with the version number of your specific installation
```

If you are not seeing that, then it's possible that Zarf was not installed onto your $PATH, [this $PATH guide](https://zwbetz.com/how-to-add-a-binary-to-your-path-on-macos-linux-windows/) should help with that.

---

## Where to next?

Depending on how familiar you are with Kubernetes, DevOps, and Zarf, let's find what set of information would be most useful to you!

- If you want to dive straight into Zarf, you can find examples and guides on the [Walkthroughs](./13-walkthroughs/index.md)](./13-walkthroughs/index.md) page.

- More information about the Zarf CLI is available on the [Zarf](./4-user-guide/1-the-zarf-cli/index.md) CLI](./4-user-guide/1-the-zarf-cli/index.md) page, or by browsing through the help descriptions of all the commands available through `zarf --help`.

- More information about the packages that Zarf creates and deploy is available in the [Understanding Zarf Packages](./4-user-guide/2-zarf-packages/1-zarf-packages.md) page.

- If you want to take a step back and better understand the problem Zarf is trying to solve, you can find more context on the [Understand](./1-understand-the-basics.md) the Basics](./1-understand-the-basics.md) and [Core Concepts](./2-core-concepts.md) pages.
