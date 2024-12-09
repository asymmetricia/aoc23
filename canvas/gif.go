package canvas

import (
	"fmt"
	"image"
	"image/gif"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/asymmetricia/aoc23/aoc"
)

const DefaultFrameTiming = 2

// RenderGif creates a GIF from the stack of canvases. (TODO: motion blending?)
func RenderGif(canvases []*Canvas, filename string, log logrus.FieldLogger) {
	anim := &gif.GIF{}

	err := render(canvases, log, func(frame renderedFrame, last bool) error {
		anim.Image = append(anim.Image, frame.img)
		anim.Delay = append(anim.Delay, int(frame.timing))
		anim.Disposal = append(anim.Disposal, gif.DisposalNone)
		return nil
	})

	if err != nil {
		log.Fatalf("error rendering GIF: %w", err)
	}

	anim.Delay[len(anim.Delay)-2] = 500

	aoc.SaveGIF(anim, filename, log)
	log.Printf("rendered %d frames", len(canvases))
}

type renderedFrame struct {
	id     int
	img    *image.Paletted
	timing float32
}

func render(canvases []*Canvas, log logrus.FieldLogger, encode func(frame renderedFrame, last bool) error) error {
	width, height := Bounds(canvases)

	type RenderReq struct {
		id     int
		canvas *Canvas
	}

	canvasCh := make(chan RenderReq)
	frames := make(chan renderedFrame)

	go func() {
		last := time.Now()
		for i, canvas := range canvases {
			canvasCh <- RenderReq{i, canvas}
			if time.Since(last) >= 5*time.Second {
				last = time.Now()
				log.Printf("rendering frame %d of %d", i, len(canvases))
			}
		}
		close(canvasCh)
	}()

	wg := &sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			for req := range canvasCh {
				frames <- renderedFrame{req.id, req.canvas.RenderRect(width, height), req.canvas.Timing}
			}
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(frames)
	}()

	next := 0
	frameCache := map[int]renderedFrame{}
	var last = time.Now()
	for {
		frame, ok := <-frames
		if !ok {
			break
		}

		frameCache[frame.id] = frame

		for {
			f, ok := frameCache[next]
			if !ok {
				break
			}

			if time.Since(last) >= 5*time.Second {
				last = time.Now()
				log.Printf("encoding frame %d of %d", next, len(canvases))
			}

			if err := encode(f, next == len(canvases)-1); err != nil {
				return fmt.Errorf("encoding frame %d: %w", next, err)
			}

			delete(frameCache, next)
			next++
		}
	}

	log.Printf("rendered %d frames", len(canvases))
	return nil
}

// RenderMp4 creates an MP4 from the stack of canvases. In this mode, timing is
// interpreted as the number of times to duplicate the frame *divided by two*
// (so a default frame timing of 2 will appear once). (TODO: motion blending?)
func RenderMp4(canvases []*Canvas, filename string, log logrus.FieldLogger) error {
	const framerate = 60
	enc, err := aoc.NewMP4Encoder(filename, framerate, log)
	if err != nil {
		return fmt.Errorf("error creating MP4 encoder: %w", err)
	}

	err = render(canvases, log, func(f renderedFrame, last bool) error {
		var count int
		if f.timing == 0 {
			count = DefaultFrameTiming
		} else {
			count = int(f.timing / 100 * framerate)
			if count < 1 {
				count = 1
			}
		}
		if last {
			count = framerate * 5
		}

		for i := 0; i < count; i++ {
			if err := enc.Encode(f.img); err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return err
	}

	if err := enc.Close(); err != nil {
		return fmt.Errorf("error closing MP4 encoder: %w", err)
	}

	return nil
}
