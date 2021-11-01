package shopify

import (
	"context"
	"fmt"
	"time"

	"github.com/es-hs/go-shopify-graphql/graphql"

	"github.com/sirupsen/logrus"
)

type WebhookService interface {
	NewWebhookSubcription(topic WebhookTopic, input WebhookTopicSubscription) (output WebhookSubscriptionCreatePayload)
	GetAllWebhookSubcription() (output []*WebhookSubscription, err error)
	DeleteWebhook(webhookID string) (output WebhookSubscriptionDeletePayload, err error)
}

type WebhookServiceOp struct {
	client *Client
}

type WebhookSubscriptionCreatePayload struct {
	// The list of errors that occurred from executing the mutation.
	UserErrors []UserErrors `json:"userErrors,omitempty"`
	// The webhook subscription that was created.
	WebhookSubscription WebhookSubscription `json:"webhookSubscription,omitempty"`
}

// Return type for `webhookSubscriptionDelete` mutation.
type WebhookSubscriptionDeletePayload struct {
	// The ID of the deleted webhook subscription.
	DeletedWebhookSubscriptionID *graphql.String `json:"deletedWebhookSubscriptionId,omitempty"`
	// The list of errors that occurred from executing the mutation.
	UserErrors []*UserErrors `json:"userErrors,omitempty"`
}

type WebhookSubscriptionConnection struct {
	// A list of edges.
	Edges []*WebhookSubscriptionEdge `json:"edges,omitempty"`
	// Information to aid in pagination.
	PageInfo *PageInfo `json:"pageInfo,omitempty"`
}

type WebhookSubscriptionEdge struct {
	// A cursor for use in pagination.
	Cursor graphql.String `json:"cursor,omitempty"`
	// The item at the end of WebhookSubscriptionEdge.
	Node *WebhookSubscription `json:"node,omitempty"`
}

// A webhook subscription is a persisted data object created by an app using the REST Admin API or GraphQL Admin API.
// It describes the topic that the app wants to receive, and a destination where Shopify should send webhooks of the specified topic.
// When an event for a given topic occurs, the webhook subscription sends a relevant payload to the destination.
// Learn more about the [webhooks system](https://shopify.dev/tutorials/manage-webhooks).
type WebhookSubscription struct {
	// The destination URI to which the webhook subscription will send a message when an event occurs.
	CallbackURL graphql.String `json:"callbackUrl,omitempty"`
	// The date and time when the webhook subscription was created.
	CreatedAt time.Time `json:"createdAt,omitempty"`
	// The endpoint to which the webhook subscription will send events. todo check this
	Endpoint WebhookSubscriptionEndpoint `json:"endpoint,omitempty"`
	// The format in which the webhook subscription should send the data.
	Format WebhookSubscriptionFormat `json:"format,omitempty"`
	// A globally-unique identifier.
	ID graphql.ID `json:"id,omitempty"`
	// An optional array of top-level resource fields that should be serialized and sent in the webhook message. If null, then all fields will be sent.
	IncludeFields []graphql.String `json:"includeFields,omitempty"`
	// The ID of the corresponding resource in the REST Admin API.
	LegacyResourceID graphql.String `json:"legacyResourceId,omitempty"`
	// The list of namespaces for any metafields that should be included in the webhook subscription.
	MetafieldNamespaces []graphql.String `json:"metafieldNamespaces,omitempty"`
	// The type of event that triggers the webhook. The topic determines when the webhook subscription sends a webhook, as well as what class of data object that webhook contains.
	Topic WebhookSubscriptionTopic `json:"topic,omitempty"`
	// The date and time when the webhook subscription was updated.
	UpdatedAt graphql.String `json:"updatedAt,omitempty"`
}

type WebhookSubscriptionEndpoint struct {
	WebhookHTTPEndpoint        `graphql:"... on WebhookHttpEndpoint"`
	WebhookEventBridgeEndpoint `graphql:"... on WebhookEventBridgeEndpoint"`
}

// Amazon EventBridge event source.
type WebhookEventBridgeEndpoint struct {
	// ARN of this EventBridge event source.
	Arn graphql.String `json:"arn,omitempty"`
}

// HTTP endpoint where POST requests will be made to.
type WebhookHTTPEndpoint struct {
	// URL of webhook endpoint to deliver webhooks to.
	CallbackURL graphql.String `json:"callbackUrl,omitempty"`
}

