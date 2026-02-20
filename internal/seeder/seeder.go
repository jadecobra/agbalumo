package seeder

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jadecobra/agbalumo/internal/domain"
)

// SeedAll inserts all predefined data into the repository.
func SeedAll(ctx context.Context, repo domain.ListingStore) {
	seedGroup(ctx, repo, "Businesses", businesses)
	seedGroup(ctx, repo, "Services", services)
	seedGroup(ctx, repo, "Products", products)
	seedGroup(ctx, repo, "Jobs", jobs)
	seedGroup(ctx, repo, "Requests", requests)
	seedGroup(ctx, repo, "Food", food)
	seedGroup(ctx, repo, "Events", events)
}

// EnsureSeeded checks if the database is empty, and if so, seeds it.
func EnsureSeeded(ctx context.Context, repo domain.ListingStore) {
	listings, err := repo.FindAll(ctx, "", "", true, 1, 0)
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

// Data
var businesses = []domain.Listing{
	{Title: "Lagos Import Export", OwnerOrigin: "Nigeria", Type: domain.Business, Anchor: "Import/Export", Description: "Major distributor of Nigerian goods in Texas.", City: "Dallas", Address: "1001 Main St, Dallas, TX 75202", ContactEmail: "contact@lagosie.com"},
	{Title: "Accra Market Square", OwnerOrigin: "Ghana", Type: domain.Business, Anchor: "Grocery", Description: "Authentic Ghanaian groceries and spices.", City: "Plano", Address: "2002 Legacy Dr, Plano, TX 75024", ContactWhatsApp: "+155510101"},
	{Title: "Dakar Textiles", OwnerOrigin: "Senegal", Type: domain.Business, Anchor: "Textiles", Description: "Premium fabrics imported from Senegal.", City: "Fort Worth", Address: "3003 Commerce St, Fort Worth, TX 76102", ContactPhone: "+155510102"},
	{Title: "Monrovia Money Transfer", OwnerOrigin: "Liberia", Type: domain.Business, Anchor: "Finance", Description: "Secure money transfer services to Liberia.", City: "Arlington", Address: "4004 Cooper St, Arlington, TX 76010", ContactEmail: "money@monrovia.com"},
	{Title: "Freetown Logistics", OwnerOrigin: "Sierra Leone", Type: domain.Business, Anchor: "Logistics", Description: "Shipping and freight services to West Africa.", City: "Irving", Address: "5005 O'Connor Blvd, Irving, TX 75039", ContactWhatsApp: "+155510103"},
	{Title: "Bamako Braid Shop", OwnerOrigin: "Mali", Type: domain.Business, Anchor: "Beauty", Description: "Expert hair braiding and styling.", City: "Garland", Address: "6006 Garland Ave, Garland, TX 75040", ContactPhone: "+155510104"},
	{Title: "Cotonou Car Dealership", OwnerOrigin: "Benin", Type: domain.Business, Anchor: "Automotive", Description: "Used cars for export or local sale.", City: "Grand Prairie", Address: "7007 Belt Line Rd, Grand Prairie, TX 75052", ContactEmail: "cars@cotonou.com"},
	{Title: "Lome Electronics", OwnerOrigin: "Togo", Type: domain.Business, Anchor: "Electronics", Description: "Repair and sales of phones and laptops.", City: "Richardson", Address: "8008 Campbell Rd, Richardson, TX 75080", ContactWhatsApp: "+155510105"},
	{Title: "Conakry Construction Co.", OwnerOrigin: "Guinea", Type: domain.Business, Anchor: "Construction", Description: "Commercial and residential construction.", City: "McKinney", Address: "9009 Eldorado Pkwy, McKinney, TX 75070", ContactPhone: "+155510106"},
	{Title: "Ouaga Office Supplies", OwnerOrigin: "Burkina Faso", Type: domain.Business, Anchor: "Retail", Description: "Office furniture and supplies.", City: "Mesquite", Address: "1010 Town East Blvd, Mesquite, TX 75150", ContactEmail: "supplies@ouaga.com"},
}

var services = []domain.Listing{
	{Title: "Yoruba Language Tutors", OwnerOrigin: "Nigeria", Type: domain.Service, Anchor: "Education", Description: "Learn to speak Yoruba fluently.", City: "Dallas", Address: "1101 Ross Ave, Dallas, TX 75204", ContactWhatsApp: "+155520201"},
	{Title: "Ashanti Web Design", OwnerOrigin: "Ghana", Type: domain.Service, Anchor: "Tech", Description: "Professional websites for small businesses.", City: "Frisco", Address: "1202 Dallas Pkwy, Frisco, TX 75034", ContactEmail: "web@ashanti.com"},
	{Title: "Wolof Interpreters", OwnerOrigin: "Senegal", Type: domain.Service, Anchor: "Translation", Description: "Official interpretation services.", City: "Dallas", Address: "1303 Pacific Ave, Dallas, TX 75201", ContactPhone: "+155520202"},
	{Title: "Monrovia Legal Aid", OwnerOrigin: "Liberia", Type: domain.Service, Anchor: "Legal", Description: "Immigration and family law assistance.", City: "Fort Worth", Address: "1404 Main St, Fort Worth, TX 76102", ContactEmail: "legal@monrovia.com"},
	{Title: "Freetown Tax Prep", OwnerOrigin: "Sierra Leone", Type: domain.Service, Anchor: "Finance", Description: "Tax preparation for individuals and businesses.", City: "Arlington", Address: "1505 Division St, Arlington, TX 76011", ContactWhatsApp: "+155520203"},
	{Title: "Bamako Moving Services", OwnerOrigin: "Mali", Type: domain.Service, Anchor: "Moving", Description: "Affordable local and long-distance moving.", City: "Plano", Address: "1606 Preston Rd, Plano, TX 75093", ContactPhone: "+155520204"},
	{Title: "Cotonou Cleaning Crew", OwnerOrigin: "Benin", Type: domain.Service, Anchor: "Cleaning", Description: "Residential and commercial cleaning.", City: "Irving", Address: "1707 MacArthur Blvd, Irving, TX 75061", ContactEmail: "clean@cotonou.com"},
	{Title: "Lome IT Support", OwnerOrigin: "Togo", Type: domain.Service, Anchor: "Tech", Description: "Computer repair and network setup.", City: "Richardson", Address: "1808 Arapaho Rd, Richardson, TX 75080", ContactWhatsApp: "+155520205"},
	{Title: "Conakry Caregivers", OwnerOrigin: "Guinea", Type: domain.Service, Anchor: "Health", Description: "Home health aides for seniors.", City: "Garland", Address: "1909 Northwest Hwy, Garland, TX 75041", ContactPhone: "+155520206"},
	{Title: "Niamey Notary Public", OwnerOrigin: "Niger", Type: domain.Service, Anchor: "Legal", Description: "Mobile notary services.", City: "Mesquite", Address: "2010 Galloway Ave, Mesquite, TX 75149", ContactEmail: "notary@niamey.com"},
}

var products = []domain.Listing{
	{Title: "Handmade Shea Butter", OwnerOrigin: "Ghana", Type: domain.Product, Anchor: "Beauty", Description: "Raw, unrefined shea butter.", City: "Dallas", Address: "2101 Elm St, Dallas, TX 75201", ContactWhatsApp: "+155530301"},
	{Title: "Nigerian Ankara Fabrics", OwnerOrigin: "Nigeria", Type: domain.Product, Anchor: "Fashion", Description: "Colorful Ankara prints by the yard.", City: "Houston", Address: "2202 Westheimer Rd, Houston, TX 77006", ContactEmail: "ankara@naija.com"},
	{Title: "Senegalese Baskets", OwnerOrigin: "Senegal", Type: domain.Product, Anchor: "Home Decor", Description: "Handwoven baskets for storage and decor.", City: "Plano", Address: "2303 Park Blvd, Plano, TX 75074", ContactPhone: "+155530302"},
	{Title: "Liberian Coffee Beans", OwnerOrigin: "Liberia", Type: domain.Product, Anchor: "Food", Description: "Robusta coffee beans direct from Monrovia.", City: "Fort Worth", Address: "2404 7th St, Fort Worth, TX 76107", ContactEmail: "coffee@liberia.com"},
	{Title: "Sierra Leone Diamonds (Art.)", OwnerOrigin: "Sierra Leone", Type: domain.Product, Anchor: "Jewelry", Description: "Artificial diamond jewelry.", City: "Arlington", Address: "2505 Collins St, Arlington, TX 76011", ContactWhatsApp: "+155530303"},
	{Title: "Malian Mud Cloth", OwnerOrigin: "Mali", Type: domain.Product, Anchor: "Fashion", Description: "Authentic Bogolanfini fabric.", City: "Irving", Address: "2606 Irving Blvd, Irving, TX 75060", ContactPhone: "+155530304"},
	{Title: "Benin Bronze Art", OwnerOrigin: "Benin", Type: domain.Product, Anchor: "Art", Description: "Replica bronze statues and art.", City: "Dallas", Address: "2707 Flora St, Dallas, TX 75201", ContactEmail: "art@benin.com"},
	{Title: "Togo Cocoa Butter", OwnerOrigin: "Togo", Type: domain.Product, Anchor: "Beauty", Description: "Pure cocoa butter for skin care.", City: "Garland", Address: "2808 Centerville Rd, Garland, TX 75043", ContactWhatsApp: "+155530305"},
	{Title: "Guinea Pepper Spice", OwnerOrigin: "Guinea", Type: domain.Product, Anchor: "Food", Description: "Spicy peppers for cooking.", City: "Grand Prairie", Address: "2909 Main St, Grand Prairie, TX 75050", ContactPhone: "+155530306"},
	{Title: "Niger Leather Bags", OwnerOrigin: "Niger", Type: domain.Product, Anchor: "Fashion", Description: "Handcrafted leather goods.", City: "McKinney", Address: "3010 Virginia Pkwy, McKinney, TX 75071", ContactEmail: "bags@niger.com"},
	{Title: "Ivory Coast Cocoa Powder", OwnerOrigin: "Cote d'Ivoire", Type: domain.Product, Anchor: "Food", Description: "Premium cocoa powder for baking.", City: "Frisco", Address: "3111 Main St, Frisco, TX 75033", ContactWhatsApp: "+155530307"},
}

var jobs = []domain.Listing{
	{Title: "Experienced Chef Needed", OwnerOrigin: "Nigeria", Type: domain.Job, Anchor: "Restaurant", Description: "Looking for a chef specializing in Nigerian cuisine.", City: "Dallas", Address: "3201 Greenville Ave, Dallas, TX 75206", ContactEmail: "jobs@chef.com"},
	{Title: "Sales Associate", OwnerOrigin: "Ghana", Type: domain.Job, Anchor: "Retail", Description: "Part-time sales associate for fabric store.", City: "Plano", Address: "3302 Parker Rd, Plano, TX 75023", ContactPhone: "+155540401"},
	{Title: "Delivery Driver", OwnerOrigin: "Senegal", Type: domain.Job, Anchor: "Logistics", Description: "Driver needed for local deliveries. Must have own van.", City: "Fort Worth", Address: "3403 Main St, Fort Worth, TX 76102", ContactWhatsApp: "+155540402"},
	{Title: "Hair Stylist", OwnerOrigin: "Liberia", Type: domain.Job, Anchor: "Beauty", Description: "Braider needed for busy salon.", City: "Arlington", Address: "3504 Cooper St, Arlington, TX 76015", ContactEmail: "salon@jobs.com"},
	{Title: "Java Developer", OwnerOrigin: "Sierra Leone", Type: domain.Job, Anchor: "Tech", Description: "Junior Java developer for startup.", City: "Irving", Address: "3605 Las Colinas Blvd, Irving, TX 75039", ContactPhone: "+155540403"},
	{Title: "Event Planner Assistant", OwnerOrigin: "Mali", Type: domain.Job, Anchor: "Events", Description: "Help organize engaging community events.", City: "Dallas", Address: "3706 Lemmon Ave, Dallas, TX 75219", ContactWhatsApp: "+155540404"},
	{Title: "Auto Mechanic", OwnerOrigin: "Benin", Type: domain.Job, Anchor: "Automotive", Description: "Experienced mechanic for Toyota/Honda shop.", City: "Garland", Address: "3807 Forest Ln, Garland, TX 75042", ContactEmail: "mechanic@jobs.com"},
	{Title: "French Teacher", OwnerOrigin: "Togo", Type: domain.Job, Anchor: "Education", Description: "Part-time French teacher for kids.", City: "Richardson", Address: "3908 Coit Rd, Richardson, TX 75080", ContactPhone: "+155540405"},
	{Title: "Construction Worker", OwnerOrigin: "Guinea", Type: domain.Job, Anchor: "Construction", Description: "General labor for construction site.", City: "McKinney", Address: "4009 Central Expy, McKinney, TX 75070", ContactWhatsApp: "+155540406"},
	{Title: "Accountant", OwnerOrigin: "Burkina Faso", Type: domain.Job, Anchor: "Finance", Description: "Bookkeeper needed for small business.", City: "Mesquite", Address: "4110 Motley Dr, Mesquite, TX 75150", ContactEmail: "finance@jobs.com"},
}

var requests = []domain.Listing{
	{Title: "Looking for Zobo leaves", OwnerOrigin: "Nigeria", Type: domain.Request, Anchor: "Sourcing", Description: "Need bulk Zobo (Hibiscus) leaves for an event.", City: "Dallas", Address: "4201 Ross Ave, Dallas, TX 75204", ContactEmail: "req1@example.com"},
	{Title: "Graphic Designer needed", OwnerOrigin: "Ghana", Type: domain.Request, Anchor: "Services", Description: "Need a logo for a new startup.", City: "Carrollton", Address: "4302 Josey Ln, Carrollton, TX 75006", ContactWhatsApp: "+155550501"},
	{Title: "Anyone flying to Senegal?", OwnerOrigin: "Senegal", Type: domain.Request, Anchor: "Travel", Description: "Need to send a small document package urgently.", City: "Euless", Address: "4403 Airport Fwy, Euless, TX 76040", ContactEmail: "req2@example.com"},
	{Title: "French Tutor for kids", OwnerOrigin: "Cote d'Ivoire", Type: domain.Request, Anchor: "Education", Description: "Looking for a tutor for my 2 kids.", City: "Suburbs", Address: "4504 Main St, Lewisville, TX 75067", ContactPhone: "+155550502"},
	{Title: "Apartment to rent", OwnerOrigin: "Gambia", Type: domain.Request, Anchor: "Housing", Description: "Looking for a 1 bedroom near the university.", City: "Denton", Address: "4605 University Dr, Denton, TX 76201", ContactWhatsApp: "+155550503"},
	{Title: "Used Toyota Camry", OwnerOrigin: "Togo", Type: domain.Request, Anchor: "Auto", Description: "Looking to buy a reliable commuter car.", City: "Grapevine", Address: "4706 Main St, Grapevine, TX 76051", ContactEmail: "req3@example.com"},
	{Title: "Web Developer partner", OwnerOrigin: "Niger", Type: domain.Request, Anchor: "Tech", Description: "Seeking a technical co-founder.", City: "Irving", Address: "4807 MacArthur Blvd, Irving, TX 75061", ContactPhone: "+155550504"},
	{Title: "Caterer for Naming Ceremony", OwnerOrigin: "Sierra Leone", Type: domain.Request, Anchor: "Events", Description: "Need catering for 50 people next Saturday.", City: "Allen", Address: "4908 McDermott Dr, Allen, TX 75013", ContactWhatsApp: "+155550505"},
	{Title: "Roommate Wanted", OwnerOrigin: "Liberia", Type: domain.Request, Anchor: "Housing", Description: "Share a 2-bedroom apartment in Arlington.", City: "Arlington", Address: "5009 Cooper St, Arlington, TX 76017", ContactEmail: "req4@example.com"},
	{Title: "Shipping Barrel to Mali", OwnerOrigin: "Mali", Type: domain.Request, Anchor: "Logistics", Description: "Who has the best rates for barrels?", City: "Dallas", Address: "5110 Buckner Blvd, Dallas, TX 75228", ContactPhone: "+155550506"},
}

var food = []domain.Listing{
	{Title: "Mama Put Dallas", OwnerOrigin: "Nigeria", Type: domain.Food, Anchor: "Restaurant", Description: "Home cooked meals: Pounded yam, Egusi, Suya.", City: "Dallas", Address: "5201 Forest Ln, Dallas, TX 75243", ContactEmail: "mama@put.com"},
	{Title: "Accra Chop Bar", OwnerOrigin: "Ghana", Type: domain.Food, Anchor: "Restaurant", Description: "Best Banku and Tilapia in town.", City: "Arlington", Address: "5302 Collins St, Arlington, TX 76014", ContactWhatsApp: "+155560601"},
	{Title: "Dakar Delights Catering", OwnerOrigin: "Senegal", Type: domain.Food, Anchor: "Catering", Description: "Catering for weddings and parties. Thieboudienne specialty.", City: "Plano", Address: "5403 Legacy Dr, Plano, TX 75024", ContactPhone: "+155560602"},
	{Title: "Monrovia Bakery", OwnerOrigin: "Liberia", Type: domain.Food, Anchor: "Bakery", Description: "Fresh rice bread and cassava cake daily.", City: "Fort Worth", Address: "5504 Berry St, Fort Worth, TX 76110", ContactEmail: "bread@monrovia.com"},
	{Title: "Freetown Kitchen", OwnerOrigin: "Sierra Leone", Type: domain.Food, Anchor: "Restaurant", Description: "Cassava leaves and Plasas like home.", City: "Garland", Address: "5605 Broadway Blvd, Garland, TX 75043", ContactWhatsApp: "+155560603"},
	{Title: "Bamako Grill", OwnerOrigin: "Mali", Type: domain.Food, Anchor: "Restaurant", Description: "Grilled fish and chicken with couscous.", City: "Irving", Address: "5706 Belt Line Rd, Irving, TX 75063", ContactPhone: "+155560604"},
	{Title: "Cotonou Cuisine", OwnerOrigin: "Benin", Type: domain.Food, Anchor: "Restaurant", Description: "Authentic Beninese flavors.", City: "Mesquite", Address: "5807 Gus Thomasson Rd, Mesquite, TX 75150", ContactEmail: "food@cotonou.com"},
	{Title: "Lome Snacks", OwnerOrigin: "Togo", Type: domain.Food, Anchor: "Snacks", Description: "Street food snacks and drinks.", City: "Richardson", Address: "5908 Plano Rd, Richardson, TX 75081", ContactWhatsApp: "+155560605"},
	{Title: "Conakry Cafe", OwnerOrigin: "Guinea", Type: domain.Food, Anchor: "Cafe", Description: "Coffee and pastries with a Guinean twist.", City: "Lewisville", Address: "6009 Main St, Lewisville, TX 75067", ContactPhone: "+155560606"},
	{Title: "Ouaga Spicy Pot", OwnerOrigin: "Burkina Faso", Type: domain.Food, Anchor: "Restaurant", Description: "Spicy stews and grilled meats.", City: "Grand Prairie", Address: "6110 Carrier Pkwy, Grand Prairie, TX 75052", ContactEmail: "spicy@ouaga.com"},
}

var events = []domain.Listing{
	{Title: "Naija Independence Day Bash", OwnerOrigin: "Nigeria", Type: domain.Event, Anchor: "Festival", Description: "Celebrating Nigeria @ 66. Food, Music, Dance.", City: "Dallas", Address: "6201 Fair Park, Dallas, TX 75210", ContactEmail: "indep@naija.com"},
	{Title: "Ghana Fest Texas", OwnerOrigin: "Ghana", Type: domain.Event, Anchor: "Festival", Description: "Cultural showcase of Ghanaian heritage.", City: "Arlington", Address: "6302 Levitt Pavilion, Arlington, TX 76010", ContactWhatsApp: "+155570701"},
	{Title: "Senegal Music Night", OwnerOrigin: "Senegal", Type: domain.Event, Anchor: "Concert", Description: "Live Mbalax music featuring top artists.", City: "Dallas", Address: "6403 Deep Ellum, Dallas, TX 75226", ContactPhone: "+155570702"},
	{Title: "Liberian Community Picnic", OwnerOrigin: "Liberia", Type: domain.Event, Anchor: "Community", Description: "Family fun day at the park.", City: "Fort Worth", Address: "6504 Trinity Park, Fort Worth, TX 76107", ContactEmail: "picnic@liberia.com"},
	{Title: "Sierra Leone Charity Gala", OwnerOrigin: "Sierra Leone", Type: domain.Event, Anchor: "Charity", Description: "Fundraiser for schools in Freetown.", City: "Plano", Address: "6605 Legacy Dr, Plano, TX 75024", ContactWhatsApp: "+155570703"},
	{Title: "Mali Art Exhibition", OwnerOrigin: "Mali", Type: domain.Event, Anchor: "Art", Description: "Showcasing contemporary Malian artists.", City: "Dallas", Address: "6706 Flora St, Dallas, TX 75201", ContactPhone: "+155570704"},
	{Title: "Benin Cultural Day", OwnerOrigin: "Benin", Type: domain.Event, Anchor: "Culture", Description: "Dance performances and food tasting.", City: "Irving", Address: "6807 Las Colinas Blvd, Irving, TX 75039", ContactEmail: "culture@benin.com"},
	{Title: "Togo Independence Party", OwnerOrigin: "Togo", Type: domain.Event, Anchor: "Party", Description: "Music and dancing all night.", City: "Garland", Address: "6908 Naaman Forest Blvd, Garland, TX 75040", ContactWhatsApp: "+155570705"},
	{Title: "Guinea Fashion Week", OwnerOrigin: "Guinea", Type: domain.Event, Anchor: "Fashion", Description: "Fashion show featuring Guinean designers.", City: "Dallas", Address: "7009 Market Center Blvd, Dallas, TX 75207", ContactPhone: "+155570706"},
	{Title: "Burkina Faso Film Screening", OwnerOrigin: "Burkina Faso", Type: domain.Event, Anchor: "Film", Description: "Screening of award-winning FESPACO films.", City: "Richardson", Address: "7110 Alamo Drafthouse, Richardson, TX 75080", ContactEmail: "film@burkina.com"},
}
