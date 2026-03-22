# Design System Components

This document contains standardized guidelines for reusable component designs, such as cards, modals, buttons, and form inputs.

---

## 1. Shapes & Border Radius

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

## 2. Buttons

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

## 3. Form Inputs

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

## 4. Animations & Transitions

| Pattern | Classes |
|---------|---------|
| Button press | `active:scale-[0.98]` |
| Hover scale | `hover:scale-105 transition-transform` |
| Color transition | `transition-colors` |
| All transitions | `transition-all duration-300` |
| Featured card zoom | `group-hover:scale-110 transition-transform duration-700` |
| Modal open | `open:animate-in open:fade-in open:zoom-in-95` |
| Loading spinner | `animate-spin rounded-full h-12 w-12 border-b-2 border-primary` |
