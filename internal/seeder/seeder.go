package seeder

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jadecobra/agbalumo/internal/domain"
)

type seedSource struct {
	group    string
	listings []domain.Listing
}

var allSeeds = []seedSource{
	{"Businesses", businesses},
	{"Services", services},
	{"Products", products},
	{"Jobs", jobs},
	{"Requests", requests},
	{"Food", food},
	{"Events", events},
}

// SeedAll inserts all predefined data into the repository.
func SeedAll(ctx context.Context, repo domain.ListingStore) {
	for _, source := range allSeeds {
		seedGroup(ctx, repo, source.group, source.listings)
	}
}

// EnsureSeeded checks if the database is empty, and if so, seeds it.
func EnsureSeeded(ctx context.Context, repo domain.ListingStore) {
	// Check if already seeded via FindAll
	listings, _, err := repo.FindAll(ctx, "", "", "", "", "", true, 1, 0)
	if err != nil {
		slog.Error("Failed to check existing listings", "error", err)
		return
	}

	if len(listings) == 0 {
		slog.Info("Database empty. Seeding data...")
		SeedAll(ctx, repo)
	}
}

func seedGroup(ctx context.Context, repo domain.ListingStore, name string, listings []domain.Listing) {
	slog.Info("Seeding", "group", name)
	for _, l := range listings {
		l.ID = uuid.New().String()
		l.CreatedAt = time.Now()
		l.IsActive = true
		if l.Type == domain.Request || l.Type == domain.Event {
			l.Deadline = time.Now().Add(30 * 24 * time.Hour)
		}

		if err := repo.Save(ctx, l); err != nil {
			slog.Error("Error saving listing", "title", l.Title, "error", err)
		} else {
			fmt.Printf("Saved: %s\n", l.Title)
		}
	}
}

func listing(title, origin string, lType domain.Category, anchor, desc, city, addr, contact string) domain.Listing {
	l := domain.Listing{
		Title:       title,
		OwnerOrigin: origin,
		Type:        lType,
		Anchor:      anchor,
		Description: desc,
		City:        city,
		Address:     addr,
	}
	// Basic heurestic for contact fields
	if strings.Contains(contact, "@") {
		l.ContactEmail = contact
	} else if strings.HasPrefix(contact, "+") {
		if len(contact) < 15 {
			l.ContactPhone = contact
		} else {
			l.ContactWhatsApp = contact
		}
	} else if contact != "" {
		l.WebsiteURL = contact
	}
	return l
}

// Data
var businesses = []domain.Listing{
	listing("Lagos Import Export", "Nigeria", domain.Business, "Import/Export", "Major distributor of Nigerian goods in Texas.", "Dallas", "1001 Main St, Dallas, TX 75202", "contact@lagosie.com"),
	listing("Accra Market Square", "Ghana", domain.Business, "Grocery", "Authentic Ghanaian groceries and spices.", "Plano", "2002 Legacy Dr, Plano, TX 75024", "+155510101"),
	listing("Dakar Textiles", "Senegal", domain.Business, "Textiles", "Premium fabrics imported from Senegal.", "Fort Worth", "3003 Commerce St, Fort Worth, TX 76102", "+155510102"),
	listing("Monrovia Money Transfer", "Liberia", domain.Business, "Finance", "Secure money transfer services to Liberia.", "Arlington", "4004 Cooper St, Arlington, TX 76010", "money@monrovia.com"),
	listing("Freetown Logistics", "Sierra Leone", domain.Business, "Logistics", "Shipping and freight services to West Africa.", "Irving", "5005 O'Connor Blvd, Irving, TX 75039", "+155510103"),
	listing("Bamako Braid Shop", "Mali", domain.Business, "Beauty", "Expert hair braiding and styling.", "Garland", "6006 Garland Ave, Garland, TX 75040", "+155510104"),
	listing("Cotonou Car Dealership", "Benin", domain.Business, "Automotive", "Used cars for export or local sale.", "Grand Prairie", "7007 Belt Line Rd, Grand Prairie, TX 75052", "cars@cotonou.com"),
	listing("Lome Electronics", "Togo", domain.Business, "Electronics", "Repair and sales of phones and laptops.", "Richardson", "8008 Campbell Rd, Richardson, TX 75080", "+155510105"),
	listing("Conakry Construction Co.", "Guinea", domain.Business, "Construction", "Commercial and residential construction.", "McKinney", "9009 Eldorado Pkwy, McKinney, TX 75070", "+155510106"),
	listing("Ouaga Office Supplies", "Burkina Faso", domain.Business, "Retail", "Office furniture and supplies.", "Mesquite", "1010 Town East Blvd, Mesquite, TX 75150", "supplies@ouaga.com"),
}

