# Specs

This document describes the specifications and internals of the `rdoc` data
structure. The `rdoc` is an implementation of an operation-based JSON CRDT as
presented in [1][Martin Kleppmann, Alastair R. Beresford work](https://arxiv.org/abs/1608.03960).

## Interfaces

The `rdoc` data structure exposes 3 interfaces with different
responsibilities and scopes:

- **User interface**: the high level interface for the user to interact with
  `rdoc` documents. It specifies how the user initiates documents and perform
modifications in it such as modifying, creating and deleting nodes in the JSON
document. The user interface also provides reading interface for the user to get
the current state of the document.

- **Document interface**: the private interface which is called by proxy and not
  explicitly by the user. This interface manages the document state and applies
the operations requested by the user interface. The document interface works at
an operations level. Its responsibility is to apply remote and local operations
against the current local state.

The `rdoc` exposes the user interface publicly. The document interface is 
private and responsible to manage the JSON document internally so that the CRDT 
properties as described by [1] are achieved.

### User interface

The user interface is exposed to the library user to mutate and read the
document state.

**API**:

```go
func rdoc.Init() *Doc
```` 

Initiates a new document. 

```go
func (*doc) Serialize() ([]byte, error)
```

Returns the JSON encoding of the current state of the document. Internally, the
`Serialize` struct method will call `encoding.Marshal(v interface)`

// TODO: define the user interface for mutating the document

**Examples**:

```go

doc := rdoc.Init()

// perform mutations on document - define interface

bdoc := doc.Serialize()
```

### Document interface

The document interface is the API which is used to internally modify and read
the document structure.

**API**:

```go
func (*Doc) ApplyOperation(operation op) (*Doc, error)
```

Traverses the document until the node pointed by the operation `cursor`
and applies the operation `mutation`. This function hides all the complexity of
mutating the CRDT document coming from the metadata management. It returns a
pointer to the document state after the mutation.

```go
func (*Doc) ApplyRemoteOperation(operation op) (*Doc, error)
````

Verifies if the remote operation has been applied before and if all dependencies
are fulfilled. If both conditions hold true, apply the operation by calling
`ApplyOperation(remoteOp)`. Otherwise, buffer the operation until all
dependencies are fulfilled or discard the operation if it has been applied
previously. 

**Examples:**

```go
doc := rdoc.Init() // initializes a new rdoc

deps := []string{...}
opID := "..."
mutation := operation.NewMutation(...)
cursor := operation.NewCursor(...)

op := operation.New(id, deps, cursor, mutation)

doc2 := doc.ApplyOperation(op)

remoteOp := operation.Operation{}
receiveRemoteOperation(&remoteOp) // receives remote operation

doc3 := doc=4.ApplyRemoteOperation(remoteOp)
```

## Document management specifications

This section defines how an operation, document and auxiliary types are 
represented and what are their methods. An operation identifies uniquely a 
mutation in the local document.

A document may contain 3 different types: a list, a map and a registers. Each of
the types may be empty or not. Each of the list and map values are pointers to 
document nodes. A register is a multi-value register represented by a hash map 
in which keys are the values of the registers and values the operations that
assigned the register value.

**Document data structure**:

```go
type Doc struct {
	id string
	deps []string
	head *Node
}
```

**Node**:

```go
type Node struct {
  key  interface{}
  deps *arraylist.List
  hmap *hashmap.Map
  list *arraylist.List
  reg  *hashmap.Map
}
```

A `Node`` may contain a no-empty map, list and/or register. In
order to ensure the properties required by the JSON CRDT, a node in the document
tree may be at the same time a map, list and/or register (see paper for more
information). Thus, a `Doc` keeps a pointer for storing information for each
type.

Each node keeps a list with IDs of dependencies. A dependency is an operation ID
that the node depends on.

### Operation type

An operation uniquely identifies a mutation in the local document and can be
shared across peers in the network. Each operation contains a cursor, a 
`mutation`, `id` and `dependencies`:

**Data structure**:

```go
type Operation struct {
	id string
	deps []string
	cursor []CursorElement
	mutation Mutation
}	
```

**API**:

A `cursor` is a set of cursor elements, each specifying the type of the node and
the respective key of the node to traverse to. A `cursor` is an ordered list of
`CursorElements`.

```go
cursor := []CursorElement{
	CursorElement{Key: "root", Type: MapT},
  CursorElement{Key: "other", Type: MapT},
  CursorElement{Key: "other-new", Type: MapT},
  CursorElement{Key: 0, Type: ListT},
}

// cursor points at root.other.other-new[0]
```

A `mutation` describes what is the modification to apply to the node selected by
the `cursor`.

```go
func NewMutation(type string, k interface, v interface) (Mutation, error)
````

Returns a new `mutation`. The type input can be one of `INSERT`, `DELETE`.
`ASSIGN`. The `k` and `v` are the key and value, respectively, to apply.

```go
func NewOperation(cursor Cursor, mut Mutation, deps []string, doc Doc) (error, Operation)
```

Creates a new operation based on a previously created `cursor`, `mutation`. The
`operation` id is calculated based on the `Doc`'s Id and current Lamport
clock state of the document. Thus the document itself must be passed as an
argument. `deps` is a list of strings which are Lamport Ids that describe which
operations the `operation` to construct depends upon.

## Algorithm for applying operation on a document 

## Document immutability

The document structure and interface does not have any public methods which
mutate the data structure. Every time the `(*Doc) ApplyOperation()` method is
called, a new document is generated (deep copy of document where to apply the
operation) and the operation is applied in the copied document structure. The
new document is returned as part of the `(*Doc) ApplyOperation()` method.

This design ensures that a document is immutable at the user level.

// TODO: how to efficiently perform deepcopy on document structure
