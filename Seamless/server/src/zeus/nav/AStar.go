package nav

import "errors"
import "math"
import "container/heap"

// AStarNode A*节点
type AStarNode struct {
	meshNode *Node
	parent   *AStarNode
	G        float32
	F        float32
}

func newAStarNode(n *Node, p *AStarNode) *AStarNode {
	node := &AStarNode{
		meshNode: n,
		parent:   p,
		G:        0,
		F:        0,
	}
	return node
}

func (n *AStarNode) getID() int32 {
	return n.meshNode.id
}

// HeapNode 极小堆节点
type HeapNode struct {
	value *AStarNode
	index int
}

func newHeapNode(n *AStarNode) *HeapNode {
	return &HeapNode{
		n,
		0,
	}
}

// HeapNodePool 极小堆节点池
type HeapNodePool struct {
	pool []*HeapNode
}

func newHeapNodePool() *HeapNodePool {
	return &HeapNodePool{
		make([]*HeapNode, 0, 20),
	}
}

func (pq *HeapNodePool) Len() int { return len(pq.pool) }

func (pq *HeapNodePool) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return pq.pool[i].value.F < pq.pool[j].value.F
}

func (pq *HeapNodePool) Swap(i, j int) {
	pq.pool[i], pq.pool[j] = pq.pool[j], pq.pool[i]
	pq.pool[i].index = i
	pq.pool[j].index = j
}

func (pq *HeapNodePool) Push(x interface{}) {
	n := len(pq.pool)
	item := x.(*HeapNode)
	item.index = n
	pq.pool = append(pq.pool, item)
}

func (pq *HeapNodePool) Pop() interface{} {
	old := pq.pool
	n := len(old)
	item := old[n-1]
	item.index = -1 // for safety
	pq.pool = old[0 : n-1]
	return item
}

// OpenList 开放列表
type OpenList struct {
	heap    *HeapNodePool
	nodeMap map[int32]*HeapNode
}

func newOpenList() *OpenList {
	return &OpenList{
		heap:    newHeapNodePool(),
		nodeMap: make(map[int32]*HeapNode),
	}
}

func (h *OpenList) put(n *AStarNode) {
	hn := newHeapNode(n)
	heap.Push(h.heap, hn)
	h.nodeMap[n.getID()] = hn
}

func (h *OpenList) popMin() (*AStarNode, error) {

	if len(h.nodeMap) == 0 {
		return nil, errors.New("no node")
	}

	hn := heap.Pop(h.heap).(*HeapNode)
	delete(h.nodeMap, hn.value.getID())

	return hn.value, nil
}

func (h *OpenList) get(id int32) *AStarNode {

	hn, ok := h.nodeMap[id]
	if !ok {
		return nil
	}

	return hn.value
}

func (h *OpenList) remove(id int32) {
	hn, ok := h.nodeMap[id]
	if !ok {
		return
	}

	delete(h.nodeMap, id)
	heap.Remove(h.heap, hn.index)
}

// AStar A*寻路
type AStar struct {
	mesh      *Mesh
	openList  *OpenList
	closeList map[int32]*AStarNode
}

func newAStar(mesh *Mesh) *AStar {
	return &AStar{
		mesh:      mesh,
		openList:  newOpenList(),
		closeList: make(map[int32]*AStarNode),
	}
}

// FindPath 寻路
func (a *AStar) FindPath(srcNode, destNode *Node) ([]*Node, error) {

	n := newAStarNode(srcNode, nil)
	a.openList.put(n)

	for {
		n, err := a.openList.popMin()

		if err != nil {
			return nil, errors.New("find path failed ")
		}

		a.closeList[n.getID()] = n

		bs := n.meshNode.borders

		for _, b := range bs {
			mn := a.mesh.getNeighbour(n.meshNode, b)
			if mn == nil {
				continue
			}

			//如果当前已经在close列表中
			if _, ok := a.closeList[mn.id]; ok {
				continue
			}

			// 此处G值取两个多边形的中间的距离
			// H 值取 当前点到最终点的x 轴与 Z轴的绝对值的和
			nb := newAStarNode(mn, n)
			nb.G = n.G + b.Cost
			h := math.Abs(float64(destNode.center.X-mn.center.X)) + math.Abs(float64(destNode.center.Z-mn.center.Z))
			nb.F = nb.G + float32(h)

			//如果找到路
			if mn == destNode {
				return a.getPath(nb), nil
			}

			// 如果已经存在openlist中
			if on := a.openList.get(nb.getID()); on != nil {
				//如果结节该结点的G值比较小的话
				if nb.G < on.G {
					a.openList.remove(nb.getID())
					a.openList.put(nb)
				}

			} else {
				a.openList.put(nb)
			}
		}
	}
}

func (a *AStar) getPath(n *AStarNode) []*Node {

	path := make([]*Node, 0, 10)

	path = append(path, n.meshNode)

	for {
		n = n.parent
		if n != nil {
			path = append(path, n.meshNode)
		} else {
			break
		}
	}

	//reverse
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}

	return path
}
