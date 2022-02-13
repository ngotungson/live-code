package mindx

import (
	"errors"
	"flag"
	"fmt"
	"mindx/internal/app/mindx/api"
	"mindx/internal/app/mindx/database"
	"os"
	"os/signal"
	"path"
	"sync"
	"syscall"

	"github.com/spf13/viper"
)

var (
	configFile   = ".config"
	printVersion = flag.Bool("version", false, "Print the version and exit")
	gracefulExit = errors.New("graceful exit")
)

type App struct {
	Name    string
	Version string
	done    chan struct{}
}

func newApp(name, version string) *App {
	return &App{
		Name:    name,
		Version: version,
		done:    make(chan struct{}),
	}
}

func Run(name, version string) error {
	return handleError(newApp(name, version).launch())
}

// HandleSignals manages OS signals that ask the api/daemon to stop.
// The stopFunction should break the loop in the Beat so that
// the api shut downs gracefully.
func HandleSignals(stopFunction func()) {
	var callback sync.Once

	// On ^C or SIGTERM, gracefully stop the sniffer
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigc
		callback.Do(stopFunction)
	}()
}

func (c *App) launch() error {
	if err := c.configure(); err != nil {
		return err
	}

	// Setup DB
	err := database.Init()
	if err != nil {
		return fmt.Errorf("error connect to database: %v", err)
	}

	//Init API
	api.InitAPI()
	////Init RPC

	// Blocks progressing. As soon as channel is closed, all defer statements come into play
	<-c.done
	return nil
}

func (c *App) Stop() {
	close(c.done)
}

func setDefaults() {
}

//configure setups default config, reads config file and init logger
func (c *App) configure() error {
	setDefaults()
	var err error
	flag.Parse()

	if *printVersion {
		fmt.Printf("%s version %s \n",
			c.Name, c.Version)
		return gracefulExit
	}

	base := path.Base(configFile)
	ext := path.Ext(configFile)
	dir := path.Dir(configFile)
	configName := base[:len(base)-len(ext)]

	viper.SetConfigName(configName)
	viper.AddConfigPath(dir)
	err = viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("error reading config file: %v", err)
	}
	return nil
}

func handleError(err error) error {
	if err == nil || err == gracefulExit {
		return nil
	}

	fmt.Fprintf(os.Stderr, "Exiting: %v\n", err)
	return err
}
