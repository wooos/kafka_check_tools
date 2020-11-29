package main

func main() {
	cmd := newRootCommand()

	_ = cmd.Execute()
}
