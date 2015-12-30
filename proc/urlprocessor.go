package proc

type UrlProcessor interface {
	Process(in <-chan string) <-chan string
}
