package examples

type PS1 struct {
	Value int64
	PS2   *PS2
}

func (this *PS1) Nothing() {

}

func (this *PS1) Call(array []int64) {
}
func (PS1) Params(array []int64, values ...int64) {

}

type PS2 struct {
	Value int64
}
