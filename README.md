# `awsweb`

[![Build Status](https://travis-ci.org/glassechidna/awsweb.svg?branch=master)](https://travis-ci.org/glassechidna/awsweb)

`awsweb` is a tiny CLI tool that makes it easier to hop between AWS accounts and
profiles without going through the regular username + password, switch role,
switch region dance. It uses the credentials in `~/.aws` and optionally user-provided
or stored MFA credentials.

## Installation

You can download the latest build of `awsweb` from the [project's Github Releases][github-releases]
page. Download the binary for your platform and place it somewhere in your `PATH`.

[github-releases]: https://github.com/glassechidna/awsweb/releases

## Setup

Your `~/.aws/config` should look something like this:

```
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
[mycompany]
aws_access_key_id = AKIA...
aws_secret_access_key = qGrg....
```

## Usage

You can then do `awsweb --mfa-secret SOMESECRET browser mycompany-prod` and a browser
window will pop up. Or `eval "$(awsweb --mfa-secret SOMESECRET env mycompany-prod)"`.
Alternatively you can store your MFA secret in `~/.awsweb.yml` in the following format:

```yaml
mfa-secret: SOMESECRET
```

### Powershell

If you're using Powershell, you can do:

```
$Cmd = (awsweb env --shell powershell mycompany-prod) | Out-String
Invoke-Expression $Cmd
```
