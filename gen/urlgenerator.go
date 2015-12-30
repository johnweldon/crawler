package gen

type UrlGenerator interface {
	Start() <-chan string
}
