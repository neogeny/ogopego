package main

import (
	"os"

	"github.com/neogeny/ogopego/util/heartbeat"
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

// CliHeartbeat implements the heartbeat.Heartbeat interface using a mpb.Bar.
type CliHeartbeat struct {
	bar      *mpb.Bar
	lastMark int
}

// NewCliHeartbeat returns a heartbeat implementation backed by a mpb.Bar.
func NewCliHeartbeat(bar *mpb.Bar) *CliHeartbeat {
	return &CliHeartbeat{bar: bar}
}

// Tick updates the progress bar with the given mark.
func (h *CliHeartbeat) Tick(mark, _ int) {
	if mark > h.lastMark {
		h.bar.SetCurrent(int64(mark))
		h.lastMark = mark
	}
}

// LoadProgress manages a spinner-style progress tracker for load operations.
type LoadProgress struct {
	bar *mpb.Bar
	hb  *CliHeartbeat
}

// NewLoadProgress creates a spinner-style progress tracker for long-running
// load operations and returns a handle providing heartbeat integration.
func NewLoadProgress(p *mpb.Progress, msg string) *LoadProgress {
	bar := p.New(0,
		mpb.SpinnerStyle("⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"),
		mpb.AppendDecorators(decor.Name(msg)),
		mpb.BarRemoveOnComplete(),
	)
	return &LoadProgress{bar: bar, hb: NewCliHeartbeat(bar)}
}

// Heartbeat returns the heartbeat interface for the load progress.
func (lp *LoadProgress) Heartbeat() heartbeat.Heartbeat {
	return lp.hb
}

// Finish marks the load progress as complete.
func (lp *LoadProgress) Finish() {
	lp.bar.SetTotal(0, true)
	lp.bar.Wait()
}

// FileProgress manages a progress bar for individual file processing.
type FileProgress struct {
	bar    *mpb.Bar
	hb     *CliHeartbeat
	length int
}

// NewFileProgress creates a file progress bar with a name and returns a
// handle providing heartbeat integration.
func NewFileProgress(p *mpb.Progress, name string) *FileProgress {
	bar := p.New(0,
		mpb.BarStyle().
			Lbound(" ").
			Rbound(" ").
			Filler("--").
			Padding("  ").
			Tip("-"),
		mpb.PrependDecorators(
			decor.Name(name, decor.WC{W: 40, C: decor.DindentRight}),
		),
		mpb.AppendDecorators(
			decor.Percentage(),
			//decor.Elapsed(decor.ET_STYLE_GO),
		),
		//mpb.BarRemoveOnComplete(),
	)
	return &FileProgress{bar: bar, hb: NewCliHeartbeat(bar)}
}

// Heartbeat returns the heartbeat interface for the file progress.
func (fp *FileProgress) Heartbeat() heartbeat.Heartbeat {
	return fp.hb
}

// SetLength sets the total length for the file progress bar.
func (fp *FileProgress) SetLength(length int) {
	fp.length = length
	fp.bar.SetTotal(int64(length), false)
}

// Success marks the file progress as successful.
func (fp *FileProgress) Success() {
	fp.bar.SetCurrent(1 + int64(fp.length))
	fp.bar.Abort(false)
}

// Fail marks the file progress as failed.
func (fp *FileProgress) Fail() {
	fp.bar.Abort(false)
}

// ProgressUI manages the overall progress display for the CLI.
type ProgressUI struct {
	p     *mpb.Progress
	files *mpb.Bar
}

// NewProgressUI constructs the terminal progress UI used by the CLI and
// returns a manager for creating per-file and load progress handles.
func NewProgressUI(total int) *ProgressUI {
	p := mpb.New(mpb.WithOutput(os.Stderr))

	files := p.New(int64(total),
		mpb.BarStyle().
			Lbound(" ").
			Rbound(" ").
			Filler("--").
			Padding("  ").
			Tip("-"),
		mpb.PrependDecorators(
			decor.CountersNoUnit("%d/%d files"),
		),
		mpb.BarRemoveOnComplete(),
	)

	return &ProgressUI{p: p, files: files}
}

// Loading creates and returns a new LoadProgress instance.
func (ui *ProgressUI) Loading(msg string) *LoadProgress {
	return NewLoadProgress(ui.p, msg)
}

// AddFile creates and returns a new FileProgress instance.
func (ui *ProgressUI) AddFile(name string) *FileProgress {
	return NewFileProgress(ui.p, name)
}

// IncFiles increments the count of processed files.
func (ui *ProgressUI) IncFiles() {
	ui.files.Increment()
}

// Finish marks the overall progress UI as complete.
func (ui *ProgressUI) Finish() {
	ui.files.SetTotal(ui.files.Current()+1, true)
	ui.p.Wait()
}
