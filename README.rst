spksrc-npc

Setup Development Environment
-----------------------------
Docker

* Fork and clone spksrc: ``git clone https://github.com/cnlh/spksrc.git ~/spksrc``
* Install Docker on your host OS: `Docker installation`_. A wget-based alternative for linux: `Install Docker with wget`_.
* Download the spksrc docker container: ``docker pull ffdfgdfg/spksrc``
* Run the container with ``docker run -it -v ~/spksrc:/spksrc ffdfgdfg/spksrc /bin/bash``


Build Npc
-------
* ``cd spksrc``

* ``make setup``

* ``cd spksrc/spk/npc``

* ``make``

