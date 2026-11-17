package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"syscall"
	"time"
)

// openSource returns the IQ sample stream (unsigned 8-bit interleaved I/Q)
// and whether it is a live receiver.
func openSource(ctx context.Context, cfg config) (io.ReadCloser, bool, error) {
	switch cfg.input {
	case "rtl":
		r, err := startRTL(ctx, cfg)
		return r, true, err
	case "-":
		return os.Stdin, true, nil
	}
	f, err := os.Open(cfg.input)
	if err != nil {
		return nil, false, err
	}
	return struct {
		io.Reader
		io.Closer
	}{bufio.NewReaderSize(f, 1<<20), f}, false, nil
}

// rtlSource is a running rtl_sdr process streaming samples on its stdout.
type rtlSource struct {
	cmd *exec.Cmd
	out io.ReadCloser
}

func startRTL(ctx context.Context, cfg config) (*rtlSource, error) {
	args := []string{
		"-d", cfg.device,
		"-f", strconv.FormatInt(int64(cfg.freq), 10),
		"-s", strconv.FormatInt(int64(cfg.rate), 10),
		"-g", cfg.gain,
		"-p", strconv.Itoa(cfg.ppm),
		"-",
	}
	cmd := exec.CommandContext(ctx, "rtl_sdr", args...)
	cmd.Stderr = os.Stderr
	// SIGTERM lets rtl_sdr close the dongle cleanly (and hand it back to the
	// kernel driver it detached).
	cmd.Cancel = func() error { return cmd.Process.Signal(syscall.SIGTERM) }
	cmd.WaitDelay = 3 * time.Second
	out, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("starting rtl_sdr (is the rtl-sdr package installed?): %w", err)
	}
	log.Printf("rtl_sdr started: %.3f MHz, %.0f samples/s, gain %s", cfg.freq/1e6, cfg.rate, cfg.gain)
	return &rtlSource{cmd: cmd, out: out}, nil
}

func (r *rtlSource) Read(p []byte) (int, error) {
	n, err := r.out.Read(p)
	if err == io.EOF {
		// rtl_sdr only stops by itself when the dongle fails or goes away.
		if werr := r.cmd.Wait(); werr != nil {
			return n, fmt.Errorf("rtl_sdr exited: %w", werr)
		}
		return n, fmt.Errorf("rtl_sdr exited")
	}
	return n, err
}

func (r *rtlSource) Close() error {
	if r.cmd.ProcessState == nil {
		_ = r.cmd.Process.Signal(syscall.SIGTERM)
		_ = r.cmd.Wait()
	}
	return nil
}
