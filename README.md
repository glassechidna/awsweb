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

Or, via homebrew:
```
brew tap glassechidna/taps
brew install awsweb
```

## Setup

Your `~/.aws/config` should *already* look something like this:

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

You should also add the MFA serial number to your _config_ file in the "source" account, e.g.

```
[mycompany]
region = ap-southeast-2
mfa_serial = arn:aws:iam::0987654321:mfa/aidan.steele@example.com
```

## Usage

You can then do `awsweb browser mycompany-prod` and a browser window will pop up.  
Or `eval "$(awsweb env mycompany-prod)"` to set `AWS_*` environment variables for profile *mycompany-prod*.  
Also `awsweb set mycompany-prod` will update the default profile (in `~/.aws/config` and `~/.aws/credentials`) with the temporary credentials.

Temporary credentials last for 1 hour and are cached in `~/.aws/awswebcache`

### Handy bash aliases

```
alias ab='awsweb browser'

awsenv() {
  eval $(awsweb env "$@")
}

alias ae='awsenv'
```
eg: `ae mycompany-prod` to set `AWS-*` environment variables for profile *mycompany-prod*.

### Powershell

If you're using Powershell, you can do:

```
$Cmd = (awsweb env --shell powershell mycompany-prod) | Out-String
Invoke-Expression $Cmd
```
