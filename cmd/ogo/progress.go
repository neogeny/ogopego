package main

import (
	"github.com/neogeny/ogopego/util/heartbeat"
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

type CliHeartbeat struct {
	bar      *mpb.Bar
	lastMark int
}

func NewCliHeartbeat(bar *mpb.Bar) *CliHeartbeat {
	return &CliHeartbeat{bar: bar}
}

func (h *CliHeartbeat) Tick(mark, _ int) {
	if mark > h.lastMark {
		h.bar.SetCurrent(int64(mark))
		h.lastMark = mark
	}
}

type LoadProgress struct {
	bar *mpb.Bar
	hb  *CliHeartbeat
}

func NewLoadProgress(p *mpb.Progress, msg string) *LoadProgress {
	bar := p.New(0,
		mpb.SpinnerStyle("⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏").Build(),
		mpb.AppendDecorators(decor.Name(msg)),
		mpb.BarRemoveOnComplete(),
	)
	return &LoadProgress{bar: bar, hb: NewCliHeartbeat(bar)}
}

func (lp *LoadProgress) Heartbeat() heartbeat.Heartbeat {
	return lp.hb
}

func (lp *LoadProgress) Finish() {
	lp.bar.SetTotal(0, true)
	lp.bar.Wait()
}

type FileProgress struct {
	bar *mpb.Bar
	hb  *CliHeartbeat
}

func NewFileProgress(p *mpb.Progress, name string) *FileProgress {
	bar := p.New(0,
		mpb.BarStyle().Build(),
		mpb.PrependDecorators(
			decor.Name(name, decor.WC{W: 40, C: decor.DidentRight}),
		),
		mpb.AppendDecorators(
			decor.Percentage(),
			decor.Elapsed(decor.ET_STYLE_GO),
		),
		mpb.BarRemoveOnComplete(),
	)
	return &FileProgress{bar: bar, hb: NewCliHeartbeat(bar)}
}

func (fp *FileProgress) Heartbeat() heartbeat.Heartbeat {
	return fp.hb
}

func (fp *FileProgress) SetLength(length int) {
	fp.bar.SetTotal(int64(length), false)
}

func (fp *FileProgress) Success() {
	fp.bar.SetTotal(0, true)
	fp.bar.Wait()
}

func (fp *FileProgress) Fail() {
	fp.bar.Abort(true)
}

type ProgressUI struct {
	p     *mpb.Progress
	files *mpb.Bar
}

func NewProgressUI(total int) *ProgressUI {
	p := mpb.New(mpb.WithOutput(nil))

	files := p.New(int64(total),
		mpb.BarStyle().Build(),
		mpb.PrependDecorators(
			decor.CountersNoUnit("%d/%d files"),
		),
		mpb.BarRemoveOnComplete(),
	)

	return &ProgressUI{p: p, files: files}
}

func (ui *ProgressUI) Loading(msg string) *LoadProgress {
	return NewLoadProgress(ui.p, msg)
}

func (ui *ProgressUI) AddFile(name string) *FileProgress {
	return NewFileProgress(ui.p, name)
}

func (ui *ProgressUI) IncFiles() {
	ui.files.Increment()
}

func (ui *ProgressUI) Finish() {
	ui.files.SetTotal(0, true)
	ui.p.Wait()
}
