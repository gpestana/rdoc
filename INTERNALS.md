# crdt-json implementation details

This document describes documentations details of the json-crdt library. You
don't need to learn those when using it but its a good help if you'd like to
contribute or learning about the internals of the library.

## JSON document

A JSON CRDT document is represented by a collection of `Nodes` linked with each
other. 

A document is a data structure which maintains a set with all the operations
ID applied to the document and a set with pointers to the first nodes in the
document. A document is identified by an unique ID. 

A document exposes an high level API for the user to mutate the local state,
merge different documents and print the document in an adapted JSON format.

The document is implemented in `/json.go`.

Each node maintains a list of pointers to the next nodes, a set of
dependencies (operation IDs which relied on the node's existence) and a `type`.
A node may be of type `mapT`, `listT` or `registerT`. Each type has a value,
which means different things depending on the type. In the case of `registerT -
which is a multi-value register - a value is a set of either `int` or `string`
The set should not be empty. The value of a `mapT` is the key of the map and it
should never bt empty. The value of a `listT` is always nil, since the content
of a list are nodes which are maintained in the `links` set.

The `node` data type is implemented in `/node.go` and the different data
types are implemented under `/types/`.

## Operations

The crdt-json is an operation-based JSON crdt implementation. An operation is
represented by:

```
op {	
	ID       : string (Lamport timestamp ID)
	deps     : set of casual dependencies, represented by strings of Lamport
timestamp IDs
	cursor   : path to the node where the mutation will be applied, represented as
an array of string
	mutation : mutation to apply to node, represented as a data structure	 
}

This data structure is defined in `/operation/operation.go`.

Internally, an operation is applied to the current document state. This happens
using the primitive `doc.apply(operation)`

Applying an operation to a document consists of the following steps:

	1) Make sure the operations has not been applied already

This can be checked by inspecting the set in the document which contains all the
operation IDs applied so far.

	2) Make sure that all the operations dependencies are met

This can be checked by comparing the operations applied to the document and the
operations dependencies. If not all the operation dependencies are respected,
the operation is buffered until that's the case.

	3) Select the node to apply mutation

Traverse the document from the root until the node in which the mutation should
be applied, as described on the operation's cursor. If the cursor contains node 
values that do not exist in the current document, new nodes must be created.

	4) Apply the mutation

