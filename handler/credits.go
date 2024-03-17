package handler

import (
	"fmt"
	"net/http"
	"os"

	"com.github.denisbytes.goimageai/db"
	"com.github.denisbytes.goimageai/view/credits"
	"github.com/go-chi/chi/v5"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/checkout/session"
)

func HandleCreditsIndex(w http.ResponseWriter, r *http.Request) error {
	return credits.Index().Render(r.Context(), w)
}

func HandleStripeCheckoutPost(w http.ResponseWriter, r *http.Request) error {
	stripe.Key = os.Getenv("STRIPE_API_KEY")
	checkoutParams := &stripe.CheckoutSessionParams{
		SuccessURL: stripe.String("http://localhost:3000/checkout/success/{CHECKOUT_SESSION_ID}"),
		CancelURL:  stripe.String("http://localhost:3000/checkout/cancel"),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(chi.URLParam(r, "productID")),
				Quantity: stripe.Int64(1),
			},
		},
		Mode: stripe.String(string(stripe.CheckoutSessionModePayment)),
	}
	s, err := session.New(checkoutParams)
	if err != nil {
		return err
	}
	http.Redirect(w, r, s.URL, http.StatusSeeOther)
	return nil
}

func HandleStripeCheckoutSuccess(w http.ResponseWriter, r *http.Request) error {
	user := GetAuthenticatedUser(r)
	sessionID := chi.URLParam(r, "sessionID")
	stripe.Key = os.Getenv("STRIPE_API_KEY")
	sess, err := session.Get(sessionID, nil)
	if err != nil {
		return nil
	}

	lineItemsParams := stripe.CheckoutSessionListLineItemsParams{}
	lineItemsParams.Session = stripe.String(sess.ID)
	iter := session.ListLineItems(&lineItemsParams)
	iter.Next()
	item := iter.LineItem()
	priceID := item.Price.ID

	switch priceID {
	case os.Getenv("100_CREDITS_STRIPE_API"):
		user.Account.Credits += 100
	case os.Getenv("250_CREDITS_STRIPE_API"):
		user.Account.Credits += 250
	case os.Getenv("550_CREDITS_STRIPE_API"):
		user.Account.Credits += 550
	default:
		return fmt.Errorf("invalid price ID: %s", priceID)
	}
	if err := db.UpdateAccount(&user.Account); err != nil {
		return err
	}
	http.Redirect(w, r, "/generate", http.StatusSeeOther)
	return nil
}

func HandleStripeCheckoutCancel(w http.ResponseWriter, r *http.Request) error {
	return nil
}
