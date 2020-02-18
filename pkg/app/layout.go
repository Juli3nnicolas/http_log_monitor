package app

import (
	"context"
	"time"

	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/container/grid"
	"github.com/mum4k/termdash/linestyle"
	"github.com/mum4k/termdash/widgets/barchart"
	"github.com/mum4k/termdash/widgets/text"
	"github.com/mum4k/termdash/widgets/textinput"
)

// redrawInterval is how often termdash redraws the screen.
const redrawInterval = 250 * time.Millisecond

// widgets holds the widgets used by this demo.
type widgets struct {
	alertMessage *text.Text
	ratesMsg     *text.Text
	mostHits     *text.Text
	httpCodes100 *text.Text
	httpCodes200 *text.Text
	httpCodes300 *text.Text
	httpCodes400 *text.Text
	httpCodes500 *text.Text
	reqPerSec    *barchart.BarChart
}

// newWidgets creates all widgets used by this demo.
func newWidgets(ctx context.Context, c *container.Container) (*widgets, error) {
	reqPerSec, err := newBarChart(ctx)
	if err != nil {
		return nil, err
	}

	alertMessage, err := newTextLabel(alertMessageHeader)
	if err != nil {
		return nil, err
	}

	ratesMsg, err := newTextLabel(formatRateMsg(rateMsgContent{
		frameDuration: 1,
		maxReqPSec:    0,
		avgReqPSec:    0,
		nbSuccesses:   0,
		nbFailures:    0,
	}))
	if err != nil {
		return nil, err
	}

	mostHits, err := newTextLabel(mostHitsNoTraffic)
	if err != nil {
		return nil, err
	}

	httpCodes100, err := newTextLabel(httpCodes100Header)
	if err != nil {
		return nil, err
	}

	httpCodes200, err := newTextLabel(httpCodes200Header)
	if err != nil {
		return nil, err
	}

	httpCodes300, err := newTextLabel(httpCodes300Header)
	if err != nil {
		return nil, err
	}

	httpCodes400, err := newTextLabel(httpCodes400Header)
	if err != nil {
		return nil, err
	}

	httpCodes500, err := newTextLabel(httpCodes500Header)
	if err != nil {
		return nil, err
	}

	return &widgets{
		alertMessage: alertMessage,
		ratesMsg:     ratesMsg,
		mostHits:     mostHits,
		httpCodes100: httpCodes100,
		httpCodes200: httpCodes200,
		httpCodes300: httpCodes300,
		httpCodes400: httpCodes400,
		httpCodes500: httpCodes500,
		reqPerSec:    reqPerSec,
	}, nil
}

// gridLayout prepares container options that represent the desired screen layout.
// This function demonstrates the use of the grid builder.
// gridLayout() and contLayout() demonstrate the two available layout APIs and
// both produce equivalent layouts for layoutType layoutAll.
func gridLayout(w *widgets) ([]container.Option, error) {
	builder := grid.New()
	builder.Add(
		grid.RowHeightPerc(8,
			grid.Widget(w.alertMessage,
				container.Border(linestyle.Light),
				container.BorderTitle("Alert:"),
				container.BorderTitleAlignLeft(),
			),
		),
		grid.RowHeightPerc(92,
			grid.ColWidthPerc(70,
				grid.RowHeightPerc(10,
					grid.Widget(w.ratesMsg,
						container.Border(linestyle.Light),
						container.BorderTitle("Rates"),
						container.BorderTitleAlignLeft(),
					),
				),
				grid.RowHeightPerc(80,
					grid.Widget(w.reqPerSec,
						container.Border(linestyle.Light),
						container.BorderTitle("Req/s"),
						container.BorderTitleAlignLeft(),
					),
				),
			),
			grid.ColWidthPerc(30,
				grid.RowHeightPerc(50,
					grid.Widget(w.mostHits,
						container.Border(linestyle.Light),
						container.BorderTitle("Most hits"),
						container.BorderTitleAlignLeft(),
					),
				),
				// HTTP error codes - add a container
				grid.RowHeightPercWithOpts(50,
					[]container.Option{
						container.Border(linestyle.Light),
						container.BorderTitle("HTTP codes"),
						container.BorderTitleAlignLeft(),
					},
					grid.ColWidthPerc(20,
						grid.Widget(w.httpCodes100,
							container.Border(linestyle.None),
						),
					),
					grid.ColWidthPerc(20,
						grid.Widget(w.httpCodes200,
							container.Border(linestyle.None),
						),
					),
					grid.ColWidthPerc(20,
						grid.Widget(w.httpCodes300,
							container.Border(linestyle.None),
						),
					),
					grid.ColWidthPerc(20,
						grid.Widget(w.httpCodes400,
							container.Border(linestyle.None),
						),
					),
					grid.ColWidthPerc(20,
						grid.Widget(w.httpCodes500,
							container.Border(linestyle.None),
						),
					),
				),
			),
		),
	)

	gridOpts, err := builder.Build()
	if err != nil {
		return nil, err
	}
	return gridOpts, nil
}

// newBarChart returns a BarcChart that displays random values on multiple bars.
func newBarChart(ctx context.Context) (*barchart.BarChart, error) {
	bc, err := barchart.New(
		barchart.BarColors([]cell.Color{
			cell.ColorNumber(33),
			cell.ColorNumber(33),
			cell.ColorNumber(33),
			cell.ColorNumber(33),
			cell.ColorNumber(33),
			cell.ColorNumber(33),
			cell.ColorNumber(33),
			cell.ColorNumber(33),
			cell.ColorNumber(33),
			cell.ColorNumber(33),
			cell.ColorNumber(33),
			cell.ColorNumber(33),
		}),
		barchart.ValueColors([]cell.Color{
			cell.ColorYellow,
			cell.ColorYellow,
			cell.ColorYellow,
			cell.ColorYellow,
			cell.ColorYellow,
			cell.ColorYellow,
			cell.ColorYellow,
			cell.ColorYellow,
			cell.ColorYellow,
			cell.ColorYellow,
			cell.ColorYellow,
			cell.ColorYellow,
		}),
		barchart.ShowValues(),
	)
	if err != nil {
		return nil, err
	}

	return bc, nil
}

// newTextInput creates a new TextInput field that changes the text on the
// SegmentDisplay.
func newTextInput(text, inputPlaceHolder string, onSubmit func(input string) error) (*textinput.TextInput, error) {
	input, err := textinput.New(
		textinput.Label(text, cell.FgColor(cell.ColorWhite)),
		textinput.MaxWidthCells(20),
		textinput.PlaceHolder(inputPlaceHolder),
		textinput.OnSubmit(onSubmit),
	)
	if err != nil {
		return nil, err
	}
	return input, err
}

// Adds a label (raw text)
func newTextLabel(msg string) (*text.Text, error) {
	txt, err := text.New()
	txt.Write(msg, text.WriteCellOpts(cell.FgColor(cell.ColorWhite)))
	if err != nil {
		return nil, err
	}
	return txt, err
}
