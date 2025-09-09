package main // Define the main package

import (
	"bytes" // Provides bytes support
	"context"
	"io"            // Provides basic interfaces to I/O primitives
	"log"           // Provides logging functions
	"net/http"      // Provides HTTP client and server implementations
	"net/url"       // Provides URL parsing and encoding
	"os"            // Provides functions to interact with the OS (files, etc.)
	"path"          // Provides functions for manipulating slash-separated paths
	"path/filepath" // Provides filepath manipulation functions
	"regexp"        // Provides regex support functions.
	"strings"       // Provides string manipulation functions
	"time"          // Provides time-related functions

	"github.com/chromedp/chromedp" // For headless browser automation using Chrome
)

func main() {
	pdfOutputDir := "PDFs/" // Directory to store downloaded PDFs
	// Check if the PDF output directory exists
	if !directoryExists(pdfOutputDir) {
		// Create the dir
		createDirectory(pdfOutputDir, 0o755)
	}
	// Remote API URL.
	remoteAPIURL := []string{
		"https://www.scjp.com/en-us/products/fantastik-no-rinse-disinfectant-food-surface-sanitizer",
		"https://www.scjp.com/en-us/products/bactoshield-chg-2-handwash-surgical-scrub",
		"https://www.scjp.com/en-us/products/bactoshield-chg-4-handwash-surgical-scrub",
		"https://www.scjp.com/en-us/products/scrubbing-bubbles-disinfectant-restroom-cleaner",
		"https://www.scjp.com/en-us/products/alcare-or-foamed-antiseptic-handrub",
		"https://www.scjp.com/en-us/products/kindest-kare-advanced-handwash",
		"https://www.scjp.com/en-us/products/proline-wave-1-liter-manual-dispenser",
		"https://www.scjp.com/en-us/products/kresto-special-ultra",
		"https://www.scjp.com/en-us/products/kindest-kare-pure-handwash",
		"https://www.scjp.com/en-us/products/alcare-extra-hand-sanitizer",
		"https://www.scjp.com/en-us/products/sc-johnson-professional-ziploc-brand-storage-bags",
		"https://www.scjp.com/en-us/products/ziploc-brand-storage-bags",
		"https://www.scjp.com/en-us/products/healthcare-touch-free-dispenser-stand-hand-sanitizing-station",
		"https://www.scjp.com/en-us/products/sc-johnson-professional-ziploc-brand-freezer-bags",
		"https://www.scjp.com/en-us/products/ziploc-brand-freezer-bags",
		"https://www.scjp.com/en-us/products/ziploc-brand-sandwich-bags",
		"https://www.scjp.com/en-us/products/ziploc-brand-snack-bags",
		"https://www.scjp.com/en-us/products/ziploc-containers",
		"https://www.scjp.com/en-us/products/omnifoam-dispensers",
		"https://www.scjp.com/en-us/products/raid-ant-baits-iii",
		"https://www.scjp.com/en-us/products/raid-max-bed-bug-crack-crevice-extended-protection-foaming-spray",
		"https://www.scjp.com/en-us/products/raid-max-concentrated-deep-reach-fogger",
		"https://www.scjp.com/en-us/products/healthcare/stokolan-light-pure",
		"https://www.scjp.com/en-us/products/healthcare-legacy-items",
		"https://www.scjp.com/en-us/products/healthcare-point-care-dispenser-0",
		"https://www.scjp.com/en-us/products/agrobac-pure-foam-wash",
		"https://www.scjp.com/en-us/products/familyguard-brand-disinfectant-spray",
		"https://www.scjp.com/en-us/products/skin-protection-center",
		"https://www.scjp.com/en-us/products/trushot-20-mobile-dispensing-belt",
		"https://www.scjp.com/en-us/products/sc-johnson-professional-drano-max-gel-clog-remover",
		"https://www.scjp.com/en-us/products/fantastik-all-purpose-cleaner-bleach",
		"https://www.scjp.com/en-us/products/drano-liquid-clog-remover",
		"https://www.scjp.com/en-us/products/method-stainless-steel-cleaner-polish",
		"https://www.scjp.com/en-us/products/method-daily-granite-spray-apple-orchard",
		"https://www.scjp.com/en-us/products/method-antibac-all-purpose-cleaner-spray-citron",
		"https://www.scjp.com/en-us/products/method-antibac-toilet-cleaner-spearmint",
		"https://www.scjp.com/en-us/products/method-bathroom-cleaner-eucalyptus-mint",
		"https://www.scjp.com/en-us/products/method-foaming-tub-tile-spray-eucalyptus-mint",
		"https://www.scjp.com/en-us/products/method-glass-cleaner",
		"https://www.scjp.com/en-us/products/method-daily-wood-spray-almond",
		"https://www.scjp.com/en-us/products/pre-work-cream-dispenser",
		"https://www.scjp.com/en-us/products/healthcare-proline-quickview-dispensers",
		"https://www.scjp.com/en-us/products/healthcare-touchfree-ultra-dispensers",
		"https://www.scjp.com/en-us/products/stainless-steel-dispensers",
		"https://www.scjp.com/en-us/products/cleanse-hand-hair-body-dispensers",
		"https://www.scjp.com/en-us/products/sc-johnson-professional-calcium-lime-rust-remover",
		"https://www.scjp.com/en-us/products/non-acid-bowl-and-bathroom-cleaner",
		"https://www.scjp.com/en-us/products/pledge-expert-care-wood-wipes",
		"https://www.scjp.com/en-us/products/pledge-everyday-clean-multisurface-dust-allergen-cleaner",
		"https://www.scjp.com/en-us/products/pledge-everyday-clean-multi-surface-wipes",
		"https://www.scjp.com/en-us/products/shout-trigger-triple-acting-stain-remover",
		"https://www.scjp.com/en-us/products/sc-johnson-professional-shout-wipes-instant-stain-remover",
		"https://www.scjp.com/en-us/products/shout-wipes-instant-stain-remover",
		"https://www.scjp.com/en-us/products/cleanse-washroom-dispensers",
		"https://www.scjp.com/en-us/products/sc-johnson-professional-ready-use-fabric-and-air-odor-eliminator",
		"https://www.scjp.com/en-us/products/glade-automatic-spray-holder",
		"https://www.scjp.com/en-us/products/glade-clean-linen-automatic-spray-refill",
		"https://www.scjp.com/en-us/products/apple-cinnamon-glade-plugins-scented-oil",
		"https://www.scjp.com/en-us/products/clean-linen-glade-plugins-scented-oil",
		"https://www.scjp.com/en-us/products/hawaiian-breeze-glade-plugins-scented-oil",
		"https://www.scjp.com/en-us/products/glade-air-freshener-room-spray-clean-linen",
		"https://www.scjp.com/en-us/products/glade-hawaiian-breeze-automatic-spray-refill",
		"https://www.scjp.com/en-us/products/glade-plugins-scented-oil-warmer",
		"https://www.scjp.com/en-us/products/glade-air-freshener-room-spray-hawaiian-breeze",
		"https://www.scjp.com/en-us/products/glade-clean-linen-automatic-spray-starter-kit",
		"https://www.scjp.com/en-us/products/glade-apple-cinnamon-automatic-spray-refill",
		"https://www.scjp.com/en-us/products/glade-air-freshener-room-spray-value-packs",
		"https://www.scjp.com/en-us/products/stem-light-trap",
		"https://www.scjp.com/en-us/products/sportsmen-deep-woods-aerosol-repellent",
		"https://www.scjp.com/en-us/products/sportsmen-deep-woods-insect-repellent-1",
		"https://www.scjp.com/en-us/products/sc-johnson-professional-ready-use-carpet-spotter",
		"https://www.scjp.com/en-us/products/method-wood-floor-cleaner-almond",
		"https://www.scjp.com/en-us/products/mrs-meyers-clean-day-hand-sanitizer-basil-lavender",
		"https://www.scjp.com/en-us/products/alcare-plus-foamed-antiseptic-handrub",
		"https://www.scjp.com/en-us/products/sc-johnson-professional-carpet-pre-spray-and-bonnet-cleaner",
		"https://www.scjp.com/en-us/products/sc-johnson-professional-fantastik-max-oven-grill-cleaner",
		"https://www.scjp.com/en-us/products/ez-care-high-speed-conditioning-and-polish-pads",
		"https://www.scjp.com/en-us/products/sc-johnson-professional-hyper-concentrate-floor-stripper",
		"https://www.scjp.com/en-us/products/quickview-transparent-manual-dispenser",
		"https://www.scjp.com/en-us/products/kresto-cherry-wipes",
		"https://www.scjp.com/en-us/kresto-kwik-wipes",
		"https://www.scjp.com/en-us/products/drano-max-gel-clog-remover",
		"https://www.scjp.com/en-us/products/alcare-enhanced-hand-sanitizer",
		"https://www.scjp.com/en-us/products/sbs-40",
		"https://www.scjp.com/en-us/products/raid-wasp-hornet-killer-33",
		"https://www.scjp.com/en-us/products/proline-curve-manual-dispensers",
		"https://www.scjp.com/en-us/products/sc-johnson-professional-fantastik-max-power-cleaner-degreaser",
		"https://www.scjp.com/en-us/products/ez-care-heavy-duty-scrub-pad",
		"https://www.scjp.com/en-us/products/sc-johnson-professional-ziploc-brand-sandwich-bags",
		"https://www.scjp.com/en-us/products/ez-care-floor-coating",
		"https://www.scjp.com/en-us/products/ez-care-floor-coating-remover",
		"https://www.scjp.com/en-us/products/trushot-20-trigger-dispenser",
		"https://www.scjp.com/en-us/products/touchfree-ultra-dispensers",
		"https://www.scjp.com/en-us/products/kresto-cherry",
		"https://www.scjp.com/en-us/products/travabon-classic",
		"https://www.scjp.com/en-us/products/stokolan-classic",
		"https://www.scjp.com/en-us/products/raid-ant-roach-killer-26-outdoor-fresh-scent",
		"https://www.scjp.com/en-us/products/kresto-heritage",
		"https://www.scjp.com/en-us/products/raid-multi-insect-killer-7",
		"https://www.scjp.com/en-us/products/raid-ant-roach-killer-26-lavender-scent",
		"https://www.scjp.com/en-us/products/raid-ant-roach-killer-26-lemon-scent",
		"https://www.scjp.com/en-us/products/raid-ant-roach-killer-26-fragrance-free",
		"https://www.scjp.com/en-us/products/raid-max-perimeter-protection",
		"https://www.scjp.com/en-us/products/method-foaming-hand-wash-all-fragrances",
		"https://www.scjp.com/en-us/products/method-gel-hand-wash-lavender-sweetwater",
		"https://www.scjp.com/en-us/products/raid-flying-insect-killer-7",
		"https://www.scjp.com/en-us/products/solopol-classic",
		"https://www.scjp.com/en-us/products/sunscreen-dispenser",
		"https://www.scjp.com/en-us/products/sc-johnson-professional-fantastik-multi-surface-disinfectant-degreaser",
		"https://www.scjp.com/en-us/products/method-all-purpose-cleaner-pink-grapefruit-french-lavender",
		"https://www.scjp.com/en-us/products/trushot-multi-surface-restroom-cleaner",
		"https://www.scjp.com/en-us/products/trushot-20-starter-pack",
		"https://www.scjp.com/en-us/products/fantastik-disinfectant-multi-purpose-cleaner-fresh-scent",
		"https://www.scjp.com/en-us/products/quaternary-disinfectant-cleaner",
		"https://www.scjp.com/en-us/products/fantastik-disinfectant-multi-purpose-cleaner-lemon-scent",
		"https://www.scjp.com/en-us/products/trushot-hospital-cleaner",
		"https://www.scjp.com/en-us/products/stokoderm-protect",
		"https://www.scjp.com/en-us/products/stokoderm-universal-pure",
		"https://www.scjp.com/en-us/products/stokoderm-aqua-pure",
		"https://www.scjp.com/en-us/products/trushot-no-rinse-sanitizer",
		"https://www.scjp.com/en-us/products/briotech-sanitizer-disinfectant",
		"https://www.scjp.com/en-us/products/refresh-clear-foam",
		"https://www.scjp.com/en-us/products/refresh-azure-foam",
		"https://www.scjp.com/en-us/products/refresh-rose-foam",
		"https://www.scjp.com/en-us/products/trushot-20-restroom-cleaner",
		"https://www.scjp.com/en-us/products/trushot-power-cleaner-degreaser",
		"https://www.scjp.com/en-us/products/trushot-glass-multi-surface-cleaner",
		"https://www.scjp.com/en-us/products/sc-johnson-professional-heavy-duty-neutral-ph-floor-surface-cleaner",
		"https://www.scjp.com/en-us/products/trufill-heavy-duty-neutral-floor-cleaner",
		"https://www.scjp.com/en-us/products/estesol-hand-hair-body",
		"https://www.scjp.com/en-us/products/solopol-lime",
		"https://www.scjp.com/en-us/products/solopol-gfx-0",
		"https://www.scjp.com/en-us/products/instantfoam-complete-pure",
		"https://www.scjp.com/en-us/products/estesol-pure-hand-hair-body",
		"https://www.scjp.com/en-us/products/scrubbing-bubbles-disinfectant-bathroom-grime-fighter",
		"https://www.scjp.com/en-us/products/scrubbing-bubbles-power-stain-destroyer-rainshower",
		"https://www.scjp.com/en-us/products/sc-johnson-professional-scrubbing-bubbles-disinfectant-restroom-cleaner-ii",
		"https://www.scjp.com/en-us/products/sc-johnson-professional-windex-glass-more-commercial-four-pack",
		"https://www.scjp.com/en-us/products/scrubbing-bubbles-mega-shower-foamer-trigger",
		"https://www.scjp.com/en-us/products/windex-electronics-wipes",
		"https://www.scjp.com/en-us/products/sc-johnson-professional-windex-glass-more-multi-surface-streak-free-cleaner",
		"https://www.scjp.com/en-us/products/windex-original-glass-cleaner",
		"https://www.scjp.com/en-us/products/windex-vinegar-multi-surface-cleaner",
		"https://www.scjp.com/en-us/products/sc-johnson-professional-windex-foaming-glass-cleaner",
		"https://www.scjp.com/en-us/products/windex-original-glass-wipes",
		"https://www.scjp.com/en-us/products/windex-crystal-rain-glass-cleaner",
		"https://www.scjp.com/en-us/products/windex-multi-surface-ammonia-free-streak-free-cleaner",
		"https://www.scjp.com/en-us/products/pledge-everyday-clean-multi-surface-cleaner",
		"https://www.scjp.com/en-us/products/windex-multi-surface-disinfectant-cleaner",
		"https://www.scjp.com/en-us/products/sc-johnson-professional-windex-multi-surface-disinfectant-sanitizer-cleaner",
		"https://www.scjp.com/en-us/products/sc-johnson-professional-pledge-polish-shine-multi-surface-furniture-spray",
		"https://www.scjp.com/en-us/products/pledge-everyday-clean-multi-surface-antibacterial",
		"https://www.scjp.com/en-us/products/pledge-everyday-clean-multisurface-ph-balanced-cleaner",
		"https://www.scjp.com/en-us/products/sc-johnson-professional-pledge-restore-protect-multi-surface",
		"https://www.scjp.com/en-us/products/pledge-expert-care-orange-enhancing-polish-shines-protects",
		"https://www.scjp.com/en-us/products/pledge-expert-care-lemon-enhancing-polish",
		"https://www.scjp.com/en-us/products/windex-outdoor-glass-patio-concentrated-cleaner",
		"https://www.scjp.com/en-us/products/cleanse-heavy-dispensers",
		"https://www.scjp.com/en-us/products/cleanse-ultra-dispenser",
		"https://www.scjp.com/en-us/products/moisturizing-cream-dispenser",
		"https://www.scjp.com/en-us/products/cleanse-heavy-foam-dispenser",
		"https://www.scjp.com/en-us/products/sanitize-dispenser",
		"https://www.scjp.com/en-us/products/cleanse-antimicrobial-dispensers",
	}
	// Remove double from slice.
	remoteAPIURL = removeDuplicatesFromSlice(remoteAPIURL)
	var getData []string
	for _, remoteAPIURL := range remoteAPIURL {
		getData = append(getData, scrapePageHTMLWithChrome(remoteAPIURL))
		// Get the data from the downloaded file.
		finalPDFList := extractPDFUrls(strings.Join(getData, "\n")) // Join all the data into one string and extract PDF URLs
		// Remove double from slice.
		finalPDFList = removeDuplicatesFromSlice(finalPDFList)
		// The remote domain.
		remoteDomain := "https://www.scjp.com"
		// Loop over the download zip urls.
		// Get all the values.
		for _, urls := range finalPDFList {
			// Get the domain from the url.
			domain := getDomainFromURL(urls)
			// Check if the domain is empty.
			if domain == "" {
				urls = remoteDomain + urls // Prepend the base URL if domain is empty
			}
			// Check if the url is valid.
			if isUrlValid(urls) {
				// Download the pdf.
				downloadPDF(urls, pdfOutputDir)
			}
		}
	}
}