var services = []domain.Listing{
	listing("Yoruba Language Tutors", "Nigeria", domain.Service, "Education", "Learn to speak Yoruba fluently.", "Dallas", "1101 Ross Ave, Dallas, TX 75204", "+155520201"),
	listing("Ashanti Web Design", "Ghana", domain.Service, "Tech", "Professional websites for small businesses.", "Frisco", "1202 Dallas Pkwy, Frisco, TX 75034", "web@ashanti.com"),
	listing("Wolof Interpreters", "Senegal", domain.Service, "Translation", "Official interpretation services.", "Dallas", "1303 Pacific Ave, Dallas, TX 75201", "+155520202"),
	listing("Monrovia Legal Aid", "Liberia", domain.Service, "Legal", "Immigration and family law assistance.", "Fort Worth", "1404 Main St, Fort Worth, TX 76102", "legal@monrovia.com"),
	listing("Freetown Tax Prep", "Sierra Leone", domain.Service, "Finance", "Tax preparation for individuals and businesses.", "Arlington", "1505 Division St, Arlington, TX 76011", "+155520203"),
	listing("Bamako Moving Services", "Mali", domain.Service, "Moving", "Affordable local and long-distance moving.", "Plano", "1606 Preston Rd, Plano, TX 75093", "+155520204"),
	listing("Cotonou Cleaning Crew", "Benin", domain.Service, "Cleaning", "Residential and commercial cleaning.", "Irving", "1707 MacArthur Blvd, Irving, TX 75061", "clean@cotonou.com"),
	listing("Lome IT Support", "Togo", domain.Service, "Tech", "Computer repair and network setup.", "Richardson", "1808 Arapaho Rd, Richardson, TX 75080", "+155520205"),
	listing("Conakry Caregivers", "Guinea", domain.Service, "Health", "Home health aides for seniors.", "Garland", "1909 Northwest Hwy, Garland, TX 75041", "+155520206"),
	listing("Niamey Notary Public", "Niger", domain.Service, "Legal", "Mobile notary services.", "Mesquite", "2010 Galloway Ave, Mesquite, TX 75149", "notary@niamey.com"),
}

var products = []domain.Listing{
	listing("Handmade Shea Butter", "Ghana", domain.Product, "Beauty", "Raw, unrefined shea butter.", "Dallas", "2101 Elm St, Dallas, TX 75201", "+155530301"),
	listing("Nigerian Ankara Fabrics", "Nigeria", domain.Product, "Fashion", "Colorful Ankara prints by the yard.", "Houston", "2202 Westheimer Rd, Houston, TX 77006", "ankara@naija.com"),
	listing("Senegalese Baskets", "Senegal", domain.Product, "Home Decor", "Handwoven baskets for storage and decor.", "Plano", "2303 Park Blvd, Plano, TX 75074", "+155530302"),
	listing("Liberian Coffee Beans", "Liberia", domain.Product, "Food", "Robusta coffee beans direct from Monrovia.", "Fort Worth", "2404 7th St, Fort Worth, TX 76107", "coffee@liberia.com"),
	listing("Sierra Leone Diamonds (Art.)", "Sierra Leone", domain.Product, "Jewelry", "Artificial diamond jewelry.", "Arlington", "2505 Collins St, Arlington, TX 76011", "+155530303"),
	listing("Malian Mud Cloth", "Mali", domain.Product, "Fashion", "Authentic Bogolanfini fabric.", "Irving", "2606 Irving Blvd, Irving, TX 75060", "+155530304"),
	listing("Benin Bronze Art", "Benin", domain.Product, "Art", "Replica bronze statues and art.", "Dallas", "2707 Flora St, Dallas, TX 75201", "art@benin.com"),
	listing("Togo Cocoa Butter", "Togo", domain.Product, "Beauty", "Pure cocoa butter for skin care.", "Garland", "2808 Centerville Rd, Garland, TX 75043", "+155530305"),
	listing("Guinea Pepper Spice", "Guinea", domain.Product, "Food", "Spicy peppers for cooking.", "Grand Prairie", "2909 Main St, Grand Prairie, TX 75050", "+155530306"),
	listing("Niger Leather Bags", "Niger", domain.Product, "Fashion", "Handcrafted leather goods.", "McKinney", "3010 Virginia Pkwy, McKinney, TX 75071", "bags@niger.com"),
	listing("Ivory Coast Cocoa Powder", "Cote d'Ivoire", domain.Product, "Food", "Premium cocoa powder for baking.", "Frisco", "3111 Main St, Frisco, TX 75033", "+155530307"),
}

