# Static Assets

This directory contains static assets for the web application.

## Structure

```
static/
├── images/          # Image files (PNG, JPG, SVG, etc.)
├── css/            # CSS files (if needed)
├── js/             # JavaScript files (if needed)
└── favicon.ico     # Favicon
```

## Image Guidelines

- Use descriptive filenames (e.g., `logo.png`, `hero-banner.jpg`)
- Optimize images for web (compress PNG/JPG files)
- Use appropriate formats:
  - PNG for logos, icons, graphics with transparency
  - JPG for photographs
  - SVG for scalable graphics and icons
- Keep file sizes reasonable for web performance

## Usage in Templates

Reference images in your HTML templates like this:

```html
<img src="/static/images/logo.png" alt="Company Logo" />
<img src="/static/images/hero-banner.jpg" alt="Hero Banner" />
```

## Roland Images

The application uses specific Roland-themed images:

- `roland-standing.png` - Logo image displayed at the top of the home page
- `roland-phone.png` - Displayed on the home page below the title
- `roland-maintenance.png` - Displayed on the 404 error page

These images should be placed in the `static/images/` directory and will be automatically styled with the appropriate CSS classes.