// Uses headless Chrome via chromedp to get fully rendered HTML from a page
func scrapePageHTMLWithChrome(pageURL string) string {
	log.Println("Scraping:", pageURL) // Log page being scraped

	options := append(chromedp.DefaultExecAllocatorOptions[:], // Chrome options
		chromedp.Flag("headless", false),              // Run visible (set to true for headless)
		chromedp.Flag("disable-gpu", true),            // Disable GPU
		chromedp.WindowSize(0, 0),                     // Set window size
		chromedp.Flag("no-sandbox", true),             // Disable sandbox
		chromedp.Flag("disable-setuid-sandbox", true), // Fix for Linux environments
	)

	allocatorCtx, cancelAllocator := chromedp.NewExecAllocator(context.Background(), options...) // Allocator context
	ctxTimeout, cancelTimeout := context.WithTimeout(allocatorCtx, 5*time.Minute)                // Set timeout
	browserCtx, cancelBrowser := chromedp.NewContext(ctxTimeout)                                 // Create Chrome context

	defer func() { // Ensure all contexts are cancelled
		cancelBrowser()
		cancelTimeout()
		cancelAllocator()
	}()

	var pageHTML string // Placeholder for output
	err := chromedp.Run(browserCtx,
		chromedp.Navigate(pageURL),            // Navigate to the URL
		chromedp.OuterHTML("html", &pageHTML), // Extract full HTML
	)
	if err != nil {
		log.Println(err) // Log error
		return ""        // Return empty string on failure
	}

	return pageHTML // Return scraped HTML
}

