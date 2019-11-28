# spksrc-npc

Setup Development Environment
-----------------------------
Docker

* Fork and clone spksrc: ``git clone https://github.com/cnlh/spksrc.git ~/spksrc``
* Install Docker on your host OS: `Docker installation`_. A wget-based alternative for linux: `Install Docker with wget`_.
* Download the spksrc docker container: ``docker pull synocommunity/spksrc``
* Run the container with ``docker run -it -v ~/spksrc:/spksrc synocommunity/spksrc /bin/bash``


Build Npc
-------
* ```cd spksrc```

* ```make setup```

* ```cd spksrc/spk/npc```

* ```make```

