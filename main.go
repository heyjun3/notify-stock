package main

func main() {
	client := NewHTTPClient()
	client.FetchCurrentStock()
}
