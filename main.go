package main

import (
	"github.com/gotk3/gotk3/gtk"
	"time"
	"log"
	"net"
	"fmt"
	//"strconv"
	"sync"
	//"os"
)
var goroutines = 256
var i int
var msgChan = make(chan string)
var USAGE = `
The Crawl Help Page:
	LEGEND:
	{} = mandatory option
	[] = optional option 
	- = short flag
	-- = long flag
	CRAWL {start} [{-g} {int}]
`

func dialip(end int) (net.Conn, error) {
	//return net.DialTimeout("tcp", fmt.Sprintf("192.168.1.%d:8080", end), 500 * time.Millisecond)
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("192.168.1.%d:8080", end), 1000*time.Millisecond)
	return conn, err

}

func start() {
	var wg sync.WaitGroup
	for w := 0; w < goroutines; w++ {
		wg.Add(1)
		go dial(&wg)
	}
	wg.Wait()
}

func dial(wg *sync.WaitGroup) {
	for ; i <= 256; i++ {
		var worked = true
		conn, err := dialip(i)
		if err != nil {
			worked = false
		}
		if worked == true {
			var ntime = time.Now()
			_ = conn.SetReadDeadline(ntime.Add(10000 * time.Millisecond))
			var text = make([]byte, 100)
			_, err = conn.Read(text)
			if err != nil {
				continue
			}
			msgChan <-fmt.Sprintf("message: %s", string(text))
			conn.Close()
		}
	}
	wg.Done()
}
func crawl() {
	/*var action string
	if len(os.Args) >= 2 {
		action = os.Args[1]
	}
	for q := 2; q < len(os.Args); q++ {
		switch os.Args[q] {
		case "-g":
			if len(os.Args) < q {
				fmt.Println(USAGE)
				os.Exit(1)
			} else {
				//goroutines = int(os.Args[q+1])
				if len(os.Args) <= q+1 {
					fmt.Println(USAGE)
					os.Exit(1)
				}
				if s, err := strconv.Atoi(os.Args[q+1]); err == nil {
					goroutines = s
				} else {
					fmt.Println(USAGE)
					os.Exit(1)
				}
				q++
			}
		}

	}
	if goroutines > 1020 {
		fmt.Println("too many goroutines, might cause problems on some systems, I would reccomend less than 1020")
		fmt.Println("goroutines greater than 255 are pointless because this is a subnet scanner")
	}
	switch action {
	case "start":
		start()
	default:
		fmt.Println(USAGE)
	}
}
*/
	start()
}

func main() {
	// Initialize GTK without parsing any command line arguments.
	gtk.Init(nil)

	// Create a new toplevel window, set its title, and connect it to the
	// "destroy" signal to exit the GTK main loop when it is destroyed.
	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("Unable to create window:", err)
	}
	win.SetTitle("Simple Example")
	win.Connect("destroy", func() {
		gtk.MainQuit()
	})

	// Create a new label widget to show in the window.
	l, err := gtk.LabelNew("")
	if err != nil {
		log.Fatal("Unable to create label:", err)
	}

	// Add the label to the window.
	win.Add(l)
	go func() {
		for {
			select {
			case b := <-msgChan:
				l.SetLabel(b)
			default:
				time.Sleep(50*time.Millisecond)
			}
		}
	}()
	go crawl()

	// Set the default window size.
	win.SetDefaultSize(800, 600)

	// Recursively show all widgets contained in this window.
	win.ShowAll()

	// Begin executing the GTK main loop.  This blocks until
	// gtk.MainQuit() is run. 
	gtk.Main()
}
