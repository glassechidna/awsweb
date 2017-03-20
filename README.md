# `awsweb`

[![Build Status](https://travis-ci.org/glassechidna/awsweb.svg?branch=master)](https://travis-ci.org/glassechidna/awsweb)

`awsweb` is a tiny CLI tool that makes it easier to hop between AWS accounts and profiles without going through the 
regular username + password, switch role, switch region dance. It uses the credentials in `~/.aws` and optionally 
user-provided or stored MFA credentials.

## Usage

Your `~/.aws/config` should look something like this:

```
[default]
region = ap-southeast-2

[mycompany]
region = ap-southeast-2

[profile mycompany-prod]
role_arn = arn:aws:iam::1234567890:role/Developer
source_profile = mycompany
region = us-east-1
mfa_serial = arn:aws:iam::0987654321:mfa/aidan.steele@example.com
```

Your `~/.aws/credentials` will look like this:

```
[default]
aws_access_key_id = AKIA...
aws_secret_access_key = qxxg....

[mycompany]
aws_access_key_id = AKIA...
aws_secret_access_key = qGrg....
```

You can then do `awsweb --mfa-secret SOMESECRET mycompany-prod` and a browser window will pop up. Alternatively you
can store your MFA secret in `~/.awsweb.yml` in the following format:

```yaml
mfa-secret: SOMESECRET
```
