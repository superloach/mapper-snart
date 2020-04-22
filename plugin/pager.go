package main

import (
	"errors"
	"fmt"

	dgw "github.com/Necroforger/dgwidgets"
	dg "github.com/bwmarrin/discordgo"
)

var (
	ErrIndexOutOfBounds = errors.New("index out of bounds")
	ErrNilMessage       = errors.New("nil message")
)

type Pager struct {
	Pages []*dg.MessageEmbed
	Index int

	Widget *dgw.Widget

	Ses *dg.Session
}

func NewPager(ses *dg.Session, channelID string, userID string) *Pager {
	p := &Pager{}

	p.Ses = ses
	p.Pages = make([]*dg.MessageEmbed, 0)

	w := dgw.NewWidget(ses, channelID, nil)
	w.UserWhitelist = []string{userID}
	p.Widget = w

	return p
}

func (p *Pager) Spawn() error {
	_f := "(*Pager).Spawn"

	defer p.Close(nil, nil)

	p.Widget.Handle("\u25C0", p.PreviousPage)
	p.Widget.Handle("\u25B6", p.NextPage)
	p.Widget.Handle("\u2705", p.Close)

	page, err := p.Page()
	if err != nil {
		err = fmt.Errorf("page: %w", err)
		Log.Error(_f, err)
		return err
	}
	p.Widget.Embed = page

	return p.Widget.Spawn()
}

func (p *Pager) Add(embeds ...*dg.MessageEmbed) {
	p.Pages = append(p.Pages, embeds...)
}

func (p *Pager) Page() (*dg.MessageEmbed, error) {
	if p.Index < 0 || p.Index >= len(p.Pages) {
		return nil, ErrIndexOutOfBounds
	}

	return p.Pages[p.Index], nil
}

func (p *Pager) NextPage(w *dgw.Widget, r *dg.MessageReaction) {
	_f := "(*Pager).NextPage"

	if p.Index+1 >= 0 && p.Index+1 < len(p.Pages) {
		p.Index++
	} else {
		p.Index = 0
	}

	err := p.Update()
	if err != nil {
		err = fmt.Errorf("update: %w", err)
		Log.Error(_f, err)
		return
	}
}

func (p *Pager) PreviousPage(w *dgw.Widget, r *dg.MessageReaction) {
	_f := "(*Pager).PreviousPage"

	if p.Index-1 >= 0 && p.Index-1 < len(p.Pages) {
		p.Index--
	} else {
		p.Index = len(p.Pages) - 1
	}

	err := p.Update()
	if err != nil {
		err = fmt.Errorf("update: %w", err)
		Log.Error(_f, err)
		return
	}
}

func (p *Pager) Close(w *dgw.Widget, r *dg.MessageReaction) {
	_f := "(*Pager).Close"

	page, err := p.Page()
	if err != nil {
		err = fmt.Errorf("page: %w", err)
		Log.Error(_f, err)
		return
	}

	page.Color = 0x77B255

	err = p.Update()
	if err != nil {
		err = fmt.Errorf("update: %w", err)
		Log.Error(_f, err)
		return
	}

	err = p.Ses.MessageReactionsRemoveAll(p.Widget.ChannelID, p.Widget.Message.ID)
	if err != nil {
		err = fmt.Errorf("remove reacts %#v %#v: %w", p.Widget.ChannelID, p.Widget.Message.ID, err)
		Log.Error(_f, err)
		return
	}
}

func (p *Pager) Update() error {
	if p.Widget.Message == nil {
		return ErrNilMessage
	}

	page, err := p.Page()
	if err != nil {
		return err
	}

	_, err = p.Widget.UpdateEmbed(page)
	return err
}
