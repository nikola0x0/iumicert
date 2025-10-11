# IU-MiCert Issuer Dashboard - Design Document

## Overview
This document outlines the redesign of the IU-MiCert Issuer Dashboard, transforming it from a top-navigation layout to a modern sidebar-based dashboard with improved UX and professional educational aesthetics.

## Design Philosophy

### Core Principles
1. **Educational Professional**: Clean, trustworthy design suitable for academic institutions
2. **Simplicity First**: Minimal clutter, clear information hierarchy
3. **Desktop-First**: Optimized for desktop workflows, mobile shows access warning
4. **Modern & Light**: Contemporary design with light theme and subtle shadows
5. **Rounded & Soft**: Generous border radius (12-16px) for cards, buttons use softer shadows and rounded corners for a friendly, approachable feel

### Color Palette

#### Primary Colors
- **Navy Blue** (`#1e40af` / blue-800): Primary brand color for headers, active states
- **Royal Blue** (`#2563eb` / blue-600): Interactive elements, CTAs, links
- **Soft Blue** (`#3b82f6` / blue-500): Cards, stat cards with gradient overlays
- **Light Blue** (`#eff6ff` / blue-50): Subtle backgrounds, hover states
- **Pale Blue** (`#dbeafe` / blue-100): Secondary card backgrounds

#### Accent Colors
- **Success Green** (`#16a34a` / green-600): Success states, verified status
- **Light Green** (`#dcfce7` / green-100): Success card backgrounds
- **Warning Amber** (`#f59e0b` / amber-500): Warning states, pending actions
- **Light Amber** (`#fef3c7` / amber-100): Warning card backgrounds
- **Error Red** (`#dc2626` / red-600): Error states, destructive actions
- **Light Red** (`#fee2e2` / red-100): Error card backgrounds
- **Neutral Gray** (`#6b7280` / gray-500): Secondary text, borders

#### Background & Text
- **Main Background**: `#f8fafc` (slate-50) - softer than pure gray
- **Card Background**: `#ffffff` (white) with soft shadows
- **Primary Text**: `#0f172a` (slate-900)
- **Secondary Text**: `#64748b` (slate-500)

#### Shadows
- **Card Shadow**: `0 1px 3px 0 rgb(0 0 0 / 0.1), 0 1px 2px -1px rgb(0 0 0 / 0.1)` (shadow-md)
- **Hover Shadow**: `0 10px 15px -3px rgb(0 0 0 / 0.1), 0 4px 6px -4px rgb(0 0 0 / 0.1)` (shadow-lg)
- **Stat Card Shadow**: Soft, colorful shadows matching card theme (e.g., blue-500/10)

### Typography
- **Headings**: System font stack (Inter/SF Pro on macOS)
- **Body**: System font stack
- **Monospace**: 'SF Mono', 'Monaco', 'Courier New' for code/data

