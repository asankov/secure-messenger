# Secure Messaging System

Secure Messaging System is a CLI tool that allows you to send and receive encrypted messages.

## Usage

The `secure-messenger` CLI exposes a few command that allow you to send and receive encrypted messages.

### generate-key

In order to send and receive encrypted messages you need a secret key.

You can generate one via the `generate-key` command:

```console
secure-messenger generate-key
```

By default, the key will be save into the OS Keychain.
This is the most secure option, as the Keychain can have native encryption and access control.

All commands support using the key from the keychain, so you never really have to see it and use it manually.

If, however, you do need to see it or have it in another location, you can retrieve it via the `retrieve-key` command:

```console
$ secure-messenger retrieve-key
Y1Ywl91T6xHTSwuaJL3tCwc7taHYjmTT
```

or you can output it to `stdout` during the generation:

```console
$ secure-messenger generate-key --output-to-stdout
Consider using the keychain to store the key, so that it don't get lost or exposed.
1nteEHHUjDUMoilUoquO8gSAfI7Btd5A
```

you can also redirect the output and save it into a file:

```console
$ secure-messenger generate-key --output-to-stdout > key.file
Consider using the keychain to store the key, so that it don't get lost or exposed.
$ cat key.file
9ekJni2mWM0BUhnfvg2eq5F7xcAaqPzM
```

You can also generate 16-byte, 24-byte and 32-byte long keys.

```console
$ secure-messenger generate-key --output-to-stdout --key-size=16
Consider using the keychain to store the key, so that it don't get lost or exposed.
iPUML2/JXJ96NrW+

$ secure-messenger generate-key --output-to-stdout --key-size=24
Consider using the keychain to store the key, so that it don't get lost or exposed.
h/LVOSopnqb/e2maMxAjIQ==

$ secure-messenger generate-key --output-to-stdout --key-size=32
Consider using the keychain to store the key, so that it don't get lost or exposed.
9313MMyVGPuuX29Pz6N5rWSVGWSn9CN1
```

You need to share this key with the person you want to exchange messages with.

