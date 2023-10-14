package main

import "github.com/progrium/macdriver/macos/appkit"

type Option func(appkit.Window)

func WithDelegate(delegate appkit.PWindowDelegate) Option {
	return func(w appkit.Window) {
		w.SetDelegate(delegate)
	}
}

func WithStyleMask(mask appkit.WindowStyleMask) Option {
	return func(w appkit.Window) {
		w.SetStyleMask(mask)
	}
}

func NewWindow(title string, contentView appkit.IView, opts ...Option) appkit.Window {
	controller := appkit.NewViewController()
	controller.SetView(contentView)
	w := appkit.Window_WindowWithContentViewController(controller)
	w.Center()
	w.SetTitle(title)
	w.SetTitlebarAppearsTransparent(true)
	w.SetStyleMask(appkit.ClosableWindowMask |
		appkit.TitledWindowMask |
		appkit.WindowStyleMaskFullSizeContentView |
		appkit.WindowStyleMaskUnifiedTitleAndToolbar)
	w.SetLevel(appkit.MainMenuWindowLevel)

	for _, op := range opts {
		op(w)
	}
	return w
}
