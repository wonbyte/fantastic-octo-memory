import { Stack } from 'expo-router';
import { QueryClient } from '@tanstack/react-query';
import { PersistQueryClientProvider } from '@tanstack/react-query-persist-client';
import AsyncStorage from '@react-native-async-storage/async-storage';
import { AuthProvider } from '../src/contexts/AuthContext';
import { ThemeProvider } from '../src/contexts/ThemeContext';

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      retry: 1,
      refetchOnWindowFocus: false,
      gcTime: 1000 * 60 * 60 * 24, // 24 hours
      staleTime: 1000 * 60 * 5, // 5 minutes
    },
  },
});

// Create async storage persister manually
const asyncStoragePersister = {
  persistClient: async (client: unknown) => {
    await AsyncStorage.setItem('REACT_QUERY_OFFLINE_CACHE', JSON.stringify(client));
  },
  restoreClient: async () => {
    const cached = await AsyncStorage.getItem('REACT_QUERY_OFFLINE_CACHE');
    return cached ? JSON.parse(cached) : undefined;
  },
  removeClient: async () => {
    await AsyncStorage.removeItem('REACT_QUERY_OFFLINE_CACHE');
  },
};

export default function RootLayout() {
  return (
    <PersistQueryClientProvider
      client={queryClient}
      persistOptions={{ persister: asyncStoragePersister }}
    >
      <ThemeProvider>
        <AuthProvider>
          <Stack screenOptions={{ headerShown: false }}>
            <Stack.Screen name="index" />
            <Stack.Screen name="(auth)" />
            <Stack.Screen name="(main)" />
          </Stack>
        </AuthProvider>
      </ThemeProvider>
    </PersistQueryClientProvider>
  );
}
