# Agbalumo Brand Guidelines

> **Purpose**: Single source of truth for all visual elements across every page and modal. Every new template, component, or page MUST reference this guide.

---

## 1. Color Palette

All colors are defined in [tailwind.config.js](file:///Users/johnnyblase/gym/agbalumo/tailwind.config.js) and must be used via their Tailwind token names. **Never hardcode hex values in templates.**

### Brand Colors

| Token | Hex | Usage | Tailwind Class |
|-------|-----|-------|----------------|
| **primary** | `#FF5E0E` | CTAs, active states, links, logo accent | `bg-primary`, `text-primary` |
| **secondary** | `#2D5A27` | "Ask" buttons, success, green accents | `bg-secondary`, `text-secondary` |

### Surface & Background

| Token | Hex | Usage | Tailwind Class |
|-------|-----|-------|----------------|
| **background-light** | `#FFF2EB` | Page background (light mode) | `bg-background-light` |
| **background-dark** | `#23160f` | Page background (dark mode) | `bg-background-dark` |
| **surface-light** | `#ffffff` | Card/modal surface (light mode) | `bg-surface-light` / `bg-white` |
| **surface-dark** | `#2f221c` | Card/modal surface (dark mode) | `bg-surface-dark` |

### Text

| Token | Hex | Usage | Tailwind Class |
|-------|-----|-------|----------------|
| **text-main** | `#181310` | Primary body text | `text-text-main` |
| **text-sub** | `#6d4c41` | Secondary/label text | `text-text-sub` |

### Neutral Scale

Use **`stone-*`** for all neutral grays. **Never use `gray-*`** (generic Tailwind gray).

| Use Case | Light | Dark |
|----------|-------|------|
| Borders | `border-stone-200` | `border-stone-700` |
| Muted text | `text-stone-400` / `text-stone-500` | `text-stone-400` |
| Hover surfaces | `bg-stone-50` / `bg-stone-100` | `bg-stone-800` |
| Dividers | `border-stone-100` | `border-stone-800` |

### Category Accent Gradients (Listing Cards Without Images)

| Type | Gradient |
|------|----------|
| Business | `from-indigo-500 to-indigo-700` |
| Service | `from-emerald-500 to-emerald-700` |
| Food | `from-orange-500 to-orange-700` |
| Product | (default) `from-stone-400 to-stone-600` |
| Event | `from-purple-500 to-purple-700` |
| Job | `from-blue-500 to-blue-700` |
| Request | `from-teal-500 to-teal-700` |

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

## 4. Shapes & Border Radius

| Element | Radius | Class |
|---------|--------|-------|
| Cards | 24px | `rounded-3xl` |
| Modals | 24px | `rounded-3xl` |
| Buttons (pill) | 9999px | `rounded-full` |
| Form inputs | 12px | `rounded-xl` |
| Badges/pills | 9999px | `rounded-full` |
| Feature cards | 24px | `rounded-3xl` |
| About page tiles | 24px | `rounded-3xl` |
| Admin cards | 8px | `rounded-lg` ⚠️ **inconsistent** |

> **Rule**: User-facing surfaces use `rounded-3xl`. Form inputs use `rounded-xl`. Buttons and badges use `rounded-full`.

---

## 5. Buttons

### Primary CTA (Orange)

```html
class="bg-primary hover:bg-orange-600 text-white font-bold rounded-full px-5 h-10
       shadow-md shadow-orange-500/20 transition-all active:scale-[0.98]"
```

Used for: **Post**, **Save Changes**, **Submit**, **Upload**

### Secondary CTA (Green)

```html
class="bg-secondary/10 hover:bg-secondary/20 text-secondary rounded-full px-4 h-10
       font-semibold transition-all"
```

Used for: **Ask**, **Sign In** (logged out)

### Ghost/Outline

```html
class="border border-stone-200 dark:border-stone-700 rounded-full px-5 py-2.5
       font-bold text-stone-600 hover:bg-stone-50 transition-all"
```

Used for: **Cancel**, **Join the Directory**

### Filter Pill (Active)

```html
class="rounded-full bg-stone-900 text-white px-4 h-8 text-xs font-bold uppercase"
```

### Filter Pill (Inactive)

```html
class="rounded-full bg-white border border-stone-200 px-4 h-8 text-xs font-semibold uppercase
       hover:bg-stone-50 text-text-main"
```

### Danger

```html
class="bg-red-600 text-white rounded-full px-3 py-1 text-xs font-bold"
```

Used for: **Delete**

### Contact Action Buttons

| Contact | Color |
|---------|-------|
| WhatsApp | `bg-green-600 hover:bg-green-700` |
| Call | `bg-blue-600 hover:bg-blue-700` |
| Email | `bg-stone-700 hover:bg-stone-800` |
| Website | `bg-primary hover:bg-orange-600` |

---

## 6. Backgrounds & Surfaces

### Page Background

```html
body: bg-background-light dark:bg-background-dark
```

### Glassmorphic Header

```html
class="backdrop-blur-md bg-background-light/80 dark:bg-background-dark/80
       border-b border-orange-100/50 dark:border-stone-800/50"
```

### Card Surface

```html
class="bg-white dark:bg-surface-dark border border-stone-100 dark:border-stone-800
       rounded-3xl shadow-sm hover:shadow-md"
```

### Modal Surface

```html
class="bg-background-light dark:bg-background-dark rounded-3xl
       border border-orange-100 dark:border-stone-800 shadow-soft"
```

### Modal Backdrop

```html
class="backdrop:bg-black/50 backdrop:animate-in backdrop:fade-in"
```

### Decorative Orb (Modal)

```html
class="absolute top-0 right-0 w-32 h-32 bg-primary/10 rounded-full blur-3xl
       -z-10 transform translate-x-10 -translate-y-10"
```

---

## 7. Shadows

| Token | Value | Usage |
|-------|-------|-------|
| `shadow-soft` | `0 4px 20px -2px rgba(0,0,0,0.05)` | Cards, modals |
| `shadow-sm` | Tailwind default | Badges, nav items |
| `shadow-md` | Tailwind default | Hover states |
| `shadow-lg` | Tailwind default | Active buttons |
| `shadow-2xl` | Tailwind default | Modals |
| `shadow-orange-500/20` | Orange glow | Primary CTA buttons |

---

## 8. Animations & Transitions

| Pattern | Classes |
|---------|---------|
| Button press | `active:scale-[0.98]` |
| Hover scale | `hover:scale-105 transition-transform` |
| Color transition | `transition-colors` |
| All transitions | `transition-all duration-300` |
| Featured card zoom | `group-hover:scale-110 transition-transform duration-700` |
| Modal open | `open:animate-in open:fade-in open:zoom-in-95` |
| Loading spinner | `animate-spin rounded-full h-12 w-12 border-b-2 border-primary` |

---

## 9. Form Inputs

### Standard Input

```html
class="w-full bg-white dark:bg-surface-dark border border-stone-200 dark:border-stone-700
       rounded-xl px-4 py-3 text-sm
       focus:ring-2 focus:ring-primary/20 focus:border-primary outline-none
       transition-all placeholder:text-stone-400 text-text-main dark:text-white"
```

### Select

Same as input, add `appearance-none`.

### Textarea

Same as input, add `resize-none`.

### Label

```html
class="text-xs font-bold uppercase tracking-wider text-text-sub dark:text-stone-400 ml-1"
```

---

## 10. Known Inconsistencies ⚠️

These deviations from the brand system should be corrected over time:

| File | Issue |
|------|-------|
| [admin_dashboard.html](file:///Users/johnnyblase/gym/agbalumo/ui/templates/admin_dashboard.html) | Uses `text-gray-*`, `bg-gray-*`, `border-gray-*` instead of `stone-*` tokens. Uses `rounded-lg` not `rounded-3xl`. |
| [admin_login.html](file:///Users/johnnyblase/gym/agbalumo/ui/templates/admin_login.html) | Loads Tailwind via CDN (`cdn.tailwindcss.com`) and duplicates the config inline (with `text-sub` as `#8d6d5e` not `#6d4c41`). Should use the compiled `output.css` and extend `base.html`. |
| [admin_users.html](file:///Users/johnnyblase/gym/agbalumo/ui/templates/admin_users.html) | Likely uses `gray-*` (same pattern as dashboard). |
| [admin_listings.html](file:///Users/johnnyblase/gym/agbalumo/ui/templates/admin_listings.html) | Likely uses `gray-*`. |
| [error.html](file:///Users/johnnyblase/gym/agbalumo/ui/templates/error.html) | Hardcoded hex values (`bg-[#FFF2EB]`, `text-[#FF5E0E]`, `bg-[#2D5A27]`) instead of tokens. Loads Tailwind CDN instead of `output.css`. |
| [about.html](file:///Users/johnnyblase/gym/agbalumo/ui/templates/about.html) | "Explore Listings" button uses `bg-stone-900` instead of brand primary/secondary. "Services" icon uses `text-blue-500` and "Products" uses `text-green-500` — not brand colors. |

---

## 11. Reference Screenshots

![Agbalumo Homepage](file:///Users/johnnyblase/.gemini/antigravity/brain/62c1ce97-7538-478d-8440-ffa2a449b107/homepage_screenshot_1771707186729.png)

![Agbalumo About Page](file:///Users/johnnyblase/.gemini/antigravity/brain/62c1ce97-7538-478d-8440-ffa2a449b107/about_page_screenshot_1771707193543.png)
