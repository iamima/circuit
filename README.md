The Circuit
===========

For a visual introduction to The Circuit, dive into the 
[GopherCon 2014 Slides](https://docs.google.com/presentation/d/1ooedstHs8_ow-eHY7z8MCV_1m65gSaiB6Q1ruz3j7Hk/edit#slide=id.g26f183bd0_00).

The circuit is a tool for executing and synchronizing UNIX processes across entire clusters
by means of a file-system interface. 

The circuit represents any cluster of UNIX machines in Hoare's _Communicating
Sequential Processes_ model, hereafter CSP for short.

The circuit's CSP model provides three simple _elements_ that suffice for
unconstrained general-purpose programmability and orchestration of cloud
applications across clusters of any size. These elements are:

* _Channel:_ An ordered, message-oriented communication primitive
* _Selection:_ A synchronization mechanism for waiting on multiple event sources
* _Process:_ A primitive for executing, monitoring and synchronizing with UNIX processes

Build
-----

The Circuit comprises one small binary. It can be built for Linux and Darwin.

Given that the [Go Language](http://golang.org) compiler is installed,
you can build and install the circuit binary with one line:

	go get github.com/gocircuit/circuit/cmd/circuit

Run
---

Prepare a local directory that can be FUSE-mounted by your user. 
For instance, `/circuit` is a good choice.

To run the circuit agent, pick a public IP address and port for it to
listen on, and start it like so

	circuit -a 10.20.30.7:11022 -m /circuit

Among a few other things, the circuit agent will print its own circuit URL.
It should look like this:

	…
	circuit://10.20.30.7:11022/78517/R56e7a2a0d47a7b5d
	…

Copy it. We will need it to tell the next circuit agent to join this one.

Log onto another machine and similarly start a circuit agent there, as well.
This time, use the `-j` option to tell the new agent to join the
circuit of the first one:

	circuit -a 10.20.30.5:11088 -m /circuit -j circuit://10.20.30.7:11022/78517/R56e7a2a0d47a7b5d

You now have two mutually-aware circuit agents, running on two different hosts in your cluster.
You can join any number of additional hosts to the circuit environment in a similar fashion,
even billions:

The circuit uses a modern [expander graph](http://en.wikipedia.org/wiki/Expander_graph)-based
algorithm for presence awareness and ordered communication, which is genuinely distributed;
It uses communication and connectivity sparingly, hardly leaving a footprint when idle.

Explore
-------

On any host with a running circuit agent, go to the local circuit mount directory

	cd /circuit
	ls

Each of its subdirectories corresponds to a live circuit agent. Navigate into
any one of them and explore the file system. Each directory is equipped with a
`help` file to guide you.

Learn more
----------

To stay up to date with new developments, documentation and articles, follow
The Circuit Project on Twitter [@gocircuit](https://twitter.com/gocircuit) or
me [@maymounkov](https://twitter.com/maymounkov).
