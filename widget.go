package mapper

import (
	"errors"
	"fmt"
	"time"

	dgw "github.com/Necroforger/dgwidgets"
	dg "github.com/bwmarrin/discordgo"
)

var (
	// ErrIndexOutOfBounds indicates that an index did not correspond to a page.
	ErrIndexOutOfBounds = errors.New("index out of bounds")

	// ErrNilMessage indicates that there was no message to update.
	ErrNilMessage = errors.New("nil message")
)

const duration = 2*time.Minute + 30*time.Second

// Widget wraps the state of an embed pager.
type Widget struct {
	Pages []*dg.MessageEmbed
	Index int

	Timer  *time.Timer
	Widget *dgw.Widget
	Ses    *dg.Session
}

// NewWidget creates a Widget with sensible defaults.
func NewWidget(ses *dg.Session, channelID string, userID string) *Widget {
	w := &Widget{}

	w.Ses = ses
	w.Pages = make([]*dg.MessageEmbed, 0)

	w.Widget = dgw.NewWidget(ses, channelID, nil)
	w.Widget.UserWhitelist = []string{userID}

	return w
}

// Spawn adds handlers for the Widget.
func (w *Widget) Spawn() {
	_f := "(*Widget).Spawn"

	defer w.Close(nil, nil)

	err := w.Widget.Handle("\u25C0", w.PreviousPage)
	if err != nil {
		err = fmt.Errorf("handle \u25C0: %w", err)
		Log.Warn(_f, err)

		return
	}

	err = w.Widget.Handle("\u25B6", w.NextPage)
	if err != nil {
		err = fmt.Errorf("handle \u25B6: %w", err)
		Log.Warn(_f, err)

		return
	}

	err = w.Widget.Handle("\u2705", w.Close)
	if err != nil {
		err = fmt.Errorf("handle \u2705: %w", err)
		Log.Warn(_f, err)

		return
	}

	page, err := w.Page()
	if err != nil {
		err = fmt.Errorf("page: %w", err)
		Log.Warn(_f, err)

		return
	}

	w.Widget.Embed = page
	w.Timer = time.NewTimer(duration)

	go w.Expire()

	err = w.Widget.Spawn()
	if err != nil {
		err = fmt.Errorf("widget spawn: %w", err)
		Log.Warn(_f, err)
	}
}

// Expire closes the Widget after its Timer returns.
func (w *Widget) Expire() {
	<-w.Timer.C
	w.Close(nil, nil)
}

// Add adds embed pages to the Widget.
func (w *Widget) Add(embeds ...*dg.MessageEmbed) {
	w.Pages = append(w.Pages, embeds...)
}

// Page returns the Widget's current embed page.
func (w *Widget) Page() (*dg.MessageEmbed, error) {
	if w.Index < 0 || w.Index >= len(w.Pages) {
		return nil, ErrIndexOutOfBounds
	}

	return w.Pages[w.Index], nil
}

// NextPage is a handler for the right arrow, which advances the Widget by 1 page.
func (w *Widget) NextPage(_ *dgw.Widget, r *dg.MessageReaction) {
	_f := "(*Widget).NextPage"

	if w.Index+1 >= 0 && w.Index+1 < len(w.Pages) {
		w.Index++
	} else {
		w.Index = 0
	}

	err := w.Update()
	if err != nil {
		err = fmt.Errorf("update: %w", err)
		Log.Warn(_f, err)

		return
	}
}

// PreviousPage is a handler for the left arrow, which retracts the Widget by 1 page.
func (w *Widget) PreviousPage(_ *dgw.Widget, r *dg.MessageReaction) {
	_f := "(*Widget).PreviousPage"

	if w.Index-1 >= 0 && w.Index-1 < len(w.Pages) {
		w.Index--
	} else {
		w.Index = len(w.Pages) - 1
	}

	err := w.Update()
	if err != nil {
		err = fmt.Errorf("update: %w", err)
		Log.Warn(_f, err)

		return
	}
}

// Close if a handler for the check mark, which shuts down the Widget.
func (w *Widget) Close(_ *dgw.Widget, r *dg.MessageReaction) {
	_f := "(*Widget).Close"

	page, err := w.Page()
	if err != nil {
		err = fmt.Errorf("page: %w", err)
		Log.Warn(_f, err)

		return
	}

	page.Color = 0x77B255

	err = w.Update()
	if err != nil {
		err = fmt.Errorf("update: %w", err)
		Log.Warn(_f, err)

		return
	}

	err = w.Ses.MessageReactionsRemoveAll(w.Widget.ChannelID, w.Widget.Message.ID)
	if err != nil {
		err = fmt.Errorf("remove reacts %#v %#v: %w", w.Widget.ChannelID, w.Widget.Message.ID, err)
		Log.Warn(_f, err)

		return
	}

	w.Widget.Close <- true
}

// Update ensures that the current page is displayed correctly.
func (w *Widget) Update() error {
	if w.Widget.Message == nil {
		return ErrNilMessage
	}

	w.Timer.Reset(duration)

	page, err := w.Page()
	if err != nil {
		return err
	}

	_, err = w.Widget.UpdateEmbed(page)
	if err != nil {
		return err
	}

	return nil
}