For more info on how to share the key securely, go to [Key Exchange](#key-exchange).

### encrypt

Once you have a key you can use it to generate encrypted messages.

```console
$ secure-messenger encrypt --sender-id=abc123 --receiver-id=987 --payload="Hello from Planet Earth"
6Vq4bYb+9yK+wB8P22evR+FTYUZG6zksB7tAzqgnpv66Z2y9f1fJ6UcXZtKPc3Sm9SkwiBg/fXTqLPokvw178WxAqqa3JtRdvUGRr4Ksp/ABXF06IyX48EaIhDAivM4sACYYiditNkLoGyz0b3685yFgMxLc1K7f0Ce13dDuYQ==
```

In order to generate a message you need to input 4 arguments:

- `sender-id` - this is the ID of the sender (your ID)
- `receiver-id` - this is the ID of the receiver
- `payload` - the payload you want to send

The secret key will automatically be retrieve from the keychain if it exists there.

If you want to provide the secret key yourself you can use the `--secret-key` and `--secret-key-file` arguments.

The output of the command is the encrypted message, which you can safely share with the other side, without worrying that someone
might intercept it and read the contents.

### decrypt

Once you have receive a message and you have the key you can use the `secure-messenger` CLI to decrypt it and read it.

```console
$ secure-messenger decrypt UIqwQxl9ntJc9SIaJidJOrx6QrgCrJbr7Jy8rkS3BjVIE6TofnO4ljW6mIc4Eo8CoM/w9rSYNaRvTbtArEyMxebrVJrH0xcKxZhJnAKC3A83EXL+rfh9+wNki6DH/aKqt1XnzoajK6lH1Bep2O74oR8aRNwgIRYN4R9GdtPYvQ==

{"senderId":"abc123","receiverId":"987","payload":"Hello from Planet Earth","timestamp":1700317704}
```

The secret key will automatically be retrieve from the keychain if it exists there.

If you want to provide the secret key yourself you can use the `--secret-key` and `--secret-key-file` arguments.

The output of the command is the encrypted message, which you can safely share with the other side, without worrying that someone
might intercept it and read the contents.

TODO: passing long encrypted message as a CLI arguments is not a good UX.
Support a better way like reading the messages from file or from stdin.

Voila. You have exchange an encrypted message.

### Key Exchange

In order to securely exchange the secret key with the other party, you can again use the CLI.

The two command to do that are `exchange-key` and `exchange-key-server`.

The other side (let's call them Bob, and we are Alice) need to run the `exchange-key-server` command.
This will start a web server that will listen for a request that will start the key exchange process.

```console
$(bob) secure-messenger exchange-key-server
Starting exchange key server on [:8080]
```

Then we (Alice) can use the `exchange-key` command to send a request to Bob's server that will initiate the key-exchange process.

```console
$(alice) secure-messenger exchange-key --remote-addr=http://localhost:8080 --secret-key-file=key.file
successfully exchanged secret key
```

Meanwhile Bob sees this:

```console
$(bob) secure-messenger exchange-key-server
Starting exchange key server on [:8080]
time=2023-11-18T18:28:09.647+02:00 level=INFO msg="Retrieved secret key"
time=2023-11-18T18:28:09.648+02:00 level=INFO msg="Secret key stored" location=keychain
```

This means that the key-exchange was successful and now Bob has the secret key stored in its keychain.

These command use the [Diffie-Helmman algorithm](https://en.wikipedia.org/wiki/Diffie%E2%80%93Hellman_key_exchange) for secure exchanging secret data over an untrusted network.

After this is completed, both Alice and Bob have the secret key and can start exchange messages.

## Design decision

### App structure

#### CLI

I have created this program as a [Cobra](https://github.com/spf13/cobra) CLI app.
I decided to use Cobra, because it is the de-facto standard for building Go CLI apps, I have used it in my work, so am familiar with it, and it has everything I needed for the task.

Most of the command are simple one-off tasks that run, do something and then exit.

The only exception is the `exchange-key-server` command that starts a web-server and listens for requests.
The command does not exit, unless explicitly stopped by the user.

I used the [cobra-cli generator](https://github.com/spf13/cobra-cli/blob/main/README.md) to generate the boilerplate for the commands.
All the commands are in the [`cmd`](./cmd) directory.

#### Internal

The other source folder I have is the [`internal`](./internal/) one.

In Go, the `internal` package is special as it does not allow an external project to import this code.
I this this, because I don't think that this code should be imported by another project, so I don't want to have it public
and to have to deal with API backwards-compatibility guarantees.

The packages inside the `internal` folder are:

- `crypto` - deals with all the crypthography.
In `crypto/exchange/listener.go` there is a bit of code that deals with HTTP stuff, which maybe is not the best place for it, but I'll let it roll for now and refactor it later.
- `messages` - domain model for the messages that are being exchanged
- `secretstore` - package for the different types of stores used to store the keys

Most of the code in `internal` is covered by unit tests.


### Crypthography

#### Message encryption

For the message encryption I used [AES](https://en.wikipedia.org/wiki/Advanced_Encryption_Standard).
This type of encryption is considered secure and used widely in the industry.

For the mode of operation I chose [Galois/Counter Mode (GCM)](https://en.wikipedia.org/wiki/Galois/Counter_Mode).
It is supported out-of-the-box by the Go stdlib, it is fast, and it provides both authenticity and confidentiality.

#### Storing the keys

By default, all CLI commands use the OS keychain for storing and reading the secret key.

This is the most secure option, because the OS keychain (depending on the OS) can support native encryption and access control.

If needed, the command also support storing(and reading) the key from a file or from a CLI argument.
This is not recommended, because files can be read by anyone and CLI arguments are visible in the shell history.
The command output warnings when these options are used.

### Building and running the code locally

To build the code locally, you need to have the Go toolchain installed on your machine.
The version of Go we are using is 1.21.

To build the code:

```console
go build -o secure-messenger main.go
```

This wil build the secure-messenger CLI binary for your OS.
You can now use it:

```console
$ ./secure-messenger generate-key --output-to-stdout
Consider using the keychain to store the key, so that it don't get lost or exposed.
fiV2Lo4ZAnrzDNyZfA5hY4+Fe+xReiZy
```

To test the code you can again use the Go tool:

```console
$ go test ./...
?       github.com/asankov/secure-messenger     [no test files]
?       github.com/asankov/secure-messenger/cmd [no test files]
?       github.com/asankov/secure-messenger/internal/crypto/exchange    [no test files]
?       github.com/asankov/secure-messenger/internal/secretstore        [no test files]
ok      github.com/asankov/secure-messenger/internal/crypto     0.750s
ok      github.com/asankov/secure-messenger/internal/messages   (cached)
```

### Future Improvements

A software can always be improved and new features can always be added.

Possible new features for the Secure Messenger are:

- support more key-exchange algorithms (for example, Public-Private Key crypthography)
- support for building and running the CLI (and server) as a container
- support for running it like a server and being able to exchange and listen for messages
- a simple UI for sending and receiving messages (after the server support is implemented)
- right now, you can generate encrypted messages via the CLI, but you cannot send them,
so a command like `secure-messenger send-message --remote-addr=xxx` would be a good UX improvement
