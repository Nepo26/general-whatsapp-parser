// Command remote is a chromedp example demonstrating how to connect to an
// existing Chrome DevTools instance using a remote WebSocket URL.
package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"log"
	"os"
	"regexp"
	"time"
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

	if err := chromedp.Run(ctx,
		chromedp.Navigate(targetUrl),
	); err != nil {
		log.Fatalf("Couldn't open whatsapp web: %v", err)
	}

	if !checkIfIsLoggedIn(ctx) {
		log.Printf("Waiting for the QRCODE")
		// Wait for the QRCODE to be visible
		qrcodeSelector := "//div[@data-testid='qrcode']"
		var buf []byte
		if err := chromedp.Run(ctx,
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
		if loggedIn := checkIfIsLoggedIn(ctx); loggedIn != true {
			log.Fatal("Couldn't login to whatsapp")
		}
		log.Printf("Whatsapp Started")

	}

	// Chat Titles
	// //div[@data-testid='chat-list']//div[@data-testid='cell-frame-title']//span[starts-with(text(),"Os ex-gordinhos")]
	if err := goToChat(ctx, "Os ex-gordinhos"); err != nil {
		log.Fatal(err)
	}
	if err := clickChatInfo(ctx); err != nil {
		log.Fatal(err)
	}

	if err := getGroupMembers(ctx); err != nil {
		log.Fatal(err)
	}
}

func checkIfIsLoggedIn(ctx context.Context) bool {
	// Wait for Whatsapp Web (APP) to be visible
	log.Printf("checking if whatsapp is logged in...")
	appSelector := "//header[@data-testid='chatlist-header']"
	if err := chromedp.Run(ctx, chromedp.WaitVisible(appSelector)); err != nil {
		log.Printf("It's NOT logged in.")

		return false
	}

	log.Printf("It's logged in.")
	return true
}

func goToChat(ctx context.Context, name string) error {
	filterChatByText := fmt.Sprintf("//div[@data-testid='chat-list']//div[@data-testid='cell-frame-title']//span[starts-with(text(),'%s')]", name)

	log.Printf("click on Chat with name '%s'", name)
	if err := chromedp.Run(ctx,
		chromedp.WaitVisible(filterChatByText),
		chromedp.Click(filterChatByText, chromedp.NodeVisible),
	); err != nil {
		return fmt.Errorf("couldn't click on chat: %v", err)
	}

	return nil
}

func clickChatInfo(ctx context.Context) error {
	filterChatHeader := "//header[@data-testid='conversation-header']"

	log.Printf("Click on Chat Info")
	if err := chromedp.Run(ctx,
		chromedp.WaitVisible(filterChatHeader),
		chromedp.Click(filterChatHeader),
	); err != nil {
		return fmt.Errorf("couldn't click on header", err)
	}

	return nil
}

func getGroupMembers(ctx context.Context) error {
	filterMoreMembersButton := "//div[@data-testid='group-info-participants-section']//div[@role='button' and @data-ignore-capture='any']"

	// Get Profile Name of every member
	// OBS.: To get every Member there is a need toclick "show more"
	//filterMemberNames := "//div[@data-testid='group-info-participants-section']//div[@data-testid='cell-frame-title']//span/text()"

	// Each Profile Box
	filterMemberBox := "//div[@data-testid='group-info-participants-section']//div[@role='list']//div[@data-testid='cell-frame-container']"
	//filterMemberBox := "div[data-testid='group-info-participants-section'] div[role='list'] div[data-testid='cell-frame-container']"

	log.Printf("Getting group members")
	var memberNodes []*cdp.Node
	if err := chromedp.Run(ctx,
		chromedp.WaitVisible(filterMoreMembersButton),
		chromedp.Click(filterMoreMembersButton),
		chromedp.WaitVisible(filterMemberBox),
		chromedp.Nodes(filterMemberBox, &memberNodes, chromedp.BySearch),
	); err != nil {
		return fmt.Errorf("couldn't get to member profiles: %v", err)
	}

	log.Printf("All nodes: %v", memberNodes)
	for i := range memberNodes {
		filterUniqueMemberBox := fmt.Sprintf("(//div[@data-testid='group-info-participants-section']//div[@role='list']//div[@data-testid='cell-frame-container'])[%d]", i+1)

		profile, err := getMemberInfo(ctx, filterUniqueMemberBox)
		if err != nil {
			log.Fatalf("Couldn't get unique member: %v", err)
			return err
		}

		log.Print(profile)
	}

	return nil
}

func getMemberInfo(ctx context.Context, filter string) (Person, error) {
	filterCloseButton := "//div[@data-testid='btn-closer-drawer']"
	filterProfileName := "//span[@data-testid='contact-info-subtitle']"
	//filterProfilePhoneNumber := "//span[@data-testid='contact-info-subtitle']/../following-sibling::div/span/span"
	filterProfilePhoneNumber := "//span[@data-testid='contact-info-subtitle']/../following-sibling::div/span/span | //div[@data-testid='section-about-and-phone-number']/following-sibling::div//span/span"
	filterProfilePicture := "//span[@data-testid='contact-info-subtitle']/../../preceding-sibling::div//img | //span[@data-testid='contact-info-subtitle']/../../../..//img"
	filterMoreMembersButton := "//div[@data-testid='group-info-participants-section']//div[@role='button' and @data-ignore-capture='any']"

	var profileName string
	var profilePhoneNumber string
	var profilePictureNodes []*cdp.Node
	if err := chromedp.Run(ctx,
		chromedp.WaitVisible(filter),
		chromedp.Click(filter),

		chromedp.WaitVisible(filterProfileName),
		chromedp.Text(filterProfileName, &profileName),

		chromedp.WaitVisible(filterProfilePhoneNumber),
		chromedp.Text(filterProfilePhoneNumber, &profilePhoneNumber),

		chromedp.WaitVisible(filterProfilePicture),
		chromedp.Nodes(filterProfilePicture, &profilePictureNodes),

		chromedp.Click(filterCloseButton),
		chromedp.Click(filterMoreMembersButton),
	); err != nil {
		return Person{
			Name:            "",
			ProfilePhotoURL: "",
			PhoneNumber:     "",
		}, fmt.Errorf("couldn't get member info: %v", err)
	}
	var profilePictureUrl string = profilePictureNodes[0].AttributeValue("src")

	// Currently if a name contains a number it means that you don't have it on your contacts yet.
	if regexp.MustCompile(`\d`).MatchString(profileName) {
		return Person{
			Name:            profilePhoneNumber,
			ProfilePhotoURL: profilePictureUrl,
			PhoneNumber:     profileName,
		}, nil
	}

	return Person{
		Name:            profileName,
		ProfilePhotoURL: profilePictureUrl,
		PhoneNumber:     profilePhoneNumber,
	}, nil
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

	// Temporarily Initialize Empty Person to stop errors
	var person = []Person{
		{
			Name:            "",
			ProfilePhotoURL: "",
			PhoneNumber:     "",
		},
	}

	return person, nil
}
