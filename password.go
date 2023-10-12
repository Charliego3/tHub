package main

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/progrium/macdriver/dispatch"
	"github.com/progrium/macdriver/helper/action"
	"github.com/progrium/macdriver/helper/layout"
	"github.com/progrium/macdriver/helper/widgets"
	"github.com/progrium/macdriver/macos/appkit"
	"github.com/progrium/macdriver/macos/foundation"
	"github.com/progrium/macdriver/objc"
)

type PasswordType string

func (t PasswordType) addHistory2Combox(g *generatePwdItem) {
	v := g.combox.StringValue()
	if len(v) == 0 {
		return
	}

	g.combox.InsertItemWithObjectValueAtIndex(foundation.String_StringWithString(v), 0)
	if g.combox.NumberOfItems() > 10 {
		g.combox.RemoveItemAtIndex(10)
	}
}

const (
	PasswordManual          PasswordType = "Manual"
	PasswordLetterAndNumber PasswordType = "Letter And Number"
	PasswordNumber          PasswordType = "Numbers"
	PasswordRandom          PasswordType = "Random"
	PasswordFips181         PasswordType = "Compliant with FIPS-181"
)

type generatePwdItem struct {
	appkit.MenuItem
	w      appkit.Window
	popup  appkit.PopUpButton
	combox appkit.ComboBox
	regbtn appkit.Button
	slider appkit.Slider
	length int
	r      *rand.Rand
}

func getGeneratePasswordItem(menu appkit.StatusItem) *generatePwdItem {
	item := appkit.NewMenuItem()
	item.SetImage(getSymbolImage("key.fill"))
	g := &generatePwdItem{MenuItem: item, length: 12}
	item.SetTitle("Generate Password")
	item.SetAllowsKeyEquivalentWhenHidden(false)
	item.SetKeyEquivalentModifierMask(appkit.EventModifierFlagCommand)
	item.SetKeyEquivalent("p")
	target, selecor := action.Wrap(g.showGenerateWindow)
	item.SetAction(selecor)
	item.SetTarget(target)
	return g
}

func (g *generatePwdItem) showGenerateWindow(_ objc.Object) {
	if !g.w.IsNil() {
		return
	}

	const width = 350
	g.popup = appkit.NewPopUpButton()
	g.popup.SetControlSize(appkit.ControlSizeSmall)
	g.popup.AddItemWithTitle(string(PasswordManual))
	g.popup.Menu().AddItem(appkit.MenuItem_SeparatorItem())
	g.popup.AddItemsWithTitles([]string{
		string(PasswordLetterAndNumber),
		string(PasswordNumber),
		string(PasswordRandom),
		string(PasswordFips181),
	})
	g.popup.SelectItemAtIndex(4)
	genTarget, genSelector := action.Wrap(g.gen)
	g.popup.SetTarget(genTarget)
	g.popup.SetAction(genSelector)

	g.combox = appkit.NewComboBox()
	g.combox.SetBezelStyle(appkit.TextFieldSquareBezel)
	g.combox.SetControlSize(appkit.ControlSizeSmall)
	g.combox.SetEditable(false)
	g.combox.SetSelectable(true)

	g.slider = appkit.NewSlider()
	g.slider.SetContinuous(true)
	g.slider.SetNumberOfTickMarks(31 - 8)
	g.slider.SetTickMarkPosition(appkit.TickMarkBelow)
	g.slider.SetControlSize(appkit.ControlSizeSmall)
	g.slider.SetAllowsTickMarkValuesOnly(true)
	g.slider.SetMinValue(8)
	g.slider.SetMaxValue(31)
	g.slider.SetIntValue(g.length)
	text := appkit.NewLabel(strconv.Itoa(g.slider.IntValue()))
	sliderTarget, sliderSelector := action.Wrap(func(sender objc.Object) {
		dispatch.MainQueue().DispatchAsync(func() {
			g.length = g.slider.IntValue()
			text.SetStringValue(strconv.Itoa(g.length))
			g.gen(sender)
		})
	})

	g.slider.SetTarget(sliderTarget)
	g.slider.SetAction(sliderSelector)
	sliderView := appkit.GridView_GridViewWithNumberOfColumnsRows(2, 0)
	sliderView.SetTranslatesAutoresizingMaskIntoConstraints(false)
	sliderView.SetContentHuggingPriorityForOrientation(750.0, appkit.LayoutConstraintOrientationHorizontal)
	sliderView.SetContentHuggingPriorityForOrientation(750.0, appkit.LayoutConstraintOrientationVertical)
	text.SetContentCompressionResistancePriorityForOrientation(appkit.LayoutPriorityRequired, appkit.LayoutConstraintOrientationVertical)
	sliderView.AddRowWithViews([]appkit.IView{
		g.slider,
		text,
	})
	sliderView.ColumnAtIndex(1).SetXPlacement(appkit.GridCellPlacementTrailing)
	sliderView.ColumnAtIndex(1).SetWidth(16)

	form := widgets.NewFormView()
	form.AddRow("Type:", g.popup)
	form.AddRow("Suggest:", g.combox)
	form.GridView.AddRowWithViews([]appkit.IView{
		appkit.NewLabel("Length:"),
		sliderView,
	})
	form.SetTranslatesAutoresizingMaskIntoConstraints(false)
	form.SetLabelFont(appkit.Font_LabelFontOfSize(12))
	form.SetLabelControlSpacing(8)
	form.SetLabelWidth(50)
	form.SetLabelAlignment(widgets.LabelAlignmentTrailing)
	form.GridView.SetRowSpacing(10)

	okbtn := appkit.NewButtonWithTitle("OK")
	okbtn.SetBezelStyle(appkit.BezelStyleRounded)
	okbtn.SetBezelColor(appkit.Color_SystemBlueColor())
	okbtn.SetTranslatesAutoresizingMaskIntoConstraints(false)
	okTarget, okSelector := action.Wrap(func(sender objc.Object) {
		v := g.combox.StringValue()
		if len(v) > 0 {
			pasteboard := appkit.Pasteboard_GeneralPasteboard()
			pasteboard.ClearContents()
			pasteboard.SetStringForType(v, appkit.PasteboardTypeString)
		}
		g.w.Close()
	})
	okbtn.SetTarget(okTarget)
	okbtn.SetAction(okSelector)

	g.regbtn = appkit.NewButtonWithImage(getSymbolImage("arrow.triangle.2.circlepath"))
	g.regbtn.SetTranslatesAutoresizingMaskIntoConstraints(false)
	g.regbtn.SetBezelStyle(appkit.BezelStyleRounded)
	g.regbtn.SetBezelColor(appkit.Color_SystemGreenColor())
	g.regbtn.SetTarget(genTarget)
	g.regbtn.SetAction(genSelector)

	view := appkit.NewView()
	view.AddSubview(form)
	view.AddSubview(okbtn)
	view.AddSubview(g.regbtn)
	layout.SetWidth(form, width-30)
	layout.AliginCenterX(form, view)
	layout.SetMinWidth(view, width)
	layout.SetMinHeight(view, 145)
	form.TopAnchor().ConstraintEqualToAnchorConstant(view.TopAnchor(), 38).SetActive(true)
	okbtn.TopAnchor().ConstraintEqualToAnchorConstant(form.BottomAnchor(), 10).SetActive(true)
	okbtn.TrailingAnchor().ConstraintEqualToAnchorConstant(view.TrailingAnchor(), -15).SetActive(true)
	g.regbtn.TopAnchor().ConstraintEqualToAnchorConstant(form.BottomAnchor(), 10).SetActive(true)
	g.regbtn.TrailingAnchor().ConstraintEqualToAnchorConstant(okbtn.LeadingAnchor(), -10).SetActive(true)

	g.r = rand.New(rand.NewSource(time.Now().UnixNano()))
	g.gen(objc.Object{})

	controller := appkit.NewViewController()
	controller.SetView(view)
	delegate := &appkit.WindowDelegate{}
	delegate.SetWindowWillClose(func(notification foundation.Notification) {
		g.w = appkit.Window{}
		g.length = 12
	})
	g.w = appkit.Window_WindowWithContentViewController(controller)
	g.w.Center()
	g.w.SetDelegate(delegate)
	g.w.SetTitle("Generate Password")
	g.w.SetTitlebarAppearsTransparent(true)
	g.w.SetStyleMask(appkit.ClosableWindowMask |
		appkit.TitledWindowMask |
		appkit.WindowStyleMaskFullSizeContentView |
		appkit.WindowStyleMaskUnifiedTitleAndToolbar)
	g.w.SetLevel(appkit.MainMenuWindowLevel)
	g.w.MakeKeyAndOrderFront(nil)
}

