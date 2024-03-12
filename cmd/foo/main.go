package main

import (
	"context"
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/replicate/replicate-go"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	r8, err := replicate.NewClient(replicate.WithTokenFromEnv())
	if err != nil {
		log.Fatal(err)
	}

	// https://replicate.com/stability-ai/stable-diffusion
	version := "ac732df83cea7fff18b8472768c88ad041fa750ff7682a21affe81863cbe77e4"


	input := replicate.PredictionInput{
		"prompt": "an astronaut riding a horse on mars, hd, dramatic lighting",
	}

	webhook := replicate.Webhook{
		URL:    "https://webhook.site/5fe9f245-ef5c-432c-931e-6eeff587066b",
		Events: []replicate.WebhookEventType{"start", "completed"},
	}

	// Run a model and wait for its output
	output, err := r8.CreatePrediction(ctx, version, input, &webhook, false)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("output: ", output)

	// // Create a prediction
	// prediction, err := r8.CreatePrediction(ctx, version, input, &webhook, false)
	// if err != nil {
	// 	// handle error
	// }

}
