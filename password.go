package main

import (
	"github.com/progrium/macdriver/helper/action"
	"github.com/progrium/macdriver/helper/layout"
	"github.com/progrium/macdriver/helper/widgets"
	"github.com/progrium/macdriver/macos/appkit"
	"github.com/progrium/macdriver/macos/foundation"
	"github.com/progrium/macdriver/objc"
)

type generatePwdItem struct {
	appkit.MenuItem
	w appkit.Window
}

func getGeneratePasswordItem(menu appkit.StatusItem) *generatePwdItem {
	item := appkit.NewMenuItem()
	item.SetImage(getSymbolImage("key.fill"))
	g := &generatePwdItem{MenuItem: item}
	item.SetTitle("Generate Password")
	item.SetAllowsKeyEquivalentWhenHidden(false)
	item.SetKeyEquivalentModifierMask(appkit.EventModifierFlagCommand)
	item.SetKeyEquivalent("p")
	target, selecor := action.Wrap(g.showGenerateWindow())
	item.SetAction(selecor)
	item.SetTarget(target)
	return g
}

func (g *generatePwdItem) showGenerateWindow() action.Handler {
	mask := appkit.ClosableWindowMask |
		appkit.TitledWindowMask |
		appkit.WindowStyleMaskFullSizeContentView |
		appkit.WindowStyleMaskUnifiedTitleAndToolbar
	return func(objc.Object) {
		if !g.w.IsNil() {
			g.w.OrderFront(nil)
			return
		}

		const width = 330
		popup := appkit.NewPopUpButton()
		popup.SetControlSize(appkit.ControlSizeSmall)
		popup.AddItemWithTitle("手动")
		popup.Menu().AddItem(appkit.MenuItem_SeparatorItem())
		popup.AddItemsWithTitles([]string{
			"字母与数字",
			"仅数字",
			"随机",
			"符合FIPS-181",
		})
		popup.SelectItemAtIndex(4)

		combox := appkit.NewComboBox()
		combox.SetBezelStyle(appkit.TextFieldSquareBezel)
		combox.SetControlSize(appkit.ControlSizeSmall)

		slider := appkit.NewSlider()
		slider.SetContinuous(true)
		slider.SetNumberOfTickMarks(31 - 8)
		slider.SetTickMarkPosition(appkit.TickMarkBelow)
		slider.SetControlSize(appkit.ControlSizeSmall)
		slider.SetAllowsTickMarkValuesOnly(true)
		slider.SetMinValue(8)
		slider.SetMaxValue(31)
		slider.SetDoubleValue(12)

		form := widgets.NewFormView()
		form.AddRow("类型:", popup)
		form.AddRow("建议:", combox)
		form.AddRow("长度:", slider)
		form.SetTranslatesAutoresizingMaskIntoConstraints(false)
		form.SetLabelFont(appkit.Font_LabelFontOfSize(12))
		form.SetLabelControlSpacing(10)

		view := appkit.NewView()
		view.AddSubview(form)
		layout.SetWidth(form, width-30)
		layout.AliginCenterX(form, view)
		form.TopAnchor().ConstraintEqualToAnchorConstant(view.TopAnchor(), 38).SetActive(true)
		controller := appkit.NewViewController()
		controller.SetView(view)
		delegate := &appkit.WindowDelegate{}
		delegate.SetWindowWillClose(func(notification foundation.Notification) {
			g.w = appkit.Window{}
		})
		g.w = appkit.Window_WindowWithContentViewController(controller)
		g.w.Center()
		g.w.SetDelegate(delegate)
		g.w.SetTitle("Generate Password")
		g.w.SetTitlebarAppearsTransparent(true)
		g.w.SetContentSize(sizeOf(width, 150))
		g.w.SetStyleMask(mask)
		g.w.SetLevel(appkit.MainMenuWindowLevel)
		g.w.MakeKeyAndOrderFront(nil)
	}
}
