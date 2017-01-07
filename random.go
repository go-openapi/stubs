package stubs

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	randomdata "github.com/Pallinder/go-randomdata"
	"github.com/asaskevich/govalidator"
	conv "github.com/cstockton/go-conv"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/manveru/faker"
	regen "github.com/zach-klippenstein/goregen"
)

var (
	generatorAliases map[string]string
)

func init() {
	RegisterAltGenNames("state", "state-code")
	RegisterAltGenNames("country", "country-name")
	RegisterAltGenNames("latitude", "lat")
	RegisterAltGenNames("longitude", "lon")
	RegisterAltGenNames("company", "company-name")
	RegisterAltGenNames("company-slogan", "company-catch-phrase")
	RegisterAltGenNames("company-bs", "company-mission")
	RegisterAltGenNames("landline", "phone-number")
	RegisterAltGenNames("mobile", "mobile-number", "cell", "cell-phone", "gsm", "gsm-number")
	RegisterAltGenNames("user-name", "username", "login", "nickname", "nick-name")
	RegisterAltGenNames("hostname", "domainword", "domain-word", "host", "host-name")
	RegisterAltGenNames("domain", "domain-name")
	RegisterAltGenNames("credit-card", "creditcard")
	RegisterAltGenNames("ssn", "socialsecurity", "social-security", "social-security-number")
	RegisterAltGenNames("hexcolor", "hex-color", "hexcolour", "hex-colour")
	RegisterAltGenNames("rgbcolor", "rgb-color", "rgbcolour", "rgb-colour")
	RegisterAltGenNames("mac-address", "mac", "macaddress")
	RegisterAltGenNames("ipv4", "ip4")
	RegisterAltGenNames("ipv6", "ip6")
	RegisterAltGenNames("isbn10", "isbnv10")
	RegisterAltGenNames("isbn13", "isbnv13")
	RegisterAltGenNames("uuid4", "uuidv4")
	RegisterAltGenNames("uuid3", "uuidv3")
	RegisterAltGenNames("uuid5", "uuidv5")
}

// RegisterAltGenNames registers alternatives for a generator name
func RegisterAltGenNames(key string, alts ...string) {
	if generatorAliases == nil {
		generatorAliases = make(map[string]string, 300)
	}
	for _, v := range alts {
		generatorAliases[v] = key
	}
}

// DataGenerator represents a function to generate a piece of random data
type DataGenerator func(...interface{}) (interface{}, error)

func newGenerator(lang string) (*generators, error) {
	if lang == "" {
		lang = "en"
	}
	faker, err := faker.New(lang)
	if err != nil {
		return nil, err
	}
	g := &generators{
		faker: faker,
		conv:  conv.Conv{},
	}
	g.makeGenerators()
	return g, nil
}

type generators struct {
	faker *faker.Faker
	conv  conv.Converter
	gens  map[string]DataGenerator
}

func (g *generators) makeGenerators() {
	g.gens = map[string]DataGenerator{
		"characters":        g.wi(g.faker.Characters),
		"noun":              g.ws(randomdata.Noun),
		"adjective":         g.ws(randomdata.Noun),
		"word":              g.ws(func() string { return g.faker.Words(1, false)[0] }),
		"words":             g.wibs(g.faker.Words),
		"sentence":          g.wib(g.faker.Sentence),
		"sentences":         g.wibs(g.faker.Sentences),
		"paragraph":         g.wib(g.faker.Paragraph),
		"paragraphs":        g.wibs(g.faker.Paragraphs),
		"city":              g.ws(g.faker.City),
		"street-name":       g.ws(g.faker.StreetName),
		"street-address":    g.ws(g.faker.StreetAddress),
		"secondary-address": g.ws(g.faker.SecondaryAddress),
		"postcode":          g.ws(g.faker.PostCode),
		"street-suffix":     g.ws(g.faker.StreetSuffix),
		"city-suffix":       g.ws(g.faker.CitySuffix),
		"city-prefix":       g.ws(g.faker.CityPrefix),
		"state":             g.ws(g.faker.StateAbbr),
		"state-name":        g.ws(g.faker.State),
		"country":           g.ws(g.faker.Country),
		"latitude":          g.wf(g.faker.Latitude),
		"longitude":         g.wf(g.faker.Longitude),
		"company":           g.ws(g.faker.CompanyName),
		"company-suffix":    g.ws(g.faker.CompanySuffix),
		"company-slogan":    g.ws(g.faker.CompanyCatchPhrase),
		"company-bs":        g.ws(g.faker.CompanyBs),
		"landline":          g.ws(g.faker.PhoneNumber),
		"mobile":            g.ws(g.faker.CellPhoneNumber),
		"email":             g.ws(g.faker.Email),
		"free-email":        g.ws(g.faker.FreeEmail),
		"safe-email":        g.ws(g.faker.SafeEmail),
		"user-name":         g.ws(g.faker.UserName),
		"hostname":          g.ws(g.faker.DomainWord),
		"domain":            g.ws(g.faker.DomainName),
		"domain-suffix":     g.ws(g.faker.DomainSuffix),
		"ipv4":              g.ws(randomdata.IpV4Address),
		"ipv6":              g.ws(randomdata.IpV6Address),
		"ip":                g.altws(randomdata.IpV4Address, randomdata.IpV6Address),
		"name":              g.ws(g.faker.Name),
		"silly-name":        g.ws(randomdata.SillyName),
		"first-name":        g.ws(g.faker.FirstName),
		"last-name":         g.ws(g.faker.LastName),
		"name-prefix":       g.ws(g.faker.NamePrefix),
		"name-suffix":       g.ws(g.faker.NameSuffix),
		"job-title":         g.ws(g.faker.JobTitle),
		"credit-card":       g.wsp(govalidator.CreditCard),
		"isbn":              g.altwsp(govalidator.ISBN10, govalidator.ISBN13),
		"isbn10":            g.wsp(govalidator.ISBN10),
		"isbn13":            g.wsp(govalidator.ISBN13),
		"ssn":               g.wsp(govalidator.SSN),
		"hexcolor":          g.wsp(govalidator.Hexcolor),
		"rgbcolor":          g.wsp(govalidator.RGBcolor),
		"mac-address":       g.wsp("^([0-9A-Fa-f]{2}[:]){5}([0-9A-Fa-f]{2})$"),
		"uuid":              g.wsp(strfmt.UUIDPattern),
		"uuid3":             g.wsp(strfmt.UUID3Pattern),
		"uuid4":             g.wsp(strfmt.UUID4Pattern),
		"uuid5":             g.wsp(strfmt.UUID5Pattern),
	}

	/* TODO:
	* add date
	* add date-time
	* add duration
	* add integers
	* add decimals
	* add slices
	 */
}

