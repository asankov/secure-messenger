# Secure Messaging System

Secure Messaging System is a CLI tool that allows you to send and receive encrypted messages.

## Usage

The `secure-messenger` CLI exposes a few command that allow you to send and receive encrypted messages.

### generate-key

In order to send and receive encrypted messages you need a secret key.

You can generate one via the `generate-key` command:

```shell
$ secure-messenger generate-key
zTYpiKkH9jsAR1EMPAgmdoeMJCDAcPiD
```

This is a 32-byte long key that you can use to encrypt and decrypt messages.

You can also generate 16-byte and 24-byte long keys.

```shell
$ secure-messenger generate-key --key-size=16
iPUML2/JXJ96NrW+
$ secure-messenger generate-key --key-size=24
h/LVOSopnqb/e2maMxAjIQ==
```

You need to share this key with the person you want to exchange messages with.

For more info on how to share the key securely, go to [Key Exchange](#key-exchange).

### encrypt

Once you have a key you can use it to generate encrypted messages.

```shell
$ secure-messenger encrypt --sender-id=abc123 --receiver-id=987 --secret-key=zTYpiKkH9jsAR1EMPAgmdoeMJCDAcPiD --payload="Hello from Planet Earth"
6Vq4bYb+9yK+wB8P22evR+FTYUZG6zksB7tAzqgnpv66Z2y9f1fJ6UcXZtKPc3Sm9SkwiBg/fXTqLPokvw178WxAqqa3JtRdvUGRr4Ksp/ABXF06IyX48EaIhDAivM4sACYYiditNkLoGyz0b3685yFgMxLc1K7f0Ce13dDuYQ==
```

In order to generate a message you need to input 4 arguments:

- `sender-id` - this is the ID of the sender (your ID)
- `receiver-id` - this is the ID of the receiver
- `payload` - the payload you want to send
- `secret-key` - the secret key

The output of the command is the encrypted message, which you can safely share with the other side, without worrying that someone
might intercept it and read the contents.

Alternatively, you can save the key into a file and tell the `secure-messenger` CLI to read it from there.

```shell
$ secure-messenger generate-key > key.file
$ secure-messenger encrypt --sender-id=abc123 --receiver-id=987 --secret-key-file=key.file --payload="Hello from Planet Earth"
UIqwQxl9ntJc9SIaJidJOrx6QrgCrJbr7Jy8rkS3BjVIE6TofnO4ljW6mIc4Eo8CoM/w9rSYNaRvTbtArEyMxebrVJrH0xcKxZhJnAKC3A83EXL+rfh9+wNki6DH/aKqt1XnzoajK6lH1Bep2O74oR8aRNwgIRYN4R9GdtPYvQ==
```

WARNING: Neither of these ways of storing the keys are secure.
Passing the key as an argument to the CLI makes it visible into the shell history.
Saving the key into a simple files means that everyone that got access to the machine can read it.

TODO: Implement a more secure way to save the key.

### decrypt

Once you have receive a message and you have the key you can use the `secure-messenger` CLI to decrypt it and read it.

```shell
$ secure-messenger decrypt --secret-key-file=key.file UIqwQxl9ntJc9SIaJidJOrx6QrgCrJbr7Jy8rkS3BjVIE6TofnO4ljW6mIc4Eo8CoM/w9rSYNaRvTbtArEyMxebrVJrH0xcKxZhJnAKC3A83EXL+rfh9+wNki6DH/aKqt1XnzoajK6lH1Bep2O74oR8aRNwgIRYN4R9GdtPYvQ==

{"senderId":"abc123","receiverId":"987","payload":"Hello from Planet Earth","timestamp":1700317704}
```

TODO: passing long encrypted message as a CLI arguments is not a good UX.
Support a better way like reading the messages from file or from stdin.

Voila. You have exchange an encrypted message.

## Key Exchange

In order to securely exchange the secret key with the other party, you can again use the CLI.

The two command to do that are `exchange-key` and `exchange-key-server`.

The other side (let's call them Bob, and we are Alice) need to run the `exchange-key-server` command.
This will start a web server that will listen for a request that will start the key exchange process.

```shell
(bob)$ secure-messenger exchange-key-server
Starting exchange key server on [:8080]
```

Then we (Alice) can use the `exchange-key` command to send a request to Bob's server that will initiate the key-exchange process.

```shell
(alice)$ secure-messenger exchange-key --remote-addr=http://localhost:8080 --secret-key-file=key.file
successfully exchanged secret key
```

Meanwhile Bob sees this:

```shell
(bob)$ secure-messenger exchange-key-server
Starting exchange key server on [:8080]
time=2023-11-18T18:28:09.647+02:00 level=INFO msg="Retrieved secret key"
time=2023-11-18T18:28:09.648+02:00 level=INFO msg="Secret key stored" location=keychain
```

This means that the key-exchange was successful and now Bob has the secret key stored in its keychain.

These command use the Diffie-Helmman algorithm for secure exchanging secret data over an untrusted network.

After this is completed, both Alice and Bob have the secret key and can start exchange messages.
