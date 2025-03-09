package clispinner

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/briandowns/spinner"

	"github.com/ignite/cli/v29/ignite/pkg/cliui/colors"
)

var (
	simpleCharset     = spinner.CharSets[4]
	simpleRefreshRate = time.Millisecond * 300
	simpleColor       = colors.Spinner
)

type SimpleSpinner struct {
	mu       sync.Mutex
	writer   io.Writer
	charset  []string
	text     string
	prefix   string
	color    string
	active   bool
	stopChan chan struct{}
}

// newSimpleSpinner creates a new simple spinner.
func newSimpleSpinner(o Options) *SimpleSpinner {
	text := o.text
	if text == "" {
		text = DefaultText
	}

	charset := o.charset
	if len(charset) == 0 {
		charset = simpleCharset
	}

	writer := o.writer
	if writer == nil {
		writer = os.Stdout
	}

	return &SimpleSpinner{
		charset: charset,
		text:    text,
		writer:  writer,
	}
}

// SetText sets the text for the spinner.
func (s *SimpleSpinner) SetText(text string) Spinner {
	s.mu.Lock()
	s.text = text
	s.mu.Unlock()
	return s
}

// SetPrefix sets the prefix for the spinner.
func (s *SimpleSpinner) SetPrefix(prefix string) Spinner {
	s.mu.Lock()
	s.prefix = prefix
	s.mu.Unlock()
	return s
}

// SetCharset sets the charset for the spinner.
func (s *SimpleSpinner) SetCharset(charset []string) Spinner {
	s.mu.Lock()
	s.charset = charset
	s.mu.Unlock()
	return s
}

// SetColor sets the color for the spinner (if color functionality is added).
func (s *SimpleSpinner) SetColor(color string) Spinner {
	s.mu.Lock()
	s.color = color
	s.mu.Unlock()
	return s
}

// Start begins the spinner animation.
func (s *SimpleSpinner) Start() Spinner {
	s.mu.Lock()
	if s.active {
		s.mu.Unlock()
		return s // Do nothing if already active
	}
	s.active = true
	s.stopChan = make(chan struct{})

	// Extract spinner data safely within the mutex
	prefix := s.prefix
	text := s.text
	writer := s.writer
	charset := s.charset
	s.mu.Unlock()

	// Start the animation loop in a separate goroutine
	go func() {
		ticker := time.NewTicker(simpleRefreshRate)
		defer ticker.Stop()

		index := 0
		for {
			select {
			case <-s.stopChan: // Stop the spinner
				_, _ = fmt.Fprintf(writer, "\r\033[K") // Clear the spinner's line
				return
			case <-ticker.C: // Update the spinner on each tick
				s.mu.Lock()
				frame := charset[index]
				str := fmt.Sprintf("\r%s%s %s", prefix, simpleColor(frame), text)
				_, _ = fmt.Fprint(writer, str) // Update the spinner in the same line
				index++
				if index >= len(charset) {
					index = 0
				}
				s.mu.Unlock()
			}
		}
	}()
	return s
}

// Stop ends the spinner animation.
func (s *SimpleSpinner) Stop() Spinner {
	s.mu.Lock()
	if !s.active {
		s.mu.Unlock()
		return s // Do nothing if already inactive
	}
	close(s.stopChan)
	s.active = false
	s.stopChan = nil
	fmt.Print("\r") // Clear spinner line on stop
	s.mu.Unlock()
	return s
}

// IsActive returns whether the spinner is currently active.
func (s *SimpleSpinner) IsActive() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.active
}

// Writer returns the spinner writer.
func (s *SimpleSpinner) Writer() io.Writer {
	return s.writer
}
