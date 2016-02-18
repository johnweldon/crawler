package gen

type URLGenerator interface {
	Start() <-chan string
}
