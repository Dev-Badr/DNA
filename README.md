# Skytells - Decentralized Network Architecture (DNA)
We are building a network that we believe has the potential to re-engineer privacy, security and freedom on the internet.
Skytells DNA is the world's first AI-Powered Decentralized Network Architecture.


[![N|Skytells DNA](https://cdn-images-1.medium.com/max/1200/0*hoYKuIeh7LXHYE8s)](https://www.skytells.org)


## How it works?

In order to understand how Skytells DNA works, You may have to look at this graph.

[![N|Skytells DNA](https://raw.githubusercontent.com/skytells-research/DNA/master/resources/graph/toplogy.png)](https://www.skytells.org)


So, Skytells DNA is
  - Very light and easy (one similar config on all hosts)
  - Use same config for all hosts (autedetect local params) - useful with puppet etc
  - Uses AES-128, AES-192 or AES-256 encryption (note that AES-256 is **much slower** than AES-128 on most computers) + optional HMAC-SHA256 or (super secure! ) NONE encryption (just copy without modification)
  - Communicates via UDP directly to selected host (no central server)
  - Works only on Linux (uses TUN device)
  - Support of basic routing - can be used to connect several networks
  - Multithread send and receive - scaleable for big traffc
  - Due to use so_reuseport better result in case of bigger number of hosts
  - It's still in beta stage, use it on your own risk (and please use only versions marked as "release")
