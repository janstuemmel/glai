package main

import (
	"bytes"
	"io"
	"math"
	"os"
	"os/signal"
	"syscall"

	"github.com/charmbracelet/glamour"
)

func clearAltScreen() {
	os.Stdout.WriteString("\033[H\033[2J")
}

func enterAltScreen() {
	os.Stdout.WriteString("\033[?1049h")
}

func leaveAltScreen() {
	os.Stdout.WriteString("\033[?1049l")
}

func clamp(x, min, max int) int {
	if x < min {
		return min
	}
	if x > max {
		return max
	}
	return x
}

func adaptiveThreshold(bufLen int) int {
	t := 32 + int(math.Sqrt(float64(bufLen)))*2
	return clamp(t, 32, 16*1024)
}

func main() {
	enterAltScreen()
	clearAltScreen()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	renderer, _ := glamour.NewTermRenderer(
		glamour.WithWordWrap(120),
		glamour.WithAutoStyle(),
		glamour.WithEmoji(),
		glamour.WithEnvironmentConfig(),
	)

	var buf bytes.Buffer
	var out []byte
	var lastBufLen int
	tmp := make([]byte, 4096)

	go func() {
		<-sigChan
		leaveAltScreen()
		clearAltScreen()
		os.Stdout.Write(out)
		os.Exit(0)
	}()

	for {
		n, err := os.Stdin.Read(tmp)
		if n > 0 {
			buf.Write(tmp[:n])
			delta := adaptiveThreshold(buf.Len())

			if buf.Len()-lastBufLen >= delta || lastBufLen == 0 {
				clearAltScreen()
				out, _ = renderer.RenderBytes(buf.Bytes())
				os.Stdout.Write(out)
				lastBufLen = buf.Len()
			}
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			break
		}
	}

	sigChan <- os.Interrupt
	select {}
}
