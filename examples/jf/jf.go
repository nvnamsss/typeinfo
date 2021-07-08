package examples

type Layer1 struct {
	Value  int64 `tinfo:"ignore"`
	Layer2 Layer2
}

func (Layer1) AAA() int64 {
	return 0
}

func (Layer1) Contains(a int64, arr ...int64) bool {
	return true
}

type Layer2 struct {
	Value2 int64
	Layer3 Layer3
}

func (Layer2) AAA() int64 {
	return 0
}

type Layer3 struct {
	Value3 int64
	Layer4 Layer4
}

func (Layer3) AAA() int64 {
	return 0
}

type Layer4 struct {
	Value4 int64
}
