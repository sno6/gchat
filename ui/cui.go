package ui

import (
	"github.com/jroimartin/gocui"
	"github.com/sno6/gchat/chat"
)

// CUI is a wrapper around gocui.Gui.
type CUI struct {
	gui *gocui.Gui
}

// A SendCloser is anything that can take the buffer from the user input and do something with it.
// In this case, the sender will be implemented by the client who sends grpc messages over a stream.
type SendCloser interface {
	Send(string) error
	Close() error
}

// NewCUI initialises a new CUI.
func NewCUI(s SendCloser) (*CUI, error) {
	g, err := gocui.NewGui(gocui.Output256)
	if err != nil {
		return nil, err
	}
	g.SetManagerFunc(layout)
	g.Cursor = true
	if err := setKeyBindings(g, s); err != nil {
		return nil, err
	}
	return &CUI{gui: g}, nil
}

// Notify updates the chat view with a string.
// This satisfies the Notifier interface of the chat client so it can subscribe
// to events that happen and update the UI when input is received from the server.
func (c *CUI) Notify(msg *chat.Message) {
	c.gui.Execute(func(g *gocui.Gui) error {
		v, err := g.View("chat")
		if err != nil {
			return err
		}

		var c Color
		if msg.Register || msg.Close {
			c = BoldBlue
		} else {
			c = BoldWhite
		}

		Fprintf(v, c, "[%s] %s", msg.User, msg.Data)
		return nil
	})
}

// Run runs the user interface main loop until a quit signal is detected.
func (c *CUI) Run() error {
	return c.gui.MainLoop()
}

// Close closes the CUI.
func (c *CUI) Close() error {
	return c.Close()
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	// Display view that shows the input from the server or other users.
	if chatView, err := g.SetView("chat", 0, 0, maxX-1, maxY-4); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		chatView.Frame = true
		chatView.Autoscroll = true
		chatView.Title = "Chat Room"
	}

	// Input view where the user writes what they are going to send to the server.
	if inputView, err := g.SetView("input", 0, maxY-3, maxX-1, maxY-1); err != nil {
		if err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}

			inputView.Editable = true
			inputView.Title = "Input"

			if _, err := g.SetCurrentView("input"); err != nil {
				return err
			}
		}
	}
	return nil
}

func setKeyBindings(g *gocui.Gui, s SendCloser) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit(s)); err != nil {
		return err
	}
	return g.SetKeybinding("input", gocui.KeyEnter, gocui.ModNone, send(s))
}

// reset clears the buffer of a view and sets the cursor back to 0, 0.
func reset(v *gocui.View) error {
	v.Clear()
	if err := v.SetCursor(0, 0); err != nil {
		return err
	}
	return v.SetOrigin(0, 0)
}

func quit(s SendCloser) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		if err := s.Close(); err != nil {
			return err
		}
		return gocui.ErrQuit
	}
}

func send(s SendCloser) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		if err := s.Send(v.Buffer()); err != nil {
			return err
		}
		return reset(v)
	}
}
