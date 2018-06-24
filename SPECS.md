# Specs

This document describes the specifications and internals of the `rdoc` data
structure. The `rdoc` is an implementation of an operation-based JSON CRDT as
presented in [Martin Kleppmann, Alastair R. Beresford work](https://arxiv.org/abs/1608.03960).

## Interfaces

The `rdoc` data structure exposes 3 interfaces with different
responsibilities and scopes:

- **User interface**: the high level interface for the user to interact with
  `rdoc` documents. It specifies how the user initiates documents and perform
modifications in it such as modifying, creating and deleting nodes in the JSON
document. The user interface also provides reading interface for the user to get
the current state of the document.

- **Network interface**: the high level interface for the operations applied
  remotely to be applied locally. Similarly to the user interface, it is exposed
for the library user to call whenever remote operations are received.

- **Document interface**: the private interface which is called by proxy and not
  explicitly by the used. This interface manages the document state and applied
the operations requested by the used and network interfaces.

The `rdoc` exposes the user interface and network interface publicly. The
document interface is private and responsible to manage the JSON document
internally do that the CRDT properties are achieved.

### User interface

The user interface is exposed to the library user to mutate and read the
document state.

**API**:

`func Init() ()`: 

Initiates a new document. 

`func (*doc) Serialize() ([]byte, error)`

Returns the JSON encoding of the current state of the document. Internally, the
`Serialize` struct method will call `encoding.Marshal(v interface)`

// TODO: define the user interface for mutating the document


**Examples**:

```go

doc := rdoc.Init()

// perform mutations on document - define interface

bdoc := doc.Serialize()
```

### Network interface

The network interface is exposed to the library user to mutate the document
state with operations received remotely. The exposed API is the following:

**API**:

`func (*doc) ApplyRemote(operation op) (*doc, error)`

Receives a remote operation and applies it to the current state. Returns a new 
copy of the document after the operation was applied.

**Examples**:

```go
doc := rdoc.Init() // initializes a new rdoc

op := make(operation.Operation)
receiveRemoteOperation(&op) // receives remote operation

doc2 := doc.ApplyRemote(op)
```

### Document interface

The document interface is the API which is used to internally modify and read
the document structure.

**API**:

`func Init() document`

Initiates a new document 

`func (*doc) Mutate(operation) *doc`

Traverses the document until the node pointed by the operation `cursor`
and applies the operation `mutation`. This function hides all the complexity of
mutating the CRDT document coming from the metadata management. It returns a
pointer to the document state after the mutation.

`func (*doc) Get(cursor) *doc`

Traverses the document until the node pointed by the input `cursor`.

**Examples:**

```go
doc := rdoc.Init() // initializes a new rdoc

// construct operation for mutation
mutation := operation.NewMutation(...) // TODO: define how mutation is represented and APIs
cursor := operation.NewCursor(...) // TODO: define how operation is represented and APIs

op := operation.New(cursor, mutation)

doc2 := doc.Mutate(op)
doc3 := doc2.Get(cursor) // returns the new document added previsouly
```

## Document management specifications

This section defines how a operation, document and its types are represented and
what are their APIs. An operation identifies uniquely a mutation in the local
document and it is the main data structure passed from user to document when
mutations are to be applied.

A document may contain 3 different types: a list, a map and a registers. Each of
the types may be empty or not. The types' values are pointers to other
documents.

### Document representation and types

**Document data structure**:

```go
type Doc struct {
	deps []string
	map map[string]*Doc
	list []*Doc
	register map[string]interface
}
```

Each `Doc` keeps a list of operation dependencies. An operation dependency is an
operation ID in string format. 

**Types**:

A `Doc` may be of type map, list or register. In
order to ensure the properties required by the JSON CRDT, a node in the document
tree may be at the same time a map, list and/or register (see paper for more
information). Thus, a `Doc` keeps a placeholder for storing information for each
type.

The three types implement  a common `Type` interface:

```go
type Type interface {	
	func Delete()
	func Get()
	func Serialize() ([]byte, error)
}
```

### Operation type

An operation uniquely identifies a mutation in the local document and can be
shared across peers in the network. Each operation contains a `cursor`,
`mutation`, `id` and `dependencies`:

**Data structure**:

```go
type Operation struct {
	Id string
	deps []string
	cursor Cursor
	mutation Mutation
}	
```

**API**:

The `operation` interface defines how to create a new operation as well as all
the components which are part of it - namely `cursor` and `mutation`.

`func NewCursor(path []map[interface]string) (Cursor, error)`

Returns a new `cursor` or error. As input, `NewCursor` receives a list which
indicates the path from the head of the document until the position in the
document referred by the cursor. E.g:

```go
path := []map[interface]string{
		"todo", "list",
		0, "map",
		"done": register
	}

cursor := NewCursor(path)

// cursor will point at {"todo".[0]."done"}
```

`func NewMutation(type string, k interface, v interface) (Mutation, error)`

Returns a new `mutation`. The type input can be one of `INSERT`, `DELETE`.
`ASSIGN`. The `k` and `v` are the key and value, respectively, to apply.

`func NewOperation(cursor Cursor, mut Mutation, deps []string, doc Doc) (error, Operation)`

Creates a new operation based on a previously created `cursor`, `mutation`. The
`operation` id is calculated based on the main `Doc` Id and current Lamport
clock state of the document, thus the document itself must be passed as an
argument. `deps` is a list of strings which are Lamport Ids that describe which
operations the `operation` to construct depends upon.

## Algorithm for applying operation on a document 

