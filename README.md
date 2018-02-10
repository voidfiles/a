# A

A project to produce an easy to use authority server

## Contributing

1. Fork
2. Commit
3. Make sure tests pass
4. Create PR
5. Get code review

## Getting Started

This will install dependencies and cayley. You'll be able to run tests and
build the codebase afterwards.

```bash
make init
```

Run the tests

```bash
make test
```

To build

```bash
make build
```

Now, you can find `./_work/a_darmin_amd64` or something like it.

To kick the tires start with a small data set. This will download and load
the LOC childrens subjects headers.

```bash
make small_build
```


## Running locally

```
make build
make run_boltdb
```

You should now have a server running on localhost.

## Small Data Queries

http://localhost:8080/api/v1/query/subject?subject=%3Chttp://id.loc.gov/authorities/childrensSubjects/sj00001253%3E

http://localhost:8080/api/v1/query/object?object=Ceratopsians


## Fun Queries

```js
var vID = "<http://id.loc.gov/authorities/subjects/sh00002388>";
g.V(vID).OutPredicates().ForEach( function(r){
	g.V(vID).Out(r.id).ForEach( function(t){
		var node = {
		  source: r.id,
		  target: t.id
		}
		g.Emit(node)
	})
});
```

Emits this

```json
{
  "result": [
    {
      "source": "<http://www.w3.org/1999/02/22-rdf-syntax-ns#type>",
      "target": "<http://www.w3.org/2004/02/skos/core#Concept>"
    },
    {
      "source": "<http://www.w3.org/2004/02/skos/core#prefLabel>",
      "target": "Antonovych Prize"
    },
    {
      "source": "<http://www.w3.org/2008/05/skos-xl#altLabel>",
      "target": "_:bnode17585414189761283976"
    },
    {
      "source": "<http://www.w3.org/2008/05/skos-xl#altLabel>",
      "target": "_:bnode18314860344878351640"
    },
    {
      "source": "<http://www.w3.org/2008/05/skos-xl#altLabel>",
      "target": "_:bnode5083688311867135174"
    },
    {
      "source": "<http://www.w3.org/2004/02/skos/core#broader>",
      "target": "<http://id.loc.gov/authorities/subjects/sh85113526>"
    },
    {
      "source": "<http://www.w3.org/2004/02/skos/core#inScheme>",
      "target": "<http://id.loc.gov/authorities/subjects>"
    },
    {
      "source": "<http://www.w3.org/2004/02/skos/core#altLabel>",
      "target": "Nahoroda Antonovychiv"
    },
    {
      "source": "<http://www.w3.org/2004/02/skos/core#altLabel>",
      "target": "Nahoroda Fundat︠s︡iï Omeli︠a︡na i Teti︠a︡ny Antonovychiv"
    },
    {
      "source": "<http://www.w3.org/2004/02/skos/core#altLabel>",
      "target": "Nahoroda imeny Teti︠a︡ny i Omeli︠a︡na Antonovychiv"
    },
    {
      "source": "<http://www.w3.org/2004/02/skos/core#changeNote>",
      "target": "_:bnode10482368289832061814"
    },
    {
      "source": "<http://www.w3.org/2004/02/skos/core#changeNote>",
      "target": "_:bnode3260349085687493725"
    }
  ]
}
```
