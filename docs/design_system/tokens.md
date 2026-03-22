# Brand Tokens

> **Purpose**: Single source of truth for color palette, backgrounds, surfaces, and shadows.

## Color Palette

All colors are defined in [tailwind.config.js](../../tailwind.config.js) and must be used via their Tailwind token names. **Never hardcode hex values in templates.**

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

## Backgrounds & Surfaces

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

## Shadows

| Token | Value | Usage |
|-------|-------|-------|
| `shadow-soft` | `0 4px 20px -2px rgba(0,0,0,0.05)` | Cards, modals |
| `shadow-sm` | Tailwind default | Badges, nav items |
| `shadow-md` | Tailwind default | Hover states |
| `shadow-lg` | Tailwind default | Active buttons |
| `shadow-2xl` | Tailwind default | Modals |
| `shadow-orange-500/20` | Orange glow | Primary CTA buttons |
