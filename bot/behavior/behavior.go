package behavior

import (
	"encoding/xml"
	"sort"
)

type Mode int

const (
	Thread Mode = 1 + iota // 线程运行模式（一般用于压测
	Block                  // 阻塞运行模式（一般用于http调用
	Step                   // 步进运行模式（一般用于调试
)

type Tree struct {
	ID string `xml:"id"`
	Ty string `xml:"ty"`

	Wait int32 `xml:"wait"`

	Loop int32  `xml:"loop"` // 用于记录循环节点的循环x次数
	Code string `xml:"code"`

	Pos struct {
		X float32 `xml:"x"`
		Y float32 `xml:"y"`
	} `xml:"pos"`
	Alias string `xml:"alias"`

	root INod
	mode Mode

	Children []*Tree `xml:"children"`
}

func (t *Tree) GetRoot() INod {
	return t.root
}

func (t *Tree) GetMode() Mode {
	return t.mode
}

func (t *Tree) link(self INod, parent INod, mode Mode) {

	self.Init(t, parent, mode)

	//当前节点的子节点个数大于1，则根据t.pos.x的值，将当前节点的子节点进行排序从小到大
	// modified by dgh 2023-10-24
	if len(t.Children) > 1 {
		sort.Slice(t.Children, func(i, j int) bool {
			return t.Children[i].Pos.X < t.Children[j].Pos.X
		})
	}

	for k := range t.Children {
		child := NewNode(t.Children[k].Ty).(INod)
		self.AddChild(child)
		t.Children[k].link(child, self, mode)
	}
}

func Load(f []byte, mode Mode) (*Tree, error) {

	tree := &Tree{
		mode: mode,
	}

	err := xml.Unmarshal([]byte(f), &tree)
	if err != nil {
		panic(err)
	}

	tree.root = NewNode(tree.Ty).(INod)
	tree.root.Init(tree, nil, mode)

	//当前节点的子节点个数大于1，则根据k.pos.x的值，将当前节点的子节点进行排序从小到大
	// modified by dgh 2023-10-24
	if len(tree.Children) > 1 {
		sort.Slice(tree.Children, func(i, j int) bool {
			return tree.Children[i].Pos.X < tree.Children[j].Pos.X
		})
	}

	for k := range tree.Children {
		cn := NewNode(tree.Children[k].Ty).(INod)
		cn.Init(tree.Children[k], tree.root, mode)
		tree.root.AddChild(cn)

		tree.Children[k].link(cn, tree.root, mode)
	}
	return tree, nil
}

// 定义一个反序列化函数
func UnmarshalXML(tree *Tree) ([]byte, error)  {
	type behavior Tree
	return xml.MarshalIndent((*behavior)(tree), "", "\t")
}
