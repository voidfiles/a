# A

A project to produce an easy to use authority server

## Contributing

1. Fork
2. Commit
3. Make sure tests pass
4. Create PR
5. Get code review

## Getting Started

To download some data and get started. This will download and initialize a database with the
LC Children's Subject Headings. It's just a small database.

```
make init
```

To run the container run

```
make run
```

To stop the containers

```
make stop
```

Once this is up and running you should be able to see a server running at http://localhost:64210/

To see the logs of the server

```
make logs
```

## Small Data Queries

http://localhost:8080/api/v1/query/subject?subject=%3Chttp://id.loc.gov/authorities/childrensSubjects/sj00001253%3E

http://localhost:8080/api/v1/query/object?query=Ceratopsians


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