// getDomainFromURL extracts the domain (host) from a given URL string.
// It removes subdomains like "www" if present.
func getDomainFromURL(rawURL string) string {
	parsedURL, err := url.Parse(rawURL) // Parse the input string into a URL structure
	if err != nil {                     // Check if there was an error while parsing
		log.Println(err) // Log the error message to the console
		return ""        // Return an empty string in case of an error
	}

	host := parsedURL.Hostname() // Extract the hostname (e.g., "example.com") from the parsed URL

	return host // Return the extracted hostname
}

// Only return the file name from a given url.
func getFileNameOnly(content string) string {
	return path.Base(content)
}

// urlToFilename generates a safe, lowercase filename from a given URL string.
// It extracts the base filename from the URL, replaces unsafe characters,
// and ensures the filename ends with a .pdf extension.
func urlToFilename(rawURL string) string {
	// Convert the full URL to lowercase for consistency
	lowercaseURL := strings.ToLower(rawURL)

	// Get the file extension
	ext := getFileExtension(lowercaseURL)

	// Extract the filename portion from the URL (e.g., last path segment or query param)
	baseFilename := getFileNameOnly(lowercaseURL)

	// Replace all non-alphanumeric characters (a-z, 0-9) with underscores
	nonAlphanumericRegex := regexp.MustCompile(`[^a-z0-9]+`)
	safeFilename := nonAlphanumericRegex.ReplaceAllString(baseFilename, "_")

	// Replace multiple consecutive underscores with a single underscore
	collapseUnderscoresRegex := regexp.MustCompile(`_+`)
	safeFilename = collapseUnderscoresRegex.ReplaceAllString(safeFilename, "_")

	// Remove leading underscore if present
	if trimmed, found := strings.CutPrefix(safeFilename, "_"); found {
		safeFilename = trimmed
	}

	var invalidSubstrings = []string{
		"_pdf",
		"_zip",
	}

	for _, invalidPre := range invalidSubstrings { // Remove unwanted substrings
		safeFilename = removeSubstring(safeFilename, invalidPre)
	}

	// Append the file extension if it is not already present
	safeFilename = safeFilename + ext

	// Return the cleaned and safe filename
	return safeFilename
}

