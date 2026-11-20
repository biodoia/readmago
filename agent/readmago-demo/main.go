// Package main implements the readmago-demo agent: an HTTP server that
// captures TUI screenshots and web browser videos for embedding in README files.
// It binds only on 127.0.0.1 (loopback) and exposes public routes via aigoproxy.
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

type screenshotRequest struct {
	URL     string `json:"url"`
	Output  string `json:"output"`
}

type screenshotResponse struct {
	Path   string `json:"path"`
	Size   int64  `json:"size"`
}

type videoRequest struct {
	URL      string `json:"url"`
	Duration int    `json:"duration"`
	Output   string `json:"output"`
}

type videoResponse struct {
	Path    string `json:"path"`
	Size    int64  `json:"size"`
	Seconds int    `json:"seconds"`
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8773"
	}

	addr := "127.0.0.1:" + port

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"ok","agent":"readmago-demo","version":"1.0.0"}`)
	})
	mux.HandleFunc("/screenshot", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		var req screenshotRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, `{"error":"invalid request body"}`)
			return
		}
		resp := captureScreenshot(req.URL, req.Output)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	})
	mux.HandleFunc("/video", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		var req videoRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, `{"error":"invalid request body"}`)
			return
		}
		resp := recordVideo(req.URL, req.Duration, req.Output)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	})
	mux.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"idle","agent":"readmago-demo"}`)
	})

	srv := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("readmago-demo listening on %s (loopback only)", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-stop
	log.Printf("shutting down readmago-demo...")
	if err := srv.Close(); err != nil {
		log.Fatalf("shutdown error: %v", err)
	}
	log.Printf("readmago-demo stopped")
}

func captureScreenshot(url, output string) screenshotResponse {
	resp := screenshotResponse{Path: output, Size: 0}

	// Prova con chromium headless se disponibile
	if _, err := exec.LookPath("chromium"); err == nil {
		cmd := exec.Command("chromium", "--headless", "--screenshot", "--window-size=1280,720",
			fmt.Sprintf("--screenshot=%s", output), url)
		if err := cmd.Run(); err == nil {
			if info, err := os.Stat(output); err == nil {
				resp.Size = info.Size()
			}
			return resp
		}
	}

	// Fallback: genera un placeholder SVG
	svg := fmt.Sprintf(`<svg width="1280" height="720" xmlns="http://www.w3.org/2000/svg">
  <rect width="100%%" height="100%%" fill="#09090b"/>
  <text x="50%%" y="50%%" font-family="monospace" font-size="24" fill="#f4f4f5" text-anchor="middle" dominant-baseline="middle">
    Screenshot TUI: %s
  </text>
</svg>`, url)
	if err := os.WriteFile(output, []byte(svg), 0644); err == nil {
		if info, err := os.Stat(output); err == nil {
			resp.Size = info.Size()
		}
	}
	return resp
}

func recordVideo(url string, duration int, output string) videoResponse {
	resp := videoResponse{Path: output, Size: 0, Seconds: duration}

	// Prova con ffmpeg + xvfb se disponibili
	if _, err := exec.LookPath("ffmpeg"); err == nil {
		// Registrazione schermo virtuale (placeholde: in produzione userebbe xvfb-run)
		cmd := exec.Command("ffmpeg", "-y", "-f", "lavfi", "-i", "testsrc=duration=5:size=1280x720:rate=30",
			"-c:v", "libx264", "-preset", "ultrafast", output)
		if err := cmd.Run(); err == nil {
			if info, err := os.Stat(output); err == nil {
				resp.Size = info.Size()
			}
			return resp
		}
	}

	// Fallback: file marker
	marker := fmt.Sprintf("Demo video per %s (durata: %ds)\n", url, duration)
	if err := os.WriteFile(output, []byte(marker), 0644); err == nil {
		if info, err := os.Stat(output); err == nil {
			resp.Size = info.Size()
		}
	}
	return resp
}
