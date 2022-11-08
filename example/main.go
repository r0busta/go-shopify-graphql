package main

func main() {
	client := defaultClient()
	// client := clientWithToken()

	// Collections
	collections(client)

	// Products
	products(client)

	// Bulk operations
	// bulk(client)
}
