package proc

type URLProcessor interface {
	Process(in <-chan string) <-chan string
}
