# `awsweb`

[![Build Status](https://travis-ci.org/glassechidna/awsweb.svg?branch=master)](https://travis-ci.org/glassechidna/awsweb)

`awsweb` is a tiny CLI tool that makes it easier to hop between AWS accounts and profiles without going through the 
regular username + password, switch role, switch region dance. It uses the credentials in `~/.aws` and optionally 
user-provided or stored MFA credentials.
