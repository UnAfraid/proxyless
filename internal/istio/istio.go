package istio

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

func WaitForSidecar(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	httpClient := &http.Client{
		Timeout: time.Second,
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if err := checkIstioProxy(ctx, httpClient); err != nil {
				if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
					return err
				}
				time.Sleep(500 * time.Millisecond)
				continue
			}
			return nil
		}
	}
}

func checkIstioProxy(ctx context.Context, httpClient *http.Client) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:15020/healthz/ready", http.NoBody)
	if err != nil {
		return err
	}

	res, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	_, _ = io.Copy(io.Discard, res.Body)

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return fmt.Errorf("istio proxy not ready: %v", res.Status)
	}
	return nil
}