var jobs = []domain.Listing{
	listing("Experienced Chef Needed", "Nigeria", domain.Job, "Restaurant", "Looking for a chef specializing in Nigerian cuisine.", "Dallas", "3201 Greenville Ave, Dallas, TX 75206", "jobs@chef.com"),
	listing("Sales Associate", "Ghana", domain.Job, "Retail", "Part-time sales associate for fabric store.", "Plano", "3302 Parker Rd, Plano, TX 75023", "+155540401"),
	listing("Delivery Driver", "Senegal", domain.Job, "Logistics", "Driver needed for local deliveries. Must have own van.", "Fort Worth", "3403 Main St, Fort Worth, TX 76102", "+155540402"),
	listing("Hair Stylist", "Liberia", domain.Job, "Beauty", "Braider needed for busy salon.", "Arlington", "3504 Cooper St, Arlington, TX 76015", "salon@jobs.com"),
	listing("Java Developer", "Sierra Leone", domain.Job, "Tech", "Junior Java developer for startup.", "Irving", "3605 Las Colinas Blvd, Irving, TX 75039", "+155540403"),
	listing("Event Planner Assistant", "Mali", domain.Job, "Events", "Help organize engaging community events.", "Dallas", "3706 Lemmon Ave, Dallas, TX 75219", "+155540404"),
	listing("Auto Mechanic", "Benin", domain.Job, "Automotive", "Experienced mechanic for Toyota/Honda shop.", "Garland", "3807 Forest Ln, Garland, TX 75042", "mechanic@jobs.com"),
	listing("French Teacher", "Togo", domain.Job, "Education", "Part-time French teacher for kids.", "Richardson", "3908 Coit Rd, Richardson, TX 75080", "+155540405"),
	listing("Construction Worker", "Guinea", domain.Job, "Construction", "General labor for construction site.", "McKinney", "4009 Central Expy, McKinney, TX 75070", "+155540406"),
	listing("Accountant", "Burkina Faso", domain.Job, "Finance", "Bookkeeper needed for small business.", "Mesquite", "4110 Motley Dr, Mesquite, TX 75150", "finance@jobs.com"),
}

var requests = []domain.Listing{
	listing("Looking for Zobo leaves", "Nigeria", domain.Request, "Sourcing", "Need bulk Zobo (Hibiscus) leaves for an event.", "Dallas", "4201 Ross Ave, Dallas, TX 75204", "req1@example.com"),
	listing("Graphic Designer needed", "Ghana", domain.Request, "Services", "Need a logo for a new startup.", "Carrollton", "4302 Josey Ln, Carrollton, TX 75006", "+155550501"),
	listing("Anyone flying to Senegal?", "Senegal", domain.Request, "Travel", "Need to send a small document package urgently.", "Euless", "4403 Airport Fwy, Euless, TX 76040", "req2@example.com"),
	listing("French Tutor for kids", "Cote d'Ivoire", domain.Request, "Education", "Looking for a tutor for my 2 kids.", "Suburbs", "4504 Main St, Lewisville, TX 75067", "+155550502"),
	listing("Apartment to rent", "Gambia", domain.Request, "Housing", "Looking for a 1 bedroom near the university.", "Denton", "4605 University Dr, Denton, TX 76201", "+155550503"),
	listing("Used Toyota Camry", "Togo", domain.Request, "Auto", "Looking to buy a reliable commuter car.", "Grapevine", "4706 Main St, Grapevine, TX 76051", "req3@example.com"),
	listing("Web Developer partner", "Niger", domain.Request, "Tech", "Seeking a technical co-founder.", "Irving", "4807 MacArthur Blvd, Irving, TX 75061", "+155550504"),
	listing("Caterer for Naming Ceremony", "Sierra Leone", domain.Request, "Events", "Need catering for 50 people next Saturday.", "Allen", "4908 McDermott Dr, Allen, TX 75013", "+155550505"),
	listing("Roommate Wanted", "Liberia", domain.Request, "Housing", "Share a 2-bedroom apartment in Arlington.", "Arlington", "5009 Cooper St, Arlington, TX 76017", "req4@example.com"),
	listing("Shipping Barrel to Mali", "Mali", domain.Request, "Logistics", "Who has the best rates for barrels?", "Dallas", "5110 Buckner Blvd, Dallas, TX 75228", "+155550506"),
}

