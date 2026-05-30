package cli

//goland:noinspection GoRedundantImportAlias
import (
	"os"

	"github.com/fatih/color"
	"github.com/neogeny/ogopego/pkg/util/heartbeat"
	mpb "github.com/vbauerster/mpb/v8"
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
	if h == nil || h.bar == nil {
		return
	}
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
	if p == nil {
		return &LoadProgress{}
	}
	bar := p.New(0,
		mpb.SpinnerStyle("⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"),
		mpb.AppendDecorators(decor.Name(msg)),
		mpb.BarRemoveOnComplete(),
	)
	return &LoadProgress{bar: bar, hb: NewCliHeartbeat(bar)}
}

// Heartbeat returns the heartbeat interface for the load progress.
func (lp *LoadProgress) Heartbeat() heartbeat.Heartbeat {
	if lp == nil || lp.hb == nil {
		return heartbeat.NullHeartbeat{}
	}
	return lp.hb
}

// Finish marks the load progress as complete.
func (lp *LoadProgress) Finish() {
	if lp == nil || lp.bar == nil {
		return
	}
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
	if p == nil {
		return &FileProgress{}
	}
	yellow := func(s string) string { return color.YellowString(s) }
	bar := p.New(0,
		mpb.BarStyle().
			Lbound(" ").
			Rbound(" ").
			Filler("-").FillerMeta(yellow).
			Padding("  ").
			Tip("-").TipMeta(yellow),
		mpb.PrependDecorators(
			decor.Name(name, decor.WC{W: 40, C: decor.DindentRight}),
		),
		mpb.AppendDecorators(
			decor.Percentage(),
		),
		mpb.BarRemoveOnComplete(),
	)
	return &FileProgress{bar: bar, hb: NewCliHeartbeat(bar)}
}

// Heartbeat returns the heartbeat interface for the file progress.
func (fp *FileProgress) Heartbeat() heartbeat.Heartbeat {
	if fp == nil || fp.hb == nil {
		return heartbeat.NullHeartbeat{}
	}
	return fp.hb
}

// SetLength sets the total length for the file progress bar.
func (fp *FileProgress) SetLength(length int) {
	if fp == nil || fp.bar == nil {
		return
	}
	fp.length = length
	fp.bar.SetTotal(int64(length), false)
}

// Success marks the file progress as successful.
func (fp *FileProgress) Success() {
	if fp == nil || fp.bar == nil {
		return
	}
	fp.bar.SetCurrent(1 + int64(fp.length))
	fp.bar.Abort(true)
}

// Fail marks the file progress as failed.
func (fp *FileProgress) Fail() {
	if fp == nil || fp.bar == nil {
		return
	}
	fp.bar.Abort(true)
}

// ProgressUI manages the overall progress display for the CLI.
type ProgressUI struct {
	p     *mpb.Progress
	files *mpb.Bar
}

// NewProgressUI constructs the terminal progress UI used by the CLI and
// returns a manager for creating per-file and load progress handles.
// When quiet is true, all progress bar calls are no-ops.
func NewProgressUI(total int, quiet bool) *ProgressUI {
	if quiet {
		return &ProgressUI{}
	}
	p := mpb.New(mpb.WithOutput(os.Stderr))

	green := func(s string) string { return color.GreenString(s) }
	files := p.New(int64(total),
		mpb.BarStyle().
			Lbound(" ").
			Rbound(" ").
			Filler(".").FillerMeta(green).
			Padding(" ").
			Tip(".").TipMeta(green),
		mpb.PrependDecorators(
			decor.CountersNoUnit("%d/%d files"),
		),
		mpb.BarRemoveOnComplete(),
	)

	return &ProgressUI{p: p, files: files}
}

// Loading creates and returns a new LoadProgress instance.
func (ui *ProgressUI) Loading(msg string) *LoadProgress {
	if ui == nil || ui.p == nil {
		return &LoadProgress{}
	}
	return NewLoadProgress(ui.p, msg)
}

// AddFile creates and returns a new FileProgress instance.
func (ui *ProgressUI) AddFile(name string) *FileProgress {
	if ui == nil || ui.p == nil {
		return &FileProgress{}
	}
	return NewFileProgress(ui.p, name)
}

// IncFiles increments the count of processed files.
func (ui *ProgressUI) IncFiles() {
	if ui == nil || ui.files == nil {
		return
	}
	ui.files.Increment()
}

// Finish marks the overall progress UI as complete.
func (ui *ProgressUI) Finish() {
	if ui == nil || ui.p == nil {
		return
	}
	ui.files.SetTotal(ui.files.Current()+1, true)
	ui.p.Wait()
}
