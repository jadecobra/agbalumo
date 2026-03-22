# Design System: Typography

> **Purpose:** Centralized definitions for all text styling, font families, and typographic scales across the agbalumo application.

## Font Family

| Token          | Font         | Weights                 | Source       |
| -------------- | ------------ | ----------------------- | ------------ |
| `font-display` | **Lexend**   | 300, 400, 500, 600, 700 | Google Fonts |

Applied globally via `<body class="font-display">`. No other font families are used.

## Scale

| Element         | Classes                                           |
| --------------- | ------------------------------------------------- |
| Page heading    | `text-4xl md:text-5xl font-bold tracking-tight` |
| Section heading | `text-2xl font-bold`                            |
| Card title      | `text-xl font-bold uppercase leading-tight`     |
| Modal title     | `text-xl font-bold`                             |
| Body text       | `text-sm leading-relaxed`                       |
| Label           | `text-xs font-bold uppercase tracking-wider`    |
| Badge           | `text-[10px] uppercase font-bold tracking-wider`|
| Caption         | `text-[10px] font-medium`                       |
