package main

import (
	"github.com/PrathmeshAdhav2006/letsGO/internal/models"
)

// Define a templateData type to act as the holding structure for
// any dynamic data that we want to pass to our HTML templates.

type templateData struct {
	Snippet  *models.Snippet
	Snippets []*models.Snippet
}
