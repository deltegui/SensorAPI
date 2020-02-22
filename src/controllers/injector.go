package controllers

import (
	"github.com/deltegui/locomotive"
)

func Register() {
	locomotive.MapRoot(NewErrorController)
}