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
	RegisterAltGenNames("bool", "boolean")
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

// ValueGenerator represents a function to generate a piece of random data
type ValueGenerator func(GeneratorOpts) (interface{}, error)

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
	gens  map[string]ValueGenerator
}

func (g *generators) makeGenerators() {
	g.gens = map[string]ValueGenerator{
		"characters":        g.intString(g.faker.Characters),
		"noun":              g.string(randomdata.Noun),
		"adjective":         g.string(randomdata.Noun),
		"word":              g.string(func() string { return g.faker.Words(1, false)[0] }),
		"words":             g.intBoolStrings(g.faker.Words),
		"sentence":          g.intBoolString(g.faker.Sentence),
		"sentences":         g.intBoolStrings(g.faker.Sentences),
		"paragraph":         g.intBoolString(g.faker.Paragraph),
		"paragraphs":        g.intBoolStrings(g.faker.Paragraphs),
		"city":              g.string(g.faker.City),
		"street-name":       g.string(g.faker.StreetName),
		"street-address":    g.string(g.faker.StreetAddress),
		"secondary-address": g.string(g.faker.SecondaryAddress),
		"postcode":          g.string(g.faker.PostCode),
		"street-suffix":     g.string(g.faker.StreetSuffix),
		"city-suffix":       g.string(g.faker.CitySuffix),
		"city-prefix":       g.string(g.faker.CityPrefix),
		"state":             g.string(g.faker.StateAbbr),
		"state-name":        g.string(g.faker.State),
		"country":           g.string(g.faker.Country),
		"latitude":          g.float(g.faker.Latitude),
		"longitude":         g.float(g.faker.Longitude),
		"company":           g.string(g.faker.CompanyName),
		"company-suffix":    g.string(g.faker.CompanySuffix),
		"company-slogan":    g.string(g.faker.CompanyCatchPhrase),
		"company-bs":        g.string(g.faker.CompanyBs),
		"landline":          g.string(g.faker.PhoneNumber),
		"mobile":            g.string(g.faker.CellPhoneNumber),
		"email":             g.string(g.faker.Email),
		"free-email":        g.string(g.faker.FreeEmail),
		"safe-email":        g.string(g.faker.SafeEmail),
		"user-name":         g.string(g.faker.UserName),
		"hostname":          g.string(g.faker.DomainWord),
		"domain":            g.string(g.faker.DomainName),
		"domain-suffix":     g.string(g.faker.DomainSuffix),
		"ipv4":              g.string(randomdata.IpV4Address),
		"ipv6":              g.string(randomdata.IpV6Address),
		"ip":                g.altws(randomdata.IpV4Address, randomdata.IpV6Address),
		"name":              g.string(g.faker.Name),
		"silly-name":        g.string(randomdata.SillyName),
		"first-name":        g.string(g.faker.FirstName),
		"last-name":         g.string(g.faker.LastName),
		"name-prefix":       g.string(g.faker.NamePrefix),
		"name-suffix":       g.string(g.faker.NameSuffix),
		"job-title":         g.string(g.faker.JobTitle),
		"credit-card":       g.fromPattern(govalidator.CreditCard),
		"isbn":              g.altwsp(govalidator.ISBN10, govalidator.ISBN13),
		"isbn10":            g.fromPattern(govalidator.ISBN10),
		"isbn13":            g.fromPattern(govalidator.ISBN13),
		"ssn":               g.fromPattern(govalidator.SSN),
		"hexcolor":          g.fromPattern(govalidator.Hexcolor),
		"rgbcolor":          g.fromPattern(govalidator.RGBcolor),
		"mac-address":       g.fromPattern("^([0-9A-Fa-f]{2}[:]){5}([0-9A-Fa-f]{2})$"),
		"uuid":              g.fromPattern(strfmt.UUIDPattern),
		"uuid3":             g.fromPattern(strfmt.UUID3Pattern),
		"uuid4":             g.fromPattern(strfmt.UUID4Pattern),
		"uuid5":             g.fromPattern(strfmt.UUID5Pattern),
		"bool":              g.bool,
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

func (g *generators) For(opts GeneratorOpts) (ValueGenerator, bool) {
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

func (g *generators) altws(fns ...func() string) ValueGenerator {
	return func(opts GeneratorOpts) (interface{}, error) {
		idx := seedAndReturnRandom(len(fns))
		return fns[idx](), nil
	}
}

func (g *generators) altwsp(patterns ...string) ValueGenerator {
	return func(opts GeneratorOpts) (interface{}, error) {
		idx := seedAndReturnRandom(len(patterns))
		return regen.Generate(patterns[idx])
	}
}

func (g *generators) fromPattern(pattern string) ValueGenerator {
	return func(opts GeneratorOpts) (interface{}, error) {
		return regen.Generate(pattern)
	}
}

func (g *generators) stringError(fn func() (string, error)) ValueGenerator {
	return func(opts GeneratorOpts) (interface{}, error) {
		return fn()
	}
}

func (g *generators) string(fn func() string) ValueGenerator {
	return func(opts GeneratorOpts) (interface{}, error) {
		return fn(), nil
	}
}

func (g *generators) stringer(fn func() fmt.Stringer) ValueGenerator {
	return func(opts GeneratorOpts) (interface{}, error) {
		return fn().String(), nil
	}
}

func (g *generators) float(fn func() float64) ValueGenerator {
	return func(opts GeneratorOpts) (interface{}, error) {
		return fn(), nil
	}
}

func (g *generators) intString(fn func(int) string) ValueGenerator {
	return func(opts GeneratorOpts) (interface{}, error) {
		args := opts.Args()
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

func (g *generators) intBoolString(fn func(int, bool) string) ValueGenerator {
	return func(opts GeneratorOpts) (interface{}, error) {
		args := opts.Args()
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

func (g *generators) intBoolStrings(fn func(int, bool) []string) ValueGenerator {
	return func(opts GeneratorOpts) (interface{}, error) {
		args := opts.Args()
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

func (g *generators) bool(opts GeneratorOpts) (interface{}, error) {
	answer := seedAndReturnRandom(2) == 2
	return answer, nil
}
