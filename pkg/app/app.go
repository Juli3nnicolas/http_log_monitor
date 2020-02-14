package app

import "time"

// Run executes the entire application (both frontend and backend)
func Run() error {
	// Init backend
	b := Backend{}
	if err := b.init(); err != nil {
		return err
	}
	defer b.shutdown()

	// Init view
	r := renderer{}
	ctx, err := r.init()
	if err != nil {
		return err
	}
	defer r.shutdown()

	updateChan := make(chan ViewFrame)
	go r.update(updateChan, LogUpdateError(""))
	go b.run(time.Second, updateChan)

	err = r.render(ctx)
	if err != nil {
		return err
	}

	return nil
}
