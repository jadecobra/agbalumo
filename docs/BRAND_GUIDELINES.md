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

## 6. Backgrounds, Surfaces & Shadows

All surface variants, modal backdrops, and shadow definitions have been extracted to the core design tokens.

> **See [Design System Tokens](./design_system/tokens.md#backgrounds--surfaces)** for backgrounds and shadows.

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
| [about.html](file:///Users/johnnyblase/gym/agbalumo/ui/templates/about.html) | "explore" button uses `bg-stone-900` instead of brand primary/secondary. "Services" icon uses `text-blue-500` and "Products" uses `text-green-500` — not brand colors. |

---

## 11. Reference Screenshots

![agbalumo Homepage](file:///Users/johnnyblase/.gemini/antigravity/brain/62c1ce97-7538-478d-8440-ffa2a449b107/homepage_screenshot_1771707186729.png)

![agbalumo About Page](file:///Users/johnnyblase/.gemini/antigravity/brain/62c1ce97-7538-478d-8440-ffa2a449b107/about_page_screenshot_1771707193543.png)
