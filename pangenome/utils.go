package pangenome

// import "sort"

// func makeFacets(facetKVs map[string]string) (fs []*facets.Facet, err error) {
// 	if len(facetKVs) == 0 {
// 		return nil, nil
// 	}
// 	allKeys := make([]string, 0, len(facetKVs))
// 	for k := range facetKVs {
// 		allKeys = append(allKeys, k)
// 	}
// 	sort.Strings(allKeys)
// 	for _, k := range allKeys {
// 		v := facetKVs[k]
// 		typ, err := facets.ValType(v)
// 		if err != nil {
// 			return nil, err
// 		}
// 		fs = append(fs, &facets.Facet{
// 			k,
// 			[]byte(v),
// 			typ,
// 		})
// 	}
// 	return fs, nil
// }

// func addEdgeToUID(attr string, src string,
// 	dst uint64, facetKVs map[string]string) {
// 	fs, err := makeFacets(facetKVs)
// 	require.NoError(t, err)
// 	edge := &task.DirectedEdge{
// 		ValueId: dst,
// 		Label:   "edges",
// 		Attr:    attr,
// 		Entity:  src,
// 		Op:      task.DirectedEdge_SET,
// 		Facets:  fs,
// 	}
// 	addEdge(t, attr, src, edge)
// }
