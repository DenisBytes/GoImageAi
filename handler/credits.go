package handler

import (
	"net/http"
	"os"

	"com.github.denisbytes.goimageai/view/credits"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/checkout/session"
)

func HandleCreditsIndex(w http.ResponseWriter, r *http.Request) error {
	return credits.Index().Render(r.Context(), w)
}

func HandleStripeCheckoutPost(w http.ResponseWriter, r *http.Request) error {
	stripe.Key = os.Getenv("STRIPE_API_KEY")
	checkoutParams := &stripe.CheckoutSessionParams{
		SuccessURL: stripe.String(""),
		CancelURL:  stripe.String(""),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String("___"),
				Quantity: stripe.Int64(1),
			},
		},
	}
	s, err := session.New(checkoutParams)
	if err != nil {
		return err
	}
	http.Redirect(w, r, s.URL, http.StatusSeeOther)
	return nil
}
