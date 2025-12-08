import React from 'react';
import { View, Text, StyleSheet } from 'react-native';
import { useTheme } from '../../contexts/ThemeContext';
import { useNetworkStatus } from '../../hooks/useNetworkStatus';

export const OfflineIndicator: React.FC = () => {
  const { colors } = useTheme();
  const { isOffline } = useNetworkStatus();

  if (!isOffline) {
    return null;
  }

  return (
    <View
      style={[styles.container, { backgroundColor: colors.warning }]}
      accessibilityRole="alert"
      accessibilityLabel="You are offline"
    >
      <Text style={styles.text}>📵 You are offline</Text>
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    paddingVertical: 8,
    paddingHorizontal: 16,
    alignItems: 'center',
    justifyContent: 'center',
  },
  text: {
    color: '#FFFFFF',
    fontSize: 14,
    fontWeight: '600',
  },
});
