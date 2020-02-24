package main

func init() {
	CreateEnv()
}
func main() {
	defer SaveStack()
	server := SofaTennis{}
	server.run()
}
