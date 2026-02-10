package auth

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"go.temporal.io/sdk/workflow"
)

// LoginWorkflow orchestrates the login process
func LoginWorkflow(ctx workflow.Context, email string) (string, error) {
	// 1. Setup ActivityOptions
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 1 * time.Minute,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	// 2. Deterministic Randomness (Crucial)
	var code string
	err := workflow.SideEffect(ctx, func(ctx workflow.Context) interface{} {
		// Generate a random 6-digit code
		// Note: In a real production scenario, might want crypto/rand or a more robust generation
		// but for SideEffect, standard math/rand seeded with time is common pattern as long as it's inside the function
		rng := rand.New(rand.NewSource(time.Now().UnixNano()))
		return fmt.Sprintf("%06d", rng.Intn(1000000))
	}).Get(&code)
	if err != nil {
		return "", err
	}

	// 3. Send Code
	err = workflow.ExecuteActivity(ctx, SendMagicCode, email, code).Get(ctx, nil)
	if err != nil {
		return "", err
	}

	// 4. Wait for User Input (The Signal)
	var userCode string
	signalChan := workflow.GetSignalChannel(ctx, "SubmitCode")

	selector := workflow.NewSelector(ctx)

	// Case 1: Signal Received
	selector.AddReceive(signalChan, func(c workflow.ReceiveChannel, more bool) {
		c.Receive(ctx, &userCode)
	})

	// Case 2: Timeout
	timer := workflow.NewTimer(ctx, 10*time.Minute)
	selector.AddFuture(timer, func(f workflow.Future) {
		// Timer fired
	})

	// 5. Verification Logic
	selector.Select(ctx)

	if userCode == "" {
		return "", errors.New("Login timed out")
	}

	if userCode != code {
		return "", errors.New("Invalid code provided")
	}

	// Match - Generate Token
	// Mock user for now
	user := User{
		ID:    "test-user-id",
		Email: email,
		Role:  "admin",
	}

	var token string
	err = workflow.ExecuteActivity(ctx, GenerateToken, user).Get(ctx, &token)
	if err != nil {
		return "", err
	}

	return token, nil
}
