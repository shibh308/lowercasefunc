package a

func run() {
	var gopher int
	print(gopher)
}

func Run() { // want `{FuncName:Run, TargetPos:"a.go:6:3", CalledPos:\["a.go:2:9", "a.go:26:10"\]}`
	run()
	func () func () {return run}()
}
func RUN() {
	run()
}
