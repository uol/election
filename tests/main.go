package tests

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/uol/election"
	"github.com/uol/funks"
	"github.com/uol/logh"
)

//
// Does the election process using the election package
// author: rnojiri
//

func main() {

	logh.ConfigureGlobalLogger(logh.INFO, logh.CONSOLE)

	cfg := election.Config{
		ZKURL:                  []string{"zookeeper.intranet"},
		ZKElectionNodeURI:      "/master",
		ZKSlaveNodesURI:        "/slaves",
		ReconnectionTimeout:    *funks.ForceNewStringDuration("3s"),
		SessionTimeout:         *funks.ForceNewStringDuration("5s"),
		ClusterChangeCheckTime: *funks.ForceNewStringDuration("1s"),
		ClusterChangeWaitTime:  *funks.ForceNewStringDuration("1s"),
	}

	manager, err := election.New(&cfg)
	if err != nil {
		logh.Error().Err(err).Send()
		os.Exit(0)
	}

	feedbackChannel, err := manager.Start()

	go func() {
		for {
			select {
			case signal := <-*feedbackChannel:
				if signal == election.Master {
					logh.Info().Msg("master signal received")
				} else if signal == election.Slave {
					logh.Info().Msg("slave signal received")
				} else if signal == election.ClusterChanged {
					logh.Info().Msg("cluster changed signal received")
				} else if signal == election.Disconnected {
					logh.Info().Msg("disconnected signal received")
				}
			}
		}
	}()

	ci, err := manager.GetClusterInfo()
	if err != nil {
		logh.Error().Err(err).Send()
		os.Exit(0)
	}

	logh.Info().Msgf("%+v", ci)

	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)

	go func() {
		<-gracefulStop
		logh.Error().Msg("exiting...")
		manager.Disconnect()
		time.Sleep(2 * time.Second)
		os.Exit(0)
	}()

	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}
