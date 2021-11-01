package shopify

import "gopkg.in/guregu/null.v4"

type QueryRoot struct {
	// Lookup an App by ID or return the currently authenticated App.
	// App *App `json:"app,omitempty"`
	// Fetches app by handle.
	// Returns null if the app doesn't exist.
	// AppByHandle *App `json:"appByHandle,omitempty"`
	// Fetches app by apiKey.
	// Returns null if the app doesn't exist.
	// AppByKey *App `json:"appByKey,omitempty"`
	// Lookup an AppInstallation by ID or return the AppInstallation for the currently authenticated App.
	// AppInstallation *AppInstallation `json:"appInstallation,omitempty"`
	// List of app installations.
	// AppInstallations *AppInstallationConnection `json:"appInstallations,omitempty"`
	// Returns an automatic discount resource by ID.
	// AutomaticDiscount DiscountAutomatic `json:"automaticDiscount,omitempty"`
	// Returns an automatic discount resource by ID.
	// AutomaticDiscountNode *DiscountAutomaticNode `json:"automaticDiscountNode,omitempty"`
	// List of automatic discounts.
	// AutomaticDiscountNodes *DiscountAutomaticNodeConnection `json:"automaticDiscountNodes,omitempty"`
	// List of the shop's automatic discount saved searches.
	// AutomaticDiscountSavedSearches *SavedSearchConnection `json:"automaticDiscountSavedSearches,omitempty"`
	// List of automatic discounts.
	// AutomaticDiscounts *DiscountAutomaticConnection `json:"automaticDiscounts,omitempty"`
	// List of activated carrier services and which shop locations support them.
	// AvailableCarrierServices []*DeliveryCarrierServiceAndLocations `json:"availableCarrierServices,omitempty"`
	// List of available locales.
	// AvailableLocales []*Locale `json:"availableLocales,omitempty"`
	// Lookup a carrier service by ID.
	// CarrierService *DeliveryCarrierService `json:"carrierService,omitempty"`
	// Lookup a channel by ID.
	// Channel *Channel `json:"channel,omitempty"`
	// List of the active sales channels.
	// Channels *ChannelConnection `json:"channels,omitempty"`
	// Returns a code discount resource by ID.
	// CodeDiscountNode *DiscountCodeNode `json:"codeDiscountNode,omitempty"`
	// Returns a code discount identified by its code.
	// CodeDiscountNodeByCode *DiscountCodeNode `json:"codeDiscountNodeByCode,omitempty"`
	// List of code discounts. Special fields for query params:
	//  * status: active, expired, scheduled
	//  * discount_type: bogo, fixed_amount, free_shipping, percentage.
	// CodeDiscountNodes *DiscountCodeNodeConnection `json:"codeDiscountNodes,omitempty"`
	// List of the shop's code discount saved searches.
	// CodeDiscountSavedSearches *SavedSearchConnection `json:"codeDiscountSavedSearches,omitempty"`
	// Returns a Collection resource by ID.
	Collection *Collection `json:"collection,omitempty"`
	// Return a collection by its handle.
	CollectionByHandle *Collection `json:"collectionByHandle,omitempty"`
	// A list of rule conditions to define how collections with rules can be created.
	// CollectionRulesConditions []*CollectionRuleConditions `json:"collectionRulesConditions,omitempty"`
	// List of the shop's collection saved searches.
	// CollectionSavedSearches *SavedSearchConnection `json:"collectionSavedSearches,omitempty"`
	// List of collections.
	// Collections *CollectionConnection `json:"collections,omitempty"`
	// Return the AppInstallation for the currently authenticated App.
	// CurrentAppInstallation *AppInstallation `json:"currentAppInstallation,omitempty"`
	// Returns the current app's most recent BulkOperation.
	// CurrentBulkOperation *BulkOperation `json:"currentBulkOperation,omitempty"`
	// Returns a Customer resource by ID.
	// Customer *Customer `json:"customer,omitempty"`
	// Returns a CustomerPaymentMethod resource by ID.
	// CustomerPaymentMethod *CustomerPaymentMethod `json:"customerPaymentMethod,omitempty"`
	// List of the shop's customer saved searches.
	// CustomerSavedSearches *SavedSearchConnection `json:"customerSavedSearches,omitempty"`
	// List of customers.
	// Customers *CustomerConnection `json:"customers,omitempty"`
	// The paginated list of deletion events.
	// DeletionEvents *DeletionEventConnection `json:"deletionEvents,omitempty"`
	// Lookup a Delivery Profile by ID.
	// DeliveryProfile *DeliveryProfile `json:"deliveryProfile,omitempty"`
	// List of saved delivery profiles.
	// DeliveryProfiles *DeliveryProfileConnection `json:"deliveryProfiles,omitempty"`
	// The shop-wide shipping settings.
	// DeliverySettings *DeliverySetting `json:"deliverySettings,omitempty"`
	// The total number of discount codes for the shop.
	DiscountCodeCount int `json:"discountCodeCount,omitempty"`
	// Returns a bulk code creation resource by ID.
	// DiscountRedeemCodeBulkCreation *DiscountRedeemCodeBulkCreation `json:"discountRedeemCodeBulkCreation,omitempty"`
	// List of the shop's redeemed discount code saved searches.
	// DiscountRedeemCodeSavedSearches *SavedSearchConnection `json:"discountRedeemCodeSavedSearches,omitempty"`
	// Lookup a Domain by ID.
	// Domain *Domain `json:"domain,omitempty"`
	// Returns a DraftOrder resource by ID.
	// DraftOrder *DraftOrder `json:"draftOrder,omitempty"`
	// List of the shop's draft order saved searches.
	// DraftOrderSavedSearches *SavedSearchConnection `json:"draftOrderSavedSearches,omitempty"`
	// List of saved draft orders.
	// DraftOrders *DraftOrderConnection `json:"draftOrders,omitempty"`
	// A list of the shop's file saved searches.
	// FileSavedSearches *SavedSearchConnection `json:"fileSavedSearches,omitempty"`
	// A list of files.
	// Files *FileConnection `json:"files,omitempty"`
	// Returns a Fulfillment resource by ID.
	// Fulfillment *Fulfillment `json:"fulfillment,omitempty"`
	// Returns a Fulfillment order resource by ID.
	FulfillmentOrder *FulfillmentOrder `json:"fulfillmentOrder,omitempty"`
	// Returns a FulfillmentService resource by ID.
	FulfillmentService *FulfillmentService `json:"fulfillmentService,omitempty"`
	// Returns a gift card resource by ID.
	// GiftCard *GiftCard `json:"giftCard,omitempty"`
	// Returns a list of gift cards.
	// GiftCards *GiftCardConnection `json:"giftCards,omitempty"`
	// The total number of gift cards issued for the shop.
	GiftCardsCount null.String `json:"giftCardsCount,omitempty"`
	// Returns an InventoryItem resource by ID.
	InventoryItem *InventoryItem `json:"inventoryItem,omitempty"`
	// List of inventory items.
	// InventoryItems *InventoryItemConnection `json:"inventoryItems,omitempty"`
	// Returns an InventoryLevel resource by ID.
	InventoryLevel *InventoryLevel `json:"inventoryLevel,omitempty"`
	// Returns a Job resource by ID. Used to check the status of internal jobs and any applicable changes.
	// Job *Job `json:"job,omitempty"`
	// Returns an inventory Location resource by ID.
	Location *Location `json:"location,omitempty"`
	// List of active locations.
	// Locations *LocationConnection `json:"locations,omitempty"`
	// Returns a list of all origin locations available for a delivery profile.
	LocationsAvailableForDeliveryProfiles []*Location `json:"locationsAvailableForDeliveryProfiles,omitempty"`
	// Returns a list of all origin locations available for a delivery profile.
	// LocationsAvailableForDeliveryProfilesConnection *LocationConnection `json:"locationsAvailableForDeliveryProfilesConnection,omitempty"`
	// List of a campaign's marketing activities.
	// MarketingActivities *MarketingActivityConnection `json:"marketingActivities,omitempty"`
	// Returns a MarketingActivity resource by ID.
	// MarketingActivity *MarketingActivity `json:"marketingActivity,omitempty"`
	// Returns a MarketingEvent resource by ID.
	// MarketingEvent *MarketingEvent `json:"marketingEvent,omitempty"`
	// List of marketing events.
	// MarketingEvents *MarketingEventConnection `json:"marketingEvents,omitempty"`
	// Returns a metafield by ID.
	Metafield *Metafield `json:"metafield,omitempty"`
	// Returns a metafield definition by ID.
	// MetafieldDefinition *MetafieldDefinition `json:"metafieldDefinition,omitempty"`
	// All available metafield definition types.
	// MetafieldDefinitionTypes []*MetafieldDefinitionType `json:"metafieldDefinitionTypes,omitempty"`
	// List of metafield definitions.
	// MetafieldDefinitions *MetafieldDefinitionConnection `json:"metafieldDefinitions,omitempty"`
	// List of metafield namespaces and keys visible to the Storefront API.
	// MetafieldStorefrontVisibilities *MetafieldStorefrontVisibilityConnection `json:"metafieldStorefrontVisibilities,omitempty"`
	// Returns metafield storefront visibility by ID.
	// MetafieldStorefrontVisibility *MetafieldStorefrontVisibility `json:"metafieldStorefrontVisibility,omitempty"`
	// Returns a specific node by ID.
	// Node Node `json:"node,omitempty"`
	// Returns the list of nodes with the given IDs.
	// Nodes []Node `json:"nodes,omitempty"`
	// Returns an Order resource by ID.
	Order *Order `json:"order,omitempty"`
	// List of the shop's order saved searches.
	// OrderSavedSearches *SavedSearchConnection `json:"orderSavedSearches,omitempty"`
	// Returns a list of orders placed.
	// Orders *OrderConnection `json:"orders,omitempty"`
	// The list of payment terms templates eligible for all shops and users.
	// PaymentTermsTemplates []*PaymentTermsTemplate `json:"paymentTermsTemplates,omitempty"`
	// A list of price lists.
	// PriceList *PriceList `json:"priceList,omitempty"`
	// All price lists for a shop.
	// PriceLists *PriceListConnection `json:"priceLists,omitempty"`
	// Lookup a price rule by ID.
	// PriceRule *PriceRule `json:"priceRule,omitempty"`
	// List of the shop's price rule saved searches.
	// PriceRuleSavedSearches *SavedSearchConnection `json:"priceRuleSavedSearches,omitempty"`
	// List of price rules.
	// PriceRules *PriceRuleConnection `json:"priceRules,omitempty"`
	// Returns a private metafield by ID.
	// PrivateMetafield *PrivateMetafield `json:"privateMetafield,omitempty"`
	// List of private metafields.
	// PrivateMetafields *PrivateMetafieldConnection `json:"privateMetafields,omitempty"`
	// Returns a Product resource by ID.
	// Product *Product `json:"product,omitempty"`
	// Return a product by its handle.
	// ProductByHandle *Product `json:"productByHandle,omitempty"`
	// The product resource feedback for the currently authenticated app.
	// ProductResourceFeedback *ProductResourceFeedback `json:"productResourceFeedback,omitempty"`
	// List of the shop's product saved searches.
	// ProductSavedSearches *SavedSearchConnection `json:"productSavedSearches,omitempty"`
	// Returns a ProductVariant resource by ID.
	ProductVariant *ProductVariant `json:"productVariant,omitempty"`
	// List of the product variants.
	// ProductVariants *ProductVariantConnection `json:"productVariants,omitempty"`
	// List of products.
	// Products *ProductConnection `json:"products,omitempty"`
	// The list of public Admin API versions, including supported, release candidate and unstable versions.
	// PublicAPIVersions []*APIVersion `json:"publicApiVersions,omitempty"`
	// Lookup a publication by ID.
	// Publication *Publication `json:"publication,omitempty"`
	// List of the active publications.
	// Publications *PublicationConnection `json:"publications,omitempty"`
	// Returns a Refund resource by ID.
	// Refund *Refund `json:"refund,omitempty"`
	// Lookup a script tag resource by ID.
	// ScriptTag *ScriptTag `json:"scriptTag,omitempty"`
	// A list of script tags.
	// ScriptTags *ScriptTagConnection `json:"scriptTags,omitempty"`
	// Returns a Selling Plan Group resource by ID.
	// SellingPlanGroup *SellingPlanGroup `json:"sellingPlanGroup,omitempty"`
	// List Selling Plan Groups.
	// SellingPlanGroups *SellingPlanGroupConnection `json:"sellingPlanGroups,omitempty"`
	// Returns a Shop resource corresponding to access token used in request.
	// Shop *Shop `json:"shop,omitempty"`
	// List of locales available on a shop.
	// ShopLocales []*ShopLocale `json:"shopLocales,omitempty"`
	// Shopify Payments account information, including balances and payouts.
	// ShopifyPaymentsAccount *ShopifyPaymentsAccount `json:"shopifyPaymentsAccount,omitempty"`
	// All available standard metafield definition templates.
	// StandardMetafieldDefinitionTemplates *StandardMetafieldDefinitionTemplateConnection `json:"standardMetafieldDefinitionTemplates,omitempty"`
	// Returns a SubscriptionBillingAttempt by ID.
	// SubscriptionBillingAttempt *SubscriptionBillingAttempt `json:"subscriptionBillingAttempt,omitempty"`
	// Returns a Subscription Contract resource by ID.
	// SubscriptionContract *SubscriptionContract `json:"subscriptionContract,omitempty"`
	// List Subscription Contracts.
	// SubscriptionContracts *SubscriptionContractConnection `json:"subscriptionContracts,omitempty"`
	// Returns a Subscription Draft resource by ID.
	// SubscriptionDraft *SubscriptionDraft `json:"subscriptionDraft,omitempty"`
	// List of TenderTransactions associated with the Shop.
	// TenderTransactions *TenderTransactionConnection `json:"tenderTransactions,omitempty"`
	// Translatable resource.
	// TranslatableResource *TranslatableResource `json:"translatableResource,omitempty"`
	// List of translatable resources.
	// TranslatableResources *TranslatableResourceConnection `json:"translatableResources,omitempty"`
	// Returns a redirect resource by ID.
	// URLRedirect *URLRedirect `json:"urlRedirect,omitempty"`
	// Returns a redirect import resource by ID.
	// URLRedirectImport *URLRedirectImport `json:"urlRedirectImport,omitempty"`
	// A list of the shop's URL redirect saved searches.
	// URLRedirectSavedSearches *SavedSearchConnection `json:"urlRedirectSavedSearches,omitempty"`
	// A list of redirects for a shop.
	// URLRedirects *URLRedirectConnection `json:"urlRedirects,omitempty"`
	// Returns a webhook subscription by ID.
	WebhookSubscription *WebhookSubscription `json:"webhookSubscription,omitempty"`
	// List of webhook subscriptions.
	WebhookSubscriptions *WebhookSubscriptionConnection `json:"webhookSubscriptions,omitempty"`
}
