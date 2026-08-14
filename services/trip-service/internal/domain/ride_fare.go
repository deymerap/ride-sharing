package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type RideFareModel struct {
	ID              primitive.ObjectID
	UserID          string
	PackageSlug     string  // e.g. "standard", "premium", "sedan", etc.
	TotalPriceCents float64 // Total price in cents to avoid floating point issues. For example, $10.50 would be represented as 1050.
}