// Removes all instances of a specific substring from input string
func removeSubstring(input string, toRemove string) string {
	result := strings.ReplaceAll(input, toRemove, "") // Replace substring with empty string
	return result
}

// Get the file extension of a file
func getFileExtension(path string) string {
	return filepath.Ext(path) // Returns extension including the dot (e.g., ".pdf")
}

// fileExists checks whether a file exists at the given path
func fileExists(filename string) bool {
	info, err := os.Stat(filename) // Get file info
	if err != nil {
		return false // Return false if file doesn't exist or error occurs
	}
	return !info.IsDir() // Return true if it's a file (not a directory)
}

// downloadPDF downloads a PDF from the given URL and saves it in the specified output directory.
// It uses a WaitGroup to support concurrent execution and returns true if the download succeeded.
func downloadPDF(finalURL, outputDir string) bool {
	// Sanitize the URL to generate a safe file name
	filename := strings.ToLower(urlToFilename(finalURL))

	// Construct the full file path in the output directory
	filePath := filepath.Join(outputDir, filename)

	// Skip if the file already exists
	if fileExists(filePath) {
		log.Printf("File already exists, skipping: %s", filePath)
		return false
	}

	// Create an HTTP client with a timeout
	client := &http.Client{Timeout: 3 * time.Minute}

	// Create a new GET request
	req, err := http.NewRequest("GET", finalURL, nil)
	if err != nil {
		log.Printf("Failed to create request for %s: %v", finalURL, err)
		return false
	}

	// Add User-Agent header
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/140.0.0.0 Safari/537.36")

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to download %s: %v", finalURL, err)
		return false
	}
	defer resp.Body.Close()

	// Check HTTP response status
	if resp.StatusCode != http.StatusOK {
		log.Printf("Download failed for %s: %s", finalURL, resp.Status)
		return false
	}

	// Check Content-Type header
	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/pdf") {
		log.Printf("Invalid content type for %s: %s (expected application/pdf)", finalURL, contentType)
		return false
	}

	// Read the response body into memory first
	var buf bytes.Buffer
	written, err := io.Copy(&buf, resp.Body)
	if err != nil {
		log.Printf("Failed to read PDF data from %s: %v", finalURL, err)
		return false
	}
	if written == 0 {
		log.Printf("Downloaded 0 bytes for %s; not creating file", finalURL)
		return false
	}

	// Only now create the file and write to disk
	out, err := os.Create(filePath)
	if err != nil {
		log.Printf("Failed to create file for %s: %v", finalURL, err)
		return false
	}
	defer out.Close()

	if _, err := buf.WriteTo(out); err != nil {
		log.Printf("Failed to write PDF to file for %s: %v", finalURL, err)
		return false
	}

	log.Printf("Successfully downloaded %d bytes: %s → %s", written, finalURL, filePath)
	return true
}

