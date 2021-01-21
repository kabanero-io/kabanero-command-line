# Build Guide

## Installation
### OSX, Windows, Linux
You can download the binary `kabanero` command from here:
https://github.com/kabanero-io/kabanero-command-line/releases


## Generating the README.md
The project uses Cobra's built in mechanism to generate the CLI README.md.  Run the following command to generate the README:
    ` ./build/kabanero docs --docFile ./README.md`


## Travis build
The project is instrumented with Travis CI and with an appropriate `Makefile`. Most of the build logic is triggered from within the `Makefile`.

Upon commit, only the `test` and `lint` actions are executed by Travis.

## Creating a new release

In order for Travis to go all the way to `package` and `deploy`, you need to create a *new* release on Github(one that is tagged with a never seen before tag). When you create a new release, a Travis build will automatically run, and the resulting artifacts will be posted on the `Releases` page. With each new release or release candidate don't forget to indicate on the "This is a pre-release" tick box. Update notable changes in the comments section from the last release (check the commits from now to last release).

## Updating the brew install for new releases

The homebrew repository for Kabanero CLI: https://github.com/kabanero-io/homebrew-kabanero

This brew install method is similar to what the Appsody team did for the Appsody CLI. You just have to update the `url` and `sha` fields. 

`url`: This will be the url for the homebrew tar that gets generated from the travis build for a new release. 

`sha`: Download the homebrew tar locally and calculate the sha256 for it. You can use the command:  `shasum -a 256 <homebrew_tar>` 

## Building the CLI manually
You can also test the build process manually.


Prerequisites:

After setting the `GOPATH` env var correctly, just run `make <action...>` from the command line, within the same directory where `Makefile` resides. For example `make package clean` will run the `package` and then the `clean` actions.

The version number in manual builds default to 0.1.0. To change the value yourself change the `VERSION` value in the `Makefile`. 


Some of the scripts have conditional paths, because certain Linux commands behave differently on OS/X and elsewhere (fun).

### What gets produced by the build?

Here's a description of the various artifacts as you would see them in a release page:

** for OS/X

* The binary to run
* The tarball for the brew install

** for Linux

* The binary to run

** for Windows

* The windows .exe to run 

* Some other stuff that's always there

### Running the CLI
`./kabanero`


