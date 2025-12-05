# Admin UI for Cost Database Management

This directory contains the admin interface for managing cost database pricing and overrides.

## Features

### Pricing Overrides (`pricing-overrides.tsx`)

Allows users to create, view, and delete custom pricing overrides for materials and labor rates.

**Features:**
- Create custom pricing for specific material or labor categories
- Support for absolute values or percentage adjustments
- Notes field for documenting pricing decisions
- Real-time updates with React Query
- Responsive mobile-friendly UI

**Usage:**
1. Navigate to Settings > Pricing Overrides
2. Click "Add New Override"
3. Select type (Material or Labor)
4. Enter the item key (e.g., "lumber", "carpentry")
5. Enter the override value
6. Optionally check "Is Percentage Adjustment"
7. Add notes if needed
8. Click "Create"

**API Integration:**
- GET `/api/company/pricing-overrides` - Fetch all overrides
- POST `/api/company/pricing-overrides` - Create new override
- DELETE `/api/company/pricing-overrides/:id` - Delete override

### Settings Dashboard (`../settings.tsx`)

Main settings landing page with navigation to:
- Pricing Overrides
- Regional Adjustments (planned)
- Cost Database Sync (planned)

## Architecture

The admin UI uses:
- **Expo Router** for navigation
- **TanStack Query** for data fetching and caching
- **Axios** for API calls
- **NativeWind** for styling (Tailwind CSS)

## Development

### Adding New Admin Features

1. Create a new screen in `app/(main)/settings/`
2. Add navigation link in `settings.tsx`
3. Implement API calls using React Query
4. Follow existing patterns for consistency

### Styling

Uses NativeWind (Tailwind CSS for React Native):
```tsx
<View className="bg-white p-4 rounded-lg shadow-sm">
  <Text className="text-lg font-bold">Title</Text>
</View>
```

### State Management

Uses TanStack Query for server state:
```tsx
const { data, isLoading, error } = useQuery({
  queryKey: ['key'],
  queryFn: async () => {
    const response = await api.get('/endpoint');
    return response.data;
  },
});
```

## Future Enhancements

- [ ] Regional Adjustments viewer
- [ ] Cost Database sync interface
- [ ] Bulk import/export for pricing overrides
- [ ] Price history and analytics
- [ ] Admin-only role enforcement
- [ ] Audit log for pricing changes
