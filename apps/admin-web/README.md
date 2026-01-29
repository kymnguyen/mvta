# MVTA Admin Web

React-based admin interface for MVTA vehicle tracking system.

## Features

- **Vehicle Fleet View**: Grid layout showing all vehicles with key information
- **Vehicle Details**: Detailed view of individual vehicles
- **Change History**: Timeline view of all vehicle changes (location, status, mileage, fuel)
- **Real-time Updates**: React Query for efficient data fetching and caching
- **Responsive Design**: Works on desktop and mobile devices

## Tech Stack

- **React 18** - UI framework
- **TypeScript** - Type safety
- **React Router** - Navigation
- **TanStack Query (React Query)** - Data fetching and caching
- **Axios** - HTTP client
- **Vite** - Build tool
- **date-fns** - Date formatting

## Getting Started

### Install Dependencies

```bash
cd apps/admin-web
npm install
```

### Configure API Endpoint

The app proxies API requests to the tracking service. Update `vite.config.ts` if your tracking service runs on a different port:

```typescript
proxy: {
  '/api': {
    target: 'http://localhost:50002', // Change to your tracking-svc port
    changeOrigin: true
  }
}
```

### Run Development Server

```bash
npm run dev
```

The app will be available at `http://localhost:3000`

### Build for Production

```bash
npm run build
```

## Project Structure

```
src/
├── App.tsx                    # Main app component with routing
├── App.css                    # Global styles
├── main.tsx                   # Entry point
├── features/
│   └── vehicles/
│       ├── VehicleList.tsx    # Fleet view component
│       └── VehicleDetail.tsx  # Vehicle detail + history
└── shared/
    └── api/
        └── client.ts          # API client and TypeScript types
```

## API Integration

The app expects the tracking service to expose these endpoints:

- `GET /api/v1/vehicles` - List all vehicles
- `GET /api/v1/vehicles/:id` - Get vehicle details
- `GET /api/v1/vehicles/:id/history?limit=50&offset=0` - Get vehicle change history

## Features Overview

### Vehicle List Page
- Displays all vehicles in a responsive grid
- Shows key metrics: VIN, model, license, mileage, fuel level
- Status badges with color coding
- Click to view details

### Vehicle Detail Page
- Complete vehicle information
- Change history timeline with:
  - Change type indicators (created, location, status, mileage, fuel)
  - Timestamps formatted with date-fns
  - Before/after value comparison
  - Version tracking
- Visual timeline with markers

### Change Types Supported
- 🆕 Vehicle Created
- 📍 Location Updated
- 🔄 Status Changed
- 🚗 Mileage Updated
- ⛽ Fuel Updated

## Customization

### Styling
All styles are in `App.css`. Key CSS variables and colors:
- Primary: `#2c3e50`
- Accent: `#3498db`
- Success: `#27ae60`
- Danger: `#e74c3c`

### Query Configuration
Adjust React Query settings in `App.tsx`:
```typescript
const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      refetchOnWindowFocus: false,
      retry: 1,
      staleTime: 5000, // Add this for caching
    },
  },
});
```

## Development Tips

1. **Hot Reload**: Vite provides instant HMR for fast development
2. **Type Safety**: TypeScript interfaces match the Go backend DTOs
3. **Error Handling**: Add error boundaries for production
4. **Loading States**: Already handled with React Query's loading states
5. **Real-time Updates**: Can add WebSocket support for live vehicle tracking

## Next Steps

Potential enhancements:
- Add vehicle search and filtering
- Map view with vehicle locations
- Export history to CSV
- Vehicle statistics dashboard
- Real-time WebSocket updates
- Authentication/authorization
- Dark mode support
