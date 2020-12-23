package shopify

import "github.com/r0busta/graphql"

type Money string   // Serialized and truncated to 2 decimals decimal.Decimal
type Decimal string // Serialized decimal.Decimal

type MoneyV2 struct {
	Amount       Decimal      `json:"amount,omitempty"`
	CurrencyCode CurrencyCode `json:"currencyCode,omitempty"`
}

type MoneyBag struct {
	PresentmentMoney MoneyV2 `json:"presentmentMoney,omitempty"`
	ShopMoney        MoneyV2 `json:"shopMoney,omitempty"`
}

type MailingAddress struct {
	// The first line of the address. Typically the street address or PO Box number.
	Address1 graphql.String `json:"address1,omitempty"`

	// The second line of the address. Typically the number of the apartment, suite, or unit.
	Address2 graphql.String `json:"address2,omitempty"`

	// The name of the city, district, village, or town.
	City graphql.String `json:"city,omitempty"`

	// The name of the customer's company or organization.
	Company graphql.String `json:"company,omitempty"`

	// The name of the country.
	Country graphql.String `json:"country,omitempty"`

	// The two-letter code for the country of the address. For example, US.
	CountryCodeV2 CountryCode `json:"countryCodeV2,omitempty"`

	// The first name of the customer.
	FirstName graphql.String `json:"firstName,omitempty"`

	// A formatted version of the address, customized by the provided arguments.
	Formatted []graphql.String `json:"formatted,omitempty"`

	// Whether to include the customer's company in the formatted address.
	// Default value: true
	WithCompany graphql.Boolean `json:"withCompany,omitempty"`

	// Whether to include the customer's name in the formatted address.
	// Default value: false
	WithName graphql.Boolean `json:"withName,omitempty"`

	// A comma-separated list of the values for city, province, and country.
	FormattedArea graphql.String `json:"formattedArea,omitempty"`

	// Globally unique identifier.
	ID graphql.ID `json:"id,omitempty"`

	// The last name of the customer.
	LastName graphql.String `json:"lastName,omitempty"`

	// The latitude coordinate of the customer address.
	Latitude graphql.Float `json:"latitude,omitempty"`

	// The longitude coordinate of the customer address.
	Longitude graphql.Float `json:"longitude,omitempty"`

	// The full name of the customer, based on firstName and lastName.
	Name graphql.String `json:"name,omitempty"`

	// A unique phone number for the customer.
	// Formatted using E.164 standard. For example, +16135551111.
	Phone graphql.String `json:"phone,omitempty"`

	// The region of the address, such as the province, state, or district.
	Province graphql.String `json:"province,omitempty"`

	// The two-letter code for the region.
	// For example, ON.
	ProvinceCode graphql.String `json:"provinceCode,omitempty"`

	// The zip or postal code of the address.
	Zip graphql.String `json:"zip,omitempty"`
}

//CountryCode enum ISO 3166-1 alpha-2 country codes with some differences.
type CountryCode string

// CurrencyCode enum
// USD United States Dollars (USD).
// EUR Euro
// GBP British Pound
// ...
// see more at https://shopify.dev/docs/admin-api/graphql/reference/common-objects/currencycode
type CurrencyCode string

type DateTime string

type PageInfo struct {
	// Indicates if there are more pages to fetch.
	HasNextPage graphql.Boolean `json:"hasNextPage"`
	// Indicates if there are any pages prior to the current page.
	HasPreviousPage graphql.Boolean `json:"hasPreviousPage"`
}

// URL An RFC 3986 and RFC 3987 compliant URI string.
//
// Example value: "https://johns-apparel.myshopify.com".
type URL string
