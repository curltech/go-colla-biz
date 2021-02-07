package actual

import (
	"github.com/curltech/go-colla-biz/spec/entity"
	"github.com/curltech/go-colla-core/logger"
	"github.com/curltech/go-colla-core/util/convert"
	"github.com/huandu/xstrings"
	"strings"
)

const (
	Path_Seperator = "."
	Path_RootChar  = "/"
	Path_FuzzyChar = "*"
	Path_StartChar = "["
	Path_EndChar   = "]"
)

/**
 * @Description: R：从根节点开始；C：从当前节点开始；F：模糊查询
 * @author: Jason
 * @date: 2017年4月17日 下午8:24:56
 *
 */
const (
	StartPosition_Root    = "R"
	StartPosition_Current = "C"
	StartPosition_Fuzzy   = "F"
)

type Node struct {
	Kind     string
	Position int
}

type PositionPath struct {
	Current  int
	Starter  string //StartPosition.R
	Nodes    []*Node
	SpecType string
}

func NewPositionPath(path string) *PositionPath {
	pp := PositionPath{}
	pp.splitPath(path)

	return &pp
}

/**
 * @Description: 获取最后的名称
 * @param path
 * @return
 */
func GetLastKind(path string) string {
	start := strings.LastIndex(path, ".")
	end := strings.LastIndex(path, "[")
	if end <= start {
		end = -1
	}
	return xstrings.Slice(path, start, end)
}

/**
 * @Description: 获取路径
 * @param path
 * @return
 */
func GetRolePath(path string) string {
	var node string
	if !strings.Contains(path, ".") {
		node = ""
	} else {
		start := strings.Index(path, ":")
		end := strings.LastIndex(path, ".")
		node = xstrings.Slice(path, start, end)
	}

	return "R:" + node
}

/**
 * @Description: 获取最后的节点
 * @param path
 * @return
 */
func (this *PositionPath) getLastPosition(path string) int {
	start := strings.LastIndex(path, ".")
	node := xstrings.Slice(path, start, -1)
	if node == "" {
		node = path
	}
	start = strings.LastIndex(node, "[")
	node = xstrings.Slice(node, start, -1)
	if node == "" {
		return 0
	} else {
		end := strings.LastIndex(node, "]")
		node = xstrings.Slice(node, 0, end)
		if node == "" {
			return 0
		} else {
			o, _ := convert.ToObject(node, "int")

			return o.(int)
		}
	}
}

/**
 * @Description: 获取当前的节点
 * @return
 */
func (this *PositionPath) getCurrent() *Node {
	if this.Current < len(this.Nodes) && this.Current > -1 {
		return this.Nodes[this.Current]
	}

	return nil
}

/**
 * @Description: 获取最后的节点
 * @return
 */
func (this *PositionPath) getLast() *Node {
	return this.Nodes[len(this.Nodes)-1]
}

/**
 * @Description: 节点数目
 * @return
 */
func (this *PositionPath) getSize() int {
	return len(this.Nodes)
}

/**
 * @Description: 是否有下一个节点
 * @return
 */
func (this *PositionPath) hasNext() bool {
	if this.Current < len(this.Nodes)-1 {
		return true
	} else {
		return false
	}
}

/**
 * @Description: 下一个路径
 * @return
 */
func (this *PositionPath) next() *PositionPath {
	if this.hasNext() {
		this.Current++
	}

	return this
}

/**
 * @Description: 将一个路径化成路径节点的数组
 * @param path
 */
func (this *PositionPath) splitPath(path string) {
	p := path
	start := strings.Index(path, ":")
	if start > -1 {
		specType := xstrings.Slice(path, 0, start)
		if specType != "" {
			if specType == "Role" || specType == "R" {
				this.SpecType = entity.SpecType_Role
			} else if specType == "ActionResult" || specType == "A" {
				this.SpecType = entity.SpecType_Action
			} else if specType == "Property" || specType == "P" {
				this.SpecType = entity.SpecType_Property
			} else if specType == "Connection" || specType == "C" {
				this.SpecType = entity.SpecType_Connection
			}
		}
		p = xstrings.Slice(path, start+1, -1)
	}
	if this.SpecType == "" {
		this.SpecType = entity.SpecType_Property
	}
	// 默认为当前路径
	this.Starter = StartPosition_Current

	// 根路径
	if strings.HasPrefix(p, Path_RootChar) {
		this.Starter = StartPosition_Root
		p = xstrings.Slice(p, 1, -1)
	} else if strings.HasPrefix(p, Path_Seperator) {
		this.Starter = StartPosition_Current
		p = xstrings.Slice(p, 1, -1)
	} else if strings.HasPrefix(p, Path_FuzzyChar) {
		this.Starter = StartPosition_Fuzzy
		p = xstrings.Slice(p, 1, -1)
	}

	ns := strings.Split(p, Path_Seperator)
	for _, n := range ns {
		begin := strings.Index(n, Path_StartChar)
		end := strings.Index(n, Path_EndChar)
		var kind string
		var position int
		if begin > -1 && end > begin+1 {
			kind = xstrings.Slice(n, 0, begin)
			pos := xstrings.Slice(n, begin+1, end)
			o, err := convert.ToObject(pos, "int")
			if err != nil {
				logger.Sugar.Errorf("NumberFormatException")
			}
			position = o.(int)
		} else if begin > -1 && end > begin {
			kind = xstrings.Slice(n, 0, begin) // add by lxp, for path eg.
			// ./xxx/xxx[]
			position = 0
		} else {
			kind = n
			// position = 0;
		}
		if kind != "" {
			this.Nodes = append(this.Nodes, &Node{Kind: kind, Position: position})
		}
		this.Current = 0
	}
}

/**
 * @Description: 开始位置
 * @return
 */
func (this *PositionPath) startWhere() string {
	return this.Starter
}
