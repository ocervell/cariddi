/*
==========
Cariddi
==========

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see http://www.gnu.org/licenses/.

	@Repository:  https://github.com/ocervell/cariddi

	@Author:      edoardottt, https://www.edoardoottavianelli.it

	@License: https://github.com/ocervell/cariddi/blob/main/LICENSE

*/

package crawler

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"time"

	fileUtils "github.com/ocervell/cariddi/internal/file"
	sliceUtils "github.com/ocervell/cariddi/internal/slice"
	urlUtils "github.com/ocervell/cariddi/internal/url"
	"github.com/ocervell/cariddi/pkg/input"
	"github.com/ocervell/cariddi/pkg/output"
	"github.com/ocervell/cariddi/pkg/scanner"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
)

// New it's the actual crawler engine.
// It controls all the behaviours of a scan
// (event handlers, secrets, errors, extensions and endpoints scanning).
func New(target string, jsonl bool, txt string, html string, delayTime int, concurrency int,
	ignore string, ignoreTxt string, cache bool, timeout int, intensive bool, rua bool,
	proxy string, insecure bool, secretsFlag bool, secretsFile []string, plain bool, endpointsFlag bool,
	endpointsFile []string, fileType int, headers map[string]string, errorsFlag bool, infoFlag bool,
	debug bool, userAgent string) ([]string, []scanner.SecretMatched, []scanner.EndpointMatched,
	[]scanner.FileTypeMatched, []scanner.ErrorMatched, []scanner.InfoMatched) {
	// This is to avoid to insert into the crawler target regular
	// expression directories passed as input.
	var targetTemp, protocolTemp string

	// if there isn't a scheme use http.
	if !urlUtils.HasProtocol(target) {
		protocolTemp = "http"
		targetTemp = urlUtils.GetHost(protocolTemp + "://" + target)
	} else {
		protocolTemp = urlUtils.GetProtocol(target)
		targetTemp = urlUtils.GetHost(target)
	}

	if intensive {
		var err error
		targetTemp, err = urlUtils.GetRootHost(protocolTemp + "://" + targetTemp)

		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	}

	if targetTemp == "" {
		fmt.Println("The URL provided is not built in a proper way: " + target)
		os.Exit(1)
	}

	// clean target input
	target = urlUtils.RemoveProtocol(target)

	ignoreSlice := []string{}
	ignoreBool := false

	// if ignore -> produce the slice
	if ignore != "" {
		ignoreBool = true
		ignoreSlice = sliceUtils.CheckInputArray(ignore)
	}

	// if ignoreTxt -> produce the slice
	if ignoreTxt != "" {
		ignoreBool = true
		ignoreSlice = fileUtils.ReadFile(ignoreTxt)
	}

	FinalResults := []string{}
	FinalSecrets := []scanner.SecretMatched{}
	FinalEndpoints := []scanner.EndpointMatched{}
	FinalExtensions := []scanner.FileTypeMatched{}
	FinalErrors := []scanner.ErrorMatched{}
	FinalInfos := []scanner.InfoMatched{}

	// crawler creation
	c := CreateColly(delayTime, concurrency, cache, timeout, intensive, rua, proxy, insecure, userAgent, target)

	// On every request that Colly is making, print the URL it's currently visiting
	c.OnRequest(func(e *colly.Request) {
		if jsonl == false {
			fmt.Println(e.URL.String())
		}
	})

	// On every a element which has href attribute call callback
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		if len(link) != 0 && link[0] != '#' {
			visitHTMLLink(link, protocolTemp, targetTemp, target, intensive, ignoreBool, debug, ignoreSlice, &FinalResults, e, c)
		}
	})

	// On every script element which has src attribute call callback
	c.OnHTML("script[src]", func(e *colly.HTMLElement) {
		link := e.Attr("src")
		visitHTMLLink(link, protocolTemp, targetTemp, target, intensive, ignoreBool, debug, ignoreSlice, &FinalResults, e, c)
	})

	// On every link element which has href attribute call callback
	c.OnHTML("link[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		visitHTMLLink(link, protocolTemp, targetTemp, target, intensive, ignoreBool, debug, ignoreSlice, &FinalResults, e, c)
	})

	// On every iframe element which has src attribute call callback
	c.OnHTML("iframe[src]", func(e *colly.HTMLElement) {
		link := e.Attr("src")
		visitHTMLLink(link, protocolTemp, targetTemp, target, intensive, ignoreBool, debug, ignoreSlice, &FinalResults, e, c)
	})

	// On every svg element which has src attribute call callback
	c.OnHTML("svg[src]", func(e *colly.HTMLElement) {
		link := e.Attr("src")
		visitHTMLLink(link, protocolTemp, targetTemp, target, intensive, ignoreBool, debug, ignoreSlice, &FinalResults, e, c)
	})

	// On every img element which has src attribute call callback
	c.OnHTML("img[src]", func(e *colly.HTMLElement) {
		link := e.Attr("src")
		visitHTMLLink(link, protocolTemp, targetTemp, target, intensive, ignoreBool, debug, ignoreSlice, &FinalResults, e, c)
	})

	// On every from element which has action attribute call callback
	c.OnHTML("form[action]", func(e *colly.HTMLElement) {
		link := e.Attr("action")
		visitHTMLLink(link, protocolTemp, targetTemp, target, intensive, ignoreBool, debug, ignoreSlice, &FinalResults, e, c)
	})

	// Create a callback on the XPath query searching for the URLs
	c.OnXML("//url", func(e *colly.XMLElement) {
		link := e.Text
		visitXMLLink(link, protocolTemp, targetTemp, target, intensive, ignoreBool, debug, ignoreSlice, &FinalResults, e, c)
	})

	// Create a callback on the XPath query searching for the URLs
	c.OnXML("//link", func(e *colly.XMLElement) {
		link := e.Text
		visitXMLLink(link, protocolTemp, targetTemp, target, intensive, ignoreBool, debug, ignoreSlice, &FinalResults, e, c)
	})

	// Create a callback on the XPath query searching for the URLs
	c.OnXML("//href", func(e *colly.XMLElement) {
		link := e.Text
		visitXMLLink(link, protocolTemp, targetTemp, target, intensive, ignoreBool, debug, ignoreSlice, &FinalResults, e, c)
	})

	// Create a callback on the XPath query searching for the URLs
	c.OnXML("//loc", func(e *colly.XMLElement) {
		link := e.Text
		visitXMLLink(link, protocolTemp, targetTemp, target, intensive, ignoreBool, debug, ignoreSlice, &FinalResults, e, c)
	})

	// Create a callback on the XPath query searching for the URLs
	c.OnXML("//fileurl", func(e *colly.XMLElement) {
		link := e.Text
		visitXMLLink(link, protocolTemp, targetTemp, target, intensive, ignoreBool, debug, ignoreSlice, &FinalResults, e, c)
	})

	// Add headers (if needed) on each request
	if (len(headers)) > 0 {
		c.OnRequest(func(r *colly.Request) {
			for header, value := range headers {
				r.Headers.Set(header, value)
			}
		})
	}

	c.OnResponse(func(r *colly.Response) {
		minBodyLentgh := 10
		lengthOk := len(string(r.Body)) > minBodyLentgh
		secrets := []scanner.SecretMatched{}
		parameters := []scanner.Parameter{}
		errors := []scanner.ErrorMatched{}
		infos := []scanner.InfoMatched{}
		filetype := scanner.FileType{}

		// if endpoints or secrets or filetype: scan
		if endpointsFlag || secretsFlag || (1 <= fileType && fileType <= 7) || errorsFlag || infoFlag {
			// HERE SCAN FOR SECRETS
			if secretsFlag && lengthOk {
				secretsSlice := huntSecrets(secretsFile, r.Request.URL.String(), string(r.Body))
				for _, elem := range secretsSlice {
					FinalSecrets = append(FinalSecrets, elem)
					secrets = append(secrets, elem)
				}
			}
			// HERE SCAN FOR ENDPOINTS
			if endpointsFlag {
				endpointsSlice := huntEndpoints(endpointsFile, r.Request.URL.String())
				for _, elem := range endpointsSlice {
					if len(elem.Parameters) != 0 {
						FinalEndpoints = append(FinalEndpoints, elem)
						parameters = append(parameters, elem.Parameters...)
					}
				}
			}
			// HERE SCAN FOR EXTENSIONS
			if 1 <= fileType && fileType <= 7 {
				extension := huntExtensions(r.Request.URL.String(), fileType)
				if extension.URL != "" {
					FinalExtensions = append(FinalExtensions, extension)
					filetype = extension.Filetype
				}
			}
			// HERE SCAN FOR ERRORS
			if errorsFlag {
				errorsSlice := huntErrors(r.Request.URL.String(), string(r.Body))
				for _, elem := range errorsSlice {
					FinalErrors = append(FinalErrors, elem)
					errors = append(errors, elem)
				}
			}

			// HERE SCAN FOR INFOS
			if infoFlag {
				infosSlice := huntInfos(r.Request.URL.String(), string(r.Body))
				for _, elem := range infosSlice {
					FinalInfos = append(FinalInfos, elem)
					infos = append(infos, elem)
				}
			}
		}
		if jsonl == true {
			output.GetJsonString(r, secrets, parameters, filetype, errors, infos)
		}
	})

	// Start scraping on target
	path, err := urlUtils.GetPath(protocolTemp + "://" + target)
	if err == nil {
		var (
			addPath     string
			absoluteURL string
		)

		if path == "" {
			addPath = "/"
		}

		absoluteURL = protocolTemp + "://" + target + addPath + "robots.txt"
		if !ignoreBool || (ignoreBool && !IgnoreMatch(absoluteURL, ignoreSlice)) {
			err = c.Visit(absoluteURL)
			if err != nil && debug && !errors.Is(err, colly.ErrAlreadyVisited) {
				log.Println(err)
			}
		}

		absoluteURL = protocolTemp + "://" + target + addPath + "sitemap.xml"
		if !ignoreBool || (ignoreBool && !IgnoreMatch(absoluteURL, ignoreSlice)) {
			err = c.Visit(absoluteURL)
			if err != nil && debug && !errors.Is(err, colly.ErrAlreadyVisited) {
				log.Println(err)
			}
		}
	}

	err = c.Visit(protocolTemp + "://" + target)
	if err != nil && debug && !errors.Is(err, colly.ErrAlreadyVisited) {
		log.Println(err)
	}

	// Setup graceful exit
	chanC := make(chan os.Signal, 1)
	lettersNum := 23
	cCount := 0

	signal.Notify(chanC, os.Interrupt)
	rand.Seed(time.Now().UnixNano())

	go func() {
		for range chanC {
			if cCount > 0 {
				os.Exit(1)
			}

			if !plain {
				fmt.Fprint(os.Stdout, "\r")
				fmt.Println("CTRL+C pressed: Exiting")
				cCount++
			}

			c.AllowedDomains = []string{sliceUtils.RandSeq(lettersNum)}
		}
	}()

	c.Wait()

	if html != "" {
		output.FooterHTML(html)
	}

	return FinalResults, FinalSecrets, FinalEndpoints, FinalExtensions, FinalErrors, FinalInfos
}