func normalizeGeneratorName(str string) string {
	kn := strings.ToLower(str)
	if k, ok := generatorAliases[kn]; ok {
		return k
	}
	return kn
}

func (g *generators) For(opts GeneratorOpts) (DataGenerator, bool) {
	if gen, ok := g.gens[normalizeGeneratorName(opts.Name())]; ok {
		return gen, true
	}
	if gen, ok := g.gens[normalizeGeneratorName(swag.ToCommandName(opts.FieldName()))]; ok {
		return gen, true
	}
	return nil, false
}

func seedAndReturnRandom(n int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(n)
}

func (g *generators) altws(fns ...func() string) DataGenerator {
	return func(_ ...interface{}) (interface{}, error) {
		idx := seedAndReturnRandom(len(fns))
		return fns[idx](), nil
	}
}

func (g *generators) altwsp(patterns ...string) DataGenerator {
	return func(_ ...interface{}) (interface{}, error) {
		idx := seedAndReturnRandom(len(patterns))
		return regen.Generate(patterns[idx])
	}
}

func (g *generators) wsp(pattern string) DataGenerator {
	return func(_ ...interface{}) (interface{}, error) {
		return regen.Generate(pattern)
	}
}

func (g *generators) wse(fn func() (string, error)) DataGenerator {
	return func(_ ...interface{}) (interface{}, error) {
		return fn()
	}
}

func (g *generators) ws(fn func() string) DataGenerator {
	return func(_ ...interface{}) (interface{}, error) {
		return fn(), nil
	}
}

func (g *generators) wss(fn func() fmt.Stringer) DataGenerator {
	return func(_ ...interface{}) (interface{}, error) {
		return fn().String(), nil
	}
}

func (g *generators) wf(fn func() float64) DataGenerator {
	return func(_ ...interface{}) (interface{}, error) {
		return fn(), nil
	}
}

func (g *generators) wi(fn func(int) string) DataGenerator {
	return func(args ...interface{}) (interface{}, error) {
		count := 10
		if len(args) > 0 {
			i, err := g.conv.Int(args[0])
			if err != nil {
				return nil, err
			}
			count = i
		}

		return fn(count), nil
	}
}

func (g *generators) wib(fn func(int, bool) string) DataGenerator {
	return func(args ...interface{}) (interface{}, error) {
		count := 10
		if len(args) > 0 {
			i, err := g.conv.Int(args[0])
			if err != nil {
				return nil, err
			}
			count = i
		}
		var supplemental bool
		if len(args) > 1 {
			b, err := g.conv.Bool(args[1])
			if err != nil {
				return nil, err
			}
			supplemental = b
		}

		return fn(count, supplemental), nil
	}
}

func (g *generators) wibs(fn func(int, bool) []string) DataGenerator {
	return func(args ...interface{}) (interface{}, error) {
		count := 10
		if len(args) > 0 {
			i, err := g.conv.Int(args[0])
			if err != nil {
				return nil, err
			}
			count = i
		}
		var supplemental bool
		if len(args) > 1 {
			b, err := g.conv.Bool(args[1])
			if err != nil {
				return nil, err
			}
			supplemental = b
		}

		return fn(count, supplemental), nil
	}
}
