package everyoneapi // import "samhofi.us/x/everyoneapi"

// API holds everyoneapi authentication information
type API struct {
	AccountSID string
	AuthToken  string
}

func NameField(f *[]string) {
	*f = append(*f, "name")
}

func ProfileField(f *[]string) {
	*f = append(*f, "profile")
}

func CNAMField(f *[]string) {
	*f = append(*f, "cnam")
}

func GenderField(f *[]string) {
	*f = append(*f, "gender")
}

func ImageField(f *[]string) {
	*f = append(*f, "image")
}

func AddressField(f *[]string) {
	*f = append(*f, "address")
}

func LocationField(f *[]string) {
	*f = append(*f, "location")
}

func LineProviderField(f *[]string) {
	*f = append(*f, "line_provider")
}

func CarrierField(f *[]string) {
	*f = append(*f, "carrier")
}

func OriginalCarrierField(f *[]string) {
	*f = append(*f, "carrier_o")
}

func LineTypeField(f *[]string) {
	*f = append(*f, "line_type")
}

// Result holds the results of a phone number lookup
type Result struct {
	Data    Data     `json:"data"`
	Missed  []string `json:"missed"`
	Number  string   `json:"number"`
	Note    string   `json:"note"`
	Pricing Pricing  `json:"pricing"`
	Status  bool     `json:"status"`
	Type    string   `json:"type"`
}

// Carrier holds data about the carrier that is currently providing line service
type Carrier struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// CarrierO holds data about the carrier originally assigned the phone number
type CarrierO struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// ExpandedName holds the full name returned by a phone number lookup
type ExpandedName struct {
	First string `json:"first"`
	Last  string `json:"last"`
}

// Image holds the image URLs returned by a phone number lookup. These links expire after 30 days.
type Image struct {
	Cover string `json:"cover"`
	Large string `json:"large"`
	Med   string `json:"med"`
	Small string `json:"small"`
}

// LineProvider holds the consumer facing line provider (e.g. Google Voice, or MagicJack)
type LineProvider struct {
	ID       string `json:"id"`
	MmsEmail string `json:"mms_email"`
	Name     string `json:"name"`
	SmsEmail string `json:"sms_email"`
}

// Geo holds the geographical data returned by a phone number lookup
type Geo struct {
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}

// Location holds the location data returned by a phone number lookup
type Location struct {
	City  string `json:"city"`
	Geo   Geo    `json:"geo"`
	State string `json:"state"`
	Zip   string `json:"zip"`
}

// Profile holds the profile data returned by a phone number lookup
type Profile struct {
	Edu          string `json:"edu"`
	Job          string `json:"job"`
	Relationship string `json:"relationship"`
}

// Data holds the personal info fields returned by a phone number lookup
type Data struct {
	Address      *string       `json:"address"`
	Carrier      *Carrier      `json:"carrier"`
	CarrierO     *CarrierO     `json:"carrier_o"`
	Cnam         *string       `json:"cnam"`
	ExpandedName *ExpandedName `json:"expanded_name"`
	Gender       *string       `json:"gender"`
	Image        *Image        `json:"image"`
	LineProvider *LineProvider `json:"line_provider"`
	Linetype     *string       `json:"linetype"`
	Location     *Location     `json:"location"`
	Name         *string       `json:"name"`
	Profile      *Profile      `json:"profile"`
}

// Breakdown holds the pricing breakdown of a phone number lookup
type Breakdown struct {
	Address      float64 `json:"address"`
	Carrier      float64 `json:"carrier"`
	Carrier0     float64 `json:"carrier_0"`
	Cnam         float64 `json:"cnam"`
	ExpandedName int     `json:"expanded_name"`
	Gender       float64 `json:"gender"`
	Image        float64 `json:"image"`
	LineProvider float64 `json:"line_provider"`
	Linetype     float64 `json:"linetype"`
	Location     float64 `json:"location"`
	Name         float64 `json:"name"`
	Profile      float64 `json:"profile"`
}

// Pricing holds the pricing data of a phone number lookup
type Pricing struct {
	Breakdown Breakdown `json:"breakdown"`
	Total     float64   `json:"total"`
}
