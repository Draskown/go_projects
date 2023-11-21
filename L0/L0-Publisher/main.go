package main

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	stan "github.com/nats-io/stan.go"
	"github.com/sirupsen/logrus"
)

// Stan connection configuration
type StanCfg struct {
	ClusterName string
	ClientName  string
	SubjectName string
	ServerUrl   string
}

// Stan configuration constructor
//
// Accepts string values got from environment variables
func NewStanCfg(clun, clin, subn, url string) *StanCfg {
	return &StanCfg{
		ClusterName: clun,
		ClientName:  clin,
		SubjectName: subn,
		ServerUrl:   url,
	}
}

func main() {
	// Create new STAN configuration
	// from environment variables
	stanCfg := NewStanCfg(
		os.Getenv("CLUSTER_NAME"),
		os.Getenv("CLIENT_NAME"),
		os.Getenv("SUBJECT_NAME"),
		os.Getenv("SERVER_URL"),
	)

	// Connect to STAN
	sc, err := stan.Connect(
		stanCfg.ClusterName,
		stanCfg.ClientName+"_pub",
		stan.NatsURL(stanCfg.ServerUrl),
	)
	if err != nil {
		logrus.Fatalf("Error while connecting to STAN in Publisher (%s)\n", err.Error())
		return
	}

	// Open model file
	file, err := os.Open("model.json")
	if err != nil {
		logrus.Fatalf("Error while opening 'model.json' file in Publisher (%s)\n", err.Error())
		return
	}
	defer file.Close()

	// Watches for the file to change
	go func() {
		err := watchFile(file, sc, stanCfg.SubjectName)
		if err != nil {
			logrus.Fatalf("Error while reading from 'model.json' file in Publisher (%s)\n", err.Error())
		}
	}()

	// Waits for the Ctrl+C signal to suspend the main goroutine
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	// Close the STAN connection
	if err := sc.Close(); err != nil {
		logrus.Fatalf("Error while closing STAN connection in Publisher (%s)\n", err.Error())
	}
}

// Watches for the file's change in the infinite loop
// and sends its contents to the STAN subject
//
// So every time 'model.json' is changed, the message is sent
func watchFile(file *os.File, sc stan.Conn, subName string) error {
	// Get initial file statistics
	initialStat, err := file.Stat()
	if err != nil {
		return err
	}

	for {
		// Get current file statistics
		stat, err := file.Stat()
		if err != nil {
			return err
		}

		// If current statistics is different from the initial
		if stat.Size() != initialStat.Size() || stat.ModTime() != initialStat.ModTime() {
			bytes, err := io.ReadAll(file)
			if err != nil {
				fmt.Println(err.Error())
				return err
			}

			// Publish the changed file contents
			sc.Publish(subName, bytes)

			// Update initial statistics the new ones
			initialStat = stat
		}

		// Return to the start of the file
		_, err = file.Seek(0, io.SeekStart)
		if err != nil {
			return err
		}

		// One second delay
		time.Sleep(1 * time.Second)
	}
}
