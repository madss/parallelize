Parallelize
===========

This is a tool developed primarily to parallelize ad-hoc shell script for
managing databases and similar.

How to build
------------

As with all other go programs, do a

	go install

Usage
-----

You can specify an executable to run and a list of arguments like this

	parallelize fix-users.sh 1000 1001 1002 1003 1004

This will be identical to

	fix-users.sh 1000
	fix-users.sh 1001
	fix-users.sh 1002
	fix-users.sh 1003
	fix-users.sh 1004

except it will be run in parallel. The difference between `parallelize` and
just adding an `&` to run the processes in the background is that `parallelize` will limit
the number of parallel executions to a certain number (eight as default). You can specify
this number yourself with

	parallelize -n 3 /bin/sleep 1 2 3 4 5

to fit your needs.

If you need to give a script more than one argument or if you just happen to
already have  the arguments in a comma separated file, you can use that as an
input. Each record of the comma separated file will act as arguments for the
script. For example, if `users.csv` contains

	1000,tom,Tom
	1001,rob,Rob
	1002,diana,Diana
	1003,jeff,Jeff
	1004,ann,Ann

and you run

	parallelize -csv users.csv fix-users.sh

the `fix-users.sh` script will be called with

	fix-users.sh 1000 tom Tom

and so on
