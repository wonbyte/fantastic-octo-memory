import { Tabs, router } from 'expo-router';
import { TouchableOpacity, Text } from 'react-native';
import { useAuth } from '../../src/contexts/AuthContext';

export default function MainLayout() {
  const { logout } = useAuth();

  const handleLogout = async () => {
    try {
      await logout();
      router.replace('/(auth)/login');
    } catch (error) {
      console.error('Logout failed:', error);
    }
  };

  return (
    <Tabs
      screenOptions={{
        headerRight: () => (
          <TouchableOpacity onPress={handleLogout} style={{ marginRight: 16 }}>
            <Text style={{ color: '#3B82F6', fontSize: 16 }}>Logout</Text>
          </TouchableOpacity>
        ),
      }}
    >
      <Tabs.Screen
        name="projects"
        options={{
          title: 'Projects',
          headerTitle: 'Projects',
        }}
      />
    </Tabs>
  );
}
