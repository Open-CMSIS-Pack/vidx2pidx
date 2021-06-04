package main

func main() {
	cmd := NewCli()
	ExitOnError(cmd.Execute())
}
