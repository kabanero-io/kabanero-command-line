# Installation
## Mac OSX
You can download the binary `kabanero` command from here:
https://github.com/kabanero-io/kabanero-command-line/releases

# Building from Source

## Travis build
The project is instrumented with Travis CI and with an appropriate `Makefile`. Most of the build logic is triggered from within the `Makefile`.

Upon commit, only the `test` and `lint` actions are executed by Travis.

In order for Travis to go all the way to `package` and `deploy`, you need to create a *new* release (one that is tagged with a never seen before tag). When you create a new release, a Travis build with automatically run, and the resulting artifacts will be posted on the `Releases` page. 

## Manual build
You can also test the build process manually.


Prerequisites:

* docker is installed and running
* wget is installed

After setting the `GOPATH` env var correctly, just run `make <action...>` from the command line, within the same directory where `Makefile` resides. For example `make package clean` will run the `package` and then the `clean` actions.


Some of the scripts have conditional paths, because certain Linux commands behave differently on OS/X and elsewhere (fun).

### What gets produced by the build?
Quite a bit of stuff. 

Here's a description of the various artifacts as you would see them in a release page:

* The actual RPM package for RHEL/Centos (to be `yum`med or `rpm -i`)

* The binaries tarred up in the way homebrew loves them

* The plain binaries tarred up as they come out of the build process

** for OS/X

* The binary to run

** for Linux

** for Windows

* The homebrew Formula (which we should push to some git repo, once we go "public")

* The Debian package for Ubuntu-like Linux (to be `apt-get install`ed)

* Some other stuff that's always there

### Running the CLI
`./kabanero`
