iSMC-lib
========

[![GitHub license](https://img.shields.io/github/license/dkorunic/iSMC.svg)](https://github.com/xxxserxxx/iSMC/blob/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/dkorunic/iSMC)](https://goreportcard.com/report/github.com/xxxserxxx/iSMC)

## About

`iSMC` is an Apple System Management Controller (SMC) CLI tool that attempts to query SMC for a number of well known keys and determine their type and value, classifying them into temperature, power, current, voltage, fan and battery readouts.

Typically various desktop and server Apple hardware should work and most definitely all Intel-based Mac computers.

I am not maintaining this for any use other than within [gotop](https://github.com/xxxserxxx/gotop), and so provide no how-to or help.  I recommend using the upstream project if you want something user-friendly.

A word about dependencies
-------------------------

This fork removes all of the command-line code of the upstream project and strips it down to a more lightweight library.  Eliminating Cobra alone saves 9,522 LOC and 404Kb of source code.  However, Cobra *also* pulls in a bunch of other libraries, some of which are utterly unnecessary and which are nothing but bloat in a library:

- github.com/spf13/pflag
- go.mongodb.org/mongo-driver (!!)

In addition, the CLI imports logrus, which sucks in even more dependencies:

- github.com/sirupsen/logrus
- github.com/asaskevich/govalidator
- github.com/go-openapi/errors
- github.com/go-openapi/strfmt
- github.com/konsorten/go-windows-terminal-sequences
- github.com/mattn/go-runewidth

This illustrates the wise Go philosophy:

> A little copying is better than a little dependency.
