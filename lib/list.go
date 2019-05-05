package jd

import "fmt"

type jsonList []JsonNode

var _ JsonNode = jsonList(nil)

func (l jsonList) Json(metadata ...Metadata) string {
	return renderJson(l)
}

func (l1 jsonList) Equals(n JsonNode, metadata ...Metadata) bool {
	n2 := dispatch(n, metadata)
	l2, ok := n2.(jsonList)
	if !ok {
		return false
	}
	if len(l1) != len(l2) {
		return false
	}
	for i, v1 := range l1 {
		v2 := l2[i]
		if !v1.Equals(v2, metadata...) {
			return false
		}
	}
	return true
}

func (l jsonList) hashCode(metadata []Metadata) [8]byte {
	b := make([]byte, 0, len(l)*8)
	for _, v := range l {
		h := n.hashCode(metadata)
		b = append(b, h[:]...)
	}
	return hash(b)
}

func (l jsonList) Diff(n JsonNode, metadata ...Metadata) Diff {
	return l.diff(n, Path{}, metadata)
}

func (a1 jsonList) diff(n JsonNode, path Path, metadata []Metadata) Diff {
	d := make(Diff, 0)
	a2, ok := n.(jsonList)
	if !ok {
		// Different types
		e := DiffElement{
			Path:      path.clone(),
			OldValues: nodeList(a1),
			NewValues: nodeList(n),
		}
		return append(d, e)
	}
	maxLen := len(a1)
	if len(a1) < len(a2) {
		maxLen = len(a2)
	}
	for i := 0; i < maxLen; i++ {
		a1Has := i < len(a1)
		a2Has := i < len(a2)
		subPath := append(path.clone(), float64(i))
		if a1Has && a2Has {
			n1 := dispatch(a1[i])
			n2 := dispatch(a2[i])
			subDiff := n1.diff(n2, subPath, metadata)
			d = append(d, subDiff...)
		}
		if a1Has && !a2Has {
			e := DiffElement{
				Path:      subPath,
				OldValues: nodeList(a1[i]),
				NewValues: nodeList(),
			}
			d = append(d, e)
		}
		if !a1Has && a2Has {
			e := DiffElement{
				Path:      subPath,
				OldValues: nodeList(),
				NewValues: nodeList(a2[i]),
			}
			d = append(d, e)
		}
	}
	return d
}

func (l jsonList) Patch(d Diff, metadata ...Metadata) (JsonNode, error) {
	return patchAll(l, d, metadata)
}

func (l jsonList) patch(pathBehind, pathAhead Path, oldValues, newValues []JsonNode, metadata []Metadata) (JsonNode, error) {

	if len(oldValues) > 1 || len(newValues) > 1 {
		return patchErrNonSetDiff(oldValues, newValues, pathBehind)
	}
	oldValue := singleValue(oldValues)
	newValue := singleValue(newValues)
	// Base case
	if len(pathAhead) == 0 {
		if !a.Equals(oldValue) {
			return patchErrExpectValue(oldValue, a, pathBehind)
		}
		return newValue, nil
	}
	// Recursive case
	pe, ok := pathAhead[0].(float64)
	if !ok {
		return nil, fmt.Errorf(
			"Invalid path element %v. Expected float64.",
			pathAhead[0])
	}
	i := int(pe)
	var nextNode JsonNode = voidNode{}
	if len(a) > i {
		nextNode = a[i]
	}
	patchedNode, err := nextNode.patch(append(pathBehind, pe), pathAhead[1:], oldValues, newValues, metadata)
	if err != nil {
		return nil, err
	}
	if isVoid(patchedNode) {
		if i != len(a)-1 {
			return nil, fmt.Errorf(
				"Removal of a non-terminal element of an array.")
		}
		// Delete an element
		return a[:len(a)-1], nil
	}
	if i > len(a) {
		return nil, fmt.Errorf(
			"Addition beyond the terminal element of an array.")
	}
	if i == len(a) {
		// Add an element
		return append(a, patchedNode), nil
	}
	// Replace an element
	a[i] = patchedNode
	return a, nil
}