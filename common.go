package shopify

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

// CurrencyCode enum
// USD United States Dollars (USD).
// EUR Euro
// GBP British Pound
// ...
// see more at https://shopify.dev/docs/admin-api/graphql/reference/common-objects/currencycode
type CurrencyCode string

type DateTime string