func (g *generatePwdItem) gen(_ objc.Object) {
	t := PasswordType(g.popup.SelectedItem().Title())
	t.addHistory2Combox(g)
	var gtor generator
	var editable bool
	switch t {
	default:
		gtor = randomGen{}
	case PasswordManual:
		gtor = manualGen{}
		editable = true
	case PasswordLetterAndNumber:
		gtor = letterNumberGen{}
	case PasswordNumber:
		gtor = numberGen{}
	case PasswordFips181:
		gtor = fipsGen{}
	}
	gtor.gen(g)
	g.combox.SetEditable(editable)
	g.regbtn.SetEnabled(!editable)
	g.slider.SetEnabled(!editable)
}

type generator interface {
	gen(g *generatePwdItem)
}

type manualGen struct{}

func (manualGen) gen(g *generatePwdItem) {
	g.combox.SetStringValue("")
}

type letterNumberGen struct{}

func (letterNumberGen) gen(g *generatePwdItem) {
	g.combox.SetStringValue(g.genWithProps(lowerKey, upperKey, numberKey))
}

type numberGen struct{}

func (numberGen) gen(g *generatePwdItem) {
	g.combox.SetStringValue(g.genWithProps(numberKey))
}

type randomGen struct{}

func (randomGen) gen(g *generatePwdItem) {
	g.combox.SetStringValue(g.genWithProps())
}

type fipsGen struct{}

func (fipsGen) gen(g *generatePwdItem) {
	g.combox.SetStringValue(g.genWithProps(lowerKey))
}

type propKey uint

const (
	lowerKey propKey = iota
	upperKey
	numberKey
	symbolKey
)

func (g *generatePwdItem) genWithProps(keys ...propKey) string {
	v := make([]byte, g.length)
	var props []byte
	if len(keys) > 0 {
		for _, key := range keys {
			if p, ok := properties[key]; ok {
				props = append(props, p...)
			}
		}
	} else {
		for _, p := range properties {
			props = append(props, p...)
		}
	}
	for i := 0; i < g.length; i++ {
		idx := rand.Intn(len(props))
		v[i] = props[idx]
	}
	return string(v)
}

var properties = map[propKey][]byte{
	lowerKey:  {'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z'},
	upperKey:  {'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z'},
	numberKey: {'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'},
	symbolKey: {',', '.', '/', '<', '>', '?', ';', '\'', ':', '*', '-', '`', '+', '=', '[', ']', '{', '}', '|', '\\', '_', '(', ')', '&', '^', '%', '$', '#', '@', '!', '~'},
}
