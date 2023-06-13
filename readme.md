# TrueCrypt Alternative  

## Setup  

- `setup.sh.command` to generate settings

## Usage  

- `encrypt.sh.command` to encrypt the decrypted folder
- `decrypt.sh.command` to decrypt the encrypted folder
- `clean.sh.command` to delete the decrypted folder

### Dependencies

Make sure these work `cp`, `diff`, `zip`, `unzip`, `stty`. The tool is tested on MacOS but it could work in Linux or Unix.

### Installation

```bash
go install github.com/cshtdd/truecrypt@latest
```

---

## Development  

Build & test

```bash
./build.sh
```
