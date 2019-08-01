---
title: CLI Reference
path: /docs/using-kabanero/cli-commands
section: Using kabanero
---
# kabanero CLI
## kabanero-mgmt

A command line interface that can be used to manage the environment

### Synopsis

A command line interface that can be used to manage the collections that 
the environment prosents, as well as on-board the people and clusters that will be 
used in the environment to build applications.

Complete documentation is available at https://kabanero.io

### Options

```
      --config string   config file (default is $HOME/.kabaneromgmt.yaml)
      --dryrun          Turns on dry run mode
  -h, --help            help for kabanero-mgmt
  -v, --verbose         Turns on debug output and logging to a file in $HOME/.kabaneromgmt/logs
```

### SEE ALSO

* [kabanero-mgmt activate](#kabanero-mgmt-activate)	 - Activate a collection for use by the development team
* [kabanero-mgmt add](#kabanero-mgmt-add)	 - Add a collection to the indeex
* [kabanero-mgmt deactivate](#kabanero-mgmt-deactivate)	 - Prevent this collection from being shown to the development team, while not deleting it.
* [kabanero-mgmt list](#kabanero-mgmt-list)	 - List all the collections in the apphub, and optionally their status
* [kabanero-mgmt login](#kabanero-mgmt-login)	 - Will authentic you to the Kabanero instance
* [kabanero-mgmt logout](#kabanero-mgmt-logout)	 - Disconnect from Kabanero instance
* [kabanero-mgmt onboard](#kabanero-mgmt-onboard)	 - Command to onbboard a developer to the Kabanero infrastructure
* [kabanero-mgmt remove](#kabanero-mgmt-remove)	 - Remove the collection from the apphub
* [kabanero-mgmt version](#kabanero-mgmt-version)	 - Show Kabanero CLI version

###### Auto generated by spf13/cobra on 1-Aug-2019
## kabanero-mgmt login

Will authentic you to the Kabanero instance

### Synopsis


	The userid and password passed will be used
	to authenticate the user with kabanero instance.
	
	By authenticating with the Kabanero instance, 
	you will be able to manage the instance of kabanero.

```
kabanero-mgmt login userid password [kabanero-url] [flags]
```

### Examples

```

		kabanero jane janespw https://kabanero1.io
		
```

### Options

```
  -h, --help   help for login
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.kabaneromgmt.yaml)
      --dryrun          Turns on dry run mode
  -v, --verbose         Turns on debug output and logging to a file in $HOME/.kabaneromgmt/logs
```

### SEE ALSO

* [kabanero-mgmt](#kabanero-mgmt)	 - A command line interface that can be used to manage the environment

###### Auto generated by spf13/cobra on 1-Aug-2019
## kabanero-mgmt logout

Disconnect from Kabanero instance

### Synopsis


Disconnect from the instance of Kabanero that you 
have been interacting with.

```
kabanero-mgmt logout [flags]
```

### Options

```
  -h, --help   help for logout
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.kabaneromgmt.yaml)
      --dryrun          Turns on dry run mode
  -v, --verbose         Turns on debug output and logging to a file in $HOME/.kabaneromgmt/logs
```

### SEE ALSO

* [kabanero-mgmt](#kabanero-mgmt)	 - A command line interface that can be used to manage the environment

###### Auto generated by spf13/cobra on 1-Aug-2019
## kabanero-mgmt activate

Activate a collection for use by the development team

### Synopsis


A collection can be available to a development team
to use for building applications or not. Activate
will cause the collection to be shown to 
the develoopment team when they list the types of
application they can build.

```
kabanero-mgmt activate collection-name [flags]
```

### Options

```
  -h, --help   help for activate
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.kabaneromgmt.yaml)
      --dryrun          Turns on dry run mode
  -v, --verbose         Turns on debug output and logging to a file in $HOME/.kabaneromgmt/logs
```

### SEE ALSO

* [kabanero-mgmt](#kabanero-mgmt)	 - A command line interface that can be used to manage the environment

###### Auto generated by spf13/cobra on 1-Aug-2019
## kabanero-mgmt deactivate

Prevent this collection from being shown to the development team, while not deleting it.

### Synopsis


A collection can be available to a development team
to use for building applications or not. deactivate
will cause the collection to not be shown to 
the develoopment team when they list the types of
application they can build.

This would be done in the case where you have cloned the collection
and made changes for your business.  This keeps the base collection
in the apphub, and it will continue to be updated, and the 
updates will be perkolated up to your cloned collection.

```
kabanero-mgmt deactivate collection-name [flags]
```

### Options

```
  -h, --help   help for deactivate
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.kabaneromgmt.yaml)
      --dryrun          Turns on dry run mode
  -v, --verbose         Turns on debug output and logging to a file in $HOME/.kabaneromgmt/logs
```

### SEE ALSO

* [kabanero-mgmt](#kabanero-mgmt)	 - A command line interface that can be used to manage the environment

###### Auto generated by spf13/cobra on 1-Aug-2019
## kabanero-mgmt list

List all the collections in the apphub, and optionally their status

### Synopsis

A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.

```
kabanero-mgmt list [status] [flags]
```

### Options

```
  -h, --help   help for list
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.kabaneromgmt.yaml)
      --dryrun          Turns on dry run mode
  -v, --verbose         Turns on debug output and logging to a file in $HOME/.kabaneromgmt/logs
```

### SEE ALSO

* [kabanero-mgmt](#kabanero-mgmt)	 - A command line interface that can be used to manage the environment

###### Auto generated by spf13/cobra on 1-Aug-2019
## kabanero-mgmt add

Add a collection to the indeex

### Synopsis

Add the collection to the indeex

```
kabanero-mgmt add collectioon-name [flags]
```

### Options

```
  -h, --help   help for add
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.kabaneromgmt.yaml)
      --dryrun          Turns on dry run mode
  -v, --verbose         Turns on debug output and logging to a file in $HOME/.kabaneromgmt/logs
```

### SEE ALSO

* [kabanero-mgmt](#kabanero-mgmt)	 - A command line interface that can be used to manage the environment

###### Auto generated by spf13/cobra on 1-Aug-2019
## kabanero-mgmt remove

Remove the collection from the apphub

### Synopsis


Remove the collection from the apphub.

```
kabanero-mgmt remove collection-name [flags]
```

### Options

```
  -h, --help   help for remove
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.kabaneromgmt.yaml)
      --dryrun          Turns on dry run mode
  -v, --verbose         Turns on debug output and logging to a file in $HOME/.kabaneromgmt/logs
```

### SEE ALSO

* [kabanero-mgmt](#kabanero-mgmt)	 - A command line interface that can be used to manage the environment

###### Auto generated by spf13/cobra on 1-Aug-2019