// Google Cloud Pub/Sub event source.
type WebhookPubSubEndpoint struct {
	// The Google Cloud Pub/Sub project ID.
	PubSubProject graphql.String `json:"pubSubProject,omitempty"`
	// The Google Cloud Pub/Sub topic ID.
	PubSubTopic graphql.String `json:"pubSubTopic,omitempty"`
}

type WebhookSubscriptionFormat graphql.String

type WebhookSubscriptionTopic graphql.String

// Specifies the input fields for a webhook subscription.
type WebhookSubscriptionInput struct {
	// URL where the webhook subscription should send the POST request when the event occurs.
	CallbackURL graphql.String `json:"callbackUrl,omitempty"`
	// The format in which the webhook subscription should send the data.
	Format WebhookSubscriptionFormat `json:"format,omitempty"`
	// The list of fields to be included in the webhook subscription.
	IncludeFields []string `json:"includeFields,omitempty"`
	// The list of namespaces for any metafields that should be included in the webhook subscription.
	MetafieldNamespaces []string `json:"metafieldNamespaces,omitempty"`
}

type mutationWebhookCreate struct {
	WebhookCreateResult WebhookSubscriptionCreatePayload `graphql:"webhookSubscriptionCreate(topic: $topic, webhookSubscription: $webhookSubscription)" json:"webhookSubscriptionCreate"`
}

type mutationWebhookDelete struct {
	WebhookDeleteResult WebhookSubscriptionDeletePayload `graphql:"webhookSubscriptionDelete(id: $id)" json:"webhookSubscriptionCreate"`
}

type WebhookTopic struct {
	WebhookSubscriptionTopic WebhookSubscriptionTopic
}

type WebhookTopicSubscription struct {
	WebhookSubscriptionInput WebhookSubscriptionInput
}

func (w WebhookServiceOp) NewWebhookSubcription(topic WebhookTopic, input WebhookTopicSubscription) (output WebhookSubscriptionCreatePayload) {
	fmt.Println("hahahahaha")
	m := mutationWebhookCreate{}
	vars := map[string]interface{}{
		"topic":               topic.WebhookSubscriptionTopic,
		"webhookSubscription": input.WebhookSubscriptionInput,
	}
	err := w.client.gql.Mutate(context.Background(), &m, vars)
	if err != nil {
		logrus.Info("nani")
		logrus.Info(err)
		return m.WebhookCreateResult
	}

	if len(m.WebhookCreateResult.UserErrors) > 0 {
		err = fmt.Errorf("%+v", m.WebhookCreateResult.UserErrors)
		logrus.Info(err)
		return m.WebhookCreateResult
	}

	return m.WebhookCreateResult
}

func (w WebhookServiceOp) DeleteWebhook(webhookID string) (output WebhookSubscriptionDeletePayload, err error) {
	fmt.Println("hahahahaha")
	m := mutationWebhookDelete{}
	logrus.Info(webhookID)
	vars := map[string]interface{}{
		"id": webhookID,
	}
	err = w.client.gql.Mutate(context.Background(), &m, vars)
	if err != nil {
		logrus.Info("nani")
		logrus.Info(err)
		return m.WebhookDeleteResult, err
	}

	if len(m.WebhookDeleteResult.UserErrors) > 0 {
		err = fmt.Errorf("%+v", m.WebhookDeleteResult.UserErrors)
		logrus.Info(err)
	}
	return m.WebhookDeleteResult, err
}

func (w WebhookServiceOp) GetAllWebhookSubcription() (output []*WebhookSubscription, err error) {
	fmt.Println("hahahahaha")
	query := fmt.Sprintf(`{
    webhookSubscriptions(first: 20) {
      edges {
        node {
          id,
          topic,
          endpoint {
            __typename
            ... on WebhookHttpEndpoint {
              callbackUrl
            }
            ... on WebhookEventBridgeEndpoint{
              arn
            }
          }
          callbackUrl
          format
          topic
          includeFields
          createdAt
          updatedAt
        }
      }
    }
  }
  `)
	// q := fmt.Sprintf(`
	// 	{
	// 		products (first: 10, reverse: true) {
	// 			edges{
	// 				node{
	// 					%s
	// 				}
	// 			}
	// 		}
	// 	}
	// `, productBaseQuery)

	vars := map[string]interface{}{
		// "first": 10,
	}

	var op []*WebhookSubscription
	var out QueryRoot
	fmt.Println("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	err = w.client.gql.QueryString(context.Background(), query, vars, &out)
	fmt.Println(op)
	if err != nil {
		return output, err
	}
	for _, wh := range out.WebhookSubscriptions.Edges {
		op = append(op, wh.Node)
	}
	return op, nil
}