// CreateColly takes as input all the settings needed to instantiate
// a new Colly Collector object and it returns this object.
func CreateColly(delayTime int, concurrency int, cache bool, timeout int,
	intensive bool, rua bool, proxy string, insecure bool, userAgent string, target string) *colly.Collector {
	c := colly.NewCollector(
		colly.Async(true),
	)
	c.IgnoreRobotsTxt = true
	c.AllowURLRevisit = false

	err := c.Limit(
		&colly.LimitRule{
			Parallelism: concurrency,
			Delay:       time.Duration(delayTime) * time.Second,
			DomainGlob:  "*" + target,
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	// Using timeout if needed
	if timeout != input.TimeoutRequest {
		c.SetRequestTimeout(time.Second * time.Duration(timeout))
	}

	// Using cache if needed
	if cache {
		c.CacheDir = ".cariddi_cache"
	}

	// Use a Random User Agent for each request if needed
	if rua {
		extensions.RandomUserAgent(c)
	} else {
		// Avoid using the default colly user agent
		c.UserAgent = GenerateRandomUserAgent()
	}

	if userAgent != "" {
		c.UserAgent = userAgent
	}

	// Use a Proxy if needed
	if proxy != "" {
		err := c.SetProxy(proxy)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	if insecure {
		c.WithTransport(&http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		})
	}

	return c
}

// visitHTMLLink checks if the collector should visit a link or not.
func visitHTMLLink(link, protocolTemp, targetTemp, target string, intensive, ignoreBool, debug bool,
	ignoreSlice []string, finalResults *[]string, e *colly.HTMLElement, c *colly.Collector) {
	if len(link) != 0 {
		absoluteURL := urlUtils.AbsoluteURL(protocolTemp, targetTemp, e.Request.AbsoluteURL(link))
		// Visit link found on page
		// Only those links are visited which are in AllowedDomains
		if (!intensive && urlUtils.SameDomain(protocolTemp+"://"+target, absoluteURL)) ||
			(intensive && intensiveOk(targetTemp, absoluteURL, debug)) {
			if !ignoreBool || (ignoreBool && !IgnoreMatch(absoluteURL, ignoreSlice)) {
				err := c.Visit(absoluteURL)
				if !errors.Is(err, colly.ErrAlreadyVisited) {
					*finalResults = append(*finalResults, absoluteURL)

					if err != nil && debug {
						log.Println(err)
					}
				}
			}
		}
	}
}

// visitXMLLink checks if the collector should visit a link or not.
func visitXMLLink(link, protocolTemp, targetTemp, target string, intensive, ignoreBool, debug bool,
	ignoreSlice []string, finalResults *[]string, e *colly.XMLElement, c *colly.Collector) {
	if len(link) != 0 {
		absoluteURL := urlUtils.AbsoluteURL(protocolTemp, targetTemp, e.Request.AbsoluteURL(link))
		// Visit link found on page
		// Only those links are visited which are in AllowedDomains
		if (!intensive && urlUtils.SameDomain(protocolTemp+"://"+target, absoluteURL)) ||
			(intensive && intensiveOk(targetTemp, absoluteURL, debug)) {
			if !ignoreBool || (ignoreBool && !IgnoreMatch(absoluteURL, ignoreSlice)) {
				err := c.Visit(absoluteURL)
				if !errors.Is(err, colly.ErrAlreadyVisited) {
					*finalResults = append(*finalResults, absoluteURL)

					if err != nil && debug {
						log.Println(err)
					}
				}
			}
		}
	}
}

// huntSecrets hunts for secrets.
func huntSecrets(secretsFile []string, target string, body string) []scanner.SecretMatched {
	secrets := SecretsMatch(target, body, secretsFile)
	return secrets
}

// SecretsMatch checks if a body matches some secrets.
func SecretsMatch(url string, body string, secretsFile []string) []scanner.SecretMatched {
	var secrets []scanner.SecretMatched

	if len(secretsFile) == 0 {
		for _, secret := range scanner.GetSecretRegexes() {
			if matched, err := regexp.Match(secret.Regex, []byte(body)); err == nil && matched {
				re := regexp.MustCompile(secret.Regex)
				match := re.FindStringSubmatch(body)

				// Avoiding false positives
				var isFalsePositive = false

				for _, falsePositive := range secret.FalsePositives {
					if strings.Contains(strings.ToLower(match[0]), falsePositive) {
						isFalsePositive = true
						break
					}
				}

				if !isFalsePositive {
					secretFound := scanner.SecretMatched{Secret: secret, URL: url, Match: match[0]}
					secrets = append(secrets, secretFound)
				}
			}
		}
	} else {
		for _, secret := range secretsFile {
			if matched, err := regexp.Match(secret, []byte(body)); err == nil && matched {
				re := regexp.MustCompile(secret)
				match := re.FindStringSubmatch(body)
				secretScanned := scanner.Secret{Name: "CustomFromFile", Description: "", Regex: secret, Poc: ""}
				secretFound := scanner.SecretMatched{Secret: secretScanned, URL: url, Match: match[0]}
				secrets = append(secrets, secretFound)
			}
		}
	}

	return secrets
}

// huntEndpoints hunts for juicy endpoints.
func huntEndpoints(endpointsFile []string, target string) []scanner.EndpointMatched {
	endpoints := EndpointsMatch(target, endpointsFile)
	return endpoints
}

// EndpointsMatch check if an endpoint matches a juicy parameter.
func EndpointsMatch(target string, endpointsFile []string) []scanner.EndpointMatched {
	endpoints := []scanner.EndpointMatched{}
	matched := []scanner.Parameter{}
	parameters := urlUtils.RetrieveParameters(target)

	if len(endpointsFile) == 0 {
		for _, parameter := range scanner.GetJuicyParameters() {
			for _, param := range parameters {
				if strings.ToLower(param) == parameter.Parameter {
					matched = append(matched, parameter)
				}
			}
		}
		endpoints = append(endpoints, scanner.EndpointMatched{Parameters: matched, URL: target})
	} else {
		for _, parameter := range endpointsFile {
			for _, param := range parameters {
				if param == parameter {
					matched = append(matched, scanner.Parameter{Parameter: parameter, Attacks: []string{}})
				}
			}
		}
		endpoints = append(endpoints, scanner.EndpointMatched{Parameters: matched, URL: target})
	}

	return endpoints
}

// huntExtensions hunts for extensions.
func huntExtensions(target string, severity int) scanner.FileTypeMatched {
	extension := scanner.FileTypeMatched{}
	copyTarget := target

	for _, ext := range scanner.GetExtensions() {
		if ext.Severity <= severity {
			firstIndex := strings.Index(target, "?")
			if firstIndex > -1 {
				target = target[:firstIndex]
			}

			if strings.ToLower(target[len(target)-len("."+ext.Extension):]) == "."+ext.Extension {
				extension = scanner.FileTypeMatched{Filetype: ext, URL: copyTarget}
			}
		}
	}

	return extension
}

// huntErrors hunts for errors.
func huntErrors(target string, body string) []scanner.ErrorMatched {
	errorsSlice := ErrorsMatch(target, body)
	return errorsSlice
}

// ErrorsMatch checks the patterns for errors.
func ErrorsMatch(url string, body string) []scanner.ErrorMatched {
	errors := []scanner.ErrorMatched{}

	for _, errorItem := range scanner.GetErrorRegexes() {
		for _, errorRegex := range errorItem.Regex {
			if matched, err := regexp.Match(errorRegex, []byte(body)); err == nil && matched {
				re := regexp.MustCompile(errorRegex)
				match := re.FindStringSubmatch(body)
				errorFound := scanner.ErrorMatched{Error: errorItem, URL: url, Match: match[0]}
				errors = append(errors, errorFound)
			}
		}
	}

	return errors
}

// huntInfos hunts for infos.
func huntInfos(target string, body string) []scanner.InfoMatched {
	infosSlice := InfoMatch(target, body)
	return infosSlice
}

// InfoMatch checks the patterns for infos.
func InfoMatch(url string, body string) []scanner.InfoMatched {
	infos := []scanner.InfoMatched{}

	for _, infoItem := range scanner.GetInfoRegexes() {
		for _, infoRegex := range infoItem.Regex {
			if matched, err := regexp.Match(infoRegex, []byte(body)); err == nil && matched {
				re := regexp.MustCompile(infoRegex)
				match := re.FindStringSubmatch(body)
				infoFound := scanner.InfoMatched{Info: infoItem, URL: url, Match: match[0]}
				infos = append(infos, infoFound)
			}
		}
	}

	return infos
}

// RetrieveBody retrieves the body (in the response) of a url.
func RetrieveBody(target string) string {
	sb, err := GetRequest(target)
	if err == nil && sb != "" {
		return sb
	}

	return ""
}

// IgnoreMatch checks if the URL should be ignored or not.
func IgnoreMatch(url string, ignoreSlice []string) bool {
	for _, ignore := range ignoreSlice {
		if strings.Contains(url, ignore) {
			return true
		}
	}

	return false
}

// intensiveOk checks if a given url can be crawled
// in intensive mode (if the 2nd level domain matches with
// the inputted target).
func intensiveOk(target string, urlInput string, debug bool) bool {
	root, err := urlUtils.GetRootHost(urlInput)
	if err != nil {
		if debug {
			fmt.Println(err.Error() + ": " + urlInput)
		}

		return false
	}

	return root == target
}
