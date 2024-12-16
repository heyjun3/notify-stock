package main

func main() {
	client := NewHTTPClient()
	client.FetchCurrentStock("^N225")
}
