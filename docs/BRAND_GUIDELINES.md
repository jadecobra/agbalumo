# agbalumo Brand Guidelines

> **Purpose**: Single source of truth for all visual elements across every page and modal. Every new template, component, or page MUST reference this guide.

---

## 1. Color Palette & Tokens

All colors, neutral scales, and gradient tokens have been extracted to a dedicated token system. 

> **See [Design System Tokens](./design_system/tokens.md)** for the complete color palette.

---

## 2. Typography

### Font Family

| Token | Font | Weights | Source |
|-------|------|---------|--------|
| `font-display` | **Lexend** | 300, 400, 500, 600, 700 | Google Fonts |

Applied globally via `<body class="font-display">`. No other font families are used.

### Scale

| Element | Classes |
|---------|---------|
| Page heading | `text-4xl md:text-5xl font-bold tracking-tight` |
| Section heading | `text-2xl font-bold` |
| Card title | `text-xl font-bold uppercase leading-tight` |
| Modal title | `text-xl font-bold` |
| Body text | `text-sm leading-relaxed` |
| Label | `text-xs font-bold uppercase tracking-wider` |
| Badge | `text-[10px] uppercase font-bold tracking-wider` |
| Caption | `text-[10px] font-medium` |

---

## 3. Icons

| System | Package |
|--------|---------|
| **Material Symbols Outlined** | `material-symbols-outlined` |

### Standard Sizes

| Context | Class |
|---------|-------|
| Default (buttons, nav) | `text-[20px]` |
| Inline with text | `text-[18px]` |
| Small (badges/labels) | `text-[14px]` |
| Hero/empty state | `text-3xl` to `text-4xl` |

### Category Icons

| Type | Icon |
|------|------|
| Business | `storefront` |
| Service | `handyman` / `construction` |
| Food | `restaurant` |
| Product | `shopping_bag` |
| Event | `event` |
| Job | `work` / `work_outline` |
| Request | `volunteer_activism` / `campaign` |

---

## 4. Components

All reusable component designs (cards, modals, buttons, form inputs) and their associated animations have been extracted to a dedicated components system.

> **See [Design System Components](./design_system/components.md)** for component guidelines.

---

## 5. Known Inconsistencies ⚠️

These deviations from the brand system should be corrected over time:

| File | Issue |
|------|-------|
| [admin_dashboard.html](file:///Users/johnnyblase/gym/agbalumo/ui/templates/admin_dashboard.html) | Uses `text-gray-*`, `bg-gray-*`, `border-gray-*` instead of `stone-*` tokens. Uses `rounded-lg` not `rounded-3xl`. |
| [admin_login.html](file:///Users/johnnyblase/gym/agbalumo/ui/templates/admin_login.html) | Loads Tailwind via CDN (`cdn.tailwindcss.com`) and duplicates the config inline (with `text-sub` as `#8d6d5e` not `#6d4c41`). Should use the compiled `output.css` and extend `base.html`. |
| [admin_users.html](file:///Users/johnnyblase/gym/agbalumo/ui/templates/admin_users.html) | Likely uses `gray-*` (same pattern as dashboard). |
| [admin_listings.html](file:///Users/johnnyblase/gym/agbalumo/ui/templates/admin_listings.html) | Likely uses `gray-*`. |
| [error.html](file:///Users/johnnyblase/gym/agbalumo/ui/templates/error.html) | Hardcoded hex values (`bg-[#FFF2EB]`, `text-[#FF5E0E]`, `bg-[#2D5A27]`) instead of tokens. Loads Tailwind CDN instead of `output.css`. |
| [about.html](file:///Users/johnnyblase/gym/agbalumo/ui/templates/about.html) | "explore" button uses `bg-stone-900` instead of brand primary/secondary. "Services" icon uses `text-blue-500` and "Products" uses `text-green-500` — not brand colors. |

---

## 6. Reference Screenshots

![agbalumo Homepage](file:///Users/johnnyblase/.gemini/antigravity/brain/62c1ce97-7538-478d-8440-ffa2a449b107/homepage_screenshot_1771707186729.png)

![agbalumo About Page](file:///Users/johnnyblase/.gemini/antigravity/brain/62c1ce97-7538-478d-8440-ffa2a449b107/about_page_screenshot_1771707193543.png)
