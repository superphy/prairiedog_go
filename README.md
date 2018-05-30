# prairiedog

# Introduction

Public health laboratories are currently moving to whole-genome sequence (WGS) based analyses, and require software methods for characterizing bacterial species based on WGS data.
WGS data for different strains and variants of species can be collected to create a reference genome for comparison with new samples.
These reference "pan-genomes" contain all genes, both core and accessory, and known variants.

Currently, researchers have focused on the development of software to construct pan-genomes.
While useful for comparison purposes, existing software do not readily update a constructed pan-genome or integrate variant calling [Computational Pan-Genomics Consortium, 2016] - the identification of a novel gene variants when compared against the reference pan-genome.
Furthermore, there is a lack of consensus towards indexing pan-genomes for referencing [Sirén, Valimaki, Makinen, 2014], or storing and restoring a computed pan-genome.

Graph representations of pan-genomes have emerged [Paten et al., 2017] as a possible solution to these challenges.
In contrast to linear representations, graphs can embed gene annotations and be used as a reference set.
The most common forms are a De Bruijn Graph or a Directed Acyclic Graph [Paten et al., 2017].
Both methods are based off k-mers, the separation of a genome into sequences of length k.

Existing software either require large RAM stores [Marcus, Lee, Schatz, 2014] or are not performant on update tasks [Sheikhizadeh et al., 2016].
<<<<<<< HEAD
While a consensus on indexing methods is unlikely [Sirén, Valimaki, Makinen, 2014], we aim to develop a software approach to creating and updating pan-genomes which perform on both low and high RAM systems, and focuses on supporting the update of a pan-genome as new samples are sequenced.
Our approach offloads data storage to a generic graph store, and focuses on uncompressed nodes with the goal on minimizing edges.
=======
While a consensus on indexing methods is unlikely [Sirén, Valimaki, Makinen, 2014], we aim to develop a software application for creating and updating pan-genomes which scales to both low and high RAM systems, and focuses on supporting the update of a pan-genome as new samples are sequenced.
Our approach borrows ideas and software from network analysis, namely the core genome is identified using a modified PageRank algorithm [Whang, Gleich, Dhillon, 2013] from community detection research.
>>>>>>> 767731b096a93f108e16ccc059d0d881d77354a0

# Implementation

prairiedog is a Go application which uses [Dgraph](https://github.com/dgraph-io/dgraph) for the graph layer and [Badger](https://github.com/dgraph-io/badger) for the k-mer: count mapping.
We chose Dgraph as a hedge for eventual sharding requirements depending on the size of the sampled population.
For the data structure, we use a novel combined approach integrating a De Bruijn graph (with 11-mer nodes) along with weighted directed edges representing emission probabilities as in a Li-Stephens model [Li, Stephens, 2003].
<<<<<<< HEAD
The core construction method is as follows:

```python
for kmer, nextkmer in seed/new genome:
  if kmer in graph:
    node_count+=1
    color
  else:
    create_node()

  if nextkmer not in graph:
    create_node()
    color

  if edge(kmer, nextkmer):
    edge_weight +=1
  else:
    create_edge()
```

While edge weights have sampling bias, this allows us to simulate haplotypes which has not been possible in the past [Computational Pan-Genomics Consortium, 2016].
Retrieving genes and variants can be performed as in a De Bruijn graph with the added benefit of estimating occurrence probabilities based off a given collection of genomes used in construction.
Our approach also extends to new samples and records meaningful data as more samples are sequenced.
=======
We create a pan-genome by applying the PageRank algorithm through k-mers generated from each sample of a species using a *random-walk with restart* approach [Pan et al., 2004], and also increase the weight of the edges as they are encountered.
While exposing prairiedog to sampling bias, weighted edges allows us to simulate haplotypes which has not been possible in the past [Computational Pan-Genomics Consortium, 2016].
Retrieving genes and variants can be performed as in a De Bruijn graph with the added benefit of estimating occurrence probabilities, as based off the genomes used in construction.
Our approach also extends to new samples and is more accurate as additional samples are sequenced.
>>>>>>> 767731b096a93f108e16ccc059d0d881d77354a0

It is possible we may implement a compression approach, splitting overlapping segments of an 11-mer into separate edges.
This would reduce the number of nodes, and possibly edges, but would come at the cost of indexing complete 11-mers with their paths and traversing every split node when we add new genomes.
We would have to gauge the size reduction to performance tradeoff.

# Conclusion

Completion of this project will allow users, such as reference laboratories, to characterize any bacterial species by creating a reference pan-genome which can be updated as new samples are sequenced.
A novel result will be the retrieval of haplotypes and variants with probability estimations, as based off the sampled population.
Our design focuses on integrating novel information into the pan-genome with every additional sample sequenced.

# References

[Computational Pan-Genomics Consortium, 2016]: Computational Pan-Genomics Consortium. (2016). Computational pan-genomics: status, promises and challenges. Briefings in Bioinformatics, bbw089.

[Paten et al., 2017]: Paten, B., Novak, A. M., Eizenga, J. M., & Garrison, E. (2017). Genome graphs and the evolution of genome inference. Genome research, 27(5), 665-676.

[Marcus, Lee, Schatz, 2014]: Marcus, S., Lee, H., & Schatz, M. C. (2014). SplitMEM: a graphical algorithm for pan-genome analysis with suffix skips. Bioinformatics, 30(24), 3476-3483.

[Sheikhizadeh et al., 2016]: Sheikhizadeh, S., Schranz, M. E., Akdel, M., de Ridder, D., & Smit, S. (2016). PanTools: representation, storage and exploration of pan-genomic data. Bioinformatics, 32(17), i487-i493.

[Sirén, Valimaki, Makinen, 2014]: Sirén, J., Välimäki, N., & Mäkinen, V. (2014). Indexing graphs for path queries with applications in genome research. IEEE/ACM Transactions on Computational Biology and Bioinformatics (TCBB), 11(2), 375-388.

[Whang, Gleich, Dhillon, 2013]: Whang, J. J., Gleich, D. F., & Dhillon, I. S. (2013, October). Overlapping community detection using seed set expansion. In Proceedings of the 22nd ACM international conference on Conference on information & knowledge management (pp. 2099-2108). ACM.

[Li, Stephens, 2003]: Li, N., & Stephens, M. (2003). Modeling linkage disequilibrium and identifying recombination hotspots using single-nucleotide polymorphism data. Genetics, 165(4), 2213-2233.

[Pan et al., 2004]: Pan, J. Y., Yang, H. J., Faloutsos, C., & Duygulu, P. (2004, August). Automatic multimedia cross-modal correlation discovery. In Proceedings of the tenth ACM SIGKDD international conference on Knowledge discovery and data mining (pp. 653-658). ACM.

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

# Directions

* PageRank (a modified approach [Whang, Gleich, Dhillon, 2013]) might come in handy for some operations (i.e. with de bruijn graph pulling out SNPs is not easy unless paths are bounded, so haplotypes might provide those bounds)

# Open Questions

* Most pan-genome graph algorithmns have focused on custom data structures (in the code) + compression - unsure how much of a concern this is if we can offload the data structures to established systems
* SNPs
* Forward/reverse strands
* Expansion / reference allele bias?
* Problems with circular loops
  * Sol'n?: Shortest path?
* Repeats?
* Overlap by 1 (DAG/graphtyper) or k-1 (de Bruijn)?
* Coordinate system
