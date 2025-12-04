import { Redirect } from 'expo-router';
import { useAuth } from '../src/contexts/AuthContext';
import { Loading } from '../src/components/ui/Loading';

export default function Index() {
  const { isAuthenticated, isLoading } = useAuth();

  if (isLoading) {
    return <Loading message="Loading..." />;
  }

  // Redirect to auth or main based on authentication status
  if (isAuthenticated) {
    return <Redirect href="/(main)/projects" />;
  }

  return <Redirect href="/(auth)/login" />;
}
