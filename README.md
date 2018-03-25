# crdt-json

[![Build Status](https://travis-ci.org/gpestana/crdt-json.svg?branch=master)](https://travis-ci.org/gpestana/crdt-json)

Conflict-free replicated JSON implementation in Go based on 
[Martin Kleppmann, Alastair R. Beresford
work](https://arxiv.org/abs/1608.03960).

From the paper's abstract:

> [...] an algorithm and formal semantics for a JSON data structure that 
> automatically resolves concurrent modifications such that no updates are lost,
> and such that all replicas converge towards the same state (a conflict-free 
> replicated datatype or CRDT). It supports arbitrarily nested list and map 
> types, which can be modified by insertion, deletion and assignment. The 
> algorithm performs all merging client-side and does not depend on ordering 
> guarantees from the network, making it suitable for deployment on mobile 
> devices with poor network connectivity, in peer-to-peer networks, and in 
> messaging systems with end-to-end encryption.

### Use cases for CDRTs

A good discussion and suggestions of CRDT uses can be found in the 
[research-CRDT repository maintained by IPFS](https://github.com/ipfs/research-CRDT/issues/1)

### License

MIT
