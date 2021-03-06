package app

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Juli3nnicolas/http_log_monitor/pkg/logger"
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
	alert     *task.AlertState
}

type ViewFrame struct {
	Hits  []task.Hit
	Rates task.Rates
	Codes map[uint32]uint64
	Alert task.AlertState
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

	updateReqPerSeconds := createUpdateReqPerSeconds()
	for {
		view := <-viewChan

		if err := updateHit(w, view.Hits); err != nil {
			errorHandle(err)
		}

		if err := updateRates(w, &view.Rates); err != nil {
			errorHandle(err)
		}

		if err := updateReqPerSeconds(w, &view.Rates); err != nil {
			errorHandle(err)
		}

		if err := updateCodes(w, view.Codes); err != nil {
			errorHandle(err)
		}

		if err := r.updateAlerts(&view.Alert); err != nil {
			errorHandle(err)
		}
	}
}

func updateHit(w *widgets, hits []task.Hit) error {
	var msg string

	for i := range hits {

		j, lenMethods := 0, len(hits[i].Methods)-1
		msg += hits[i].Section + ": " + strconv.Itoa(int(hits[i].Total)) + " ("
		for method, count := range hits[i].Methods {
			sep := ", "
			if j == lenMethods {
				sep = ""
			}

			msg += method + ": " + strconv.Itoa(int(count)) + sep
		}
		msg += ")\n"
	}

	if msg == "" {
		msg = mostHitsNoTraffic
	}

	return updateTextWidget(w.mostHits, msg)
}

func updateRates(w *widgets, r *task.Rates) error {
	f := &r.Frame
	g := &r.Global

	msg := formatRateMsg(rateMsgContent{
		frameDuration: f.Duration,
		maxReqPSec:    g.MaxReqPerS,
		avgReqPSec:    g.AvgReqPerS,
		nbSuccesses:   f.NbSuccess,
		nbFailures:    f.NbFailures,
	})

	return updateTextWidget(w.ratesMsg, msg)
}

func createUpdateReqPerSeconds() func(w *widgets, r *task.Rates) error {
	lastreq := make([]int, 12)
	update := func(w *widgets, r *task.Rates) error {
		lenLastReqs := len(lastreq)

		reqs := make([]int, lenLastReqs)
		copy(reqs[1:], lastreq[:lenLastReqs-1])
		reqs[0] = int(r.Frame.ReqPerS)

		lastreq = reqs
		maxReqPerS := int(r.Global.MaxReqPerS)

		// Values cannot equal to 0
		if maxReqPerS == 0 {
			maxReqPerS = 1
		}
		return w.reqPerSec.Values(lastreq[:], maxReqPerS)
	}

	return update
}

func updateCodes(w *widgets, codes map[uint32]uint64) error {
	var (
		msg100 string = "100:\n"
		msg200 string = "200:\n"
		msg300 string = "300:\n"
		msg400 string = "400:\n"
		msg500 string = "500:\n"
	)

	for code, count := range codes {
		switch {
		case code < 200:
			msg100 += httpReturnCodeLine(code, count)

		case code < 300:
			msg200 += httpReturnCodeLine(code, count)

		case code < 400:
			msg300 += httpReturnCodeLine(code, count)

		case code < 500:
			msg400 += httpReturnCodeLine(code, count)

		case code < 600:
			msg500 += httpReturnCodeLine(code, count)

		default:
			msg500 += ""
		}
	}

	if err := updateTextWidget(w.httpCodes100, msg100); err != nil {
		return err
	}
	if err := updateTextWidget(w.httpCodes200, msg200); err != nil {
		return err
	}
	if err := updateTextWidget(w.httpCodes300, msg300); err != nil {
		return err
	}
	if err := updateTextWidget(w.httpCodes400, msg400); err != nil {
		return err
	}
	return updateTextWidget(w.httpCodes500, msg500)
}

func (r *renderer) updateAlerts(alert *task.AlertState) error {
	// The alarm is has been activated so log it
	if alert.IsOn {
		r.alert = alert
		msg := formatAlertOnMsg(r.alert)
		return updateTextWidget(r.widgets.alertMessage, msg)
	}

	// The alarm has just been desactivated, log it
	if r.alert != nil && r.alert.IsOn && !alert.IsOn {
		r.alert = alert
		msg := formatAlertOffMsg(r.alert)
		return updateTextWidget(r.widgets.alertMessage, msg)
	}

	if r.alert == nil {
		msg := formatAlertInfoMsg(alert)
		return updateTextWidget(r.widgets.alertMessage, msg)
	}

	return nil
}

func httpReturnCodeLine(code uint32, count uint64) string {
	return fmt.Sprintf("%d: %d\n", code, count)
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
func LogUpdateError() func(error) {
	errLog := func(err error) { logger.Get().Errorf(err.Error()) }
	return errLog
}
