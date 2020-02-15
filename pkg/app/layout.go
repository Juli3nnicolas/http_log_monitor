package app

import (
	"context"
	"math/rand"
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
	logFilePath    *text.Text
	alertThreshold *textinput.TextInput
	alertDuration  *textinput.TextInput
	alertMessage   *text.Text
	ratesMsg       *text.Text
	mostHits       *text.Text
	httpCodes100   *text.Text
	httpCodes200   *text.Text
	httpCodes300   *text.Text
	httpCodes400   *text.Text
	httpCodes500   *text.Text
	barChart       *barchart.BarChart
}

// newWidgets creates all widgets used by this demo.
func newWidgets(ctx context.Context, c *container.Container) (*widgets, error) {
	bc, err := newBarChart(ctx)
	if err != nil {
		return nil, err
	}

	logFilePath, err := newTextLabel("Reading /tmp/access.log")
	if err != nil {
		return nil, err
	}

	alertThreshold, err := newTextInput("Threshold (req/s): ", "10", func(string) error { return nil })
	if err != nil {
		return nil, err
	}

	alertDuration, err := newTextInput("Duration (s): ", "10", func(string) error { return nil })
	if err != nil {
		return nil, err
	}

	alertMessage, err := newTextLabel("Message:")
	if err != nil {
		return nil, err
	}

	ratesMsg, err := newTextLabel("Frame: 1 s Max: 30 req/s Avg: 5 req/s Success: 5 Failure: 0")
	if err != nil {
		return nil, err
	}

	mostHits, err := newTextLabel("/instance: 5 req (GET: 3, POST: 2)\n/: 2 req (GET: 2)")
	if err != nil {
		return nil, err
	}

	httpCodes100, err := newTextLabel("100:")
	if err != nil {
		return nil, err
	}

	httpCodes200, err := newTextLabel("200:")
	if err != nil {
		return nil, err
	}

	httpCodes300, err := newTextLabel("300:")
	if err != nil {
		return nil, err
	}

	httpCodes400, err := newTextLabel("400:")
	if err != nil {
		return nil, err
	}

	httpCodes500, err := newTextLabel("500:")
	if err != nil {
		return nil, err
	}

	return &widgets{
		logFilePath:    logFilePath,
		alertThreshold: alertThreshold,
		alertDuration:  alertDuration,
		alertMessage:   alertMessage,
		ratesMsg:       ratesMsg,
		mostHits:       mostHits,
		httpCodes100:   httpCodes100,
		httpCodes200:   httpCodes200,
		httpCodes300:   httpCodes300,
		httpCodes400:   httpCodes400,
		httpCodes500:   httpCodes500,
		barChart:       bc,
	}, nil
}

// gridLayout prepares container options that represent the desired screen layout.
// This function demonstrates the use of the grid builder.
// gridLayout() and contLayout() demonstrate the two available layout APIs and
// both produce equivalent layouts for layoutType layoutAll.
func gridLayout(w *widgets) ([]container.Option, error) {
	builder := grid.New()
	builder.Add(
		grid.RowHeightPerc(5,
			grid.Widget(w.logFilePath,
				container.Border(linestyle.None),
			),
		),
		grid.RowHeightPerc(5,
			grid.ColWidthPerc(25,
				grid.ColWidthPerc(12,
					grid.Widget(w.alertThreshold,
						container.Border(linestyle.None),
					),
				),
				grid.ColWidthPerc(12,
					grid.Widget(w.alertDuration,
						container.Border(linestyle.None),
					),
				),
				grid.ColWidthPerc(75,
					grid.Widget(w.alertMessage,
						container.Border(linestyle.None),
					),
				),
			),
		),
		grid.RowHeightPerc(90,
			grid.ColWidthPerc(70,
				grid.RowHeightPerc(8,
					grid.Widget(w.ratesMsg,
						container.Border(linestyle.Light),
						container.BorderTitle("Rates"),
						container.BorderTitleAlignLeft(),
					),
				),
				grid.RowHeightPerc(92,
					grid.Widget(w.barChart,
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
				grid.RowHeightPerc(50,
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

// periodic executes the provided closure periodically every interval.
// Exits when the context expires.
func periodic(ctx context.Context, interval time.Duration, fn func() error) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if err := fn(); err != nil {
				panic(err)
			}
		case <-ctx.Done():
			return
		}
	}
}

// newBarChart returns a BarcChart that displays random values on multiple bars.
func newBarChart(ctx context.Context) (*barchart.BarChart, error) {
	bc, err := barchart.New(
		barchart.BarColors([]cell.Color{
			cell.ColorNumber(33),
			cell.ColorNumber(39),
			cell.ColorNumber(45),
			cell.ColorNumber(51),
			cell.ColorNumber(81),
			cell.ColorNumber(87),
		}),
		barchart.ValueColors([]cell.Color{
			cell.ColorBlack,
			cell.ColorBlack,
			cell.ColorBlack,
			cell.ColorBlack,
			cell.ColorBlack,
			cell.ColorBlack,
		}),
		barchart.ShowValues(),
	)
	if err != nil {
		return nil, err
	}

	const (
		bars = 6
		max  = 100
	)
	values := make([]int, bars)
	go periodic(ctx, 1*time.Second, func() error {
		for i := range values {
			values[i] = int(rand.Int31n(max + 1))
		}

		return bc.Values(values, max)
	})
	return bc, nil
}

// newTextInput creates a new TextInput field that changes the text on the
// SegmentDisplay.
func newTextInput(text, inputPlaceHolder string, onSubmit func(input string) error) (*textinput.TextInput, error) {
	input, err := textinput.New(
		textinput.Label(text, cell.FgColor(cell.ColorBlue)),
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
	txt.Write(msg, text.WriteCellOpts(cell.FgColor(cell.ColorBlue)))
	if err != nil {
		return nil, err
	}
	return txt, err
}
