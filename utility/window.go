package utility

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
	w.SetBackingType(appkit.BackingStoreBuffered)
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

func ModalAlert(w appkit.IWindow, multi bool, title, desc string, handler ...func(appkit.ModalResponse)) {
	h := func(code appkit.ModalResponse) {}
	if len(handler) > 0 {
		h = handler[0]
	}
	dialog := appkit.NewAlert()
	dialog.SetAlertStyle(appkit.AlertStyleCritical)
	dialog.SetMessageText(title)
	dialog.SetInformativeText(desc)
	if multi {
		dialog.AddButtonWithTitle("OK")
		dialog.AddButtonWithTitle("Cancel")
	}
	if w == nil || w.IsNil() {
		h(dialog.RunModal())
		return
	}
	dialog.BeginSheetModalForWindowCompletionHandler(w, h)
}