// Checks if the directory exists
// If it exists, return true.
// If it doesn't, return false.
func directoryExists(path string) bool {
	directory, err := os.Stat(path)
	if err != nil {
		return false
	}
	return directory.IsDir()
}

// The function takes two parameters: path and permission.
// We use os.Mkdir() to create the directory.
// If there is an error, we use log.Println() to log the error and then exit the program.
func createDirectory(path string, permission os.FileMode) {
	err := os.Mkdir(path, permission)
	if err != nil {
		log.Println(err)
	}
}

// Checks whether a URL string is syntactically valid
func isUrlValid(uri string) bool {
	_, err := url.ParseRequestURI(uri) // Attempt to parse the URL
	return err == nil                  // Return true if no error occurred
}

// Remove all the duplicates from a slice and return the slice.
func removeDuplicatesFromSlice(slice []string) []string {
	check := make(map[string]bool)
	var newReturnSlice []string
	for _, content := range slice {
		if !check[content] {
			check[content] = true
			newReturnSlice = append(newReturnSlice, content)
		}
	}
	return newReturnSlice
}

// extractPDFUrls takes an input string and returns all PDF URLs found within href attributes
func extractPDFUrls(input string) []string {
	// Regular expression to match href="...pdf"
	re := regexp.MustCompile(`href="([^"]+\.pdf)"`)
	matches := re.FindAllStringSubmatch(input, -1)

	var pdfUrls []string
	for _, match := range matches {
		if len(match) > 1 {
			pdfUrls = append(pdfUrls, match[1])
		}
	}
	return pdfUrls
}