var food = []domain.Listing{
	listing("Mama Put Dallas", "Nigeria", domain.Food, "Restaurant", "Home cooked meals: Pounded yam, Egusi, Suya.", "Dallas", "5201 Forest Ln, Dallas, TX 75243", "mama@put.com"),
	listing("Accra Chop Bar", "Ghana", domain.Food, "Restaurant", "Best Banku and Tilapia in town.", "Arlington", "5302 Collins St, Arlington, TX 76014", "+155560601"),
	listing("Dakar Delights Catering", "Senegal", domain.Food, "Catering", "Catering for weddings and parties. Thieboudienne specialty.", "Plano", "5403 Legacy Dr, Plano, TX 75024", "+155560602"),
	listing("Monrovia Bakery", "Liberia", domain.Food, "Bakery", "Fresh rice bread and cassava cake daily.", "Fort Worth", "5504 Berry St, Fort Worth, TX 76110", "bread@monrovia.com"),
	listing("Freetown Kitchen", "Sierra Leone", domain.Food, "Restaurant", "Cassava leaves and Plasas like home.", "Garland", "5605 Broadway Blvd, Garland, TX 75043", "+155560603"),
	listing("Bamako Grill", "Mali", domain.Food, "Restaurant", "Grilled fish and chicken with couscous.", "Irving", "5706 Belt Line Rd, Irving, TX 75063", "+155560604"),
	listing("Cotonou Cuisine", "Benin", domain.Food, "Restaurant", "Authentic Beninese flavors.", "Mesquite", "5807 Gus Thomasson Rd, Mesquite, TX 75150", "food@cotonou.com"),
	listing("Lome Snacks", "Togo", domain.Food, "Snacks", "Street food snacks and drinks.", "Richardson", "5908 Plano Rd, Richardson, TX 75081", "+155560605"),
	listing("Conakry Cafe", "Guinea", domain.Food, "Cafe", "Coffee and pastries with a Guinean twist.", "Lewisville", "6009 Main St, Lewisville, TX 75067", "+155560606"),
	listing("Ouaga Spicy Pot", "Burkina Faso", domain.Food, "Restaurant", "Spicy stews and grilled meats.", "Grand Prairie", "6110 Carrier Pkwy, Grand Prairie, TX 75052", "spicy@ouaga.com"),
}

var events = []domain.Listing{
	listing("Naija Independence Day Bash", "Nigeria", domain.Event, "Festival", "Celebrating Nigeria @ 66. Food, Music, Dance.", "Dallas", "6201 Fair Park, Dallas, TX 75210", "indep@naija.com"),
	listing("Ghana Fest Texas", "Ghana", domain.Event, "Festival", "Cultural showcase of Ghanaian heritage.", "Arlington", "6302 Levitt Pavilion, Arlington, TX 76010", "+155570701"),
	listing("Senegal Music Night", "Senegal", domain.Event, "Concert", "Live Mbalax music featuring top artists.", "Dallas", "6403 Deep Ellum, Dallas, TX 75226", "+155570702"),
	listing("Liberian Community Picnic", "Liberia", domain.Event, "Community", "Family fun day at the park.", "Fort Worth", "6504 Trinity Park, Fort Worth, TX 76107", "picnic@liberia.com"),
	listing("Sierra Leone Charity Gala", "Sierra Leone", domain.Event, "Charity", "Fundraiser for schools in Freetown.", "Plano", "6605 Legacy Dr, Plano, TX 75024", "+155570703"),
	listing("Mali Art Exhibition", "Mali", domain.Event, "Art", "Showcasing contemporary Malian artists.", "Dallas", "6706 Flora St, Dallas, TX 75201", "+155570704"),
	listing("Benin Cultural Day", "Benin", domain.Event, "Culture", "Dance performances and food tasting.", "Irving", "6807 Las Colinas Blvd, Irving, TX 75039", "culture@benin.com"),
	listing("Togo Independence Party", "Togo", domain.Event, "Party", "Music and dancing all night.", "Garland", "6908 Naaman Forest Blvd, Garland, TX 75040", "+155570705"),
	listing("Guinea Fashion Week", "Guinea", domain.Event, "Fashion", "Fashion show featuring Guinean designers.", "Dallas", "7009 Market Center Blvd, Dallas, TX 75207", "+155570706"),
	listing("Burkina Faso Film Screening", "Burkina Faso", domain.Event, "Film", "Screening of award-winning FESPACO films.", "Richardson", "7110 Alamo Drafthouse, Richardson, TX 75080", "film@burkina.com"),
}