### Icon Library
- **Primary**: Lucide React (modern, consistent)
- **Secondary**: Heroicons (when Lucide doesn't have the icon)
- **Style**: Outline/stroke style, no filled icons for consistency

## Layout Structure

```
┌──────────────────────────────────────────────────────────────┐
│  Sidebar (260px)     │  Main Content (slate-50 bg)          │
│  (white, shadow)     │                                       │
│  ╔════════════╗      │  [Wallet Connect - top right]        │
│  ║ IU-MiCert  ║      │                                       │
│  ║ Issuer     ║      │  ╭─────────────────────────────╮     │
│  ╚════════════╝      │  │  Rounded Card (shadow-lg)   │     │
│  ─────────────       │  │                             │     │
│                      │  │  Stat Cards Grid (rounded)  │     │
│  ╭──────────────╮    │  │  ╭──────╮ ╭──────╮ ╭──────╮│     │
│  │ 📄 Publish   │◀── │  │  │ Icon │ │ Icon │ │ Icon ││     │
│  ╰──────────────╯    │  │  │ Data │ │ Data │ │ Data ││     │
│                      │  │  ╰──────╯ ╰──────╯ ╰──────╯│     │
│  ┌──────────────┐    │  ╰─────────────────────────────╯     │
│  │ 🛡️ Verifier  │    │                                       │
│  └──────────────┘    │  ╭─────────────────────────────╮     │
│                      │  │  Data Table / Card List     │     │
│  ┌──────────────┐    │  │  (rounded-2xl)              │     │
│  │ 💾 Demo Data │    │  │                             │     │
│  └──────────────┘    │  ╰─────────────────────────────╯     │
│  ─────────────       │                                       │
│  v1.0.0             │                                       │
└──────────────────────────────────────────────────────────────┘

Legend:
╔╗ = Logo container
╭╮╰╯ = Rounded cards (border-radius: 16px)
◀── = Active state (gradient blue)
┌┐└┘ = Inactive nav items (hover: light gray)
```

### Component Breakdown

#### 1. Sidebar Component
**File**: `src/components/layout/Sidebar.tsx`

**Structure**:
- Fixed position (260px width for more breathing room)
- Spans full height
- White background with subtle shadow-right
- Dark text for better readability

**Elements**:
- **Logo Area** (top):
  - "IU-MiCert" text logo, 20px font, bold
  - Small subtitle "Issuer Dashboard" in gray-500, 12px
  - Padding: `py-6 px-6`
  - Border bottom: subtle gray-200

- **Navigation Links**:
  - Publish Terms (icon: FileText from Lucide)
  - Verifier (icon: ShieldCheck from Lucide)
  - Demo Data (icon: Database from Lucide)
  - Each link in rounded container (rounded-xl) with generous padding
  - Gap between links: `gap-2`

- **Active State**:
  - Background: gradient from blue-500 to blue-600
  - Text: white
  - Icon: white
  - Soft shadow
  - Border radius: `rounded-xl`

- **Inactive State**:
  - Background: transparent
  - Text: gray-600
  - Icon: gray-500
  - Hover: bg-gray-100, rounded-xl

- **Footer** (bottom):
  - Version info (gray-400 text, 11px)
  - Padding: `p-4`
  - Border top: subtle gray-200

**Behavior**:
- Active link has full rounded background with gradient
- Hover states: smooth transition to light gray background (transition-all duration-200)
- Click navigates to route using Next.js Link component
- Icons are 20px (w-5 h-5) with consistent spacing

#### 2. Main Layout Wrapper
**File**: `src/components/layout/DashboardLayout.tsx`

**Structure**:
- Wraps all dashboard pages
- Contains Sidebar + Main content area
- Main area has padding and max-width constraints
- Includes wallet connection in header

**Header**:
- Position: Sticky top
- Background: White with subtle shadow
- Right side: ConnectKit wallet button
- Left side: Page title (dynamic based on route)

#### 3. Page Components

##### A. Publish Terms Page
**Route**: `/` (home/default)
**File**: `src/app/page.tsx` → `src/components/issuer/PublishTermsPage.tsx`

**Layout**:
```
┌──────────────────────────────────────────────────────────┐
│  Terms Overview                                          │
│  ┌────────────┬────────────┬────────────┐               │
│  │ Total      │ Published  │ Pending    │               │
│  │ Terms: 7   │ Terms: 3   │ Terms: 4   │               │
│  └────────────┴────────────┴────────────┘               │
│                                                           │
│  Term List (Table)                                       │
│  ┌───────────────────────────────────────────────────┐  │
│  │ Term ID        │ Students │ Status  │ Actions    │  │
│  ├───────────────────────────────────────────────────┤  │
│  │ Semester_1_23  │ 5        │ ✓ Pub'd │ [View]     │  │
│  │ Semester_2_23  │ 5        │ Pending │ [Publish]  │  │
│  └───────────────────────────────────────────────────┘  │
└──────────────────────────────────────────────────────────┘
```

**Features**:
- Statistics cards at top
- Table with all terms
- Status badges (published/pending)
- Individual term publish action
- Real-time status updates during publishing
- Success/error toast notifications
- Etherscan link after successful publish

**Components**:
- `StatCard`: Reusable stat display
- `TermsTable`: Table with actions
- `PublishButton`: Button with loading state
- `StatusBadge`: Color-coded status indicator

##### B. Verifier Page
**Route**: `/verifier`
**File**: `src/app/verifier/page.tsx` (mostly keep current)

**Enhancements**:
- Replace emoji icons with Lucide icons:
  - Upload: `Upload` icon
  - Verified: `CheckCircle2` icon
  - Failed: `XCircle` icon
  - Folder: `Folder` / `FolderOpen` icons
  - Course: `BookOpen` icon
- Update color scheme to match new palette
- Improve drag-and-drop visual feedback
- Add breadcrumb navigation

##### C. Demo Data Page
**Route**: `/demo-data`
**File**: `src/app/demo-data/page.tsx`

**Enhancements**:
- Replace emoji icons:
  - Reset: `Trash2` icon
  - Generate: `Zap` icon
  - Add Term: `PlusCircle` icon
  - Success: `CheckCircle2` icon
  - Error: `AlertCircle` icon
- Improve card layouts with better spacing
- Add confirmation modals for destructive actions
- Progress indicators for long operations

## Component Specifications

### 1. Sidebar Navigation

```tsx
interface NavItem {
  label: string;
  href: string;
  icon: LucideIcon;
}

const navItems: NavItem[] = [
  { label: "Publish Terms", href: "/", icon: FileText },
  { label: "Verifier", href: "/verifier", icon: ShieldCheck },
  { label: "Demo Data", href: "/demo-data", icon: Database },
];
```

**Styling**:
- Background: `bg-gradient-to-b from-blue-900 to-blue-800`
- Text: `text-white`
- Active link: `bg-blue-700/20 border-l-4 border-blue-400`
- Hover: `hover:bg-blue-700/10 transition-colors`
- Icon size: 20px (w-5 h-5)
- Padding: `px-6 py-3`

### 2. Stat Cards

```tsx
interface StatCardProps {
  title: string;
  value: string | number;
  icon: LucideIcon;
  color: "blue" | "green" | "amber" | "purple";
  subtitle?: string; // optional secondary info
}
```

**Styling** (Rounded, Modern Approach):
- **Container**: `bg-white rounded-2xl shadow-lg p-6 border border-gray-100`
- **Layout**: Flexbox with icon on left or top-right badge style
- **Icon Container**:
  - Soft colored background (e.g., `bg-blue-100 text-blue-600`)
  - Large rounded container: `w-12 h-12 rounded-xl flex items-center justify-center`
  - Icon size: `w-6 h-6`
- **Title**: `text-gray-500 text-sm font-medium uppercase tracking-wide`
- **Value**: `text-3xl font-bold text-slate-900 mt-3`
- **Subtitle** (optional): `text-xs text-gray-400 mt-1`
- **Hover Effect**: Subtle lift with shadow-xl on hover (transition-all duration-300)
- **Color Variants**:
  - Blue: `bg-blue-100 text-blue-600` (icon), optional subtle blue gradient overlay
  - Green: `bg-green-100 text-green-600`
  - Amber: `bg-amber-100 text-amber-600`
  - Purple: `bg-purple-100 text-purple-600`

### 3. Status Badge

```tsx
interface StatusBadgeProps {
  status: "published" | "pending" | "error";
  label?: string;
}
```

**Styling** (More Rounded):
- Published: `bg-green-100 text-green-700 border border-green-200/50`
- Pending: `bg-amber-100 text-amber-700 border border-amber-200/50`
- Error: `bg-red-100 text-red-700 border border-red-200/50`
- Base: `px-4 py-1.5 rounded-full text-xs font-semibold`
- With icon: Include small icon (12px) before text
- Subtle shadow: `shadow-sm`

### 4. Data Table

**Styling** (Rounded Card Style):
- **Container**: `bg-white rounded-2xl shadow-lg border border-gray-100 overflow-hidden`
- **Header**: `bg-gradient-to-r from-slate-50 to-gray-50 border-b border-gray-200`
- **Header text**: `text-xs font-semibold text-slate-600 uppercase tracking-wider`
- **Row**:
  - `border-b border-gray-100/50 last:border-b-0`
  - Hover: `hover:bg-blue-50/30 transition-all duration-200`
  - Rounded corners on first/last visible row
- **Cell padding**: `px-6 py-4`
- **Action buttons**: Small rounded buttons (`rounded-lg`) with icon + text
- **Empty state**: Centered with icon and message in muted colors

**Table-less Alternative (Card List)**:
For better mobile/tablet experience, consider card-based list:
- Each row as individual rounded card (`rounded-xl`)
- Stack information vertically on smaller screens
- Add subtle shadow and hover lift effect

## Responsive Behavior

### Desktop (>1024px)
- Full sidebar visible (240px fixed)
- Main content area uses remaining space
- Tables show all columns
- Multi-column layouts for stats

### Tablet (768px - 1024px)
- Show alert: "For best experience, please use desktop"
- Optionally: Collapse sidebar to icons only

### Mobile (<768px)
- Full-screen modal overlay:
  ```
  ┌──────────────────────────┐
  │    🖥️                    │
  │                          │
  │  Desktop Access Required │
  │                          │
  │  Please access this      │
  │  dashboard from a        │
  │  desktop computer for    │
  │  the best experience.    │
  │                          │
  └──────────────────────────┘
  ```

## User Flows

### Flow 1: Publishing a Term Root

1. User lands on dashboard (Publish Terms page)
2. Sees list of terms with status
3. Clicks wallet connect if not connected
4. Selects term to publish
5. Clicks "Publish" button
6. Button shows loading state ("Publishing...")
7. Transaction submitted to blockchain
8. Toast notification: "Transaction submitted"
9. Waits for confirmation
10. Success: Status badge updates to "Published", Etherscan link shown
11. Error: Error message displayed, user can retry

### Flow 2: Verifying a Receipt

1. User clicks "Verifier" in sidebar
2. Sees upload interface
3. Drags/clicks to upload JSON file
4. File parsed and displayed
5. User clicks "Verify Entire Journey"
6. Loading state shown
7. Results displayed with color-coded badges
8. User can expand terms/courses to see details
9. Blockchain verification status shown if available

### Flow 3: Managing Demo Data

1. User clicks "Demo Data" in sidebar
2. Sees current data status
3. Option 1: Reset all data (destructive)
   - Clicks "Reset"
   - Confirmation modal appears
   - Confirms action
   - Loading state
   - Success message
4. Option 2: Generate new data
   - Selects number of students
   - Selects terms to generate
   - Clicks "Generate"
   - Progress indicator shown
   - Success message with summary

## Technical Implementation

### File Structure

```
packages/issuer/web/iumicert-issuer/src/
├── components/
│   ├── layout/
│   │   ├── DashboardLayout.tsx       # Main layout wrapper
│   │   ├── Sidebar.tsx                # Sidebar navigation
│   │   └── MobileWarning.tsx          # Mobile access warning
│   ├── ui/                            # shadcn components
│   │   ├── stat-card.tsx
│   │   ├── status-badge.tsx
│   │   ├── data-table.tsx
│   │   └── ...
│   ├── issuer/
│   │   ├── PublishTermsPage.tsx       # New main dashboard page
│   │   ├── TermsTable.tsx             # Terms list table
│   │   └── PublishButton.tsx          # Publish action button
│   └── ReceiptVerifier.tsx            # Updated verifier (remove emojis)
├── app/
│   ├── layout.tsx                     # Root layout with DashboardLayout
│   ├── page.tsx                       # Publish Terms page
│   ├── verifier/
│   │   └── page.tsx                   # Verifier page
│   └── demo-data/
│       └── page.tsx                   # Demo data page
└── lib/
    └── icons.ts                       # Icon exports (Lucide + Heroicons)
```

### shadcn Components to Use

Install these shadcn components:
```bash
npx shadcn@latest add button
npx shadcn@latest add card
npx shadcn@latest add badge
npx shadcn@latest add table
npx shadcn@latest add dialog
npx shadcn@latest add toast
npx shadcn@latest add alert
npx shadcn@latest add separator
```

### Dependencies

```json
{
  "lucide-react": "^0.400.0",
  "@heroicons/react": "^2.1.0",
  "class-variance-authority": "^0.7.0",
  "clsx": "^2.1.0",
  "tailwind-merge": "^2.3.0"
}
```

## Accessibility Considerations

1. **Keyboard Navigation**: All interactive elements accessible via keyboard
2. **ARIA Labels**: Proper labels for icons and buttons
3. **Focus Indicators**: Visible focus states on all interactive elements
4. **Color Contrast**: WCAG AA compliant contrast ratios
5. **Screen Readers**: Semantic HTML and ARIA attributes

## Performance Optimizations

1. **Code Splitting**: Lazy load heavy components (verifier, demo data)
2. **Memoization**: Memo-ize expensive renders (tables, lists)
3. **Debouncing**: Debounce search/filter inputs
4. **Optimistic UI**: Show immediate feedback before API calls complete

## Future Enhancements (Out of Scope)

1. Dark mode toggle
2. Collapsible sidebar
3. User settings page
4. Multi-language support
5. Advanced filtering and sorting
6. Bulk operations for multiple terms
7. Analytics dashboard
8. Export functionality (CSV, PDF)

## Implementation Phases

### Phase 1: Core Layout (Priority: High)
- [ ] Create DashboardLayout component
- [ ] Create Sidebar component with navigation
- [ ] Implement mobile warning screen
- [ ] Update root layout to use new structure

### Phase 2: Publish Terms Page (Priority: High)
- [ ] Create PublishTermsPage component
- [ ] Create StatCard component
- [ ] Create TermsTable component
- [ ] Implement publish functionality with new UI
- [ ] Add toast notifications

### Phase 3: UI Component Updates (Priority: Medium)
- [ ] Update Verifier page (replace emojis with icons)
- [ ] Update Demo Data page (replace emojis with icons)
- [ ] Add shadcn components
- [ ] Implement consistent styling

### Phase 4: Polish & Testing (Priority: Medium)
- [ ] Test all user flows
- [ ] Verify accessibility
- [ ] Cross-browser testing
- [ ] Performance audit
- [ ] Fix bugs and edge cases

## Design Assets

### Logo/Branding
- Text logo: "IU-MiCert" (no image initially)
- Font: Bold, 20px
- Color: White on navy blue sidebar

### Iconography Map

| Element | Current | New Icon (Lucide) |
|---------|---------|-------------------|
| Dashboard | N/A | `LayoutDashboard` |
| Publish Terms | N/A | `FileText` |
| Verifier | 🔍 | `ShieldCheck` |
| Demo Data | 🎲 | `Database` |
| Success | ✅ | `CheckCircle2` |
| Error | ❌ | `XCircle` |
| Warning | ⚠️ | `AlertTriangle` |
| Loading | ⏳ | `Loader2` (animated) |
| Upload | 📤 | `Upload` |
| Download | 📥 | `Download` |
| Folder Open | 📂 | `FolderOpen` |
| Folder Closed | 📁 | `Folder` |
| Course/Book | 📚 | `BookOpen` |
| Settings | ⚙️ | `Settings` |
| Reset | 🧹 | `Trash2` |
| Generate | 🚀 | `Zap` |
| Add | ➕ | `PlusCircle` |

## Conclusion

This design document provides a comprehensive blueprint for redesigning the IU-MiCert Issuer Dashboard. The focus is on creating a professional, modern interface suitable for educational institutions while maintaining all existing functionality and improving user experience through better visual hierarchy and interaction patterns.

The phased approach allows for incremental implementation and testing, ensuring stability throughout the redesign process.
