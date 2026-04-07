package sqlite

// ListingSelectionsSQL is the shared column selection for reading listings.
const ListingSelectionsSQL = `
	id, COALESCE(owner_id, ''), owner_origin, type, title, description,
	COALESCE(city, ''), COALESCE(address, ''), COALESCE(hours_of_operation, ''), 
	COALESCE(contact_email, ''), COALESCE(contact_phone, ''), COALESCE(contact_whatsapp, ''),
	COALESCE(website_url, ''), COALESCE(image_url, ''), created_at, deadline, is_active,
	event_start, event_end,
	COALESCE(skills, ''), job_start_date, COALESCE(job_apply_url, ''),
	COALESCE(company, ''), COALESCE(pay_range, ''), COALESCE(status, 'Approved'), featured
`

// UserSelectionsSQL is the shared column selection for reading users.
const UserSelectionsSQL = `id, google_id, email, name, avatar_url, COALESCE(role, 'User'), created_at`

// CategorySelectionsSQL is the shared column selection for reading categories.
const CategorySelectionsSQL = `id, name, claimable, is_system, active, requires_special_validation, created_at, updated_at`

// Shared SQL fragments
const (
	ListingActiveApprovedSQL = `is_active = true AND status = 'Approved'`
	ListingFilterTypeSQL     = ` AND type = ?`
)

// Shared Read Queries
const (
	ListingGetCountsSQL    = `SELECT type, COUNT(*) FROM listings WHERE ` + ListingActiveApprovedSQL + ` GROUP BY type`
	ListingGetLocationsSQL = `SELECT DISTINCT city FROM listings WHERE ` + ListingActiveApprovedSQL + ` AND city != '' ORDER BY city ASC`
	ListingTitleExistsSQL  = `SELECT EXISTS(SELECT 1 FROM listings WHERE title = ?)`
	UserGetCountSQL        = `SELECT COUNT(*) FROM users`
)

const listingColumns = `(id, owner_id, title, description, type, owner_origin, city, address, hours_of_operation, is_active, created_at, image_url, contact_email, contact_phone, contact_whatsapp, website_url, deadline, event_start, event_end, skills, job_start_date, job_apply_url, company, pay_range, status, featured)`

const listingUpsertUpdate = `ON CONFLICT(id) DO UPDATE SET
		owner_id = excluded.owner_id,
		title = excluded.title,
		description = excluded.description,
		type = excluded.type,
		owner_origin = excluded.owner_origin,
		city = excluded.city,
		address = excluded.address,
		hours_of_operation = excluded.hours_of_operation,
		is_active = excluded.is_active,
		image_url = excluded.image_url,
		contact_email = excluded.contact_email,
		contact_phone = excluded.contact_phone,
		contact_whatsapp = excluded.contact_whatsapp,
		website_url = excluded.website_url,
		deadline = excluded.deadline,
		event_start = excluded.event_start,
		event_end = excluded.event_end,
		skills = excluded.skills,
		job_start_date = excluded.job_start_date,
		job_apply_url = excluded.job_apply_url,
		company = excluded.company,
		pay_range = excluded.pay_range,
		status = excluded.status,
		featured = excluded.featured;`

// ListingUpsertSQL is the shared UPSERT query for both single and batch saves.
const ListingUpsertSQL = `INSERT INTO listings ` + listingColumns + `
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	` + listingUpsertUpdate

// CategoryUpsertSQL is the shared UPSERT query for category saving.
const CategoryUpsertSQL = `
	INSERT INTO categories (id, name, claimable, is_system, active, requires_special_validation, created_at, updated_at)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(id) DO UPDATE SET
		name = excluded.name,
		claimable = excluded.claimable,
		is_system = excluded.is_system,
		active = excluded.active,
		requires_special_validation = excluded.requires_special_validation,
		updated_at = excluded.updated_at;
	`

// ClaimUpsertSQL is the shared UPSERT query for claim saves.
const ClaimUpsertSQL = `
	INSERT INTO claim_requests (id, listing_id, listing_title, user_id, user_name, user_email, status, created_at)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(id) DO UPDATE SET
		status = excluded.status;
	`
