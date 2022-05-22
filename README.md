### gootp

gootp is a terminal-based 2FA (Two-Factor Authentication) app.

[![asciicast](https://asciinema.org/a/s9eF7EbqgnCkoLVwbv0TE0a8Y.svg)](https://asciinema.org/a/s9eF7EbqgnCkoLVwbv0TE0a8Y)

## Features
- Supported algorithms: TOTP
- Compatible with [andOTP](https://github.com/andOTP/andOTP) file format
- Allows to encrypt/decrypt andOTP files on your PC

## Installation

```shell
go install github.com/Kichiyaki/gootp@latest
```
## Examples

```shell
gootp # show OTP list
gootp -h # help for gootp
gootp -p /path/to/andotp/file/.otp_accounts.json # override default path
gootp --password xxx # specify encryption password via flag
gootp -p /path/to/andotp/file/.otp_accounts.json encrypt -o /output/.otp_accounts.json.aes # encrypt file
gootp -p /path/to/andotp/file/.otp_accounts.json decrypt -o /output/.otp_accounts.json.aes # decrypt file
```
