package app

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Juli3nnicolas/http_log_monitor/pkg/task"
	"github.com/mum4k/termdash"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/keyboard"
	"github.com/mum4k/termdash/terminal/termbox"
	"github.com/mum4k/termdash/terminal/terminalapi"
	"github.com/mum4k/termdash/widgets/text"
)

type renderer struct {
	console   *termbox.Terminal
	widgets   *widgets
	container *container.Container
	cancel    context.CancelFunc
	gridOpts  []container.Option
}

type ViewFrame struct {
	Hits  map[string]task.Hit
	Rates task.Rates
	Codes map[uint32]uint64
}

// rootID is the ID assigned to the root container.
const rootID = "root"

func (r *renderer) init() (context.Context, error) {
	t, err := termbox.New(termbox.ColorMode(terminalapi.ColorMode256))
	if err != nil {
		return nil, err
	}
	r.console = t

	r.container, err = container.New(t, container.ID(rootID))
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	r.cancel = cancel

	w, err := newWidgets(ctx, r.container)
	if err != nil {
		return nil, err
	}
	r.widgets = w

	r.gridOpts, err = gridLayout(w)
	if err != nil {
		return nil, err
	}

	return ctx, nil
}

func (r *renderer) shutdown() {
	r.console.Close()
}

// frontend runs the TUI
func (r *renderer) render(ctx context.Context) error {
	if err := r.container.Update(rootID, r.gridOpts...); err != nil {
		return err
	}

	quitter := func(k *terminalapi.Keyboard) {
		if k.Key == keyboard.KeyEsc || k.Key == keyboard.KeyCtrlC {
			r.cancel()
		}
	}
	if err := termdash.Run(ctx, r.console, r.container, termdash.KeyboardSubscriber(quitter), termdash.RedrawInterval(redrawInterval)); err != nil {
		return err
	}

	return nil
}

// This function reads from a channel to update the UI.
// It must be run asynchronously
// Parameters :
// viewChan chan ViewFrame: read to update the view
// errorHandle : called whenever an error occur
func (r *renderer) update(viewChan chan ViewFrame, errorHandle func(error)) {
	w := r.widgets
	if w == nil {
		errorHandle(fmt.Errorf("nil widget ptr"))
	}

	var lastReqPerSec [5]uint64
	for {
		view := <-viewChan

		if err := updateHit(w, view.Hits); err != nil {
			errorHandle(err)
		}

		if err := updateRates(w, &view.Rates); err != nil {
			errorHandle(err)
		}

		if err := updateReqPerSeconds(w, lastReqPerSec); err != nil {
			errorHandle(err)
		}

		if err := updateCodes(w, view.Codes); err != nil {
			errorHandle(err)
		}
	}
}

func updateHit(w *widgets, hits map[string]task.Hit) error {
	var msg string

	// Limit to the 10 firsts, order them by decreasing order
	for section, hit := range hits {

		msg += section + ": " + strconv.Itoa(int(hit.Total)) + " ("
		for method, count := range hit.Methods {
			msg += method + ": " + strconv.Itoa(int(count)) + ", "
		}
		msg += ")\n"
	}

	return updateTextWidget(w.mostHits, msg)
}

func updateRates(w *widgets, r *task.Rates) error {
	f := &r.Frame
	g := &r.Global
	msg := fmt.Sprintf("Frame: %d s Max: %d req/s Avg: %d req/s Success: %d Failure: %d",
		f.Duration, g.MaxReqPerS, g.AvgReqPerS, f.NbSuccess, f.NbFailures)

	return updateTextWidget(w.ratesMsg, msg)
}

func updateReqPerSeconds(w *widgets, lastReqPerSec [5]uint64) error {
	return nil
}

func updateCodes(w *widgets, codes map[uint32]uint64) error {
	return nil
}

func updateTextWidget(w *text.Text, msg string) error {
	const textErr string = "<Update error>"

	// It is an erreor to Write, so prevent it
	if msg == "" {
		msg = textErr
	}

	// Write and check if an unexpected error happened
	if err := w.Write(msg, text.WriteReplace()); err != nil {
		err = w.Write(textErr, text.WriteReplace())
		return err
	}
	return nil
}

// LogUpdateError writes errors to path
func LogUpdateError(path string) func(error) {
	if path == "" {
		path = "/var/log/http_log_monitor/renderupdate.log"
	}

	errLog := func(error) { fmt.Println(path) }
	return errLog
}
