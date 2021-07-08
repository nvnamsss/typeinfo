package examples

type Layer1 struct {
	Value  int64
	Layer2 Layer2
}

type Layer2 struct {
	Value2 int64
	Layer3 Layer3
}

func (Layer2) F() {

}

type Layer3 struct {
	Value3 int64
}

func (Layer3) F() {

}
