# prairiedog

# Overview

# Technical

* Edge weights
* Community detection

* 11-mers should be unique
  * Pick a start and end node: find shortest path between, read edge weights, alternate paths

```python
index = {
  (11-mer): (int),
  ...
  total: (int)
}
```

## 1. Filtering
## 2. Seeding
## 3. Seed set expansion
## 4. Propagation

## Graph Backends:

* https://github.com/snap-stanford/snap: C++
* https://github.com/indradb/indradb: Rust, require db server (postgres/rocksdb)
* https://github.com/cayleygraph/cayley: Go, can embed in a Go application or run as a server and use a Python client (https://github.com/ziyasal/pyley), client-server model not as pleasant
* https://github.com/NationalSecurityAgency/lemongraph: C with Python 2 bindings, file-based lmdb, can embed in a Python app, good interface but Python 2
* https://github.com/networkx/networkx: Python 2/3, in-memory, probably a good tool to prototype with, unsure about scalability
* https://github.com/vgteam/xg
* https://github.com/deepgraph/deepgraph: Python, based off Pandas (likely only in-memory), looks like only an analysis medium & not suitable for live graph construction
* https://github.com/bitnine-oss/agensgraph: C, runs as a server, ontop of PostgreSQL, SQL or Cypher queries
* https://github.com/dgraph-io/dgraph: Go, probably overkill

# Open Questions

* SNPs
* Forward/reverse strands
* Expansion / reference allele bias?
* Problems with circular loops
  * Sol'n?: Shortest path?
* Repeats?
* Overlap by 1 (DAG/graphtyper) or k-1 (de Bruijn)?
* Coordinate system
