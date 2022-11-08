package main

func main() {
	client := defaultClient()
	// client := clientWithToken()
	// client := clientWithVersion()

	// Collections
	collections(client)

	// Products
	products(client)

	// Bulk operations
	// bulk(client)
}
