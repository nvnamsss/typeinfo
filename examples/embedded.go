package examples

type Parent struct {
	Public  int
	private int
}

func (Parent) Meo() {

}

// Ayyooo
func (Parent) Yo() string {
	return "yo"
}

type Child struct {
	Parent
	Public  int
	private int
}
