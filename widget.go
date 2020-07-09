package mapper

import (
	"errors"
	"fmt"
	"time"

	dgw "github.com/Necroforger/dgwidgets"
	dg "github.com/bwmarrin/discordgo"
)

var (
	ErrIndexOutOfBounds = errors.New("index out of bounds")
	ErrNilMessage       = errors.New("nil message")
)

const duration = 2*time.Minute + 30*time.Second

type Widget struct {
	Pages []*dg.MessageEmbed
	Index int

	Timer  *time.Timer
	Widget *dgw.Widget
	Ses    *dg.Session
}

func NewWidget(ses *dg.Session, channelID string, userID string) *Widget {
	p := &Widget{}

	p.Ses = ses
	p.Pages = make([]*dg.MessageEmbed, 0)

	w := dgw.NewWidget(ses, channelID, nil)
	w.UserWhitelist = []string{userID}
	p.Widget = w

	return p
}

func (p *Widget) Spawn() {
	_f := "(*Widget).Spawn"

	defer p.Close(nil, nil)

	err := p.Widget.Handle("\u25C0", p.PreviousPage)
	if err != nil {
		err = fmt.Errorf("handle \u25C0: %w", err)
		Log.Warn(_f, err)

		return
	}

	err = p.Widget.Handle("\u25B6", p.NextPage)
	if err != nil {
		err = fmt.Errorf("handle \u25B6: %w", err)
		Log.Warn(_f, err)

		return
	}

	err = p.Widget.Handle("\u2705", p.Close)
	if err != nil {
		err = fmt.Errorf("handle \u2705: %w", err)
		Log.Warn(_f, err)

		return
	}

	page, err := p.Page()
	if err != nil {
		err = fmt.Errorf("page: %w", err)
		Log.Warn(_f, err)

		return
	}

	p.Widget.Embed = page
	p.Timer = time.NewTimer(duration)

	go p.Expire()

	err = p.Widget.Spawn()
	if err != nil {
		err = fmt.Errorf("widget spawn: %w", err)
		Log.Warn(_f, err)
	}
}

func (p *Widget) Expire() {
	<-p.Timer.C
	p.Close(nil, nil)
}

func (p *Widget) Add(embeds ...*dg.MessageEmbed) {
	p.Pages = append(p.Pages, embeds...)
}

func (p *Widget) Page() (*dg.MessageEmbed, error) {
	if p.Index < 0 || p.Index >= len(p.Pages) {
		return nil, ErrIndexOutOfBounds
	}

	return p.Pages[p.Index], nil
}

func (p *Widget) NextPage(w *dgw.Widget, r *dg.MessageReaction) {
	_f := "(*Widget).NextPage"

	if p.Index+1 >= 0 && p.Index+1 < len(p.Pages) {
		p.Index++
	} else {
		p.Index = 0
	}

	err := p.Update()
	if err != nil {
		err = fmt.Errorf("update: %w", err)
		Log.Warn(_f, err)

		return
	}
}

func (p *Widget) PreviousPage(w *dgw.Widget, r *dg.MessageReaction) {
	_f := "(*Widget).PreviousPage"

	if p.Index-1 >= 0 && p.Index-1 < len(p.Pages) {
		p.Index--
	} else {
		p.Index = len(p.Pages) - 1
	}

	err := p.Update()
	if err != nil {
		err = fmt.Errorf("update: %w", err)
		Log.Warn(_f, err)

		return
	}
}

func (p *Widget) Close(w *dgw.Widget, r *dg.MessageReaction) {
	_f := "(*Widget).Close"

	page, err := p.Page()
	if err != nil {
		err = fmt.Errorf("page: %w", err)
		Log.Warn(_f, err)

		return
	}

	page.Color = 0x77B255

	err = p.Update()
	if err != nil {
		err = fmt.Errorf("update: %w", err)
		Log.Warn(_f, err)

		return
	}

	err = p.Ses.MessageReactionsRemoveAll(p.Widget.ChannelID, p.Widget.Message.ID)
	if err != nil {
		err = fmt.Errorf("remove reacts %#v %#v: %w", p.Widget.ChannelID, p.Widget.Message.ID, err)
		Log.Warn(_f, err)

		return
	}

	p.Widget.Close <- true
}

func (p *Widget) Update() error {
	if p.Widget.Message == nil {
		return ErrNilMessage
	}

	p.Timer.Reset(duration)

	page, err := p.Page()
	if err != nil {
		return err
	}

	_, err = p.Widget.UpdateEmbed(page)
	if err != nil {
		return err
	}

	return nil
}
