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

## API

**Create new JSON document**

```go
import (
  jcrdt "github/gpestana/json-crdt"	
)

doc := jcrdt.Init()

fmt.Println(doc) 
// {}
```

**Document editing**

Changing the JSON document locally is triggered by calling `Change` method on
the document object with the mutation data.

```go
c := []byte(`{"todo":{"done":["read book"],"buffered":[]}}`)
doc2 := doc.Change(c)

fmt.Println(doc2)
// { todo: { done: ['read book'], buffered: []} }
```

**Document merging**

Different documents can be merged by calling the method `Merge` on one the
documents.

```go
fmt.Println(doc)
// { todo: { done: ['read book'], buffered: []} }

fmt.Println(doc1)
// { todo: { done: [], notes: ['this is a note']} }

doc2 := doc.Merge(doc, doc1)

fmt.Println(doc2)
// { todo: { done: ['read book'], buffered: [], notes: ['this is a note']} }
```

### Use cases for JSON-CDRT

A good discussion and suggestions of CRDT uses can be found in the 
[research-CRDT repository maintained by IPFS](https://github.com/ipfs/research-CRDT/issues/1)

### License

MIT
