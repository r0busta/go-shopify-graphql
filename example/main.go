package main

func main() {
	client := defaultClient()
	// client := clientWithToken()
	// client := clientWithVersion()

	// Collections
	collections(client)

	// Products
	listProducts(client)
	// createProduct(client)

	// Bulk operations
	// bulk(client)
}
