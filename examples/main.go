package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	goruntime "runtime"
	"runtime/pprof"
	"time"

	"github.com/linuxerwang/goats-html/examples/common_html"
	"github.com/linuxerwang/goats-html/examples/data"
	"github.com/linuxerwang/goats-html/examples/shelf_view_html"
	"github.com/linuxerwang/goats-html/runtime"
)

var cpuprofile = flag.String("cpuprofile", "", "Write cpu profile to file")
var benchmark = flag.Bool("benchmark", false, "Benchmark template speed")
var smallOnly = flag.Bool("small", false, "Only run small template")
var largeOnly = flag.Bool("large", false, "Only run large template")
var shelf *data.Shelf
var N = 1000000
var done chan bool

func init() {
	goruntime.GOMAXPROCS(goruntime.NumCPU())
	flag.Parse()
	shelf = data.NewBookShelf()
	runtime.InitGoats(nil)
	done = make(chan bool)
}

func main() {
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	if *benchmark {
		startTime := time.Now()

		var n int
		if !*smallOnly {
			n++
			go benchmarkShelfViews()
		} else if !*largeOnly {
			n++
			go benchmarkBookCards()
		}

		for i := 0; i < n*goruntime.GOMAXPROCS(0); i++ {
			<-done
			fmt.Println(i)
		}

		duration := time.Now().Sub(startTime).Seconds()
		fmt.Println("Duration:", duration, "seconds")
	} else {
		renderShelfView()
		renderBookCard()
	}
}

func renderShelfView() {
	args := &shelf_view_html.ShelfViewTemplateArgs{
		Shelf: shelf,
	}

	settings := runtime.TemplateSettings{
		OmitDocType: false,
	}

	var buffer bytes.Buffer
	template := shelf_view_html.NewShelfViewTemplate(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if !*benchmark {
			fmt.Println("Generated html: ", buffer.String())
		}
	} else {
		fmt.Println("Failed to render template. ", err)
	}

	if !*benchmark {
		fmt.Println()
	}
}

func benchmarkShelfView() {
	for i := 0; i < N/goruntime.GOMAXPROCS(0); i++ {
		renderShelfView()
	}
	done <- true
}

func benchmarkShelfViews() {
	fmt.Println("\nBenchmark large template...\n")
	for i := 0; i < goruntime.GOMAXPROCS(0); i++ {
		go benchmarkShelfView()
	}
}

func renderBookCard() {
	args := &common_html.BookCardTemplateArgs{
		Book:    shelf.Books[0],
		Loopvar: &runtime.LoopVar{Counter: 1},
	}

	settings := runtime.TemplateSettings{
		OmitDocType: false,
	}

	var buffer bytes.Buffer
	template := common_html.NewBookCardTemplate(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if !*benchmark {
			fmt.Println("Generated html: ", buffer.String())
		}
	} else {
		fmt.Println("Failed to render template. ", err)
	}
}

func benchmarkBookCard() {
	for i := 0; i < N/goruntime.GOMAXPROCS(0); i++ {
		renderBookCard()
	}
	done <- true
}

func benchmarkBookCards() {
	fmt.Println("\nBenchmark small template...\n")
	for i := 0; i < goruntime.GOMAXPROCS(0); i++ {
		go benchmarkBookCard()
	}
}
