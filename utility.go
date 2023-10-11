package main

import (
	"github.com/progrium/macdriver/macos/appkit"
	"github.com/progrium/macdriver/macos/foundation"
)

func getSymbolImage(name string, cfgs ...appkit.IImageSymbolConfiguration) appkit.Image {
	image := appkit.Image_ImageWithSystemSymbolNameAccessibilityDescription(name, name)
	for _, cfg := range cfgs {
		image = image.ImageWithSymbolConfiguration(cfg)
	}
	return image
}

func getImageScale(scale appkit.ImageSymbolScale) appkit.ImageSymbolConfiguration {
	return appkit.ImageSymbolConfiguration_ConfigurationWithScale(scale)
}

func sizeOf(width, height float64) foundation.Size {
	return foundation.Size{Width: width, Height: height}
}
