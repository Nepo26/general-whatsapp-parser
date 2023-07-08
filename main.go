// Command remote is a chromedp example demonstrating how to connect to an
// existing Chrome DevTools instance using a remote WebSocket URL.
package main

import (
	"context"
	"flag"
	"log"
	"os"
	"time"

	"github.com/chromedp/chromedp"
)

func main() {
	devtoolsWsURL := flag.String("devtools-ws-url", "", "DevTools WebSsocket URL")
	flag.Parse()
	if *devtoolsWsURL == "" {
		log.Fatal("must specify -devtools-ws-url")
	}

	// create allocator context for use with creating a browser context later
	allocatorContext, cancel := chromedp.NewRemoteAllocator(context.Background(), *devtoolsWsURL)
	defer cancel()

	// create context
	ctx, cancel := chromedp.NewContext(allocatorContext)
	defer cancel()

	//  ----------------------------------------------------------------------------------------------------------------

	// Set Context Timeout to 40 seconds
	ctx, cancel = context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	targetUrl := "https://web.whatsapp.com/"

	log.Printf("Waiting for the QRCODE")
	// Wait for the QRCODE to be visible
	qrcodeSelector := "//div[@data-testid='qrcode']"
	var buf []byte
	if err := chromedp.Run(ctx,
		chromedp.Navigate(targetUrl),
		chromedp.WaitVisible(qrcodeSelector),
		chromedp.Screenshot(qrcodeSelector, &buf, chromedp.NodeVisible),
	); err != nil {
		log.Fatalf("Couldn't get QRCODE: %v", err)
	}

	if err := os.WriteFile("elementScreenshot.png", buf, 0o644); err != nil {
		log.Fatal(err)
	}

	log.Printf("QRCODE is shown")

	log.Printf("Waiting whatsapp to start...")
	// Wait for Whatsapp Web (APP) to be visible
	appSelector := "//header[@data-testid='chatlist-header']"
	if err := chromedp.Run(ctx, chromedp.WaitVisible(appSelector)); err != nil {
		log.Fatalf("Couldn't get app: %v", err)
	}

	log.Printf("Whatsapp Started")

	// Chat Titles
	/// /div[@data-testid='chat-list']//div[@data-testid='cell-frame-title']//span/text()

	// run task list
	var body string
	if err := chromedp.Run(ctx,
		chromedp.Navigate(targetUrl),
		chromedp.WaitVisible("div._19vUU"),
		chromedp.OuterHTML("html", &body),
	); err != nil {
		log.Fatalf("Failed getting body of %v: %v", targetUrl, err)
	}

	log.Printf("Body of %v starts with: %v", targetUrl, body[0:100])
}

type Person struct {
	Name            string
	ProfilePhotoURL string
	PhoneNumber     string
}
type Category struct {
	Name        string
	Description string
	Identifiers []string
}
type Entry struct {
	Person   Person
	Category Category
	Time     time.Time
}

func listEveryPersonOnGroup(ctx context.Context, sect string) ([]Person, error) {
	// force max timeout of 15 seconds for retrieving and processing the data
	var cancel func()
	ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	//sel := fmt.Sprintf(`//p[text()[contains(., '%s')]]`, sect)
	//div[@data-testid='drawer-right']
	// Group's Member List
	//div[starts-with(@aria-label, 'Participants')]

	// Profile photo of members
	//div[starts-with(@aria-label, 'Part')]//img

	// Get Profile Name of every member
	// OBS.: To get every Member there is a need toclick "show more"
	//div[starts-with(@aria-label, 'Part')]//div[@data-testid="cell-frame-title"]//span/text()

	// Temporarily Initialize Empty Person to stop errors
	var person []Person = []Person{
		Person{
			Name:            "",
			ProfilePhotoURL: "",
			PhoneNumber:     "",
		},
	}

	return person, nil
}
